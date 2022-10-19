// Copyright 2022 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package ygnmi contains gNMI client library for use with a ygot Schema.
package ygnmi

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
	"google.golang.org/protobuf/encoding/prototext"

	log "github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// AnyQuery is a generic gNMI query for wildcard or non-wildcard state or config paths.
// Supported operations: Batch.
type AnyQuery[T any] interface {
	// PathStruct returns to path struct used for unmarshalling and schema validation.
	// This path must correspond to T (the parameterized type of the interface).
	PathStruct() PathStruct
	// subPaths contains the paths to subscribe to, they must be descendants of PathStruct().
	subPaths() []PathStruct
	// dirName returns the name of YANG directory schema entry.
	// For leaves, this is the parent entry.
	dirName() string
	// goStruct returns the struct that query should be unmarshalled into.
	// For leaves, this is the parent.
	goStruct() ygot.ValidatedGoStruct
	// extract is used for leaves to return the field from the parent GoStruct.
	// For non-leaves, this casts the GoStruct to the concrete type.
	extract(ygot.ValidatedGoStruct) (T, bool)
	// IsState returns if the path for this query is a state node.
	IsState() bool
	// isLeaf returns if the path for this query is a leaf.
	isLeaf() bool
	// isScalar returns whether the type (T) for this path is a pointer field (*T) in the parent GoStruct.
	isScalar() bool
	// schema returns the root schema used for unmarshalling.
	schema() *ytypes.Schema
}

// SingletonQuery is a non-wildcard gNMI query.
type SingletonQuery[T any] interface {
	AnyQuery[T]
	// IsSingleton restricts this interface to be used only where a singleton path is expected.
	IsSingleton()
}

// WildcardQuery is a wildcard gNMI query.
type WildcardQuery[T any] interface {
	AnyQuery[T]
	// IsWildcard restricts this interface to be used only where a wildcard path is expected.
	IsWildcard()
}

// ConfigQuery is a non-wildcard config gNMI query.
type ConfigQuery[T any] interface {
	AnyQuery[T]
	// IsConfig() allows this interface to be use in config funcs.
	IsConfig()
	// IsSingleton restricts this interface to be used only where a singleton path is expected.
	IsSingleton()
}

// Value contains a value received from a gNMI request and its metadata.
type Value[T any] struct {
	val     T
	present bool
	// Path is the sample's YANG path.
	Path *gpb.Path
	// Timestamp is the sample time.
	Timestamp time.Time
	// RecvTimestamp is the time the test received the sample.
	RecvTimestamp time.Time
	// ComplianceErrors contains the compliance errors encountered from an Unmarshal operation.
	ComplianceErrors *ComplianceErrors
}

// SetVal sets the value and marks it present and returns the receiver.
func (v *Value[T]) SetVal(val T) *Value[T] {
	v.present = true
	v.val = val
	return v
}

// Val returns the val and whether it is present.
func (v *Value[T]) Val() (T, bool) {
	return v.val, v.present
}

// IsPresent returns whether the value is present.
func (v *Value[T]) IsPresent() bool {
	return v.present
}

// String returns a user-readable string for the value.
func (v *Value[T]) String() string {
	path, err := ygot.PathToString(v.Path)
	if err != nil {
		path = v.Path.String()
	}
	val := fmt.Sprintf("%+v", v.val)
	if !v.present {
		val = "(not present)"
	}
	return fmt.Sprintf("path: %s\nvalue: %s", path, val)
}

// Client is used to perform gNMI requests.
type Client struct {
	gnmiC  gpb.GNMIClient
	target string
}

// ClientOption configures a client with custom options.
type ClientOption func(d *Client) error

// WithTarget sets the target of the gpb.Path for all requests made with this client.
func WithTarget(t string) ClientOption {
	return func(c *Client) error {
		c.target = t
		return nil
	}
}

// NewClient creates a new client with specified options.
func NewClient(c gpb.GNMIClient, opts ...ClientOption) (*Client, error) {
	yc := &Client{
		gnmiC: c,
	}
	for _, opt := range opts {
		if err := opt(yc); err != nil {
			return nil, err
		}
	}
	return yc, nil
}

// Option can be used modify the behavior of the gNMI requests used by the ygnmi calls (Lookup, Await, etc.).
type Option func(*opt)

type opt struct {
	useGet bool
	mode   gpb.SubscriptionMode
}

// resolveOpts applies all the options and returns a struct containing the result.
func resolveOpts(opts []Option) *opt {
	o := &opt{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithUseGet creates an option to use gnmi.Get instead of gnmi.Subscribe.
// This can only be used on Get, GetAll, Lookup, and LookupAll.
func WithUseGet() Option {
	return func(o *opt) {
		o.useGet = true
	}
}

// WithSubscriptionMode creates to an option to use input instead of the default (TARGET_DEFINED).
// This option is only relevant for Watch, WatchAll, Collect, CollectAll, Await which are STREAM subscriptions.
func WithSubscriptionMode(mode gpb.SubscriptionMode) Option {
	return func(o *opt) {
		o.mode = mode
	}
}

// Lookup fetches the value of a SingletonQuery with a ONCE subscription.
func Lookup[T any](ctx context.Context, c *Client, q SingletonQuery[T], opts ...Option) (*Value[T], error) {
	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_ONCE, resolveOpts(opts))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to path: %w", err)
	}
	data, err := receiveAll(sub, false)
	if err != nil {
		return nil, fmt.Errorf("failed to receive to data: %w", err)
	}
	val, err := unmarshalAndExtract[T](data, q, q.goStruct())
	if err != nil {
		return val, fmt.Errorf("failed to unmarshal data: %w", err)
	}
	if val.ComplianceErrors != nil {
		if q.isLeaf() {
			return val, fmt.Errorf("noncompliant data encountered while unmarshalling leaf: %v", val.ComplianceErrors)
		}
		log.V(0).Infof("noncompliant data encountered while unmarshalling: %v", val.ComplianceErrors)
	}
	return val, nil
}

// LookupConfig fetches the value of a ConfigQuery with a ONCE subscription.
// Note: This is a workaround for Go's type inference not working for this use case and may be removed in a subsequent release.
func LookupConfig[T any](ctx context.Context, c *Client, q ConfigQuery[T], opts ...Option) (*Value[T], error) {
	return Lookup[T](ctx, c, q, opts...)
}

var (
	// ErrNotPresent is returned by Get when there are no a values at a path.
	ErrNotPresent = fmt.Errorf("value not present")
	// Continue should returned by predicates to indicate the condition is not reached.
	Continue = fmt.Errorf("condition not true")
)

// Get fetches the value of a SingletonQuery with a ONCE subscription,
// returning an error that wraps ErrNotPresent if the value is not present.
// Use Lookup to get metadata and tolerate non-present data.
func Get[T any](ctx context.Context, c *Client, q SingletonQuery[T], opts ...Option) (T, error) {
	var zero T
	val, err := Lookup(ctx, c, q, opts...)
	if err != nil {
		return zero, err
	}
	ret, ok := val.Val()
	if !ok {
		return zero, fmt.Errorf("path %s: %w", val.Path.String(), ErrNotPresent)
	}
	return ret, nil
}

// GetConfig fetches the value of a ConfigQuery with a ONCE subscription,
// returning an error that wraps ErrNotPresent if the value is not present.
// Use LookupConfig to get metadata and tolerate non-present data.
// Note: This is a workaround for Go's type inference not working for this use case and may be removed in a subsequent release.
func GetConfig[T any](ctx context.Context, c *Client, q ConfigQuery[T], opts ...Option) (T, error) {
	return Get[T](ctx, c, q, opts...)
}

// Watcher represents an ongoing watch of telemetry values.
type Watcher[T any] struct {
	errCh   chan error
	lastVal *Value[T]
}

// Await waits for the watch to finish and returns the last received value
// and a boolean indicating whether the predicate evaluated to true.
// When Await returns the watcher is closed, and Await may not be called again.
func (w *Watcher[T]) Await() (*Value[T], error) {
	err, ok := <-w.errCh
	if !ok {
		return nil, fmt.Errorf("Await already called and Watcher is closed")
	}
	close(w.errCh)
	return w.lastVal, err
}

// Watch starts an asynchronous STREAM subscription, evaluating each observed value with the specified predicate.
// The predicate must return ygnmi.Continue to continue the Watch. To stop the Watch, return nil for a success
// or a non-nil error on failure. Watch can also be stopped by setting a deadline on or canceling the context.
// Calling Await on the returned Watcher waits for the subscription to complete.
// It returns the last observed value and a boolean that indicates whether that value satisfies the predicate.
func Watch[T any](ctx context.Context, c *Client, q SingletonQuery[T], pred func(*Value[T]) error, opts ...Option) *Watcher[T] {
	w := &Watcher[T]{
		errCh: make(chan error, 1),
	}

	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_STREAM, resolveOpts(opts))
	if err != nil {
		w.errCh <- err
		return w
	}

	dataCh, errCh := receiveStream[T](sub, q)
	go func() {
		// Create an intially empty GoStruct, into which all received datapoints will be unmarshalled.
		gs := q.goStruct()
		for {
			select {
			case data := <-dataCh:
				val, err := unmarshalAndExtract[T](data, q, gs)
				if err != nil {
					w.errCh <- err
					return
				}
				w.lastVal = val
				if err := pred(val); err == nil || !errors.Is(err, Continue) {
					w.errCh <- err
					return
				}
			case err := <-errCh:
				w.errCh <- err
				return
			}
		}
	}()
	return w
}

// Await observes values at Query with a STREAM subscription,
// blocking until a value that is deep equal to the specified val is received
// or the context is cancelled. To wait for a generic predicate, or to make a
// non-blocking call, use the Watch method instead.
func Await[T any](ctx context.Context, c *Client, q SingletonQuery[T], val T, opts ...Option) (*Value[T], error) {
	w := Watch(ctx, c, q, func(v *Value[T]) error {
		if v.present && reflect.DeepEqual(v.val, val) {
			return nil
		}
		return Continue
	}, opts...)
	return w.Await()
}

// Collector represents an ongoing collection of telemetry values.
type Collector[T any] struct {
	w    *Watcher[T]
	data []*Value[T]
}

// Await waits for the collection to finish and returns all received values.
// When Await returns the watcher is closed, and Await may not be called again.
// Note: the func blocks until the context is cancelled.
func (c *Collector[T]) Await() ([]*Value[T], error) {
	_, err := c.w.Await()
	return c.data, err
}

// Collect starts an asynchronous collection of the values at the query with a STREAM subscription.
// Calling Await on the return Collection waits until the context is cancelled and returns the collected values.
func Collect[T any](ctx context.Context, c *Client, q SingletonQuery[T], opts ...Option) *Collector[T] {
	collect := &Collector[T]{}
	collect.w = Watch(ctx, c, q, func(v *Value[T]) error {
		if !q.isLeaf() {
			// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#why-not-permit-type-assertions-on-values-whose-type-is-a-type-parameter
			gs, err := ygot.DeepCopy((interface{})(v.val).(ygot.GoStruct))
			if err != nil {
				return err
			}
			v.SetVal(gs.(T))
		}
		collect.data = append(collect.data, v)
		return Continue
	}, opts...)
	return collect
}

// LookupAll fetches the values of a WildcardQuery with a ONCE subscription.
// It returns an empty list if no values are present at the path.
func LookupAll[T any](ctx context.Context, c *Client, q WildcardQuery[T], opts ...Option) ([]*Value[T], error) {
	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_ONCE, resolveOpts(opts))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to path: %w", err)
	}
	data, err := receiveAll(sub, false)
	if err != nil {
		return nil, fmt.Errorf("failed to receive to data: %w", err)
	}
	p, err := resolvePath(q.PathStruct())
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	datapointGroups, sortedPrefixes, err := bundleDatapoints(data, len(p.Elem))
	if err != nil {
		return nil, fmt.Errorf("failed to bundle datapoints: %w", err)
	}
	var vals []*Value[T]
	for _, prefix := range sortedPrefixes {
		goStruct := q.goStruct()
		v, err := unmarshalAndExtract[T](datapointGroups[prefix], q, goStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}
		if v.ComplianceErrors != nil {
			log.V(0).Infof("noncompliant data encountered while unmarshalling: %v", v.ComplianceErrors)
			continue
		}
		vals = append(vals, v)
	}
	return vals, nil
}

// GetAll fetches the value of a WildcardQuery with a ONCE subscription skipping any non-present paths.
// It returns an error that wraps ErrNotPresent if no values were received.
// Use LookupAll to also get metadata containing the returned paths.
func GetAll[T any](ctx context.Context, c *Client, q WildcardQuery[T], opts ...Option) ([]T, error) {
	vals, err := LookupAll(ctx, c, q, opts...)
	if err != nil {
		return nil, err
	}
	ret := make([]T, 0, len(vals))
	for _, val := range vals {
		if v, ok := val.Val(); ok {
			ret = append(ret, v)
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("query %q: %w", q, ErrNotPresent)
	}
	return ret, nil
}

// WatchAll starts an asynchronous STREAM subscription, evaluating each observed value with the specified predicate.
// The predicate must return ygnmi.Continue to continue the Watch. To stop the Watch, return nil for a success
// or a non-nil error on failure. Watch can also be stopped by setting a deadline on or canceling the context.
// Calling Await on the returned Watcher waits for the subscription to complete.
// It returns the last observed value and a boolean that indicates whether that value satisfies the predicate.
func WatchAll[T any](ctx context.Context, c *Client, q WildcardQuery[T], pred func(*Value[T]) error, opts ...Option) *Watcher[T] {
	w := &Watcher[T]{
		errCh: make(chan error, 1),
	}
	path, err := resolvePath(q.PathStruct())
	if err != nil {
		w.errCh <- err
		return w
	}
	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_STREAM, resolveOpts(opts))
	if err != nil {
		w.errCh <- err
		return w
	}

	dataCh, errCh := receiveStream[T](sub, q)
	go func() {
		// Create a map intially empty GoStruct, into which all received datapoints will be unmarshalled based on their path prefixes.
		structs := map[string]ygot.ValidatedGoStruct{}
		for {
			select {
			case data := <-dataCh:
				datapointGroups, sortedPrefixes, err := bundleDatapoints(data, len(path.Elem))
				if err != nil {
					w.errCh <- err
					return
				}
				for _, pre := range sortedPrefixes {
					if len(datapointGroups[pre]) == 0 {
						continue
					}
					if _, ok := structs[pre]; !ok {
						structs[pre] = q.goStruct()
					}
					val, err := unmarshalAndExtract[T](data, q, structs[pre])
					if err != nil {
						w.errCh <- err
						return
					}
					w.lastVal = val
					if err := pred(val); err == nil || !errors.Is(err, Continue) {
						w.errCh <- err
						return
					}
				}
			case err := <-errCh:
				w.errCh <- err
				return
			}
		}
	}()
	return w
}

// CollectAll starts an asynchronous collection of the values at the query with a STREAM subscription.
// Calling Await on the return Collection waits until the context is cancelled to elapse and returns the collected values.
func CollectAll[T any](ctx context.Context, c *Client, q WildcardQuery[T], opts ...Option) *Collector[T] {
	collect := &Collector[T]{}
	collect.w = WatchAll(ctx, c, q, func(v *Value[T]) error {
		if !q.isLeaf() {
			// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#why-not-permit-type-assertions-on-values-whose-type-is-a-type-parameter
			gs, err := ygot.DeepCopy((interface{})(v.val).(ygot.GoStruct))
			if err != nil {
				return err
			}
			v.SetVal(gs.(T))
		}
		collect.data = append(collect.data, v)
		return Continue
	}, opts...)
	return collect
}

// Result is the result of a Set request.
type Result struct {
	// RawResponse is the raw gNMI response received from the server.
	RawResponse *gpb.SetResponse
	// Timestamp is the timestamp from the SetResponse as a native Go time struct.
	Timestamp time.Time
}

func responseToResult(resp *gpb.SetResponse) *Result {
	return &Result{
		RawResponse: resp,
		Timestamp:   time.Unix(0, resp.GetTimestamp()),
	}
}

// Update updates the configuration at the given query path with the val.
func Update[T any](ctx context.Context, c *Client, q ConfigQuery[T], val T) (*Result, error) {
	resp, path, err := set(ctx, c, q, val, updatePath)
	if err != nil {
		return nil, fmt.Errorf("Update(t) at path %s: %w", path, err)
	}
	return responseToResult(resp), nil
}

// Replace replaces the configuration at the given query path with the val.
func Replace[T any](ctx context.Context, c *Client, q ConfigQuery[T], val T) (*Result, error) {
	resp, path, err := set(ctx, c, q, val, replacePath)
	if err != nil {
		return nil, fmt.Errorf("Replace(t) at path %s: %w", path, err)
	}
	return responseToResult(resp), nil
}

// Delete deletes the configuration at the given query path.
func Delete[T any](ctx context.Context, c *Client, q ConfigQuery[T]) (*Result, error) {
	var t T
	resp, path, err := set(ctx, c, q, t, deletePath)
	if err != nil {
		return nil, fmt.Errorf("Delete(t) at path %s: %w", path, err)
	}
	return responseToResult(resp), nil
}

type batchOp struct {
	path   PathStruct
	val    interface{}
	mode   setOperation
	config bool
}

// SetBatch allows multiple Set operations (Replace, Update, Delete) to be applied as part of a single Set transaction.
// Use BatchUpdate, BatchReplace, BatchDelete to add operations, and then call the Set method to send the SetRequest.
type SetBatch struct {
	ops []*batchOp
}

// Set performs the gnmi.Set request with all queued operations.
func (sb *SetBatch) Set(ctx context.Context, c *Client) (*Result, error) {
	req := &gpb.SetRequest{}
	for _, op := range sb.ops {
		path, err := resolvePath(op.path)
		if err != nil {
			return nil, err
		}
		if err := populateSetRequest(req, path, op.val, op.mode, op.config); err != nil {
			return nil, err
		}
	}
	req.Prefix = &gpb.Path{
		Target: c.target,
	}
	log.V(1).Info(prettySetRequest(req))
	resp, err := c.gnmiC.Set(ctx, req)
	log.V(1).Infof("SetResponse:\n%s", prototext.Format(resp))
	return responseToResult(resp), err
}

// BatchUpdate stores an update operation in the SetBatch.
func BatchUpdate[T any](sb *SetBatch, q ConfigQuery[T], val T) {
	var setVal interface{} = val
	if q.isLeaf() && q.isScalar() {
		setVal = &val
	}
	sb.ops = append(sb.ops, &batchOp{
		path:   q.PathStruct(),
		val:    setVal,
		mode:   updatePath,
		config: !q.IsState(),
	})
}

// BatchReplace stores an update operation in the SetBatch.
func BatchReplace[T any](sb *SetBatch, q ConfigQuery[T], val T) {
	var setVal interface{} = val
	if q.isLeaf() && q.isScalar() {
		setVal = &val
	}
	sb.ops = append(sb.ops, &batchOp{
		path:   q.PathStruct(),
		val:    setVal,
		mode:   replacePath,
		config: !q.IsState(),
	})
}

// BatchDelete stores an update operation in the SetBatch.
func BatchDelete[T any](sb *SetBatch, q ConfigQuery[T]) {
	sb.ops = append(sb.ops, &batchOp{
		path:   q.PathStruct(),
		mode:   deletePath,
		config: !q.IsState(),
	})
}

// Batch contains a collection of paths.
// Calling State() or Config() on the batch returns a query
// that can be used to Lookup, Watch, etc on multiple paths at once.
type Batch[T ygot.ValidatedGoStruct] struct {
	root  SingletonQuery[T]
	paths []PathStruct
}

// NewBatch creates a batch object. All paths in the batch must be children of the root query.
func NewBatch[T ygot.ValidatedGoStruct](root SingletonQuery[T]) *Batch[T] {
	return &Batch[T]{
		root: root,
	}
}

// AddPaths adds the paths to the batch. Paths must be children of the root.
func (b *Batch[T]) AddPaths(paths ...PathStruct) error {
	root, _, err := ResolvePath(b.root.PathStruct())
	if err != nil {
		return err
	}
	for _, path := range paths {
		p, _, err := ResolvePath(path)
		if err != nil {
			return err
		}
		if !util.PathMatchesQuery(p, root) {
			return fmt.Errorf("root path %v is not a prefix of %v", root, p)
		}
	}
	b.paths = append(b.paths, paths...)
	return nil
}

// Query returns a Query that can be used in gNMI operations.
// The returned query is immutable, adding paths does not modify existing queries.
func (b *Batch[T]) Query() SingletonQuery[T] {
	queryPaths := make([]PathStruct, len(b.paths))
	copy(queryPaths, b.paths)
	return NewNonLeafSingletonQuery[T](
		b.root.dirName(),
		b.root.IsState(),
		b.root.PathStruct(),
		queryPaths,
		b.root.schema(),
	)
}

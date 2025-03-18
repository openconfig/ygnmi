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

	"github.com/openconfig/ygnmi/internal/logutil"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	log "github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// UntypedQuery is a generic gNMI query for wildcard or non-wildcard state or config paths.
// Supported operations: Batch.
//
// Since it is untyped, it is only suitable for representing the metadata of
// the query instead of initiating gNMI operations.
type UntypedQuery interface {
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
	// IsState returns if the path for this query is a state node.
	IsState() bool
	// isShadowPath returns if the path for this query is a shadow-path node.
	isShadowPath() bool
	// isLeaf returns if the path for this query is a leaf.
	isLeaf() bool
	// isCompressedSchema returns whether the query is for compressed ygot schema.
	isCompressedSchema() bool
	// isListContainer returns if the path for this query is for a whole list.
	isListContainer() bool
	// isScalar returns whether the type (T) for this path is a pointer field (*T) in the parent GoStruct.
	isScalar() bool
	// schema returns the root schema used for unmarshalling.
	schema() *ytypes.Schema
	// compressInfo returns information about where the path points to an
	// element that is compressed out (if applicable).
	compressInfo() *CompressionInfo
}

// AnyQuery is a generic gNMI query for wildcard or non-wildcard state or config paths.
type AnyQuery[T any] interface {
	UntypedQuery
	// extract is used for leaves to return the field from the parent GoStruct.
	// For non-leaves, this casts the GoStruct to the concrete type.
	extract(ygot.ValidatedGoStruct) (T, bool)
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
	val := "(not present)"
	if v.present {
		val = fmt.Sprintf("%+v", v.val)
	}
	return fmt.Sprintf("path: %s\nvalue: %s", path, val)
}

// Client is used to perform gNMI requests.
type Client struct {
	gnmiC           gpb.GNMIClient
	target          string
	requestLogLevel log.Level
}

// String returns a string representation of Client. This output is unstable.
func (c *Client) String() string {
	return fmt.Sprintf("ygnmi client (target: %q)", c.target)
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

// WithRequestLogLevel overrides the default info logging level (1) for gNMI
// request dumps. i.e. SetRequest, SubscribeRequest, GetRequest.
func WithRequestLogLevel(l log.Level) ClientOption {
	return func(c *Client) error {
		c.requestLogLevel = l
		return nil
	}
}

// NewClient creates a new client with specified options.
func NewClient(c gpb.GNMIClient, opts ...ClientOption) (*Client, error) {
	yc := &Client{
		gnmiC:           c,
		requestLogLevel: 1,
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
	useGet         bool
	mode           gpb.SubscriptionMode
	encoding       gpb.Encoding
	preferProto    bool
	setFallback    bool
	sampleInterval uint64
}

// resolveOpts applies all the options and returns a struct containing the result.
func resolveOpts(opts []Option) *opt {
	o := &opt{
		encoding: gpb.Encoding_PROTO,
	}
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

// WithSubscriptionMode creates an option to use input instead of the default (TARGET_DEFINED).
// This option is only relevant for Watch, WatchAll, Collect, CollectAll, Await which are STREAM subscriptions.
// The mode applies to all paths in the Subcription.
func WithSubscriptionMode(mode gpb.SubscriptionMode) Option {
	return func(o *opt) {
		o.mode = mode
	}
}

// WithSampleInterval creates an option to set the sample interval in the Subcribe request.
// NOTE: The subscription mode must be set to SAMPLE for this to have an effect.
// This option is only relevant for Watch, WatchAll, Collect, CollectAll, Await which are STREAM subscriptions.
// The mode applies to all paths in the Subcription.
func WithSampleInterval(d time.Duration) Option {
	return func(o *opt) {
		o.sampleInterval = uint64(d.Nanoseconds())
	}
}

// WithEncoding creates an option to set the Encoding for all Subscribe or Get requests.
// The default encoding is PROTO. This does not apply when using WithUseGet, whith uses JSON_IETF encoding.
func WithEncoding(enc gpb.Encoding) Option {
	return func(o *opt) {
		o.encoding = enc
	}
}

// WithSetPreferProtoEncoding creates an option that prefers encoding SetRequest using proto encoding.
// This option is only relevant for Update and Replace.
func WithSetPreferProtoEncoding() Option {
	return func(o *opt) {
		o.preferProto = true
	}
}

// WithSetFallbackEncoding creates an option that fallback to encoding SetRequests with JSON or an Any proto.
// Fallback encoding is if the parameter is neither GoStruct nor a leaf for non OpenConfig paths.
// This option is only relevant for Update and Replace.
func WithSetFallbackEncoding() Option {
	return func(o *opt) {
		o.setFallback = true
	}
}

// Lookup fetches the value of a SingletonQuery with a ONCE subscription.
func Lookup[T any](ctx context.Context, c *Client, q SingletonQuery[T], opts ...Option) (*Value[T], error) {
	resolvedOpts := resolveOpts(opts)
	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_ONCE, resolvedOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to path: %w", err)
	}
	data, err := receiveAll(sub, false)
	if err != nil {
		return nil, fmt.Errorf("failed to receive to data: %w", err)
	}
	val, err := unmarshalAndExtract[T](data, q, q.goStruct(), resolvedOpts)
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
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	w := &Watcher[T]{
		errCh: make(chan error, 1),
	}

	resolvedOpts := resolveOpts(opts)
	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_STREAM, resolvedOpts)
	if err != nil {
		cancel()
		w.errCh <- err
		return w
	}

	dataCh, errCh := receiveStream[T](ctx, sub, q)
	go func() {
		defer cancel()
		// Create an intially empty GoStruct, into which all received datapoints will be unmarshalled.
		gs := q.goStruct()
		for {
			select {
			case <-ctx.Done():
				w.errCh <- ctx.Err()
				return
			case data := <-dataCh:
				val, err := unmarshalAndExtract[T](data, q, gs, resolvedOpts)
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
	resolvedOpts := resolveOpts(opts)
	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_ONCE, resolvedOpts)
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
		v, err := unmarshalAndExtract[T](datapointGroups[prefix], q, goStruct, resolvedOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}
		if v.ComplianceErrors != nil {
			log.V(0).Infof("noncompliant data encountered while unmarshalling: %v", v.ComplianceErrors)
			if q.isLeaf() {
				continue
			}
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
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	w := &Watcher[T]{
		errCh: make(chan error, 1),
	}
	path, err := resolvePath(q.PathStruct())
	if err != nil {
		cancel()
		w.errCh <- err
		return w
	}
	resolvedOpts := resolveOpts(opts)
	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_STREAM, resolvedOpts)
	if err != nil {
		cancel()
		w.errCh <- err
		return w
	}

	dataCh, errCh := receiveStream[T](ctx, sub, q)
	go func() {
		defer cancel()
		// Create a map intially empty GoStruct, into which all received datapoints will be unmarshalled based on their path prefixes.
		structs := map[string]ygot.ValidatedGoStruct{}
		for {
			select {
			case <-ctx.Done():
				w.errCh <- ctx.Err()
				return
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
					val, err := unmarshalAndExtract[T](datapointGroups[pre], q, structs[pre], resolvedOpts)
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
func Update[T any](ctx context.Context, c *Client, q ConfigQuery[T], val T, opts ...Option) (*Result, error) {
	resp, path, err := set(ctx, c, q, val, updatePath, opts...)
	if err != nil {
		return nil, fmt.Errorf("Update(t) at path %s: %w", path, err)
	}
	return responseToResult(resp), nil
}

// Replace replaces the configuration at the given query path with the val.
func Replace[T any](ctx context.Context, c *Client, q ConfigQuery[T], val T, opts ...Option) (*Result, error) {
	resp, path, err := set(ctx, c, q, val, replacePath, opts...)
	if err != nil {
		return nil, fmt.Errorf("Replace(t) at path %s: %w", path, err)
	}
	return responseToResult(resp), nil
}

// Delete deletes the configuration at the given query path.
func Delete[T any](ctx context.Context, c *Client, q ConfigQuery[T], opts ...Option) (*Result, error) {
	var t T
	resp, path, err := set(ctx, c, q, t, deletePath, opts...)
	if err != nil {
		return nil, fmt.Errorf("Delete(t) at path %s: %w", path, err)
	}
	return responseToResult(resp), nil
}

type batchOp struct {
	path         PathStruct
	val          interface{}
	mode         setOperation
	shadowpath   bool
	isLeaf       bool
	compressInfo *CompressionInfo
}

// SetBatch allows multiple Set operations (Replace, Update, Delete) to be applied as part of a single Set transaction.
// Use BatchUpdate, BatchReplace, BatchDelete to add operations, and then call the Set method to send the SetRequest.
type SetBatch struct {
	ops []*batchOp
}

// Set performs the gnmi.Set request with all queued operations.
func (sb *SetBatch) Set(ctx context.Context, c *Client, opts ...Option) (*Result, error) {
	req := &gpb.SetRequest{}
	for _, op := range sb.ops {
		path, err := resolvePath(op.path)
		if err != nil {
			return nil, err
		}
		if err := populateSetRequest(req, path, op.val, op.mode, op.shadowpath, op.isLeaf, op.compressInfo, opts...); err != nil {
			return nil, err
		}
	}
	req.Prefix = &gpb.Path{
		Target: c.target,
	}
	logutil.LogByLine(c.requestLogLevel, prettySetRequest(req))
	resp, err := c.gnmiC.Set(ctx, req)
	log.V(c.requestLogLevel).Infof("SetResponse:\n%s", prototext.Format(resp))
	return responseToResult(resp), err
}

// BatchUpdate stores an update operation in the SetBatch.
func BatchUpdate[T any](sb *SetBatch, q ConfigQuery[T], val T) {
	var setVal interface{} = val
	if q.isLeaf() && q.isScalar() {
		setVal = &val
	}
	sb.ops = append(sb.ops, &batchOp{
		path:         q.PathStruct(),
		val:          setVal,
		mode:         updatePath,
		shadowpath:   q.isShadowPath(),
		isLeaf:       q.isLeaf(),
		compressInfo: q.compressInfo(),
	})
}

// BatchReplace stores a replace operation in the SetBatch.
func BatchReplace[T any](sb *SetBatch, q ConfigQuery[T], val T) {
	var setVal interface{} = val
	if q.isLeaf() && q.isScalar() {
		setVal = &val
	}
	sb.ops = append(sb.ops, &batchOp{
		path:         q.PathStruct(),
		val:          setVal,
		mode:         replacePath,
		shadowpath:   q.isShadowPath(),
		isLeaf:       q.isLeaf(),
		compressInfo: q.compressInfo(),
	})
}

// BatchUnionReplace stores a union_replace operation in the SetBatch.
//
// https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-union_replace.md
func BatchUnionReplace[T any](sb *SetBatch, q ConfigQuery[T], val T) {
	var setVal interface{} = val
	if q.isLeaf() && q.isScalar() {
		setVal = &val
	}
	sb.ops = append(sb.ops, &batchOp{
		path:         q.PathStruct(),
		val:          setVal,
		mode:         unionreplacePath,
		shadowpath:   q.isShadowPath(),
		isLeaf:       q.isLeaf(),
		compressInfo: q.compressInfo(),
	})
}

// BatchUnionReplaceCLI stores a CLI union_replace operation in the SetBatch.
//
//   - nos is the name of the Network operating system.
//     "_cli" is appended to it to form the origin, see
//     https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-union_replace.md#24-native-cli-configuration-cli
//   - ascii is the full CLI text.
//
// https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-union_replace.md
func BatchUnionReplaceCLI(sb *SetBatch, nos, ascii string) {
	ps := NewDeviceRootBase()
	ps.PutCustomData(OriginOverride, nos+"_cli")
	sb.ops = append(sb.ops, &batchOp{
		path:         ps,
		val:          ascii,
		mode:         unionreplacePath,
		shadowpath:   false,
		compressInfo: nil,
	})
}

// BatchDelete stores a delete operation in the SetBatch.
func BatchDelete[T any](sb *SetBatch, q ConfigQuery[T]) {
	sb.ops = append(sb.ops, &batchOp{
		path:         q.PathStruct(),
		mode:         deletePath,
		shadowpath:   q.isShadowPath(),
		isLeaf:       q.isLeaf(),
		compressInfo: q.compressInfo(),
	})
}

// Batch contains a collection of paths.
// Calling State() or Config() on the batch returns a query
// that can be used to Lookup, Watch, etc on multiple paths at once.
type Batch[T any] struct {
	root  SingletonQuery[T]
	paths []PathStruct
}

// NewBatch creates a batch object. All paths in the batch must be children of the root query.
func NewBatch[T any](root SingletonQuery[T]) *Batch[T] {
	return &Batch[T]{
		root: root,
	}
}

// AddPaths adds the paths to the batch. Paths must be children of the root.
func (b *Batch[T]) AddPaths(paths ...UntypedQuery) error {
	root, err := resolvePath(b.root.PathStruct())
	if err != nil {
		return err
	}
	var pathstructs []PathStruct
	for _, path := range paths {
		ps := path.PathStruct()
		pathstructs = append(pathstructs, ps)
		p, err := resolvePath(ps)
		if err != nil {
			return err
		}
		if !util.PathMatchesQuery(p, root) {
			return fmt.Errorf("root path %v is not a prefix of %v", root, p)
		}
	}
	b.paths = append(b.paths, pathstructs...)
	return nil
}

// Query returns a Query that can be used in gNMI operations.
// The returned query is immutable, adding paths does not modify existing queries.
func (b *Batch[T]) Query() SingletonQuery[T] {
	queryPaths := make([]PathStruct, len(b.paths))
	copy(queryPaths, b.paths)
	return NewSingletonQuery[T](
		b.root.dirName(),
		b.root.IsState(),
		b.root.isShadowPath(),
		b.root.isLeaf(),
		b.root.isScalar(),
		b.root.isCompressedSchema(),
		b.root.isListContainer(),
		b.root.PathStruct(),
		b.root.extract,
		b.root.goStruct,
		b.root.schema,
		queryPaths,
		b.root.compressInfo(),
	)
}

// WildcardBatch contains a collection of paths.
// Calling Query() on the batch returns a query
// that can be used in LookupAll, WatchAll, etc on select paths within the root path.
type WildcardBatch[T any] struct {
	root  WildcardQuery[T]
	paths []PathStruct
}

// NewWildcardBatch creates a batch object. All paths in the batch must be children of the root query.
func NewWildcardBatch[T any](root WildcardQuery[T]) *WildcardBatch[T] {
	return &WildcardBatch[T]{
		root: root,
	}
}

// AddPaths adds the paths to the batch. Paths must be children of the root.
func (b *WildcardBatch[T]) AddPaths(paths ...UntypedQuery) error {
	root, err := resolvePath(b.root.PathStruct())
	if err != nil {
		return err
	}
	var pathstructs []PathStruct
	for _, path := range paths {
		ps := path.PathStruct()
		pathstructs = append(pathstructs, ps)
		p, err := resolvePath(ps)
		if err != nil {
			return err
		}
		if !util.PathMatchesQuery(p, root) {
			return fmt.Errorf("root path %v is not a prefix of %v", root, p)
		}
	}
	b.paths = append(b.paths, pathstructs...)
	return nil
}

// Query returns a Query that can be used in gNMI operations.
// The returned query is immutable, adding paths does not modify existing queries.
func (b *WildcardBatch[T]) Query() WildcardQuery[T] {
	queryPaths := make([]PathStruct, len(b.paths))
	copy(queryPaths, b.paths)
	return NewWildcardQuery[T](
		b.root.dirName(),
		b.root.IsState(),
		b.root.isShadowPath(),
		b.root.isLeaf(),
		b.root.isScalar(),
		b.root.isCompressedSchema(),
		b.root.isListContainer(),
		b.root.PathStruct(),
		b.root.extract,
		b.root.goStruct,
		b.root.schema,
		queryPaths,
		b.root.compressInfo(),
	)
}

// NewReconciler creates a new reconciler. See Reconciler type documentation for information.
func NewReconciler[T ygot.GoStruct](c *Client, q ConfigQuery[T], opts ...Option) (*Reconciler[T], error) {
	state := configToState(q.(AnyQuery[T]))
	return &Reconciler[T]{
		rootCfg:   q,
		rootState: state,
		c:         c,
		opts:      opts,
		errCh:     make(chan error),
	}, nil
}

// ReconcilerAbortErr stops the reconcilation loop.
var ReconcilerAbortErr = fmt.Errorf("reconciler abort")

// Reconciler subscribes to a non-leaf config query, updates to the path invoke a callback function.
// The callback function accepts a GoStruct for the config and state types for the root query.
// This intended for writing a function that watches some config, checks whether config == state, and makes changes to some system until they converge.
// Reconcilers contain a root callback function invoked on every Update, and sub reconcilers that are only invoked for specific paths under the root.
// Reconciler doesn't stop until the context is cancelled or  ReconcilerAbortErr is returned by any callback func.
type Reconciler[T ygot.GoStruct] struct {
	opts      []Option
	c         *Client
	rootCfg   AnyQuery[T]
	rootState AnyQuery[T]
	errCh     chan error
	subRecs   []*subRec[T]
}

type subRec[T ygot.GoStruct] struct {
	cfg   *gpb.Path
	state *gpb.Path
	fn    func(cfg *Value[T], state *Value[T]) error
}

// AddSubReconciler adds a sub reconciler to the main reconciler. The callback function is only invokes when gNMI update matching the query are received.
// The query must be a child of the root path.
func (r *Reconciler[T]) AddSubReconciler(q UntypedQuery, fn func(cfg *Value[T], state *Value[T]) error) error {
	rootCfgPath, err := resolvePath(r.rootCfg.PathStruct())
	if err != nil {
		return err
	}

	cfgPath, err := resolvePath(q.PathStruct())
	if err != nil {
		return err
	}

	state := configToStatePS(q.PathStruct())
	statePath, err := resolvePath(state)
	if err != nil {
		return err
	}
	if !util.PathMatchesQuery(cfgPath, rootCfgPath) {
		return fmt.Errorf("root path %v is not a prefix of %v", rootCfgPath, cfgPath)
	}

	r.subRecs = append(r.subRecs, &subRec[T]{
		cfg:   cfgPath,
		state: statePath,
		fn:    fn,
	})
	return nil
}

// Start starts the reconciler.
func (r *Reconciler[T]) Start(ctx context.Context, fn func(cfg *Value[T], state *Value[T]) error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)

	resolvedOpts := resolveOpts(r.opts)
	sub, err := subscribe(ctx, r.c, r.rootCfg, gpb.SubscriptionList_STREAM, resolvedOpts)
	if err != nil {
		cancel()
		r.errCh <- err
		return
	}

	dataCh, errCh := receiveStream(ctx, sub, r.rootCfg)
	go func() {
		defer cancel()
		// Create an intially empty GoStruct, into which all received datapoints will be unmarshalled.
		cfg := r.rootCfg.goStruct()
		state := r.rootState.goStruct()
		for {
			select {
			case <-ctx.Done():
				r.errCh <- ctx.Err()
				return
			case data := <-dataCh:
				cfgVal, err := unmarshalAndExtract(data, r.rootCfg, cfg, resolvedOpts)
				if err != nil {
					r.errCh <- err
					return
				}
				stateVal, err := unmarshalAndExtract(data, r.rootState, state, resolvedOpts)
				if err != nil {
					r.errCh <- err
					return
				}

				// Pair up the datapoint so updates to /foo/state/value and /foo/config/value only call the predicate once.
				configPoints := map[string]*DataPoint{}
				statePoints := map[string]*DataPoint{}
				cfgToStatePaths := map[string]string{}
				for _, p := range data {
					if p.Sync {
						continue
					}
					pathStr, err := ygot.PathToString(p.Path)
					if err != nil {
						r.errCh <- err
						return
					}
					swapPath := swapConfigStatePath(p.Path)
					swapPathStr, err := ygot.PathToString(swapPath)
					if err != nil {
						r.errCh <- err
						return
					}
					if p.Path.Elem[len(p.Path.Elem)-2].Name == "config" {
						configPoints[pathStr] = p
						cfgToStatePaths[pathStr] = swapPathStr
					} else {
						statePoints[pathStr] = p
					}
				}

				for _, sr := range r.subRecs {
					for cfgPath, cfgPoint := range configPoints {
						if !util.PathMatchesQuery(cfgPoint.Path, sr.cfg) {
							continue
						}
						cfgVal := &Value[T]{
							val:              cfgVal.val,
							present:          cfgVal.present,
							Path:             cfgPoint.Path, // Use the datapoint path, (leaf)
							Timestamp:        cfgVal.Timestamp,
							RecvTimestamp:    cfgVal.RecvTimestamp,
							ComplianceErrors: cfgVal.ComplianceErrors,
						}
						stateVal := &Value[T]{
							val:              stateVal.val,
							present:          stateVal.present,
							Path:             swapConfigStatePath(cfgPoint.Path),
							Timestamp:        stateVal.Timestamp,
							RecvTimestamp:    stateVal.RecvTimestamp,
							ComplianceErrors: stateVal.ComplianceErrors,
						}

						delete(statePoints, cfgToStatePaths[cfgPath])

						err := sr.fn(cfgVal, stateVal)
						if errors.Is(err, ReconcilerAbortErr) {
							r.errCh <- err
							return
						} else if err != nil {
							log.Warningf("reconciler error: %v", err)
						}
					}
					for _, statePoint := range statePoints {
						if !util.PathMatchesQuery(statePoint.Path, sr.state) {
							continue
						}
						stateVal := &Value[T]{
							val:              stateVal.val,
							present:          stateVal.present,
							Path:             statePoint.Path,
							Timestamp:        stateVal.Timestamp,
							RecvTimestamp:    stateVal.RecvTimestamp,
							ComplianceErrors: stateVal.ComplianceErrors,
						}
						cfgVal := &Value[T]{
							val:     cfgVal.val,
							present: cfgVal.present,
							// This path may exist since not all state paths have config paths. This should be fine because the use case for this is reconcilation config and state
							// so the callback may need handle cases where a state path exists but the config hasn't been created.
							Path:             swapConfigStatePath(statePoint.Path),
							Timestamp:        cfgVal.Timestamp,
							RecvTimestamp:    cfgVal.RecvTimestamp,
							ComplianceErrors: cfgVal.ComplianceErrors,
						}
						err := sr.fn(cfgVal, stateVal)
						if errors.Is(err, ReconcilerAbortErr) {
							r.errCh <- err
							return
						} else if err != nil {
							log.Warningf("reconciler error: %v", err)
						}
					}
				}
				err = fn(cfgVal, stateVal)
				if errors.Is(err, ReconcilerAbortErr) {
					r.errCh <- err
					return
				} else if err != nil {
					log.Warningf("reconciler error: %v", err)
				}
			case err := <-errCh:
				r.errCh <- err
				return
			}
		}
	}()
}

// Await blocks until the reconciler exists.
func (r *Reconciler[T]) Await() error {
	return <-r.errCh
}

const (
	configPathElem = "config"
	statePathElem  = "state"
)

func swapConfigStatePath(p *gpb.Path) *gpb.Path {
	swap := proto.Clone(p).(*gpb.Path)
	if len(swap.Elem) < 2 {
		return swap
	}
	if swap.Elem[len(swap.Elem)-2].Name == configPathElem {
		swap.Elem[len(swap.Elem)-2].Name = statePathElem
	} else if swap.Elem[len(swap.Elem)-2].Name == statePathElem {
		swap.Elem[len(swap.Elem)-2].Name = configPathElem
	}
	return swap
}

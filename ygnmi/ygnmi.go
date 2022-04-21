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
	"fmt"
	"reflect"
	"time"

	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"

	log "github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// AnyQuery is a generic gNMI query for wildcard or non-wildcard state or config paths.
// Supported operations: Batch.
type AnyQuery[T any] interface {
	// pathStruct returns to path struct for this query.
	pathStruct() ygot.PathStruct
	// fieldname returns the name of YANG directory schema entry.
	// For leaves, this is the parent entry.
	dirName() string
	// goStruct returns the struct that query should be unmarshalled into.
	// For leaves, this is the parent.
	goStruct() ygot.ValidatedGoStruct
	// extract is used for leaves to return the field from the parent GoStruct.
	// For non-leaves, this casts the GoStruct to the concrete type.
	extract(ygot.ValidatedGoStruct) T
	// isState returns if the path for this query is a state node.
	isState() bool
	// isLeaf returns if the path for this query is a leaf.
	isLeaf() bool
	// schema returns the root schema used for unmarshalling.
	schema() *ytypes.Schema
}

// SingletonQuery is a non-wildcard gNMI query.
type SingletonQuery[T any] interface {
	AnyQuery[T]
	// isNonWildcard prevents this interface from being used in wildcard funcs.
	isNonWildcard()
}

// WildcardQuery is a wildcard gNMI query.
type WildcardQuery[T any] interface {
	AnyQuery[T]
	// isWildcard prevents this interface from being used in non-wildcard funcs.
	isWildcard()
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

// SetVal sets the value and marks it present.
func (v *Value[T]) SetVal(val T) {
	v.present = !reflect.ValueOf(val).IsZero()
	v.val = val
}

// Val returns the val and whether it is present.
func (v *Value[T]) Val() (T, bool) {
	return v.val, v.present
}

// IsPresent returns whether the value is present.
func (v *Value[T]) IsPresent() bool {
	return v.present
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

// Lookup fetches the value of a SingletonQuery with a ONCE subscription.
func Lookup[T any](ctx context.Context, c *Client, q SingletonQuery[T]) (*Value[T], error) {
	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_ONCE)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to path: %w", err)
	}
	data, err := receiveAll(sub, false, gpb.SubscriptionList_ONCE)
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

// Watcher represents an ongoing watch of telemetry values.
type Watcher[T any] struct {
	errCh      chan error
	lastVal    *Value[T]
	predStatus bool
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

// Watch starts an asynchronous observation of the values with a STREAM subscription, evaluating each observed value with the specified predicate.
// The subscription completes when either the predicate is true or the context is cancelled.
// Calling Await on the returned Watcher waits for the subscription to complete.
// It returns the last observed value and a boolean that indicates whether that value satisfies the predicate.
func Watch[T any](ctx context.Context, c *Client, q SingletonQuery[T], pred func(*Value[T]) bool) *Watcher[T] {
	w := &Watcher[T]{
		errCh: make(chan error, 1),
	}

	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_STREAM)
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
				if w.predStatus = pred(val); w.predStatus {
					w.errCh <- nil
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

// Lookup fetches the values of a WildcardQuery with a ONCE subscription.
// It returns an empty list if no values are present at the path.
func LookupAll[T any](ctx context.Context, c *Client, q WildcardQuery[T]) ([]*Value[T], error) {
	sub, err := subscribe[T](ctx, c, q, gpb.SubscriptionList_ONCE)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to path: %w", err)
	}
	data, err := receiveAll(sub, false, gpb.SubscriptionList_ONCE)
	if err != nil {
		return nil, fmt.Errorf("failed to receive to data: %w", err)
	}
	p, _, errs := ygot.ResolvePath(q.pathStruct())
	if err := errsToErr(errs); err != nil {
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

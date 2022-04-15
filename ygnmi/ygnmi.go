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
	"time"

	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// AnyQuery is a generic gNMI query for wildcard or non-wildcard state or config paths.
// Supported operations: Batch.
type AnyQuery[T any] interface {
	// pathStruct returns to path struct for this query.
	pathStruct() ygot.PathStruct
	// fieldname returns the name of YANG directory schema entry.
	// For leaves, this is the parent entry.
	fieldName() string
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

type SingletonQuery[T any] interface {
	AnyQuery[T]
	// isNonWildcard prevents this interface from being used in wildcard funcs.
	isNonWildcard()
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
	v.val = val
	v.present = true
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

// Lookup fetches the value at the query with a ONCE subscription.
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
	return val, nil
}

// Get fetches the value at the query with a ONCE subscription.
// It returns an error if there is no value at the path.
func Get[T any](ctx context.Context, c *Client, q SingletonQuery[T]) (T, error) {
	val, err := Lookup(ctx, c, q)
	v, ok := val.Val()
	if !ok {
		return v, fmt.Errorf("no value received; lookup err: %v", err)
	}
	return v, err
}

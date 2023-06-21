// Copyright 2022 Google LLC
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

package ygnmi

import (
	"fmt"
	"reflect"

	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// ExtractFn is the type for the func that extracts a concrete val from a GoStruct.
type ExtractFn[T any] func(ygot.ValidatedGoStruct) (T, bool)

// NewSingletonQuery creates a new SingletonQueryStruct object.
func NewSingletonQuery[T any](goStructName string, state, leaf, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schemaFn func() *ytypes.Schema, subPaths []PathStruct) *SingletonQueryStruct[T] {
	return &SingletonQueryStruct[T]{
		baseQuery: baseQuery[T]{
			goStructName,
			state,
			ps,
			leaf,
			scalar,
			schemaFn,
			extractFn,
			goStructFn,
			subPaths,
		},
	}
}

// NewConfigQuery creates a new NewLeafConfigQuery object.
func NewConfigQuery[T any](goStructName string, state, leaf, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schemaFn func() *ytypes.Schema, subPaths []PathStruct) *ConfigQueryStruct[T] {
	return &ConfigQueryStruct[T]{
		baseQuery: baseQuery[T]{
			goStructName,
			state,
			ps,
			leaf,
			scalar,
			schemaFn,
			extractFn,
			goStructFn,
			nil,
		},
	}
}

// NewWildcardQuery creates a new NewLeafWildcardQuery object.
func NewWildcardQuery[T any](goStructName string, state, leaf, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schemaFn func() *ytypes.Schema) *WildcardQueryStruct[T] {
	return &WildcardQueryStruct[T]{
		baseQuery: baseQuery[T]{
			goStructName,
			state,
			ps,
			leaf,
			scalar,
			schemaFn,
			extractFn,
			goStructFn,
			nil,
		},
	}
}

// SingletonQueryStruct is implementation of SingletonQuery interface.
// Note: Do not use this type directly, instead use the generated Path API.
type SingletonQueryStruct[T any] struct {
	baseQuery[T]
}

// IsSingleton prevents this struct from being used where a wildcard path is expected.
func (q *SingletonQueryStruct[T]) IsSingleton() {}

// WildcardQueryStruct is implementation of SingletonQuery interface for leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type WildcardQueryStruct[T any] struct {
	baseQuery[T]
}

// IsWildcard prevents this struct from being used where a non wildcard path is expected.
func (q *WildcardQueryStruct[T]) IsWildcard() {}

// ConfigQueryStruct is implementation of ConfigQuery interface for leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type ConfigQueryStruct[T any] struct {
	baseQuery[T]
}

// IsConfig restricts this struct to be used only where a config path is expected.
func (q *ConfigQueryStruct[T]) IsConfig() {}

// IsSingleton restricts this struct to be used only where a singleton path is expected.
func (q *ConfigQueryStruct[T]) IsSingleton() {}

type baseQuery[T any] struct {
	// goStructName is the name of the YANG directory or GoStruct which
	// contains this node.
	//
	// - For GoStructs this is the struct itself.
	// - For others this is the parent dir.
	goStructName string
	// state controls if state or config values should be unmarshalled.
	state bool
	// ps contains the path specification of the query.
	ps PathStruct
	// leaf indicates whether the query is on a leaf node.
	leaf bool
	// scalar is whether the type (T) for this path is a pointer field (*T) in the parent GoStruct.
	scalar bool
	// yschemaFn is parsed YANG schema to use when unmarshalling data.
	yschemaFn func() *ytypes.Schema
	// extractFn is used to extract the value from the containing GoStruct.
	extractFn ExtractFn[T]
	// goStructFn initializes a new GoStruct containing this node.
	goStructFn func() ygot.ValidatedGoStruct
	// queryPathStructs are the paths used to for the gNMI subscription.
	// They must be equal to or descendants of ps.
	queryPathStructs []PathStruct
}

// dirName returns the YANG schema name of the GoStruct containing this node.
func (q *baseQuery[T]) dirName() string {
	return q.goStructName
}

// IsState returns if the Query is for a state or config path.
func (q *baseQuery[T]) IsState() bool {
	return q.state
}

// PathStruct returns the path struct containing the path for the Query.
func (q *baseQuery[T]) PathStruct() PathStruct {
	return q.ps
}

// schema returns the schema used for unmarshalling.
func (q *baseQuery[T]) schema() *ytypes.Schema {
	return q.yschemaFn()
}

// String returns gNMI path as string for the query.
func (q *baseQuery[T]) String() string {
	protoPath, _, err := ResolvePath(q.ps)
	if err != nil {
		return fmt.Sprintf("invalid path: %v", err)
	}
	path, err := ygot.PathToString(protoPath)
	if err != nil {
		path = protoPath.String()
	}
	return path
}

// extract extracts the unmarshalled value from the containing GoStruct.
//
// For queries on GoStructs it simply casts the input GoStruct to the concrete
// type for the query.
//
// For other types it extracts and returns the correct child field from it.
func (q *baseQuery[T]) extract(gs ygot.ValidatedGoStruct) (T, bool) {
	if q.extractFn != nil {
		return q.extractFn(gs)
	}

	val := gs.(T)
	return val, !reflect.ValueOf(val).Elem().IsZero()
}

// goStruct returns new empty GoStruct into which gNMI notifications can be
// unmarshalled.
func (q *baseQuery[T]) goStruct() ygot.ValidatedGoStruct {
	if q.goStructFn != nil {
		return q.goStructFn()
	}

	// Get the underlying type of T (which is a pointer), deference it to get the base type.
	// Create a new instance of the base type and return it as a ValidatedGoStruct.
	var t T
	gs := reflect.New(reflect.TypeOf(t).Elem())
	return gs.Interface().(ygot.ValidatedGoStruct)
}

// isLeaf returns whether the query refers to a leaf.
func (q *baseQuery[T]) isLeaf() bool {
	return q.leaf
}

// subPaths returns the path structs used for creating the gNMI subscription.
func (q *baseQuery[T]) subPaths() []PathStruct {
	if len(q.queryPathStructs) == 0 {
		return []PathStruct{q.ps}
	}
	return q.queryPathStructs
}

// isScalar returns whether the type (T) for this path is a leaf type that is
// represented by a pointer field (*T) in the parent GoStruct, but whose
// natural type is not a pointer (e.g. YANG's string type).
func (q *baseQuery[T]) isScalar() bool {
	return q.scalar
}

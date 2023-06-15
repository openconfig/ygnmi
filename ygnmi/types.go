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

// CompressionInfo contains information about a compressed path element for a
// node which points to a path element that's compressed out.
//
// e.g. for OpenConfig's /interfaces/interface, if a path points to
// /interfaces, then CompressionInfo will be populated with the following:
// - PreRelPath: []{"openconfig-interfaces:interfaces"}
// - PostRelPath: []{"openconfig-interfaces:interface"}
type CompressionInfo struct {
	// PreRelPath is the list of qualified path elements prior to the
	// compressed-out node.
	PreRelPath []string
	// PostRelPath is the list of qualified path elements after the
	// compressed-out node.
	PostRelPath []string
}

func newBaseQuery[T any](dir string, state bool, ps PathStruct, schemaFn func() *ytypes.Schema, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct) baseQuery[T] {
	return baseQuery[T]{
		goStructName: dir,
		state:        state,
		ps:           ps,
		yschemaFn:    schemaFn,
		extractFn:    extractFn,
		goStructFn:   goStructFn,
	}
}

// NewLeafSingletonQuery creates a new LeafSingletonQuery object.
func NewLeafSingletonQuery[T any](parentDir string, state, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schemaFn func() *ytypes.Schema) *LeafSingletonQuery[T] {
	return &LeafSingletonQuery[T]{
		leafBaseQuery: leafBaseQuery[T]{
			baseQuery: newBaseQuery(
				parentDir,
				state,
				ps,
				schemaFn,
				extractFn,
				goStructFn,
			),
			scalar: scalar,
		},
	}
}

// NewNonLeafSingletonQuery creates a new NonLeafSingletonQuery object.
func NewNonLeafSingletonQuery[T any](dir string, state bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, subPaths []PathStruct, schemaFn func() *ytypes.Schema, compressInfo *CompressionInfo) *NonLeafSingletonQuery[T] {
	return &NonLeafSingletonQuery[T]{
		nonLeafBaseQuery: nonLeafBaseQuery[T]{
			baseQuery: newBaseQuery(
				dir,
				state,
				ps,
				schemaFn,
				extractFn,
				goStructFn,
			),
			queryPathStructs: subPaths,
			compInfo:         compressInfo,
		},
	}
}

// NewLeafConfigQuery creates a new NewLeafConfigQuery object.
func NewLeafConfigQuery[T any](parentDir string, state, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schemaFn func() *ytypes.Schema) *LeafConfigQuery[T] {
	return &LeafConfigQuery[T]{
		leafBaseQuery: leafBaseQuery[T]{
			baseQuery: newBaseQuery(
				parentDir,
				state,
				ps,
				schemaFn,
				extractFn,
				goStructFn,
			),
			scalar: scalar,
		},
	}
}

// NewNonLeafConfigQuery creates a new NewNonLeafConfigQuery object.
func NewNonLeafConfigQuery[T any](dir string, state bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, subPaths []PathStruct, schemaFn func() *ytypes.Schema, compressInfo *CompressionInfo) *NonLeafConfigQuery[T] {
	return &NonLeafConfigQuery[T]{
		nonLeafBaseQuery: nonLeafBaseQuery[T]{
			baseQuery: newBaseQuery(
				dir,
				state,
				ps,
				schemaFn,
				extractFn,
				goStructFn,
			),
			compInfo: compressInfo,
		},
	}
}

// NewLeafWildcardQuery creates a new NewLeafWildcardQuery object.
func NewLeafWildcardQuery[T any](parentDir string, state, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schemaFn func() *ytypes.Schema) *LeafWildcardQuery[T] {
	return &LeafWildcardQuery[T]{
		leafBaseQuery: leafBaseQuery[T]{
			baseQuery: newBaseQuery(
				parentDir,
				state,
				ps,
				schemaFn,
				extractFn,
				goStructFn,
			),
			scalar: scalar,
		},
	}
}

// NewNonLeafWildcardQuery creates a new NewNonLeafWildcardQuery object.
func NewNonLeafWildcardQuery[T any](dir string, state bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schemaFn func() *ytypes.Schema, compressInfo *CompressionInfo) *NonLeafWildcardQuery[T] {
	return &NonLeafWildcardQuery[T]{
		nonLeafBaseQuery: nonLeafBaseQuery[T]{
			baseQuery: newBaseQuery(
				dir,
				state,
				ps,
				schemaFn,
				extractFn,
				goStructFn,
			),
			compInfo: compressInfo,
		},
	}
}

// ExtractFn is the type for the func that extracts a concrete val from a GoStruct.
type ExtractFn[T any] func(ygot.ValidatedGoStruct) (T, bool)

// LeafSingletonQuery is implementation of SingletonQuery interface for leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type LeafSingletonQuery[T any] struct {
	leafBaseQuery[T]
}

// IsSingleton prevents this struct from being used where a wildcard path is expected.
func (lq *LeafSingletonQuery[T]) IsSingleton() {}

// LeafWildcardQuery is implementation of SingletonQuery interface for leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type LeafWildcardQuery[T any] struct {
	leafBaseQuery[T]
}

// IsWildcard prevents this struct from being used where a non wildcard path is expected.
func (lq *LeafWildcardQuery[T]) IsWildcard() {}

// baseQuery contains common fields for query objects.
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
	// yschemaFn is parsed YANG schema to use when unmarshalling data.
	yschemaFn func() *ytypes.Schema
	// extractFn extracts the value from the containing GoStruct.
	//
	// If this value is nil, then it is assumed that the baseQuery refers
	// to a GoStruct and the value itself can be returned.
	extractFn ExtractFn[T]
	// goStructFn initializes a new GoStruct able to contain the given
	// node.
	//
	// If this value is nil, then it is assumed that the baseQuery refers
	// to a GoStruct and the type itself can be returned.
	goStructFn func() ygot.ValidatedGoStruct
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

// goStruct returns new initialized GoStruct, gNMI notifications can be
// unmarshalled into this struct.
func (q *baseQuery[T]) goStruct() ygot.ValidatedGoStruct {
	if q.goStructFn != nil {
		return q.goStructFn()
	}
	var t T
	gs := reflect.New(reflect.TypeOf(t).Elem())
	return gs.Interface().(ygot.ValidatedGoStruct)
}

// extract extracts the value from the containing GoStruct, returning whether
// the value exists.
func (q *baseQuery[T]) extract(gs ygot.ValidatedGoStruct) (T, bool) {
	if q.extractFn != nil {
		return q.extractFn(gs)
	}
	val := gs.(T)
	return val, !reflect.ValueOf(val).Elem().IsZero()
}

type leafBaseQuery[T any] struct {
	baseQuery[T]
	// scalar is whether the type (T) for this path is a pointer field (*T) in the parent GoStruct.
	scalar bool
}

// String returns gNMI path as string for the query.
func (lq *leafBaseQuery[T]) String() string {
	protoPath, _, err := ResolvePath(lq.ps)
	if err != nil {
		return fmt.Sprintf("invalid path: %v", err)
	}
	path, err := ygot.PathToString(protoPath)
	if err != nil {
		path = protoPath.String()
	}
	return path
}

// isLeaf returns true, as this Query type is only for leaves.
func (lq *leafBaseQuery[T]) isLeaf() bool {
	return true
}

// subPaths returns the path structs used for creating the gNMI subscription.
// A leaf subscribes to the same path returned by PathStruct().
func (lq *leafBaseQuery[T]) subPaths() []PathStruct {
	return []PathStruct{lq.ps}
}

// isScalar returns whether the type (T) for this path is a pointer field (*T) in the parent GoStruct.
func (lq *leafBaseQuery[T]) isScalar() bool {
	return lq.scalar
}

// schema returns the schema used for unmarshalling.
func (lq *leafBaseQuery[T]) compressInfo() *CompressionInfo {
	return nil
}

// NonLeafSingletonQuery is implementation of SingletonQuery interface for non-leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type NonLeafSingletonQuery[T any] struct {
	nonLeafBaseQuery[T]
}

// IsSingleton prevents this struct from being used where a wildcard path is expected.
func (lq *NonLeafSingletonQuery[T]) IsSingleton() {}

// NonLeafWildcardQuery is implementation of SingletonQuery interface for non-leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type NonLeafWildcardQuery[T any] struct {
	nonLeafBaseQuery[T]
}

// IsWildcard prevents this struct from being used where a non-wildcard path is expected.
func (lq *NonLeafWildcardQuery[T]) IsWildcard() {}

type nonLeafBaseQuery[T any] struct {
	baseQuery[T]
	// queryPathStructs are the paths used to for the gNMI subscription.
	// They must be equal to or descendants of ps.
	queryPathStructs []PathStruct
	// compInfo stores compression information when the node points to a
	// path that's compressed out in the generated code.
	compInfo *CompressionInfo
}

// String returns gNMI path as string for the query.
func (lq *nonLeafBaseQuery[T]) String() string {
	protoPath, _, err := ResolvePath(lq.ps)
	if err != nil {
		return fmt.Sprintf("invalid path: %v", err)
	}
	path, err := ygot.PathToString(protoPath)
	if err != nil {
		path = protoPath.String()
	}
	return path
}

// isLeaf returns false, as this Query type is only for non-leaves.
func (lq *nonLeafBaseQuery[T]) isLeaf() bool {
	return false
}

// isScalar returns false, as non-leafs are always non-scalar objects.
func (lq *nonLeafBaseQuery[T]) isScalar() bool {
	return false
}

// subPaths returns the path structs used for creating the gNMI subscription.
func (lq *nonLeafBaseQuery[T]) subPaths() []PathStruct {
	if len(lq.queryPathStructs) == 0 {
		return []PathStruct{lq.ps}
	}
	return lq.queryPathStructs
}

// schema returns the schema used for unmarshalling.
func (lq *nonLeafBaseQuery[T]) compressInfo() *CompressionInfo {
	return lq.compInfo
}

// LeafConfigQuery is implementation of ConfigQuery interface for leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type LeafConfigQuery[T any] struct {
	leafBaseQuery[T]
}

// IsConfig restricts this struct to be used only where a config path is expected.
func (lq *LeafConfigQuery[T]) IsConfig() {}

// IsSingleton restricts this struct to be used only where a singleton path is expected.
func (lq *LeafConfigQuery[T]) IsSingleton() {}

// NonLeafConfigQuery is implementation of ConfigQuery interface for non-leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type NonLeafConfigQuery[T any] struct {
	nonLeafBaseQuery[T]
}

// IsConfig restricts this struct to be used only where a config path is expected.
func (nlq *NonLeafConfigQuery[T]) IsConfig() {}

// IsSingleton restricts this struct to be used only where a singleton path is expected.
func (nlq *NonLeafConfigQuery[T]) IsSingleton() {}

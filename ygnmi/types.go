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

// NewLeafSingletonQuery creates a new LeafSingletonQuery object.
func NewLeafSingletonQuery[T any](parentDir string, state, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schemaFn func() *ytypes.Schema) *LeafSingletonQuery[T] {
	return &LeafSingletonQuery[T]{
		leafBaseQuery: leafBaseQuery[T]{
			parentDir:  parentDir,
			state:      state,
			ps:         ps,
			scalar:     scalar,
			extractFn:  extractFn,
			goStructFn: goStructFn,
			yschemaFn:  schemaFn,
		},
	}
}

// NewNonLeafSingletonQuery creates a new NonLeafSingletonQuery object.
func NewNonLeafSingletonQuery[T ygot.ValidatedGoStruct](dir string, state bool, ps PathStruct, subPaths []PathStruct, schemaFn func() *ytypes.Schema) *NonLeafSingletonQuery[T] {
	return &NonLeafSingletonQuery[T]{
		nonLeafBaseQuery: nonLeafBaseQuery[T]{
			dir:              dir,
			state:            state,
			ps:               ps,
			yschemaFn:        schemaFn,
			queryPathStructs: subPaths,
		},
	}
}

// NewLeafConfigQuery creates a new NewLeafConfigQuery object.
func NewLeafConfigQuery[T any](parentDir string, state, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schemaFn func() *ytypes.Schema) *LeafConfigQuery[T] {
	return &LeafConfigQuery[T]{
		leafBaseQuery: leafBaseQuery[T]{
			parentDir:  parentDir,
			state:      state,
			ps:         ps,
			scalar:     scalar,
			extractFn:  extractFn,
			goStructFn: goStructFn,
			yschemaFn:  schemaFn,
		},
	}
}

// NewNonLeafConfigQuery creates a new NewNonLeafConfigQuery object.
func NewNonLeafConfigQuery[T ygot.ValidatedGoStruct](dir string, state bool, ps PathStruct, subPaths []PathStruct, schemaFn func() *ytypes.Schema) *NonLeafConfigQuery[T] {
	return &NonLeafConfigQuery[T]{
		nonLeafBaseQuery: nonLeafBaseQuery[T]{
			dir:       dir,
			state:     state,
			ps:        ps,
			yschemaFn: schemaFn,
		},
	}
}

// NewLeafWildcardQuery creates a new NewLeafWildcardQuery object.
func NewLeafWildcardQuery[T any](parentDir string, state, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schemaFn func() *ytypes.Schema) *LeafWildcardQuery[T] {
	return &LeafWildcardQuery[T]{
		leafBaseQuery: leafBaseQuery[T]{
			parentDir:  parentDir,
			state:      state,
			ps:         ps,
			scalar:     scalar,
			extractFn:  extractFn,
			goStructFn: goStructFn,
			yschemaFn:  schemaFn,
		},
	}
}

// NewNonLeafWildcardQuery creates a new NewNonLeafWildcardQuery object.
func NewNonLeafWildcardQuery[T ygot.ValidatedGoStruct](dir string, state bool, ps PathStruct, schemaFn func() *ytypes.Schema) *NonLeafWildcardQuery[T] {
	return &NonLeafWildcardQuery[T]{
		nonLeafBaseQuery: nonLeafBaseQuery[T]{
			dir:       dir,
			state:     state,
			ps:        ps,
			yschemaFn: schemaFn,
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

type leafBaseQuery[T any] struct {
	// parent is name of the YANG directory which contains this leaf.
	parentDir string
	// state controls if state or config values should be unmarshalled.
	state bool
	// ps contains the path of the query.
	ps PathStruct
	// extractFn gets the leaf node from the parent GoStruct.
	extractFn ExtractFn[T]
	// goStructFn initializes a new GoStruct for the given path.
	goStructFn func() ygot.ValidatedGoStruct
	// yschemaFn is parsed YANG schema to use when unmarshalling data.
	yschemaFn func() *ytypes.Schema
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

// extract takes the parent GoStruct and returns the correct child field from it.
func (lq *leafBaseQuery[T]) extract(gs ygot.ValidatedGoStruct) (T, bool) {
	return lq.extractFn(gs)
}

// dirName returns the YANG schema name of the parent GoStruct.
func (lq *leafBaseQuery[T]) dirName() string {
	return lq.parentDir
}

// goStruct returns the parent struct of the leaf node.
func (lq *leafBaseQuery[T]) goStruct() ygot.ValidatedGoStruct {
	return lq.goStructFn()
}

// isLeaf returns true, as this Query type is only for leaves.
func (lq *leafBaseQuery[T]) isLeaf() bool {
	return true
}

// IsState returns if the Query is for a state or config path.
func (lq *leafBaseQuery[T]) IsState() bool {
	return lq.state
}

// PathStruct returns the path struct containing the path for the Query.
func (lq *leafBaseQuery[T]) PathStruct() PathStruct {
	return lq.ps
}

// subPaths returns the path structs used for creating the gNMI subscription.
// A leaf subscribes to the same path returned by PathStruct().
func (lq *leafBaseQuery[T]) subPaths() []PathStruct {
	return []PathStruct{lq.ps}
}

// schema returns the schema used for unmarshalling.
func (lq *leafBaseQuery[T]) schema() *ytypes.Schema {
	return lq.yschemaFn()
}

// isScalar returns whether the type (T) for this path is a pointer field (*T) in the parent GoStruct.
func (lq *leafBaseQuery[T]) isScalar() bool {
	return lq.scalar
}

// NonLeafSingletonQuery is implementation of SingletonQuery interface for non-leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type NonLeafSingletonQuery[T ygot.ValidatedGoStruct] struct {
	nonLeafBaseQuery[T]
}

// IsSingleton prevents this struct from being used where a wildcard path is expected.
func (lq *NonLeafSingletonQuery[T]) IsSingleton() {}

// NonLeafWildcardQuery is implementation of SingletonQuery interface for non-leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type NonLeafWildcardQuery[T ygot.ValidatedGoStruct] struct {
	nonLeafBaseQuery[T]
}

// IsWildcard prevents this struct from being used where a non-wildcard path is expected.
func (lq *NonLeafWildcardQuery[T]) IsWildcard() {}

type nonLeafBaseQuery[T ygot.ValidatedGoStruct] struct {
	dir   string
	state bool
	// ps is the path used for unmarshalling and schema validation.
	ps PathStruct
	// queryPathStructs are the paths used to for the gNMI subscription.
	// They must be equal to or descendants of ps.
	queryPathStructs []PathStruct
	yschemaFn        func() *ytypes.Schema
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

// extract casts the input GoStruct to the concrete type for the query.
// As non-leaves structs are always GoStructs, a simple cast is sufficient.
func (lq *nonLeafBaseQuery[T]) extract(gs ygot.ValidatedGoStruct) (T, bool) {
	val := gs.(T)
	return val, !reflect.ValueOf(val).Elem().IsZero()
}

// dirName returns the YANG schema directory name, used to unmarshal values.
func (lq *nonLeafBaseQuery[T]) dirName() string {
	return lq.dir
}

// goStruct returns new initialized GoStruct, gNMI notifications can be unmarshalled into this struct.
func (lq *nonLeafBaseQuery[T]) goStruct() ygot.ValidatedGoStruct {
	// Get the underlying type of T (which is a pointer), deference it to get the base type.
	// Create a new instance of the base type and return it as a ValidatedGoStruct.
	var t T
	gs := reflect.New(reflect.TypeOf(t).Elem())
	return gs.Interface().(ygot.ValidatedGoStruct)
}

// isLeaf returns false, as this Query type is only for non-leaves.
func (lq *nonLeafBaseQuery[T]) isLeaf() bool {
	return false
}

// isScalar returns false, as non-leafs are always non-scalar objects.
func (lq *nonLeafBaseQuery[T]) isScalar() bool {
	return false
}

// IsState returns if the Query is for a state or config path.
func (lq *nonLeafBaseQuery[T]) IsState() bool {
	return lq.state
}

// PathStruct returns the path struct containing the path for the Query.
func (lq *nonLeafBaseQuery[T]) PathStruct() PathStruct {
	return lq.ps
}

// subPaths returns the path structs used for creating the gNMI subscription.
func (lq *nonLeafBaseQuery[T]) subPaths() []PathStruct {
	if len(lq.queryPathStructs) == 0 {
		return []PathStruct{lq.ps}
	}
	return lq.queryPathStructs
}

// schema returns the schema used for unmarshalling.
func (lq *nonLeafBaseQuery[T]) schema() *ytypes.Schema {
	return lq.yschemaFn()
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
type NonLeafConfigQuery[T ygot.ValidatedGoStruct] struct {
	nonLeafBaseQuery[T]
}

// IsConfig restricts this struct to be used only where a config path is expected.
func (nlq *NonLeafConfigQuery[T]) IsConfig() {}

// IsSingleton restricts this struct to be used only where a singleton path is expected.
func (nlq *NonLeafConfigQuery[T]) IsSingleton() {}

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
	"reflect"

	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// NewLeafSingletonQuery creates a new LeafSingletonQuery object.
func NewLeafSingletonQuery[T any](parentDir string, state, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schema *ytypes.Schema) *LeafSingletonQuery[T] {
	return &LeafSingletonQuery[T]{
		leafBaseQuery: leafBaseQuery[T]{
			parentDir:  parentDir,
			state:      state,
			ps:         ps,
			scalar:     scalar,
			extractFn:  extractFn,
			goStructFn: goStructFn,
			yschema:    schema,
		},
	}
}

// NewNonLeafSingletonQuery creates a new NonLeafSingletonQuery object.
func NewNonLeafSingletonQuery[T ygot.ValidatedGoStruct](dir string, state bool, ps PathStruct, schema *ytypes.Schema) *NonLeafSingletonQuery[T] {
	return &NonLeafSingletonQuery[T]{
		nonLeafBaseQuery: nonLeafBaseQuery[T]{
			dir:     dir,
			state:   state,
			ps:      ps,
			yschema: schema,
		},
	}
}

// NewLeafConfigQuery creates a new NewLeafConfigQuery object.
func NewLeafConfigQuery[T any](parentDir string, state, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schema *ytypes.Schema) *LeafConfigQuery[T] {
	return &LeafConfigQuery[T]{
		leafBaseQuery: leafBaseQuery[T]{
			parentDir:  parentDir,
			state:      state,
			ps:         ps,
			scalar:     scalar,
			extractFn:  extractFn,
			goStructFn: goStructFn,
			yschema:    schema,
		},
	}
}

// NewNonLeafConfigQuery creates a new NewNonLeafConfigQuery object.
func NewNonLeafConfigQuery[T ygot.ValidatedGoStruct](dir string, state bool, ps PathStruct, schema *ytypes.Schema) *NonLeafConfigQuery[T] {
	return &NonLeafConfigQuery[T]{
		nonLeafBaseQuery: nonLeafBaseQuery[T]{
			dir:     dir,
			state:   state,
			ps:      ps,
			yschema: schema,
		},
	}
}

// NewLeafWildcardQuery creates a new NewLeafWildcardQuery object.
func NewLeafWildcardQuery[T any](parentDir string, state, scalar bool, ps PathStruct, extractFn ExtractFn[T], goStructFn func() ygot.ValidatedGoStruct, schema *ytypes.Schema) *LeafWildcardQuery[T] {
	return &LeafWildcardQuery[T]{
		leafBaseQuery: leafBaseQuery[T]{
			parentDir:  parentDir,
			state:      state,
			ps:         ps,
			scalar:     scalar,
			extractFn:  extractFn,
			goStructFn: goStructFn,
			yschema:    schema,
		},
	}
}

// NewNonLeafWildcardQuery creates a new NewNonLeafWildcardQuery object.
func NewNonLeafWildcardQuery[T ygot.ValidatedGoStruct](dir string, state bool, ps PathStruct, schema *ytypes.Schema) *NonLeafWildcardQuery[T] {
	return &NonLeafWildcardQuery[T]{
		nonLeafBaseQuery: nonLeafBaseQuery[T]{
			dir:     dir,
			state:   state,
			ps:      ps,
			yschema: schema,
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

// isSingleton prevents this struct from being used where a wildcard path is expected.
func (lq *LeafSingletonQuery[T]) isSingleton() {}

// LeafWildcardQuery is implementation of SingletonQuery interface for leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type LeafWildcardQuery[T any] struct {
	leafBaseQuery[T]
}

// isWildcard prevents this struct from being used where a non wildcard path is expected.
func (lq *LeafWildcardQuery[T]) isWildcard() {}

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
	// yschema is parsed YANG schema to use when unmarshalling data.
	yschema *ytypes.Schema
	// scalar is whether the type (T) for this path is a pointer field (*T) in the parent GoStruct.
	scalar bool
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

// isState returns if the Query is for a state or config path.
func (lq *leafBaseQuery[T]) isState() bool {
	return lq.state
}

// pathStruct returns the path struct containing the path for the Query.
func (lq *leafBaseQuery[T]) pathStruct() PathStruct {
	return lq.ps
}

// schema returns the schema used for unmarshalling.
func (lq *leafBaseQuery[T]) schema() *ytypes.Schema {
	return lq.yschema
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

// isSingleton prevents this struct from being used where a wildcard path is expected.
func (lq *NonLeafSingletonQuery[T]) isSingleton() {}

// NonLeafWildcardQuery is implementation of SingletonQuery interface for non-leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type NonLeafWildcardQuery[T ygot.ValidatedGoStruct] struct {
	nonLeafBaseQuery[T]
}

// isNonWildcard prevents this struct from being used where a non-wildcard path is expected.
func (lq *NonLeafWildcardQuery[T]) isWildcard() {}

type nonLeafBaseQuery[T ygot.ValidatedGoStruct] struct {
	dir     string
	state   bool
	ps      PathStruct
	yschema *ytypes.Schema
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

// isState returns if the Query is for a state or config path.
func (lq *nonLeafBaseQuery[T]) isState() bool {
	return lq.state
}

// pathStruct returns the path struct containing the path for the Query.
func (lq *nonLeafBaseQuery[T]) pathStruct() PathStruct {
	return lq.ps
}

// schema returns the schema used for unmarshalling.
func (lq *nonLeafBaseQuery[T]) schema() *ytypes.Schema {
	return lq.yschema
}

// LeafConfigQuery is implementation of ConfigQuery interface for leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type LeafConfigQuery[T any] struct {
	leafBaseQuery[T]
}

// isConfig restricts this struct to be used only where a config path is expected.
func (lq *LeafConfigQuery[T]) isConfig() {}

// isSingleton restricts this struct to be used only where a singleton path is expected.
func (lq *LeafConfigQuery[T]) isSingleton() {}

// NonLeafConfigQuery is implementation of ConfigQuery interface for non-leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type NonLeafConfigQuery[T ygot.ValidatedGoStruct] struct {
	nonLeafBaseQuery[T]
}

// isConfig restricts this struct to be used only where a config path is expected.
func (nlq *NonLeafConfigQuery[T]) isConfig() {}

// isSingleton restricts this struct to be used only where a singleton path is expected.
func (nlq *NonLeafConfigQuery[T]) isSingleton() {}

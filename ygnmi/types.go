package ygnmi

import (
	"reflect"

	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

func NewLeafSingletonQuery[T any](parentDir string, state bool, ps ygot.PathStruct, extractFn func(ygot.ValidatedGoStruct) T, goStructFn func() ygot.ValidatedGoStruct, schema *ytypes.Schema) *LeafSingletonQuery[T] {
	return &LeafSingletonQuery[T]{
		leafBaseQuery: leafBaseQuery[T]{
			parentDir:  parentDir,
			state:      state,
			ps:         ps,
			extractFn:  extractFn,
			goStructFn: goStructFn,
			yschema:    schema,
		},
	}
}

// LeafSingletonQuery is implementation of SingletonQuery interface for leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type LeafSingletonQuery[T any] struct {
	leafBaseQuery[T]
}

// isNonWildcard prevents this struct from being used where a wildcard path is expected.
func (lq *LeafSingletonQuery[T]) isNonWildcard() {}

// LeafSingletonQuery is implementation of SingletonQuery interface for leaf nodes.
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
	ps ygot.PathStruct
	// extractFn gets the leaf node from the parent GoStruct.
	extractFn func(ygot.ValidatedGoStruct) T
	// goStructFn initializes a new GoStruct for the given path.
	goStructFn func() ygot.ValidatedGoStruct
	// yschema is parsed YANG schema to use when unmarshalling data.
	yschema *ytypes.Schema
}

// extract takes the parent GoStruct and returns the correct child field from it.
func (lq *leafBaseQuery[T]) extract(gs ygot.ValidatedGoStruct) T {
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
func (lq *leafBaseQuery[T]) pathStruct() ygot.PathStruct {
	return lq.ps
}

// schema returns the schema used for unmarshalling.
func (lq *leafBaseQuery[T]) schema() *ytypes.Schema {
	return lq.yschema
}

// NonLeafSingletonQuery is implementation of SingletonQuery interface for non-leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type NonLeafSingletonQuery[T ygot.ValidatedGoStruct] struct {
	nonLeafBaseQuery[T]
}

// isNonWildcard prevents this struct from being used where a wildcard path is expected.
func (lq *NonLeafSingletonQuery[T]) isNonWildcard() {}

// NonLeafSingletonQuery is implementation of SingletonQuery interface for non-leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type NonLeafWildcardQuery[T ygot.ValidatedGoStruct] struct {
	nonLeafBaseQuery[T]
}

// isNonWildcard prevents this struct from being used where a non-wildcard path is expected.
func (lq *NonLeafWildcardQuery[T]) isWildcard() {}

type nonLeafBaseQuery[T ygot.ValidatedGoStruct] struct {
	dir     string
	state   bool
	ps      ygot.PathStruct
	yschema *ytypes.Schema
}

// extract casts the input GoStruct to the concrete type for the query.
// As non-leaves structs are always GoStructs, a simple cast is sufficient.
func (lq *nonLeafBaseQuery[T]) extract(gs ygot.ValidatedGoStruct) T {
	return gs.(T)
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

// isState returns if the Query is for a state or config path.
func (lq *nonLeafBaseQuery[T]) isState() bool {
	return lq.state
}

// pathStruct returns the path struct containing the path for the Query.
func (lq *nonLeafBaseQuery[T]) pathStruct() ygot.PathStruct {
	return lq.ps
}

// schema returns the schema used for unmarshalling.
func (lq *nonLeafBaseQuery[T]) schema() *ytypes.Schema {
	return lq.yschema
}

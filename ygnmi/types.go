package ygnmi

import (
	"reflect"

	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// LeafSingletonQuery is implementation of SingletonQuery interface for leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type LeafSingletonQuery[T any] struct {
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

func (lq *LeafSingletonQuery[T]) extract(gs ygot.ValidatedGoStruct) T {
	return lq.extractFn(gs)
}

func (lq *LeafSingletonQuery[T]) dirName() string {
	return lq.parentDir
}

func (lq *LeafSingletonQuery[T]) goStruct() ygot.ValidatedGoStruct {
	return lq.goStructFn()
}

func (lq *LeafSingletonQuery[T]) isLeaf() bool {
	return true
}

func (lq *LeafSingletonQuery[T]) isState() bool {
	return lq.state
}

func (lq *LeafSingletonQuery[T]) pathStruct() ygot.PathStruct {
	return lq.ps
}

func (lq *LeafSingletonQuery[T]) schema() *ytypes.Schema {
	return lq.yschema
}

func (lq *LeafSingletonQuery[T]) isNonWildcard() {}

// LeafSingletonQuery is implementation of SingletonQuery interface for non-leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type NonLeafSingletonQuery[T ygot.ValidatedGoStruct] struct {
	dir     string
	state   bool
	ps      ygot.PathStruct
	yschema *ytypes.Schema
}

func (lq *NonLeafSingletonQuery[T]) extract(gs ygot.ValidatedGoStruct) T {
	return gs.(T)
}

func (lq *NonLeafSingletonQuery[T]) fieldName() string {
	return lq.dir
}

func (lq *NonLeafSingletonQuery[T]) goStruct() ygot.ValidatedGoStruct {
	// Get the underlying type of T (which is a pointer), deference it to get the base type.
	// Create a new instance of the base type and return it as a ValidatedGoStruct.
	var noop T
	val := reflect.ValueOf(noop)
	gs := reflect.New(val.Type().Elem())
	return gs.Interface().(ygot.ValidatedGoStruct)
}

func (lq *NonLeafSingletonQuery[T]) isLeaf() bool {
	return false
}

func (lq *NonLeafSingletonQuery[T]) isState() bool {
	return lq.state
}

func (lq *NonLeafSingletonQuery[T]) pathStruct() ygot.PathStruct {
	return lq.ps
}

func (lq *NonLeafSingletonQuery[T]) schema() *ytypes.Schema {
	return lq.yschema
}

func (lq *NonLeafSingletonQuery[T]) isNonWildcard() {}

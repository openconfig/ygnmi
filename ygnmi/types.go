package ygnmi

import (
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

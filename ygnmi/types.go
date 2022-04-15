package ygnmi

import (
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// LeafSingletonQuery is implementation of SingletonQuery interface for leaf nodes.
// Note: Do not use this type directly, instead use the generated Path API.
type LeafSingletonQuery[T any] struct {
	parentDir  string
	state      bool
	ps         ygot.PathStruct
	extractFn  func(ygot.ValidatedGoStruct) T
	goStructFn func() ygot.ValidatedGoStruct
	yschema    *ytypes.Schema
}

func (lq *LeafSingletonQuery[T]) extract(gs ygot.ValidatedGoStruct) T {
	return lq.extractFn(gs)
}

func (lq *LeafSingletonQuery[T]) fieldName() string {
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

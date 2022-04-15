package ygnmi

import (
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

type LeafQuery[T any] struct {
	parentDir  string
	state      bool
	ps         ygot.PathStruct
	extractFn  func(ygot.ValidatedGoStruct) T
	goStructFn func() ygot.ValidatedGoStruct
	yschema    *ytypes.Schema
}

func (lq *LeafQuery[T]) extract(gs ygot.ValidatedGoStruct) T {
	return lq.extractFn(gs)
}

func (lq *LeafQuery[T]) fieldName() string {
	return lq.parentDir
}

func (lq *LeafQuery[T]) goStruct() ygot.ValidatedGoStruct {
	return lq.goStructFn()
}

func (lq *LeafQuery[T]) isLeaf() bool {
	return true
}

func (lq *LeafQuery[T]) isState() bool {
	return lq.state
}

func (lq *LeafQuery[T]) pathStruct() ygot.PathStruct {
	return lq.ps
}

func (lq *LeafQuery[T]) schema() *ytypes.Schema {
	return lq.yschema
}

func (lq *LeafQuery[T]) isNonWildcard() {}

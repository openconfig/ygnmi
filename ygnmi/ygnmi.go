package ygnmi

import (
	"time"

	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// AnyQuery is a generic gNMI query for wildcard or non-wildcard state or config paths.
// Supported operations: Batch.
type AnyQuery[T any] interface {
	pathStruct() ygot.PathStruct
	fieldName() string
	goStruct() ygot.ValidatedGoStruct
	validate(ygot.ValidatedGoStruct) T
	reverseShadowPaths() bool
	isLeaf() bool
	schema() *ytypes.Schema
}

type Value[T any] struct {
	val              T
	present          bool
	Path             *gpb.Path         // Path is the sample's YANG path.
	Timestamp        time.Time         // Timestamp is the sample time.
	RecvTimestamp    time.Time         // Timestamp is the time the test received the sample.
	ComplianceErrors *ComplianceErrors // ComplianceErrors contains the compliance errors encountered from an Unmarshal operation.
}

func (v *Value[T]) SetVal(val T) {
	v.val = val
	v.present = true
}

func (v *Value[T]) Val() (T, bool) {
	return v.val, v.present
}
func (v *Value[T]) IsPresent() bool {
	return v.present
}

// Copyright 2022 Google Inc.
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

// Package ygnmi contains gNMI client library for use with a ygot Schema.
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

// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package modelb is a generated package which contains definitions
of structs which generate gNMI paths for a YANG schema.

This package was generated by ygnmi version: (devel): (ygot: v0.29.8)
using the following YANG input files:
  - ../pathgen/testdata/yang/openconfig-simple.yang
  - ../pathgen/testdata/yang/openconfig-withlistval.yang
  - ../pathgen/testdata/yang/openconfig-nested.yang

Imported modules were sourced from:
*/
package modelb

import (
	oc "github.com/openconfig/ygnmi/exampleoc"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// WithKey1 sets Model_MultiKeyPathAny's key "key1" to the specified value.
// Key1: uint32
func (n *Model_MultiKeyPathAny) WithKey1(Key1 uint32) *Model_MultiKeyPathAny {
	ygnmi.ModifyKey(n.NodePath, "key1", Key1)
	return n
}

// WithKey2 sets Model_MultiKeyPathAny's key "key2" to the specified value.
// Key2: uint64
func (n *Model_MultiKeyPathAny) WithKey2(Key2 uint64) *Model_MultiKeyPathAny {
	ygnmi.ModifyKey(n.NodePath, "key2", Key2)
	return n
}

// Model_MultiKey_Key1Path represents the /openconfig-withlistval/model/b/multi-key/state/key1 YANG schema element.
type Model_MultiKey_Key1Path struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Model_MultiKey_Key1PathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key/state/key1 YANG schema element.
type Model_MultiKey_Key1PathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "state/key1"
//	Path from root:       "/model/b/multi-key/state/key1"
func (n *Model_MultiKey_Key1Path) State() ygnmi.SingletonQuery[uint32] {
	return ygnmi.NewSingletonQuery[uint32](
		"Model_MultiKey",
		true,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"state", "key1"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint32, bool) {
			ret := gs.(*oc.Model_MultiKey).Key1
			if ret == nil {
				var zero uint32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_MultiKey) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "state/key1"
//	Path from root:       "/model/b/multi-key/state/key1"
func (n *Model_MultiKey_Key1PathAny) State() ygnmi.WildcardQuery[uint32] {
	return ygnmi.NewWildcardQuery[uint32](
		"Model_MultiKey",
		true,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"state", "key1"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint32, bool) {
			ret := gs.(*oc.Model_MultiKey).Key1
			if ret == nil {
				var zero uint32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_MultiKey) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
	)
}

// Config returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "config/key1"
//	Path from root:       "/model/b/multi-key/config/key1"
func (n *Model_MultiKey_Key1Path) Config() ygnmi.ConfigQuery[uint32] {
	return ygnmi.NewConfigQuery[uint32](
		"Model_MultiKey",
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"config", "key1"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint32, bool) {
			ret := gs.(*oc.Model_MultiKey).Key1
			if ret == nil {
				var zero uint32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_MultiKey) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// Config returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "config/key1"
//	Path from root:       "/model/b/multi-key/config/key1"
func (n *Model_MultiKey_Key1PathAny) Config() ygnmi.WildcardQuery[uint32] {
	return ygnmi.NewWildcardQuery[uint32](
		"Model_MultiKey",
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"config", "key1"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint32, bool) {
			ret := gs.(*oc.Model_MultiKey).Key1
			if ret == nil {
				var zero uint32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_MultiKey) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
	)
}

// Model_MultiKey_Key2Path represents the /openconfig-withlistval/model/b/multi-key/state/key2 YANG schema element.
type Model_MultiKey_Key2Path struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Model_MultiKey_Key2PathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key/state/key2 YANG schema element.
type Model_MultiKey_Key2PathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "state/key2"
//	Path from root:       "/model/b/multi-key/state/key2"
func (n *Model_MultiKey_Key2Path) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"Model_MultiKey",
		true,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"state", "key2"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Model_MultiKey).Key2
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_MultiKey) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "state/key2"
//	Path from root:       "/model/b/multi-key/state/key2"
func (n *Model_MultiKey_Key2PathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"Model_MultiKey",
		true,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"state", "key2"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Model_MultiKey).Key2
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_MultiKey) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
	)
}

// Config returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "config/key2"
//	Path from root:       "/model/b/multi-key/config/key2"
func (n *Model_MultiKey_Key2Path) Config() ygnmi.ConfigQuery[uint64] {
	return ygnmi.NewConfigQuery[uint64](
		"Model_MultiKey",
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"config", "key2"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Model_MultiKey).Key2
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_MultiKey) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// Config returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "config/key2"
//	Path from root:       "/model/b/multi-key/config/key2"
func (n *Model_MultiKey_Key2PathAny) Config() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"Model_MultiKey",
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"config", "key2"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Model_MultiKey).Key2
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_MultiKey) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
	)
}

// Model_MultiKeyPath represents the /openconfig-withlistval/model/b/multi-key YANG schema element.
type Model_MultiKeyPath struct {
	*ygnmi.NodePath
}

// Model_MultiKeyPathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key YANG schema element.
type Model_MultiKeyPathAny struct {
	*ygnmi.NodePath
}

// Model_MultiKeyPathMap represents the /openconfig-withlistval/model/b/multi-key YANG schema element.
type Model_MultiKeyPathMap struct {
	*ygnmi.NodePath
}

// Model_MultiKeyPathMapAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key YANG schema element.
type Model_MultiKeyPathMapAny struct {
	*ygnmi.NodePath
}

// Key1 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "*/key1"
//	Path from root:       "/model/b/multi-key/*/key1"
func (n *Model_MultiKeyPath) Key1() *Model_MultiKey_Key1Path {
	ps := &Model_MultiKey_Key1Path{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "key1"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Key1 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "*/key1"
//	Path from root:       "/model/b/multi-key/*/key1"
func (n *Model_MultiKeyPathAny) Key1() *Model_MultiKey_Key1PathAny {
	ps := &Model_MultiKey_Key1PathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "key1"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Key2 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "*/key2"
//	Path from root:       "/model/b/multi-key/*/key2"
func (n *Model_MultiKeyPath) Key2() *Model_MultiKey_Key2Path {
	ps := &Model_MultiKey_Key2Path{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "key2"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Key2 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "*/key2"
//	Path from root:       "/model/b/multi-key/*/key2"
func (n *Model_MultiKeyPathAny) Key2() *Model_MultiKey_Key2PathAny {
	ps := &Model_MultiKey_Key2PathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "key2"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

func binarySliceToFloatSlice(in []oc.Binary) []float32 {
	converted := make([]float32, 0, len(in))
	for _, binary := range in {
		converted = append(converted, ygot.BinaryToFloat32(binary))
	}
	return converted
}

// State returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPath) State() ygnmi.SingletonQuery[*oc.Model_MultiKey] {
	return ygnmi.NewSingletonQuery[*oc.Model_MultiKey](
		"Model_MultiKey",
		true,
		false,
		false,
		true,
		false,
		n,
		nil,
		nil,
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPathAny) State() ygnmi.WildcardQuery[*oc.Model_MultiKey] {
	return ygnmi.NewWildcardQuery[*oc.Model_MultiKey](
		"Model_MultiKey",
		true,
		false,
		false,
		true,
		false,
		n,
		nil,
		nil,
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPath) Config() ygnmi.ConfigQuery[*oc.Model_MultiKey] {
	return ygnmi.NewConfigQuery[*oc.Model_MultiKey](
		"Model_MultiKey",
		false,
		false,
		false,
		true,
		false,
		n,
		nil,
		nil,
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPathAny) Config() ygnmi.WildcardQuery[*oc.Model_MultiKey] {
	return ygnmi.NewWildcardQuery[*oc.Model_MultiKey](
		"Model_MultiKey",
		false,
		false,
		false,
		true,
		false,
		n,
		nil,
		nil,
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPathMap) State() ygnmi.SingletonQuery[map[oc.Model_MultiKey_Key]*oc.Model_MultiKey] {
	return ygnmi.NewSingletonQuery[map[oc.Model_MultiKey_Key]*oc.Model_MultiKey](
		"Model",
		true,
		false,
		false,
		true,
		true,
		n,
		func(gs ygot.ValidatedGoStruct) (map[oc.Model_MultiKey_Key]*oc.Model_MultiKey, bool) {
			ret := gs.(*oc.Model).MultiKey
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		&ygnmi.CompressionInfo{
			PreRelPath:  []string{"openconfig-withlistval:b"},
			PostRelPath: []string{"openconfig-withlistval:multi-key"},
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPathMapAny) State() ygnmi.WildcardQuery[map[oc.Model_MultiKey_Key]*oc.Model_MultiKey] {
	return ygnmi.NewWildcardQuery[map[oc.Model_MultiKey_Key]*oc.Model_MultiKey](
		"Model",
		true,
		false,
		false,
		true,
		true,
		n,
		func(gs ygot.ValidatedGoStruct) (map[oc.Model_MultiKey_Key]*oc.Model_MultiKey, bool) {
			ret := gs.(*oc.Model).MultiKey
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		&ygnmi.CompressionInfo{
			PreRelPath:  []string{"openconfig-withlistval:b"},
			PostRelPath: []string{"openconfig-withlistval:multi-key"},
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPathMap) Config() ygnmi.ConfigQuery[map[oc.Model_MultiKey_Key]*oc.Model_MultiKey] {
	return ygnmi.NewConfigQuery[map[oc.Model_MultiKey_Key]*oc.Model_MultiKey](
		"Model",
		false,
		false,
		false,
		true,
		true,
		n,
		func(gs ygot.ValidatedGoStruct) (map[oc.Model_MultiKey_Key]*oc.Model_MultiKey, bool) {
			ret := gs.(*oc.Model).MultiKey
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		&ygnmi.CompressionInfo{
			PreRelPath:  []string{"openconfig-withlistval:b"},
			PostRelPath: []string{"openconfig-withlistval:multi-key"},
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPathMapAny) Config() ygnmi.WildcardQuery[map[oc.Model_MultiKey_Key]*oc.Model_MultiKey] {
	return ygnmi.NewWildcardQuery[map[oc.Model_MultiKey_Key]*oc.Model_MultiKey](
		"Model",
		false,
		false,
		false,
		true,
		true,
		n,
		func(gs ygot.ValidatedGoStruct) (map[oc.Model_MultiKey_Key]*oc.Model_MultiKey, bool) {
			ret := gs.(*oc.Model).MultiKey
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		&ygnmi.CompressionInfo{
			PreRelPath:  []string{"openconfig-withlistval:b"},
			PostRelPath: []string{"openconfig-withlistval:multi-key"},
		},
	)
}
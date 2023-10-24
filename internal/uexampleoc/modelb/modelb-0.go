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

This package was generated by ygnmi version: (devel): (ygot: v0.29.12)
using the following YANG input files:
  - ../../pathgen/testdata/yang/openconfig-simple.yang
  - ../../pathgen/testdata/yang/openconfig-withlistval.yang
  - ../../pathgen/testdata/yang/openconfig-nested.yang

Imported modules were sourced from:
*/
package modelb

import (
	oc "github.com/openconfig/ygnmi/internal/uexampleoc"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// OpenconfigWithlistval_Model_BPath represents the /openconfig-withlistval/model/b YANG schema element.
type OpenconfigWithlistval_Model_BPath struct {
	*ygnmi.NodePath
	ygnmi.ConfigQuery[*oc.OpenconfigWithlistval_Model_B]
}

// OpenconfigWithlistval_Model_BPathAny represents the wildcard version of the /openconfig-withlistval/model/b YANG schema element.
type OpenconfigWithlistval_Model_BPathAny struct {
	*ygnmi.NodePath
	ygnmi.WildcardQuery[*oc.OpenconfigWithlistval_Model_B]
}

// MultiKeyAny (list):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "multi-key"
//	Path from root:       "/model/b/multi-key"
func (n *OpenconfigWithlistval_Model_BPath) MultiKeyAny() *OpenconfigWithlistval_Model_B_MultiKeyPathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKeyPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"multi-key"},
			map[string]interface{}{"key1": "*", "key2": "*"},
			n,
		),
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey](
		"OpenconfigWithlistval_Model_B_MultiKey",
		true,
		false,
		false,
		false,
		false,
		ps,
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

	return ps
}

// MultiKeyAny (list):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "multi-key"
//	Path from root:       "/model/b/multi-key"
func (n *OpenconfigWithlistval_Model_BPathAny) MultiKeyAny() *OpenconfigWithlistval_Model_B_MultiKeyPathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKeyPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"multi-key"},
			map[string]interface{}{"key1": "*", "key2": "*"},
			n,
		),
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey](
		"OpenconfigWithlistval_Model_B_MultiKey",
		true,
		false,
		false,
		false,
		false,
		ps,
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

	return ps
}

// WithKey1 sets OpenconfigWithlistval_Model_B_MultiKeyPathAny's key "key1" to the specified value.
// Key1: uint32
func (n *OpenconfigWithlistval_Model_B_MultiKeyPathAny) WithKey1(Key1 uint32) *OpenconfigWithlistval_Model_B_MultiKeyPathAny {
	ygnmi.ModifyKey(n.NodePath, "key1", Key1)
	return n
}

// WithKey2 sets OpenconfigWithlistval_Model_B_MultiKeyPathAny's key "key2" to the specified value.
// Key2: uint64
func (n *OpenconfigWithlistval_Model_B_MultiKeyPathAny) WithKey2(Key2 uint64) *OpenconfigWithlistval_Model_B_MultiKeyPathAny {
	ygnmi.ModifyKey(n.NodePath, "key2", Key2)
	return n
}

// MultiKey (list):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "multi-key"
//	Path from root:       "/model/b/multi-key"
//
//	Key1: uint32
//	Key2: uint64
func (n *OpenconfigWithlistval_Model_BPath) MultiKey(Key1 uint32, Key2 uint64) *OpenconfigWithlistval_Model_B_MultiKeyPath {
	ps := &OpenconfigWithlistval_Model_B_MultiKeyPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"multi-key"},
			map[string]interface{}{"key1": Key1, "key2": Key2},
			n,
		),
	}
	ps.ConfigQuery = ygnmi.NewConfigQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey](
		"OpenconfigWithlistval_Model_B_MultiKey",
		true,
		false,
		false,
		false,
		false,
		ps,
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

	return ps
}

// MultiKey (list):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "multi-key"
//	Path from root:       "/model/b/multi-key"
//
//	Key1: uint32
//	Key2: uint64
func (n *OpenconfigWithlistval_Model_BPathAny) MultiKey(Key1 uint32, Key2 uint64) *OpenconfigWithlistval_Model_B_MultiKeyPathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKeyPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"multi-key"},
			map[string]interface{}{"key1": Key1, "key2": Key2},
			n,
		),
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey](
		"OpenconfigWithlistval_Model_B_MultiKey",
		true,
		false,
		false,
		false,
		false,
		ps,
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

	return ps
}

// MultiKeyMap (list):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "multi-key"
//	Path from root:       "/model/b/multi-key"
func (n *OpenconfigWithlistval_Model_BPath) MultiKeyMap() *OpenconfigWithlistval_Model_B_MultiKeyPathMap {
	ps := &OpenconfigWithlistval_Model_B_MultiKeyPathMap{
		NodePath: ygnmi.NewNodePath(
			[]string{"multi-key"},
			map[string]interface{}{},
			n,
		),
	}
	ps.ConfigQuery = ygnmi.NewConfigQuery[map[oc.OpenconfigWithlistval_Model_B_MultiKey_Key]*oc.OpenconfigWithlistval_Model_B_MultiKey](
		"OpenconfigWithlistval_Model_B",
		true,
		false,
		false,
		false,
		true,
		ps,
		func(gs ygot.ValidatedGoStruct) (map[oc.OpenconfigWithlistval_Model_B_MultiKey_Key]*oc.OpenconfigWithlistval_Model_B_MultiKey, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B).MultiKey
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B) },
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

	return ps
}

// MultiKeyMap (list):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "multi-key"
//	Path from root:       "/model/b/multi-key"
func (n *OpenconfigWithlistval_Model_BPathAny) MultiKeyMap() *OpenconfigWithlistval_Model_B_MultiKeyPathMapAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKeyPathMapAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"multi-key"},
			map[string]interface{}{},
			n,
		),
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[map[oc.OpenconfigWithlistval_Model_B_MultiKey_Key]*oc.OpenconfigWithlistval_Model_B_MultiKey](
		"OpenconfigWithlistval_Model_B",
		true,
		false,
		false,
		false,
		true,
		ps,
		func(gs ygot.ValidatedGoStruct) (map[oc.OpenconfigWithlistval_Model_B_MultiKey_Key]*oc.OpenconfigWithlistval_Model_B_MultiKey, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B).MultiKey
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B) },
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

	return ps
}

func binarySliceToFloatSlice(in []oc.Binary) []float32 {
	converted := make([]float32, 0, len(in))
	for _, binary := range in {
		converted = append(converted, ygot.BinaryToFloat32(binary))
	}
	return converted
}

// OpenconfigWithlistval_Model_B_MultiKey_Key1Path represents the /openconfig-withlistval/model/b/multi-key/key1 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_Key1Path struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.ConfigQuery[uint32]
}

// OpenconfigWithlistval_Model_B_MultiKey_Key1PathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key/key1 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_Key1PathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.WildcardQuery[uint32]
}

// OpenconfigWithlistval_Model_B_MultiKey_Key2Path represents the /openconfig-withlistval/model/b/multi-key/key2 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_Key2Path struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.ConfigQuery[uint64]
}

// OpenconfigWithlistval_Model_B_MultiKey_Key2PathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key/key2 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_Key2PathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.WildcardQuery[uint64]
}

// OpenconfigWithlistval_Model_B_MultiKeyPath represents the /openconfig-withlistval/model/b/multi-key YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKeyPath struct {
	*ygnmi.NodePath
	ygnmi.ConfigQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey]
}

// OpenconfigWithlistval_Model_B_MultiKeyPathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKeyPathAny struct {
	*ygnmi.NodePath
	ygnmi.WildcardQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey]
}

// OpenconfigWithlistval_Model_B_MultiKeyPathMap represents the /openconfig-withlistval/model/b/multi-key YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKeyPathMap struct {
	*ygnmi.NodePath
	ygnmi.ConfigQuery[map[oc.OpenconfigWithlistval_Model_B_MultiKey_Key]*oc.OpenconfigWithlistval_Model_B_MultiKey]
}

// OpenconfigWithlistval_Model_B_MultiKeyPathMapAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKeyPathMapAny struct {
	*ygnmi.NodePath
	ygnmi.WildcardQuery[map[oc.OpenconfigWithlistval_Model_B_MultiKey_Key]*oc.OpenconfigWithlistval_Model_B_MultiKey]
}

// Config (container):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "config"
//	Path from root:       "/model/b/multi-key/config"
func (n *OpenconfigWithlistval_Model_B_MultiKeyPath) Config() *OpenconfigWithlistval_Model_B_MultiKey_ConfigPath {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_ConfigPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"config"},
			map[string]interface{}{},
			n,
		),
	}
	ps.ConfigQuery = ygnmi.NewConfigQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey_Config](
		"OpenconfigWithlistval_Model_B_MultiKey_Config",
		true,
		false,
		false,
		false,
		false,
		ps,
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

	return ps
}

// Config (container):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "config"
//	Path from root:       "/model/b/multi-key/config"
func (n *OpenconfigWithlistval_Model_B_MultiKeyPathAny) Config() *OpenconfigWithlistval_Model_B_MultiKey_ConfigPathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_ConfigPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"config"},
			map[string]interface{}{},
			n,
		),
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey_Config](
		"OpenconfigWithlistval_Model_B_MultiKey_Config",
		true,
		false,
		false,
		false,
		false,
		ps,
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

	return ps
}

// Key1 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key1"
//	Path from root:       "/model/b/multi-key/key1"
func (n *OpenconfigWithlistval_Model_B_MultiKeyPath) Key1() *OpenconfigWithlistval_Model_B_MultiKey_Key1Path {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_Key1Path{
		NodePath: ygnmi.NewNodePath(
			[]string{"key1"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.ConfigQuery = ygnmi.NewConfigQuery[uint32](
		"OpenconfigWithlistval_Model_B_MultiKey",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint32, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey).Key1
			if ret == nil {
				var zero uint32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey) },
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

	return ps
}

// Key1 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key1"
//	Path from root:       "/model/b/multi-key/key1"
func (n *OpenconfigWithlistval_Model_B_MultiKeyPathAny) Key1() *OpenconfigWithlistval_Model_B_MultiKey_Key1PathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_Key1PathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"key1"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[uint32](
		"OpenconfigWithlistval_Model_B_MultiKey",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint32, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey).Key1
			if ret == nil {
				var zero uint32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey) },
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

	return ps
}

// Key2 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key2"
//	Path from root:       "/model/b/multi-key/key2"
func (n *OpenconfigWithlistval_Model_B_MultiKeyPath) Key2() *OpenconfigWithlistval_Model_B_MultiKey_Key2Path {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_Key2Path{
		NodePath: ygnmi.NewNodePath(
			[]string{"key2"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.ConfigQuery = ygnmi.NewConfigQuery[uint64](
		"OpenconfigWithlistval_Model_B_MultiKey",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey).Key2
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey) },
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

	return ps
}

// Key2 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key2"
//	Path from root:       "/model/b/multi-key/key2"
func (n *OpenconfigWithlistval_Model_B_MultiKeyPathAny) Key2() *OpenconfigWithlistval_Model_B_MultiKey_Key2PathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_Key2PathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"key2"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[uint64](
		"OpenconfigWithlistval_Model_B_MultiKey",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey).Key2
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey) },
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

	return ps
}

// State (container):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "state"
//	Path from root:       "/model/b/multi-key/state"
func (n *OpenconfigWithlistval_Model_B_MultiKeyPath) State() *OpenconfigWithlistval_Model_B_MultiKey_StatePath {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_StatePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"state"},
			map[string]interface{}{},
			n,
		),
	}
	ps.SingletonQuery = ygnmi.NewSingletonQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey_State](
		"OpenconfigWithlistval_Model_B_MultiKey_State",
		true,
		false,
		false,
		false,
		false,
		ps,
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

	return ps
}

// State (container):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "state"
//	Path from root:       "/model/b/multi-key/state"
func (n *OpenconfigWithlistval_Model_B_MultiKeyPathAny) State() *OpenconfigWithlistval_Model_B_MultiKey_StatePathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_StatePathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"state"},
			map[string]interface{}{},
			n,
		),
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey_State](
		"OpenconfigWithlistval_Model_B_MultiKey_State",
		true,
		false,
		false,
		false,
		false,
		ps,
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

	return ps
}

// OpenconfigWithlistval_Model_B_MultiKey_Config_Key1Path represents the /openconfig-withlistval/model/b/multi-key/config/key1 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_Config_Key1Path struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.ConfigQuery[uint32]
}

// OpenconfigWithlistval_Model_B_MultiKey_Config_Key1PathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key/config/key1 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_Config_Key1PathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.WildcardQuery[uint32]
}

// OpenconfigWithlistval_Model_B_MultiKey_Config_Key2Path represents the /openconfig-withlistval/model/b/multi-key/config/key2 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_Config_Key2Path struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.ConfigQuery[uint64]
}

// OpenconfigWithlistval_Model_B_MultiKey_Config_Key2PathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key/config/key2 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_Config_Key2PathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.WildcardQuery[uint64]
}

// OpenconfigWithlistval_Model_B_MultiKey_ConfigPath represents the /openconfig-withlistval/model/b/multi-key/config YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_ConfigPath struct {
	*ygnmi.NodePath
	ygnmi.ConfigQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey_Config]
}

// OpenconfigWithlistval_Model_B_MultiKey_ConfigPathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key/config YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_ConfigPathAny struct {
	*ygnmi.NodePath
	ygnmi.WildcardQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey_Config]
}

// Key1 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key1"
//	Path from root:       "/model/b/multi-key/config/key1"
func (n *OpenconfigWithlistval_Model_B_MultiKey_ConfigPath) Key1() *OpenconfigWithlistval_Model_B_MultiKey_Config_Key1Path {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_Config_Key1Path{
		NodePath: ygnmi.NewNodePath(
			[]string{"key1"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.ConfigQuery = ygnmi.NewConfigQuery[uint32](
		"OpenconfigWithlistval_Model_B_MultiKey_Config",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint32, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey_Config).Key1
			if ret == nil {
				var zero uint32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey_Config) },
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

	return ps
}

// Key1 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key1"
//	Path from root:       "/model/b/multi-key/config/key1"
func (n *OpenconfigWithlistval_Model_B_MultiKey_ConfigPathAny) Key1() *OpenconfigWithlistval_Model_B_MultiKey_Config_Key1PathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_Config_Key1PathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"key1"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[uint32](
		"OpenconfigWithlistval_Model_B_MultiKey_Config",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint32, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey_Config).Key1
			if ret == nil {
				var zero uint32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey_Config) },
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

	return ps
}

// Key2 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key2"
//	Path from root:       "/model/b/multi-key/config/key2"
func (n *OpenconfigWithlistval_Model_B_MultiKey_ConfigPath) Key2() *OpenconfigWithlistval_Model_B_MultiKey_Config_Key2Path {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_Config_Key2Path{
		NodePath: ygnmi.NewNodePath(
			[]string{"key2"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.ConfigQuery = ygnmi.NewConfigQuery[uint64](
		"OpenconfigWithlistval_Model_B_MultiKey_Config",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey_Config).Key2
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey_Config) },
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

	return ps
}

// Key2 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key2"
//	Path from root:       "/model/b/multi-key/config/key2"
func (n *OpenconfigWithlistval_Model_B_MultiKey_ConfigPathAny) Key2() *OpenconfigWithlistval_Model_B_MultiKey_Config_Key2PathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_Config_Key2PathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"key2"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[uint64](
		"OpenconfigWithlistval_Model_B_MultiKey_Config",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey_Config).Key2
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey_Config) },
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

	return ps
}

// OpenconfigWithlistval_Model_B_MultiKey_State_Key1Path represents the /openconfig-withlistval/model/b/multi-key/state/key1 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_State_Key1Path struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.SingletonQuery[uint32]
}

// OpenconfigWithlistval_Model_B_MultiKey_State_Key1PathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key/state/key1 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_State_Key1PathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.WildcardQuery[uint32]
}

// OpenconfigWithlistval_Model_B_MultiKey_State_Key2Path represents the /openconfig-withlistval/model/b/multi-key/state/key2 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_State_Key2Path struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.SingletonQuery[uint64]
}

// OpenconfigWithlistval_Model_B_MultiKey_State_Key2PathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key/state/key2 YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_State_Key2PathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
	ygnmi.WildcardQuery[uint64]
}

// OpenconfigWithlistval_Model_B_MultiKey_StatePath represents the /openconfig-withlistval/model/b/multi-key/state YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_StatePath struct {
	*ygnmi.NodePath
	ygnmi.SingletonQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey_State]
}

// OpenconfigWithlistval_Model_B_MultiKey_StatePathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key/state YANG schema element.
type OpenconfigWithlistval_Model_B_MultiKey_StatePathAny struct {
	*ygnmi.NodePath
	ygnmi.WildcardQuery[*oc.OpenconfigWithlistval_Model_B_MultiKey_State]
}

// Key1 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key1"
//	Path from root:       "/model/b/multi-key/state/key1"
func (n *OpenconfigWithlistval_Model_B_MultiKey_StatePath) Key1() *OpenconfigWithlistval_Model_B_MultiKey_State_Key1Path {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_State_Key1Path{
		NodePath: ygnmi.NewNodePath(
			[]string{"key1"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.SingletonQuery = ygnmi.NewSingletonQuery[uint32](
		"OpenconfigWithlistval_Model_B_MultiKey_State",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint32, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey_State).Key1
			if ret == nil {
				var zero uint32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey_State) },
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

	return ps
}

// Key1 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key1"
//	Path from root:       "/model/b/multi-key/state/key1"
func (n *OpenconfigWithlistval_Model_B_MultiKey_StatePathAny) Key1() *OpenconfigWithlistval_Model_B_MultiKey_State_Key1PathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_State_Key1PathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"key1"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[uint32](
		"OpenconfigWithlistval_Model_B_MultiKey_State",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint32, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey_State).Key1
			if ret == nil {
				var zero uint32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey_State) },
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

	return ps
}

// Key2 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key2"
//	Path from root:       "/model/b/multi-key/state/key2"
func (n *OpenconfigWithlistval_Model_B_MultiKey_StatePath) Key2() *OpenconfigWithlistval_Model_B_MultiKey_State_Key2Path {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_State_Key2Path{
		NodePath: ygnmi.NewNodePath(
			[]string{"key2"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.SingletonQuery = ygnmi.NewSingletonQuery[uint64](
		"OpenconfigWithlistval_Model_B_MultiKey_State",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey_State).Key2
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey_State) },
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

	return ps
}

// Key2 (leaf):
//
//	Defining module:      "openconfig-withlistval"
//	Instantiating module: "openconfig-withlistval"
//	Path from parent:     "key2"
//	Path from root:       "/model/b/multi-key/state/key2"
func (n *OpenconfigWithlistval_Model_B_MultiKey_StatePathAny) Key2() *OpenconfigWithlistval_Model_B_MultiKey_State_Key2PathAny {
	ps := &OpenconfigWithlistval_Model_B_MultiKey_State_Key2PathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"key2"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	ps.WildcardQuery = ygnmi.NewWildcardQuery[uint64](
		"OpenconfigWithlistval_Model_B_MultiKey_State",
		true,
		true,
		true,
		false,
		false,
		ps,
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.OpenconfigWithlistval_Model_B_MultiKey_State).Key2
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.OpenconfigWithlistval_Model_B_MultiKey_State) },
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

	return ps
}

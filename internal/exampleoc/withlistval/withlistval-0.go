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

/*
Package withlistval is a generated package which contains definitions
of structs which generate gNMI paths for a YANG schema. The generated paths are
based on a compressed form of the schema.

This package was generated by ygnmi version: (devel): (ygot: v0.21.0)
using the following YANG input files:
	- ../../pathgen/testdata/yang/openconfig-simple.yang
	- ../../pathgen/testdata/yang/openconfig-withlistval.yang
	- ../../pathgen/testdata/yang/openconfig-nested.yang
Imported modules were sourced from:
*/
package withlistval

import (
	oc "github.com/openconfig/ygnmi/internal/exampleoc"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// ModelPath represents the /openconfig-withlistval/model YANG schema element.
type ModelPath struct {
	*ygnmi.NodePath
}

// ModelPathAny represents the wildcard version of the /openconfig-withlistval/model YANG schema element.
type ModelPathAny struct {
	*ygnmi.NodePath
}

// MultiKeyAny (list):
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "b/multi-key"
// 	Path from root:       "/model/b/multi-key"
func (n *ModelPath) MultiKeyAny() *Model_MultiKeyPathAny {
	return &Model_MultiKeyPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"b", "multi-key"},
			map[string]interface{}{"key1": "*", "key2": "*"},
			n,
		),
	}
}

// MultiKeyAny (list):
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "b/multi-key"
// 	Path from root:       "/model/b/multi-key"
func (n *ModelPathAny) MultiKeyAny() *Model_MultiKeyPathAny {
	return &Model_MultiKeyPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"b", "multi-key"},
			map[string]interface{}{"key1": "*", "key2": "*"},
			n,
		),
	}
}

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

// SingleKeyAny (list):
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "a/single-key"
// 	Path from root:       "/model/a/single-key"
//
// 	Key (wildcarded): string
func (n *ModelPath) SingleKeyAny() *Model_SingleKeyPathAny {
	return &Model_SingleKeyPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"a", "single-key"},
			map[string]interface{}{"key": "*"},
			n,
		),
	}
}

// SingleKeyAny (list):
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "a/single-key"
// 	Path from root:       "/model/a/single-key"
//
// 	Key (wildcarded): string
func (n *ModelPathAny) SingleKeyAny() *Model_SingleKeyPathAny {
	return &Model_SingleKeyPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"a", "single-key"},
			map[string]interface{}{"key": "*"},
			n,
		),
	}
}

// SingleKey (list):
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "a/single-key"
// 	Path from root:       "/model/a/single-key"
//
// 	Key: string
func (n *ModelPath) SingleKey(Key string) *Model_SingleKeyPath {
	return &Model_SingleKeyPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"a", "single-key"},
			map[string]interface{}{"key": Key},
			n,
		),
	}
}

// SingleKey (list):
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "a/single-key"
// 	Path from root:       "/model/a/single-key"
//
// 	Key: string
func (n *ModelPathAny) SingleKey(Key string) *Model_SingleKeyPathAny {
	return &Model_SingleKeyPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"a", "single-key"},
			map[string]interface{}{"key": Key},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *ModelPath) State() ygnmi.SingletonQuery[*oc.Model] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.Model](
		"Model",
		true,
		n,
		nil,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *ModelPathAny) State() ygnmi.WildcardQuery[*oc.Model] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Model](
		"Model",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *ModelPath) Config() ygnmi.ConfigQuery[*oc.Model] {
	return ygnmi.NewNonLeafConfigQuery[*oc.Model](
		"Model",
		false,
		n,
		nil,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *ModelPathAny) Config() ygnmi.WildcardQuery[*oc.Model] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Model](
		"Model",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
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
func (n *Model_MultiKeyPath) State() ygnmi.SingletonQuery[*oc.Model_MultiKey] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.Model_MultiKey](
		"Model_MultiKey",
		true,
		n,
		nil,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPathAny) State() ygnmi.WildcardQuery[*oc.Model_MultiKey] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Model_MultiKey](
		"Model_MultiKey",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPath) Config() ygnmi.ConfigQuery[*oc.Model_MultiKey] {
	return ygnmi.NewNonLeafConfigQuery[*oc.Model_MultiKey](
		"Model_MultiKey",
		false,
		n,
		nil,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Model_MultiKeyPathAny) Config() ygnmi.WildcardQuery[*oc.Model_MultiKey] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Model_MultiKey](
		"Model_MultiKey",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "state/key1"
// 	Path from root:       "/model/b/multi-key/state/key1"
func (n *Model_MultiKey_Key1Path) State() ygnmi.SingletonQuery[uint32] {
	return ygnmi.NewLeafSingletonQuery[uint32](
		"Model_MultiKey",
		true,
		true,
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
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "state/key1"
// 	Path from root:       "/model/b/multi-key/state/key1"
func (n *Model_MultiKey_Key1PathAny) State() ygnmi.WildcardQuery[uint32] {
	return ygnmi.NewLeafWildcardQuery[uint32](
		"Model_MultiKey",
		true,
		true,
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
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "config/key1"
// 	Path from root:       "/model/b/multi-key/config/key1"
func (n *Model_MultiKey_Key1Path) Config() ygnmi.ConfigQuery[uint32] {
	return ygnmi.NewLeafConfigQuery[uint32](
		"Model_MultiKey",
		false,
		true,
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
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "config/key1"
// 	Path from root:       "/model/b/multi-key/config/key1"
func (n *Model_MultiKey_Key1PathAny) Config() ygnmi.WildcardQuery[uint32] {
	return ygnmi.NewLeafWildcardQuery[uint32](
		"Model_MultiKey",
		false,
		true,
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
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "state/key2"
// 	Path from root:       "/model/b/multi-key/state/key2"
func (n *Model_MultiKey_Key2Path) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewLeafSingletonQuery[uint64](
		"Model_MultiKey",
		true,
		true,
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
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "state/key2"
// 	Path from root:       "/model/b/multi-key/state/key2"
func (n *Model_MultiKey_Key2PathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewLeafWildcardQuery[uint64](
		"Model_MultiKey",
		true,
		true,
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
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "config/key2"
// 	Path from root:       "/model/b/multi-key/config/key2"
func (n *Model_MultiKey_Key2Path) Config() ygnmi.ConfigQuery[uint64] {
	return ygnmi.NewLeafConfigQuery[uint64](
		"Model_MultiKey",
		false,
		true,
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
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "config/key2"
// 	Path from root:       "/model/b/multi-key/config/key2"
func (n *Model_MultiKey_Key2PathAny) Config() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewLeafWildcardQuery[uint64](
		"Model_MultiKey",
		false,
		true,
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
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
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

// Model_MultiKeyPath represents the /openconfig-withlistval/model/b/multi-key YANG schema element.
type Model_MultiKeyPath struct {
	*ygnmi.NodePath
}

// Model_MultiKeyPathAny represents the wildcard version of the /openconfig-withlistval/model/b/multi-key YANG schema element.
type Model_MultiKeyPathAny struct {
	*ygnmi.NodePath
}

// Key1 corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
func (n *Model_MultiKeyPath) Key1() *Model_MultiKey_Key1Path {
	return &Model_MultiKey_Key1Path{
		parent: n,
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "key1"},
			map[string]interface{}{},
			n,
		),
	}
}

// Key1 corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
func (n *Model_MultiKeyPathAny) Key1() *Model_MultiKey_Key1PathAny {
	return &Model_MultiKey_Key1PathAny{
		parent: n,
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "key1"},
			map[string]interface{}{},
			n,
		),
	}
}

// Key2 corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
func (n *Model_MultiKeyPath) Key2() *Model_MultiKey_Key2Path {
	return &Model_MultiKey_Key2Path{
		parent: n,
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "key2"},
			map[string]interface{}{},
			n,
		),
	}
}

// Key2 corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
func (n *Model_MultiKeyPathAny) Key2() *Model_MultiKey_Key2PathAny {
	return &Model_MultiKey_Key2PathAny{
		parent: n,
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "key2"},
			map[string]interface{}{},
			n,
		),
	}
}

// Model_SingleKey_KeyPath represents the /openconfig-withlistval/model/a/single-key/state/key YANG schema element.
type Model_SingleKey_KeyPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Model_SingleKey_KeyPathAny represents the wildcard version of the /openconfig-withlistval/model/a/single-key/state/key YANG schema element.
type Model_SingleKey_KeyPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
func (n *Model_SingleKeyPath) State() ygnmi.SingletonQuery[*oc.Model_SingleKey] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.Model_SingleKey](
		"Model_SingleKey",
		true,
		n,
		nil,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *Model_SingleKeyPathAny) State() ygnmi.WildcardQuery[*oc.Model_SingleKey] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Model_SingleKey](
		"Model_SingleKey",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Model_SingleKeyPath) Config() ygnmi.ConfigQuery[*oc.Model_SingleKey] {
	return ygnmi.NewNonLeafConfigQuery[*oc.Model_SingleKey](
		"Model_SingleKey",
		false,
		n,
		nil,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Model_SingleKeyPathAny) Config() ygnmi.WildcardQuery[*oc.Model_SingleKey] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Model_SingleKey](
		"Model_SingleKey",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "state/key"
// 	Path from root:       "/model/a/single-key/state/key"
func (n *Model_SingleKey_KeyPath) State() ygnmi.SingletonQuery[string] {
	return ygnmi.NewLeafSingletonQuery[string](
		"Model_SingleKey",
		true,
		true,
		ygnmi.NewNodePath(
			[]string{"state", "key"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Model_SingleKey).Key
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_SingleKey) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "state/key"
// 	Path from root:       "/model/a/single-key/state/key"
func (n *Model_SingleKey_KeyPathAny) State() ygnmi.WildcardQuery[string] {
	return ygnmi.NewLeafWildcardQuery[string](
		"Model_SingleKey",
		true,
		true,
		ygnmi.NewNodePath(
			[]string{"state", "key"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Model_SingleKey).Key
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_SingleKey) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "config/key"
// 	Path from root:       "/model/a/single-key/config/key"
func (n *Model_SingleKey_KeyPath) Config() ygnmi.ConfigQuery[string] {
	return ygnmi.NewLeafConfigQuery[string](
		"Model_SingleKey",
		false,
		true,
		ygnmi.NewNodePath(
			[]string{"config", "key"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Model_SingleKey).Key
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_SingleKey) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "config/key"
// 	Path from root:       "/model/a/single-key/config/key"
func (n *Model_SingleKey_KeyPathAny) Config() ygnmi.WildcardQuery[string] {
	return ygnmi.NewLeafWildcardQuery[string](
		"Model_SingleKey",
		false,
		true,
		ygnmi.NewNodePath(
			[]string{"config", "key"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Model_SingleKey).Key
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_SingleKey) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "state/value"
// 	Path from root:       "/model/a/single-key/state/value"
func (n *Model_SingleKey_ValuePath) State() ygnmi.SingletonQuery[int64] {
	return ygnmi.NewLeafSingletonQuery[int64](
		"Model_SingleKey",
		true,
		true,
		ygnmi.NewNodePath(
			[]string{"state", "value"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int64, bool) {
			ret := gs.(*oc.Model_SingleKey).Value
			if ret == nil {
				var zero int64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_SingleKey) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "state/value"
// 	Path from root:       "/model/a/single-key/state/value"
func (n *Model_SingleKey_ValuePathAny) State() ygnmi.WildcardQuery[int64] {
	return ygnmi.NewLeafWildcardQuery[int64](
		"Model_SingleKey",
		true,
		true,
		ygnmi.NewNodePath(
			[]string{"state", "value"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int64, bool) {
			ret := gs.(*oc.Model_SingleKey).Value
			if ret == nil {
				var zero int64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_SingleKey) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "config/value"
// 	Path from root:       "/model/a/single-key/config/value"
func (n *Model_SingleKey_ValuePath) Config() ygnmi.ConfigQuery[int64] {
	return ygnmi.NewLeafConfigQuery[int64](
		"Model_SingleKey",
		false,
		true,
		ygnmi.NewNodePath(
			[]string{"config", "value"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int64, bool) {
			ret := gs.(*oc.Model_SingleKey).Value
			if ret == nil {
				var zero int64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_SingleKey) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// 	Defining module:      "openconfig-withlistval"
// 	Instantiating module: "openconfig-withlistval"
// 	Path from parent:     "config/value"
// 	Path from root:       "/model/a/single-key/config/value"
func (n *Model_SingleKey_ValuePathAny) Config() ygnmi.WildcardQuery[int64] {
	return ygnmi.NewLeafWildcardQuery[int64](
		"Model_SingleKey",
		false,
		true,
		ygnmi.NewNodePath(
			[]string{"config", "value"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int64, bool) {
			ret := gs.(*oc.Model_SingleKey).Value
			if ret == nil {
				var zero int64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Model_SingleKey) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Model_SingleKey_ValuePath represents the /openconfig-withlistval/model/a/single-key/state/value YANG schema element.
type Model_SingleKey_ValuePath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Model_SingleKey_ValuePathAny represents the wildcard version of the /openconfig-withlistval/model/a/single-key/state/value YANG schema element.
type Model_SingleKey_ValuePathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Model_SingleKeyPath represents the /openconfig-withlistval/model/a/single-key YANG schema element.
type Model_SingleKeyPath struct {
	*ygnmi.NodePath
}

// Model_SingleKeyPathAny represents the wildcard version of the /openconfig-withlistval/model/a/single-key YANG schema element.
type Model_SingleKeyPathAny struct {
	*ygnmi.NodePath
}

// Key corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
func (n *Model_SingleKeyPath) Key() *Model_SingleKey_KeyPath {
	return &Model_SingleKey_KeyPath{
		parent: n,
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "key"},
			map[string]interface{}{},
			n,
		),
	}
}

// Key corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
func (n *Model_SingleKeyPathAny) Key() *Model_SingleKey_KeyPathAny {
	return &Model_SingleKey_KeyPathAny{
		parent: n,
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "key"},
			map[string]interface{}{},
			n,
		),
	}
}

// Value corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
func (n *Model_SingleKeyPath) Value() *Model_SingleKey_ValuePath {
	return &Model_SingleKey_ValuePath{
		parent: n,
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "value"},
			map[string]interface{}{},
			n,
		),
	}
}

// Value corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
func (n *Model_SingleKeyPathAny) Value() *Model_SingleKey_ValuePathAny {
	return &Model_SingleKey_ValuePathAny{
		parent: n,
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "value"},
			map[string]interface{}{},
			n,
		),
	}
}

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
Package simple is a generated package which contains definitions
of structs which generate gNMI paths for a YANG schema. The generated paths are
based on a compressed form of the schema.

This package was generated by ygnmi version: (devel): (ygot: v0.20.0)
using the following YANG input files:
  - ../../pathgen/testdata/yang/openconfig-simple.yang
  - ../../pathgen/testdata/yang/openconfig-withlistval.yang

Imported modules were sourced from:
*/
package simple

import (
	"reflect"

	oc "github.com/openconfig/ygnmi/internal/exampleoc"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// ParentPath represents the /openconfig-simple/parent YANG schema element.
type ParentPath struct {
	*ygot.NodePath
}

// ParentPathAny represents the wildcard version of the /openconfig-simple/parent YANG schema element.
type ParentPathAny struct {
	*ygot.NodePath
}

// Child (container):
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "child"
// Path from root: "/parent/child"
func (n *ParentPath) Child() *Parent_ChildPath {
	return &Parent_ChildPath{
		NodePath: ygot.NewNodePath(
			[]string{"child"},
			map[string]interface{}{},
			n,
		),
	}
}

// Child (container):
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "child"
// Path from root: "/parent/child"
func (n *ParentPathAny) Child() *Parent_ChildPathAny {
	return &Parent_ChildPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"child"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *ParentPath) State() ygnmi.SingletonQuery[*oc.Parent] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.Parent](
		"Parent",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *ParentPathAny) State() ygnmi.WildcardQuery[*oc.Parent] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Parent](
		"Parent",
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
func (n *ParentPath) Config() ygnmi.ConfigQuery[*oc.Parent] {
	return ygnmi.NewNonLeafConfigQuery[*oc.Parent](
		"Parent",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *ParentPathAny) Config() ygnmi.WildcardQuery[*oc.Parent] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Parent](
		"Parent",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Parent_ChildPath represents the /openconfig-simple/parent/child YANG schema element.
type Parent_ChildPath struct {
	*ygot.NodePath
}

// Parent_ChildPathAny represents the wildcard version of the /openconfig-simple/parent/child YANG schema element.
type Parent_ChildPathAny struct {
	*ygot.NodePath
}

// Parent_Child_FourPath represents the /openconfig-simple/parent/child/state/four YANG schema element.
type Parent_Child_FourPath struct {
	parent ygot.PathStruct
}

// Parent_Child_FourPathAny represents the wildcard version of the /openconfig-simple/parent/child/state/four YANG schema element.
type Parent_Child_FourPathAny struct {
	parent ygot.PathStruct
}

// Parent_Child_OnePath represents the /openconfig-simple/parent/child/state/one YANG schema element.
type Parent_Child_OnePath struct {
	parent ygot.PathStruct
}

// Parent_Child_OnePathAny represents the wildcard version of the /openconfig-simple/parent/child/state/one YANG schema element.
type Parent_Child_OnePathAny struct {
	parent ygot.PathStruct
}

// Parent_Child_ThreePath represents the /openconfig-simple/parent/child/state/three YANG schema element.
type Parent_Child_ThreePath struct {
	parent ygot.PathStruct
}

// Parent_Child_ThreePathAny represents the wildcard version of the /openconfig-simple/parent/child/state/three YANG schema element.
type Parent_Child_ThreePathAny struct {
	parent ygot.PathStruct
}

// Parent_Child_TwoPath represents the /openconfig-simple/parent/child/state/two YANG schema element.
type Parent_Child_TwoPath struct {
	parent ygot.PathStruct
}

// Parent_Child_TwoPathAny represents the wildcard version of the /openconfig-simple/parent/child/state/two YANG schema element.
type Parent_Child_TwoPathAny struct {
	parent ygot.PathStruct
}

// Four corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildPath) Four() *Parent_Child_FourPath {
	return &Parent_Child_FourPath{
		parent: n,
	}
}

// Four corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildPathAny) Four() *Parent_Child_FourPathAny {
	return &Parent_Child_FourPathAny{
		parent: n,
	}
}

// One corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildPath) One() *Parent_Child_OnePath {
	return &Parent_Child_OnePath{
		parent: n,
	}
}

// One corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildPathAny) One() *Parent_Child_OnePathAny {
	return &Parent_Child_OnePathAny{
		parent: n,
	}
}

// Three corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildPath) Three() *Parent_Child_ThreePath {
	return &Parent_Child_ThreePath{
		parent: n,
	}
}

// Three corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildPathAny) Three() *Parent_Child_ThreePathAny {
	return &Parent_Child_ThreePathAny{
		parent: n,
	}
}

// Two corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildPath) Two() *Parent_Child_TwoPath {
	return &Parent_Child_TwoPath{
		parent: n,
	}
}

// Two corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildPathAny) Two() *Parent_Child_TwoPathAny {
	return &Parent_Child_TwoPathAny{
		parent: n,
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *Parent_ChildPath) State() ygnmi.SingletonQuery[*oc.Parent_Child] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.Parent_Child](
		"Parent_Child",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *Parent_ChildPathAny) State() ygnmi.WildcardQuery[*oc.Parent_Child] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Parent_Child](
		"Parent_Child",
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
func (n *Parent_ChildPath) Config() ygnmi.ConfigQuery[*oc.Parent_Child] {
	return ygnmi.NewNonLeafConfigQuery[*oc.Parent_Child](
		"Parent_Child",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Parent_ChildPathAny) Config() ygnmi.WildcardQuery[*oc.Parent_Child] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Parent_Child](
		"Parent_Child",
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
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "state/four"
// Path from root: "/parent/child/state/four"
func (n *Parent_Child_FourPath) State() ygnmi.SingletonQuery[oc.Binary] {
	return ygnmi.NewLeafSingletonQuery[oc.Binary](
		"Parent_Child",
		true,
		false,
		ygot.NewNodePath(
			[]string{"state", "four"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (oc.Binary, bool) {
			ret := gs.(*oc.Parent_Child).Four
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "state/four"
// Path from root: "/parent/child/state/four"
func (n *Parent_Child_FourPathAny) State() ygnmi.WildcardQuery[oc.Binary] {
	return ygnmi.NewLeafWildcardQuery[oc.Binary](
		"Parent_Child",
		true,
		false,
		ygot.NewNodePath(
			[]string{"state", "four"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (oc.Binary, bool) {
			ret := gs.(*oc.Parent_Child).Four
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "config/four"
// Path from root: "/parent/child/config/four"
func (n *Parent_Child_FourPath) Config() ygnmi.ConfigQuery[oc.Binary] {
	return ygnmi.NewLeafConfigQuery[oc.Binary](
		"Parent_Child",
		false,
		false,
		ygot.NewNodePath(
			[]string{"config", "four"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (oc.Binary, bool) {
			ret := gs.(*oc.Parent_Child).Four
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "config/four"
// Path from root: "/parent/child/config/four"
func (n *Parent_Child_FourPathAny) Config() ygnmi.WildcardQuery[oc.Binary] {
	return ygnmi.NewLeafWildcardQuery[oc.Binary](
		"Parent_Child",
		false,
		false,
		ygot.NewNodePath(
			[]string{"config", "four"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (oc.Binary, bool) {
			ret := gs.(*oc.Parent_Child).Four
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "state/one"
// Path from root: "/parent/child/state/one"
func (n *Parent_Child_OnePath) State() ygnmi.SingletonQuery[string] {
	return ygnmi.NewLeafSingletonQuery[string](
		"Parent_Child",
		true,
		true,
		ygot.NewNodePath(
			[]string{"state", "one"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Parent_Child).One
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "state/one"
// Path from root: "/parent/child/state/one"
func (n *Parent_Child_OnePathAny) State() ygnmi.WildcardQuery[string] {
	return ygnmi.NewLeafWildcardQuery[string](
		"Parent_Child",
		true,
		true,
		ygot.NewNodePath(
			[]string{"state", "one"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Parent_Child).One
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "config/one"
// Path from root: "/parent/child/config/one"
func (n *Parent_Child_OnePath) Config() ygnmi.ConfigQuery[string] {
	return ygnmi.NewLeafConfigQuery[string](
		"Parent_Child",
		false,
		true,
		ygot.NewNodePath(
			[]string{"config", "one"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Parent_Child).One
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "config/one"
// Path from root: "/parent/child/config/one"
func (n *Parent_Child_OnePathAny) Config() ygnmi.WildcardQuery[string] {
	return ygnmi.NewLeafWildcardQuery[string](
		"Parent_Child",
		false,
		true,
		ygot.NewNodePath(
			[]string{"config", "one"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Parent_Child).One
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "state/three"
// Path from root: "/parent/child/state/three"
func (n *Parent_Child_ThreePath) State() ygnmi.SingletonQuery[oc.E_Child_Three] {
	return ygnmi.NewLeafSingletonQuery[oc.E_Child_Three](
		"Parent_Child",
		true,
		false,
		ygot.NewNodePath(
			[]string{"state", "three"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (oc.E_Child_Three, bool) {
			ret := gs.(*oc.Parent_Child).Three
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "state/three"
// Path from root: "/parent/child/state/three"
func (n *Parent_Child_ThreePathAny) State() ygnmi.WildcardQuery[oc.E_Child_Three] {
	return ygnmi.NewLeafWildcardQuery[oc.E_Child_Three](
		"Parent_Child",
		true,
		false,
		ygot.NewNodePath(
			[]string{"state", "three"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (oc.E_Child_Three, bool) {
			ret := gs.(*oc.Parent_Child).Three
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "config/three"
// Path from root: "/parent/child/config/three"
func (n *Parent_Child_ThreePath) Config() ygnmi.ConfigQuery[oc.E_Child_Three] {
	return ygnmi.NewLeafConfigQuery[oc.E_Child_Three](
		"Parent_Child",
		false,
		false,
		ygot.NewNodePath(
			[]string{"config", "three"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (oc.E_Child_Three, bool) {
			ret := gs.(*oc.Parent_Child).Three
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "config/three"
// Path from root: "/parent/child/config/three"
func (n *Parent_Child_ThreePathAny) Config() ygnmi.WildcardQuery[oc.E_Child_Three] {
	return ygnmi.NewLeafWildcardQuery[oc.E_Child_Three](
		"Parent_Child",
		false,
		false,
		ygot.NewNodePath(
			[]string{"config", "three"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (oc.E_Child_Three, bool) {
			ret := gs.(*oc.Parent_Child).Three
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "state/two"
// Path from root: "/parent/child/state/two"
func (n *Parent_Child_TwoPath) State() ygnmi.SingletonQuery[string] {
	return ygnmi.NewLeafSingletonQuery[string](
		"Parent_Child",
		true,
		true,
		ygot.NewNodePath(
			[]string{"state", "two"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Parent_Child).Two
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "state/two"
// Path from root: "/parent/child/state/two"
func (n *Parent_Child_TwoPathAny) State() ygnmi.WildcardQuery[string] {
	return ygnmi.NewLeafWildcardQuery[string](
		"Parent_Child",
		true,
		true,
		ygot.NewNodePath(
			[]string{"state", "two"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Parent_Child).Two
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Parent_Child) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// RemoteContainerPath represents the /openconfig-simple/remote-container YANG schema element.
type RemoteContainerPath struct {
	*ygot.NodePath
}

// RemoteContainerPathAny represents the wildcard version of the /openconfig-simple/remote-container YANG schema element.
type RemoteContainerPathAny struct {
	*ygot.NodePath
}

// RemoteContainer_ALeafPath represents the /openconfig-simple/remote-container/state/a-leaf YANG schema element.
type RemoteContainer_ALeafPath struct {
	parent ygot.PathStruct
}

// RemoteContainer_ALeafPathAny represents the wildcard version of the /openconfig-simple/remote-container/state/a-leaf YANG schema element.
type RemoteContainer_ALeafPathAny struct {
	parent ygot.PathStruct
}

// ALeaf corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *RemoteContainerPath) ALeaf() *RemoteContainer_ALeafPath {
	return &RemoteContainer_ALeafPath{
		parent: n,
	}
}

// ALeaf corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *RemoteContainerPathAny) ALeaf() *RemoteContainer_ALeafPathAny {
	return &RemoteContainer_ALeafPathAny{
		parent: n,
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *RemoteContainerPath) State() ygnmi.SingletonQuery[*oc.RemoteContainer] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.RemoteContainer](
		"RemoteContainer",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *RemoteContainerPathAny) State() ygnmi.WildcardQuery[*oc.RemoteContainer] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.RemoteContainer](
		"RemoteContainer",
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
func (n *RemoteContainerPath) Config() ygnmi.ConfigQuery[*oc.RemoteContainer] {
	return ygnmi.NewNonLeafConfigQuery[*oc.RemoteContainer](
		"RemoteContainer",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *RemoteContainerPathAny) Config() ygnmi.WildcardQuery[*oc.RemoteContainer] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.RemoteContainer](
		"RemoteContainer",
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
// ----------------------------------------
// Defining module: "openconfig-remote"
// Instantiating module: "openconfig-simple"
// Path from parent: "state/a-leaf"
// Path from root: "/remote-container/state/a-leaf"
func (n *RemoteContainer_ALeafPath) State() ygnmi.SingletonQuery[string] {
	return ygnmi.NewLeafSingletonQuery[string](
		"RemoteContainer",
		true,
		true,
		ygot.NewNodePath(
			[]string{"state", "a-leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.RemoteContainer).ALeaf
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RemoteContainer) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-remote"
// Instantiating module: "openconfig-simple"
// Path from parent: "state/a-leaf"
// Path from root: "/remote-container/state/a-leaf"
func (n *RemoteContainer_ALeafPathAny) State() ygnmi.WildcardQuery[string] {
	return ygnmi.NewLeafWildcardQuery[string](
		"RemoteContainer",
		true,
		true,
		ygot.NewNodePath(
			[]string{"state", "a-leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.RemoteContainer).ALeaf
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RemoteContainer) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-remote"
// Instantiating module: "openconfig-simple"
// Path from parent: "config/a-leaf"
// Path from root: "/remote-container/config/a-leaf"
func (n *RemoteContainer_ALeafPath) Config() ygnmi.ConfigQuery[string] {
	return ygnmi.NewLeafConfigQuery[string](
		"RemoteContainer",
		false,
		true,
		ygot.NewNodePath(
			[]string{"config", "a-leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.RemoteContainer).ALeaf
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RemoteContainer) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-remote"
// Instantiating module: "openconfig-simple"
// Path from parent: "config/a-leaf"
// Path from root: "/remote-container/config/a-leaf"
func (n *RemoteContainer_ALeafPathAny) Config() ygnmi.WildcardQuery[string] {
	return ygnmi.NewLeafWildcardQuery[string](
		"RemoteContainer",
		false,
		true,
		ygot.NewNodePath(
			[]string{"config", "a-leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.RemoteContainer).ALeaf
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RemoteContainer) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

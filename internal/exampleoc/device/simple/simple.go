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

This package was generated by /usr/local/google/home/dgrau/go/pkg/mod/github.com/openconfig/ygot@v0.18.0/genutil/names.go
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

// Parent represents the /openconfig-simple/parent YANG schema element.
type Parent struct {
	*ygot.NodePath
}

// ParentAny represents the wildcard version of the /openconfig-simple/parent YANG schema element.
type ParentAny struct {
	*ygot.NodePath
}

// Child (container):
// ----------------------------------------
// Defining module: "openconfig-simple"
// Instantiating module: "openconfig-simple"
// Path from parent: "child"
// Path from root: "/parent/child"
func (n *Parent) Child() *Parent_Child {
	return &Parent_Child{
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
func (n *ParentAny) Child() *Parent_ChildAny {
	return &Parent_ChildAny{
		NodePath: ygot.NewNodePath(
			[]string{"child"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *Parent) State() ygnmi.SingletonQuery[*oc.Parent] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.Parent](
		"Parent",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *ParentAny) State() ygnmi.WildcardQuery[*oc.Parent] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Parent](
		"Parent",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Parent) Config() ygnmi.ConfigQuery[*oc.Parent] {
	return ygnmi.NewNonLeafConfigQuery[*oc.Parent](
		"Parent",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *ParentAny) Config() ygnmi.WildcardQuery[*oc.Parent] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Parent](
		"Parent",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Parent_Child represents the /openconfig-simple/parent/child YANG schema element.
type Parent_Child struct {
	*ygot.NodePath
}

// Parent_ChildAny represents the wildcard version of the /openconfig-simple/parent/child YANG schema element.
type Parent_ChildAny struct {
	*ygot.NodePath
}

// Parent_Child_Four represents the /openconfig-simple/parent/child/state/four YANG schema element.
type Parent_Child_Four struct {
	parent ygot.PathStruct
}

// Parent_Child_FourAny represents the wildcard version of the /openconfig-simple/parent/child/state/four YANG schema element.
type Parent_Child_FourAny struct {
	parent ygot.PathStruct
}

// Parent_Child_One represents the /openconfig-simple/parent/child/state/one YANG schema element.
type Parent_Child_One struct {
	parent ygot.PathStruct
}

// Parent_Child_OneAny represents the wildcard version of the /openconfig-simple/parent/child/state/one YANG schema element.
type Parent_Child_OneAny struct {
	parent ygot.PathStruct
}

// Parent_Child_Three represents the /openconfig-simple/parent/child/state/three YANG schema element.
type Parent_Child_Three struct {
	parent ygot.PathStruct
}

// Parent_Child_ThreeAny represents the wildcard version of the /openconfig-simple/parent/child/state/three YANG schema element.
type Parent_Child_ThreeAny struct {
	parent ygot.PathStruct
}

// Parent_Child_Two represents the /openconfig-simple/parent/child/state/two YANG schema element.
type Parent_Child_Two struct {
	parent ygot.PathStruct
}

// Parent_Child_TwoAny represents the wildcard version of the /openconfig-simple/parent/child/state/two YANG schema element.
type Parent_Child_TwoAny struct {
	parent ygot.PathStruct
}

// Four corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_Child) Four() *Parent_Child_Four {
	return &Parent_Child_Four{
		parent: n,
	}
}

// Four corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildAny) Four() *Parent_Child_FourAny {
	return &Parent_Child_FourAny{
		parent: n,
	}
}

// One corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_Child) One() *Parent_Child_One {
	return &Parent_Child_One{
		parent: n,
	}
}

// One corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildAny) One() *Parent_Child_OneAny {
	return &Parent_Child_OneAny{
		parent: n,
	}
}

// Three corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_Child) Three() *Parent_Child_Three {
	return &Parent_Child_Three{
		parent: n,
	}
}

// Three corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildAny) Three() *Parent_Child_ThreeAny {
	return &Parent_Child_ThreeAny{
		parent: n,
	}
}

// Two corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_Child) Two() *Parent_Child_Two {
	return &Parent_Child_Two{
		parent: n,
	}
}

// Two corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *Parent_ChildAny) Two() *Parent_Child_TwoAny {
	return &Parent_Child_TwoAny{
		parent: n,
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *Parent_Child) State() ygnmi.SingletonQuery[*oc.Parent_Child] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.Parent_Child](
		"Parent_Child",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *Parent_ChildAny) State() ygnmi.WildcardQuery[*oc.Parent_Child] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Parent_Child](
		"Parent_Child",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Parent_Child) Config() ygnmi.ConfigQuery[*oc.Parent_Child] {
	return ygnmi.NewNonLeafConfigQuery[*oc.Parent_Child](
		"Parent_Child",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Parent_ChildAny) Config() ygnmi.WildcardQuery[*oc.Parent_Child] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.Parent_Child](
		"Parent_Child",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
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
func (n *Parent_Child_Four) State() ygnmi.SingletonQuery[oc.Binary] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_FourAny) State() ygnmi.WildcardQuery[oc.Binary] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_Four) Config() ygnmi.ConfigQuery[oc.Binary] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_FourAny) Config() ygnmi.WildcardQuery[oc.Binary] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_One) State() ygnmi.SingletonQuery[string] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_OneAny) State() ygnmi.WildcardQuery[string] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_One) Config() ygnmi.ConfigQuery[string] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_OneAny) Config() ygnmi.WildcardQuery[string] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_Three) State() ygnmi.SingletonQuery[oc.E_Child_Three] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_ThreeAny) State() ygnmi.WildcardQuery[oc.E_Child_Three] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_Three) Config() ygnmi.ConfigQuery[oc.E_Child_Three] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_ThreeAny) Config() ygnmi.WildcardQuery[oc.E_Child_Three] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_Two) State() ygnmi.SingletonQuery[string] {
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
			Root:       &oc.Device{},
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
func (n *Parent_Child_TwoAny) State() ygnmi.WildcardQuery[string] {
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
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// RemoteContainer represents the /openconfig-simple/remote-container YANG schema element.
type RemoteContainer struct {
	*ygot.NodePath
}

// RemoteContainerAny represents the wildcard version of the /openconfig-simple/remote-container YANG schema element.
type RemoteContainerAny struct {
	*ygot.NodePath
}

// RemoteContainer_ALeaf represents the /openconfig-simple/remote-container/state/a-leaf YANG schema element.
type RemoteContainer_ALeaf struct {
	parent ygot.PathStruct
}

// RemoteContainer_ALeafAny represents the wildcard version of the /openconfig-simple/remote-container/state/a-leaf YANG schema element.
type RemoteContainer_ALeafAny struct {
	parent ygot.PathStruct
}

// ALeaf corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *RemoteContainer) ALeaf() *RemoteContainer_ALeaf {
	return &RemoteContainer_ALeaf{
		parent: n,
	}
}

// ALeaf corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *RemoteContainerAny) ALeaf() *RemoteContainer_ALeafAny {
	return &RemoteContainer_ALeafAny{
		parent: n,
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *RemoteContainer) State() ygnmi.SingletonQuery[*oc.RemoteContainer] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.RemoteContainer](
		"RemoteContainer",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *RemoteContainerAny) State() ygnmi.WildcardQuery[*oc.RemoteContainer] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.RemoteContainer](
		"RemoteContainer",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *RemoteContainer) Config() ygnmi.ConfigQuery[*oc.RemoteContainer] {
	return ygnmi.NewNonLeafConfigQuery[*oc.RemoteContainer](
		"RemoteContainer",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *RemoteContainerAny) Config() ygnmi.WildcardQuery[*oc.RemoteContainer] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.RemoteContainer](
		"RemoteContainer",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
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
func (n *RemoteContainer_ALeaf) State() ygnmi.SingletonQuery[string] {
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
			Root:       &oc.Device{},
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
func (n *RemoteContainer_ALeafAny) State() ygnmi.WildcardQuery[string] {
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
			Root:       &oc.Device{},
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
func (n *RemoteContainer_ALeaf) Config() ygnmi.ConfigQuery[string] {
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
			Root:       &oc.Device{},
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
func (n *RemoteContainer_ALeafAny) Config() ygnmi.WildcardQuery[string] {
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
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

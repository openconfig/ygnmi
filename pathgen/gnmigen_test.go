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

package pathgen

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygot/ygen"
)

func TestGNMIGenerator(t *testing.T) {
	_, dirs, _ := getSchemaAndDirs()

	tests := []struct {
		desc           string
		pathStructName string
		dir            *ygen.Directory
		node           *NodeData
		want           string
		wantErr        string
	}{{
		desc: "field doesn't exist",
		dir:  dirs["/root-module/container"],
		node: &NodeData{
			YANGFieldName: "doesn't exist",
			IsLeaf:        true,
		},
		wantErr: "does not exist in Directory",
	}, {
		desc:           "scalar leaf without config",
		dir:            dirs["/root-module/container"],
		pathStructName: "Container_Leaf",
		node: &NodeData{
			GoTypeName:            "int32",
			LocalGoTypeName:       "int32",
			GoFieldName:           "Leaf",
			YANGFieldName:         "leaf",
			SubsumingGoStructName: "Container",
			IsLeaf:                true,
			IsScalarField:         true,
			HasDefault:            true,
			YANGPath:              "/container/leaf",
		},
		want: `
func (n *Container_Leaf) State() ygnmi.SingletonQuery[int32] {
	return ygnmi.NewLeafSingletonQuery[int32](
		"Container",
		true,
		true,
		ygot.NewNodePath(
			[]string{"leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int32, bool) { 
			ret := gs.(*oc.Container).Leaf
			return *ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Container) },
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

func (n *Container_LeafAny) State() ygnmi.WildcardQuery[int32] {
	return ygnmi.NewLeafWildcardQuery[int32](
		"Container",
		true,
		true,
		ygot.NewNodePath(
			[]string{"leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int32, bool) { 
			ret := gs.(*oc.Container).Leaf
			return *ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Container) },
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}
`,
	}, {
		desc:           "scalar leaf with config",
		dir:            dirs["/root-module/container-with-config"],
		pathStructName: "Container_Leaf",
		node: &NodeData{
			GoTypeName:            "int32",
			LocalGoTypeName:       "int32",
			GoFieldName:           "Leaf",
			YANGFieldName:         "leaf",
			SubsumingGoStructName: "Container",
			IsLeaf:                true,
			IsScalarField:         true,
			HasDefault:            true,
			YANGPath:              "/container/leaf",
		},
		want: `
func (n *Container_Leaf) State() ygnmi.SingletonQuery[int32] {
	return ygnmi.NewLeafSingletonQuery[int32](
		"Container",
		true,
		true,
		ygot.NewNodePath(
			[]string{"state", "leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int32, bool) { 
			ret := gs.(*oc.Container).Leaf
			return *ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Container) },
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

func (n *Container_LeafAny) State() ygnmi.WildcardQuery[int32] {
	return ygnmi.NewLeafWildcardQuery[int32](
		"Container",
		true,
		true,
		ygot.NewNodePath(
			[]string{"state", "leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int32, bool) { 
			ret := gs.(*oc.Container).Leaf
			return *ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Container) },
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

func (n *Container_Leaf) Config() ygnmi.ConfigQuery[int32] {
	return ygnmi.NewLeafConfigQuery[int32](
		"Container",
		false,
		true,
		ygot.NewNodePath(
			[]string{"config", "leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int32, bool) { 
			ret := gs.(*oc.Container).Leaf
			return *ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Container) },
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

func (n *Container_LeafAny) Config() ygnmi.WildcardQuery[int32] {
	return ygnmi.NewLeafWildcardQuery[int32](
		"Container",
		false,
		true,
		ygot.NewNodePath(
			[]string{"config", "leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int32, bool) { 
			ret := gs.(*oc.Container).Leaf
			return *ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Container) },
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}
`,
	}, {
		desc:           "non-scalar leaf",
		dir:            dirs["/root-module/container"],
		pathStructName: "Container_Leaf",
		node: &NodeData{
			GoTypeName:            "E_Child_Three",
			LocalGoTypeName:       "E_Child_Three",
			GoFieldName:           "Leaf",
			YANGFieldName:         "leaf",
			SubsumingGoStructName: "Container",
			IsLeaf:                true,
			IsScalarField:         false,
			HasDefault:            true,
			YANGPath:              "/container/leaf",
		},
		want: `
func (n *Container_Leaf) State() ygnmi.SingletonQuery[E_Child_Three] {
	return ygnmi.NewLeafSingletonQuery[E_Child_Three](
		"Container",
		true,
		false,
		ygot.NewNodePath(
			[]string{"leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (E_Child_Three, bool) { 
			ret := gs.(*oc.Container).Leaf
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Container) },
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

func (n *Container_LeafAny) State() ygnmi.WildcardQuery[E_Child_Three] {
	return ygnmi.NewLeafWildcardQuery[E_Child_Three](
		"Container",
		true,
		false,
		ygot.NewNodePath(
			[]string{"leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (E_Child_Three, bool) { 
			ret := gs.(*oc.Container).Leaf
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Container) },
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}
`,
	}, {
		desc:           "fake root",
		dir:            dirs["/root"],
		pathStructName: "Root",
		node: &NodeData{
			GoTypeName:            "*Root",
			LocalGoTypeName:       "Root",
			GoFieldName:           "",
			YANGFieldName:         "",
			SubsumingGoStructName: "Root",
			IsLeaf:                false,
			IsScalarField:         false,
			HasDefault:            false,
			YANGPath:              "/",
		},
		want: `
func (n *Root) State() ygnmi.SingletonQuery[*Root] {
	return ygnmi.NewNonLeafSingletonQuery[*Root](
		"Root",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

func (n *Root) Config() ygnmi.ConfigQuery[*Root] {
	return ygnmi.NewNonLeafConfigQuery[*Root](
		"Root",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Device{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}
`,
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := GNMIGenerator(tt.pathStructName, tt.dir, tt.node)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("GNMIGenerator(%q, %v, %v) returned unexpected error diff: %s", tt.pathStructName, tt.dir, tt.node, diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("GNMIGenerator(%q, %v, %v) returned unexpected error diff: %s", tt.pathStructName, tt.dir, tt.node, diff)
			}
		})
	}
}

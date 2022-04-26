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
		desc:           "scalar leaf",
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
	return &ygnmi.NewLeafSingletonQuery[int32](
		"Container",
		true,
		ygot.NewNodePath(
			[]string{"leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) int32 { 
			ret := gs.(*oc.Container).Leaf
			return *ret
		},
		func() ygot.ValidatedGoStruct { return new(oc.Container) },
		oc.GetSchema(),
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
	return &ygnmi.NewLeafSingletonQuery[E_Child_Three](
		"Container",
		true,
		ygot.NewNodePath(
			[]string{"leaf"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) E_Child_Three { 
			ret := gs.(*oc.Container).Leaf
			return ret
		},
		func() ygot.ValidatedGoStruct { return new(oc.Container) },
		oc.GetSchema(),
	)
}
`,
	}, {
		desc:           "non leaf",
		dir:            dirs["/root-module/container"],
		pathStructName: "Container",
		node: &NodeData{
			GoTypeName:            "Container",
			LocalGoTypeName:       "Container",
			GoFieldName:           "Leaf",
			YANGFieldName:         "leaf",
			SubsumingGoStructName: "Container",
			IsLeaf:                false,
			IsScalarField:         false,
			HasDefault:            true,
			YANGPath:              "/container",
		},
		want: `
func (n *Container) State() ygnmi.SingletonQuery[*Container] {
	return &ygnmi.NewNonLeafSingletonQuery[*Container](
		"Container",
		true,
		n,
		oc.GetSchema(),
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

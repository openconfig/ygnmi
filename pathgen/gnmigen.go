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
	"strings"

	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygen"
	"github.com/openconfig/ygot/ygot"
)

type gnmiStruct struct {
	PathStructName          string
	GoTypeName              string
	GoStructTypeName        string
	PathBaseTypeName        string
	SingletonTypeName       string
	GoFieldName             string
	SchemaStructPkgAccessor string
	RelPathList             string
	IsState                 bool
	MethodName              string
	IsScalar                bool
	GenerateWildcard        bool
	WildcardTypeName        string
	WildcardSuffix          string
	FakeRootName            string
}

const (
	fakeRootName = "Device"
)

// GNMIGenerator is a plugin generator for generating ygnmi query objects.
// Note: GNMIGenerator requires that PreferOperationalState be true when generating PathStructs.
// TODO(DanG100): pass schema from parent to child.
func GNMIGenerator(pathStructName string, dir *ygen.Directory, node *NodeData) (string, error) {
	tmplStruct := gnmiStruct{
		PathStructName:          pathStructName,
		GoTypeName:              node.GoTypeName,
		GoStructTypeName:        node.SubsumingGoStructName,
		PathBaseTypeName:        ygot.PathBaseTypeName,
		GoFieldName:             node.GoFieldName,
		SchemaStructPkgAccessor: "oc.",
		IsState:                 true,
		MethodName:              "State",
		SingletonTypeName:       "SingletonQuery",
		WildcardTypeName:        "WildcardQuery",
		IsScalar:                node.IsScalarField,
		GenerateWildcard:        node.YANGPath != "/", // Do not generate wildcard for the fake root.
		WildcardSuffix:          WildcardSuffix,
		FakeRootName:            fakeRootName,
	}

	tmpl := goGNMINonLeafTemplate
	if node.IsLeaf {
		relPath, err := ygen.FindSchemaPath(dir, node.YANGFieldName, false)
		if err != nil {
			return "", err
		}
		tmpl = goGNMILeafTemplate
		tmplStruct.RelPathList = `"` + strings.Join(relPath, `", "`) + `"`
	}

	var b strings.Builder
	if err := tmpl.Execute(&b, tmplStruct); err != nil {
		return "", err
	}

	if !generateConfigFunc(dir, node) {
		return b.String(), nil
	}

	tmplStruct.MethodName = "Config"
	tmplStruct.SingletonTypeName = "ConfigQuery"
	tmplStruct.IsState = false
	if node.IsLeaf {
		relPath, err := ygen.FindShadowSchemaPath(dir, node.YANGFieldName, false)
		if err != nil {
			return "", err
		}
		tmplStruct.RelPathList = `"` + strings.Join(relPath, `", "`) + `"`
	}
	if err := tmpl.Execute(&b, tmplStruct); err != nil {
		return "", err
	}
	return b.String(), nil
}

// generateConfig determines if a node should have a .Config() method.
// For leaves, it checks if the directory has a shadow-path field.
// For non-leaves, it checks if the directory or any of its descendants are config nodes.
func generateConfigFunc(dir *ygen.Directory, node *NodeData) bool {
	if node.IsLeaf {
		_, ok := dir.ShadowedFields[node.YANGFieldName]
		return ok
	}
	return util.IsConfig(dir.Entry)
}

var (
	goGNMILeafTemplate = mustTemplate("leaf-gnmi", `
func (n *{{ .PathStructName }}) {{ .MethodName }}() ygnmi.{{ .SingletonTypeName }}[{{ .GoTypeName }}] {
	return ygnmi.NewLeaf{{ .SingletonTypeName }}[{{ .GoTypeName }}](
		"{{ .GoStructTypeName }}",
		{{ .IsState }},
		{{ .IsScalar }},
		ygot.New{{ .PathBaseTypeName }}(
			[]string{ {{- .RelPathList -}} },
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) {{ .GoTypeName }} { 
			ret := gs.(*{{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}).{{ .GoFieldName }}
			{{- if .IsScalar }}
			return *ret
			{{- else }}
			return ret
			{{- end}}
		},
		func() ygot.ValidatedGoStruct { return new({{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}) },
		&ytypes.Schema{
			Root:       &{{ .SchemaStructPkgAccessor }}{{ .FakeRootName }}{},
			SchemaTree: {{ .SchemaStructPkgAccessor }}SchemaTree,
			Unmarshal:  {{ .SchemaStructPkgAccessor }}Unmarshal,
		},
	)
}

{{- if .GenerateWildcard }}

func (n *{{ .PathStructName }}{{ .WildcardSuffix }}) {{ .MethodName }}() ygnmi.{{ .WildcardTypeName }}[{{ .GoTypeName }}] {
	return ygnmi.NewLeaf{{ .WildcardTypeName }}[{{ .GoTypeName }}](
		"{{ .GoStructTypeName }}",
		{{ .IsState }},
		{{ .IsScalar }},
		ygot.New{{ .PathBaseTypeName }}(
			[]string{ {{- .RelPathList -}} },
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) {{ .GoTypeName }} { 
			ret := gs.(*{{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}).{{ .GoFieldName }}
			{{- if .IsScalar }}
			return *ret
			{{- else }}
			return ret
			{{- end}}
		},
		func() ygot.ValidatedGoStruct { return new({{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}) },
		&ytypes.Schema{
			Root:       &{{ .SchemaStructPkgAccessor }}{{ .FakeRootName }}{},
			SchemaTree: {{ .SchemaStructPkgAccessor }}SchemaTree,
			Unmarshal:  {{ .SchemaStructPkgAccessor }}Unmarshal,
		},
	)
}
{{- end }}
`)

	goGNMINonLeafTemplate = mustTemplate("non-leaf-gnmi", `
func (n *{{ .PathStructName }}) {{ .MethodName }}() ygnmi.{{ .SingletonTypeName }}[{{ .GoTypeName }}] {
	return ygnmi.NewNonLeaf{{ .SingletonTypeName }}[{{ .GoTypeName }}](
		"{{ .GoStructTypeName }}",
		{{ .IsState }},
		n,
		&ytypes.Schema{
			Root:       &{{ .SchemaStructPkgAccessor }}{{ .FakeRootName }}{},
			SchemaTree: {{ .SchemaStructPkgAccessor }}SchemaTree,
			Unmarshal:  {{ .SchemaStructPkgAccessor }}Unmarshal,
		},
	)
}

{{- if .GenerateWildcard }}

func (n *{{ .PathStructName }}{{ .WildcardSuffix }}) {{ .MethodName }}() ygnmi.{{ .WildcardTypeName }}[{{ .GoTypeName }}] {
	return ygnmi.NewNonLeaf{{ .WildcardTypeName }}[{{ .GoTypeName }}](
		"{{ .GoStructTypeName }}",
		{{ .IsState }},
		n,
		&ytypes.Schema{
			Root:       &{{ .SchemaStructPkgAccessor }}{{ .FakeRootName }}{},
			SchemaTree: {{ .SchemaStructPkgAccessor }}SchemaTree,
			Unmarshal:  {{ .SchemaStructPkgAccessor }}Unmarshal,
		},
	)
}
{{- end }}
`)
)

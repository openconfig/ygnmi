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

	"github.com/openconfig/ygot/ygen"
	"github.com/openconfig/ygot/ygot"
)

type gnmiStruct struct {
	PathStructName          string
	GoTypeName              string
	GoStructTypeName        string
	PathBaseTypeName        string
	GoFieldName             string
	SchemaStructPkgAccessor string
	RelPathList             string
	IsState                 bool
	MethodName              string
	IsScalar                bool
}

// GNMIGenerator is a plugin generator for generating ygnmi query objects
func GNMIGenerator(pathStructName string, dir *ygen.Directory, node *NodeData) (string, error) {
	tmplStruct := gnmiStruct{
		PathStructName:          pathStructName,
		GoTypeName:              node.GoTypeName,
		GoStructTypeName:        node.SubsumingGoStructName,
		PathBaseTypeName:        ygot.PathBaseTypeName,
		GoFieldName:             node.GoFieldName,
		SchemaStructPkgAccessor: "oc.",
		IsState:                 true,
		MethodName:              "State", // TODO(DanG100): add config generation
		IsScalar:                node.IsScalarField,
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
	return b.String(), nil
}

var (
	goGNMILeafTemplate = mustTemplate("leaf-gnmi", `
func (n *{{ .PathStructName }}) {{ .MethodName }}() ygnmi.SingletonQuery[{{ .GoTypeName }}] {
	return &ygnmi.NewLeafSingletonQuery[{{ .GoTypeName }}](
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
		{{ .SchemaStructPkgAccessor }}GetSchema(),
	)
}
`)

	goGNMINonLeafTemplate = mustTemplate("non-leaf-gnmi", `
func (n *{{ .PathStructName }}) {{ .MethodName }}() ygnmi.SingletonQuery[*{{ .GoStructTypeName }}] {
	return &ygnmi.NewNonLeafSingletonQuery[*{{ .GoStructTypeName }}](
		"{{ .GoStructTypeName }}",
		{{ .IsState }},
		n,
		{{ .SchemaStructPkgAccessor }}GetSchema(),
	)
}
`)
)

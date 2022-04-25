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
}

// gnmiGenerator is a extra Generator for PathStructs that generates gNMI helpers.
func gnmiGenerator(psName string, node *NodeData) (string, error) {
	// TODO(DanG100): non-leaves
	if !node.IsLeaf {
		return "", nil
	}

	relPath, err := ygen.FindSchemaPath(node.Directory, node.FieldName, false)
	if err != nil {
		return "", err
	}

	tmplStruct := gnmiStruct{
		PathStructName:          psName,
		GoTypeName:              node.GoTypeName,
		GoStructTypeName:        node.SubsumingGoStructName,
		PathBaseTypeName:        ygot.PathBaseTypeName,
		GoFieldName:             node.GoFieldName,
		SchemaStructPkgAccessor: "oc.",
		RelPathList:             `"` + strings.Join(relPath, `", "`) + `"`,
		IsState:                 true,
		MethodName:              "State",
	}
	var b strings.Builder
	if err = goGNMILeafTemplate.Execute(&b, tmplStruct); err != nil {
		return "", err
	}
	if _, ok := node.Directory.ShadowedFields[node.FieldName]; !ok {
		return b.String(), nil
	}
	tmplStruct.IsState = false
	tmplStruct.MethodName = "Config"
	if err = goGNMILeafTemplate.Execute(&b, tmplStruct); err != nil {
		return "", err
	}

	return b.String(), nil
}

var (
	goGNMILeafTemplate = mustTemplate("leaf-gnmi", `
	func (n * {{ .PathStructName }}) {{ .MethodName }}() ygnmi.SingletonQuery[ {{ .SchemaStructPkgAccessor }}{{ .GoTypeName }} ] {
		return &ygnmi.NewLeafSingletonQuery[ {{ .SchemaStructPkgAccessor }}{{ .GoTypeName }} ](
			"{{ .GoStructTypeName }}",
			{{ .IsState }},
			ygot.New{{ .PathBaseTypeName }}(
				[]string{ {{- .RelStatePathList -}} },
				nil,
				n.parent
			),
			func(gs ygot.ValidatedGoStruct) {{ .GoTypeName }} { return gs.(*{{ .GoStructTypeName }}).{{ .GoFieldName }} ) },
			func() ygot.ValidatedGoStruct { return new({{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}) },
			{{ .SchemaStructPkgAccessor }}Schema,
		)
	}
	`)
)

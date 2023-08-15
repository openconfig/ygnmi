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
	"fmt"
	"strings"
	"text/template"

	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygen"
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
	AbsPath                 string
	RelPath                 string
	DefiningModuleName      string
	InstantiatingModuleName string
	SpecialConvertFunc      string
	CompressInfo            *CompressionInfo
	IsListContainer         bool
	IsCompressedSchema      bool
}

const (
	// TODO(DanG100): pass options into custom generators and remove this.
	fakeRootName = "Root"

	wildcardQueryTypeName  = "WildcardQuery"
	singletonQueryTypeName = "SingletonQuery"
	configQueryTypeName    = "ConfigQuery"
)

var packagesSeen = map[string]bool{}

// GNMIGenerator is a plugin generator for generating ygnmi query objects for
// compressed structs.
//
// - pathStructName is the name of the PathStruct of the node.
// - node contains information of the node.
// - dir is the containing directory of the node. For leaves this is the
// parent.
//
// Note: GNMIGenerator requires that PreferOperationalState be true when generating PathStructs.
// TODO(DanG100): pass schema from parent to child.
func GNMIGenerator(pathStructName string, dir *ygen.ParsedDirectory, node *NodeData, _ bool) (string, error) {
	tmplStruct := defaultTmplStruct(pathStructName, node)
	var b strings.Builder
	if err := generateOneOff(&b, node, tmplStruct, true); err != nil {
		return "", err
	}
	if err := modifyQueryType(node, &tmplStruct); err != nil {
		return "", err
	}

	tmpl := goGNMINonLeafTemplate
	if node.IsLeaf {
		tmpl = goGNMILeafTemplate
	}

	generate := func(tmplStruct gnmiStruct, shadow bool) error {
		if node.IsLeaf {
			if err := populateTmplForLeaf(dir, node.YANGFieldName, shadow, &tmplStruct); err != nil {
				return err
			}
		}
		return tmpl.Execute(&b, &tmplStruct)
	}

	if err := generate(tmplStruct, false); err != nil {
		return "", err
	}

	if !generateConfigFunc(dir, node) {
		return b.String(), nil
	}

	tmplStruct.MethodName = "Config"
	tmplStruct.SingletonTypeName = configQueryTypeName
	tmplStruct.IsState = false
	if err := generate(tmplStruct, true); err != nil {
		return "", err
	}

	return b.String(), nil
}

// GNMIGeneratorUncompressed is a plugin generator for generating ygnmi query objects for
// uncompressed structs.
func GNMIGeneratorUncompressed(pathStructName string, dir *ygen.ParsedDirectory, node *NodeData, _ bool) (string, error) {
	var b strings.Builder
	if err := generateOneOff(&b, node, defaultTmplStruct(pathStructName, node), false); err != nil {
		return "", err
	}
	return b.String(), nil
}

func defaultTmplStruct(pathStructName string, node *NodeData) gnmiStruct {
	return gnmiStruct{
		PathStructName:          pathStructName,
		GoTypeName:              node.GoTypeName,
		GoStructTypeName:        node.SubsumingGoStructName,
		PathBaseTypeName:        ygnmi.PathBaseTypeName,
		GoFieldName:             node.GoFieldName,
		SchemaStructPkgAccessor: "oc.",
		IsState:                 true,
		IsCompressedSchema:      true,
		MethodName:              "State",
		SingletonTypeName:       singletonQueryTypeName,
		WildcardTypeName:        wildcardQueryTypeName,
		IsScalar:                node.IsScalarField,
		GenerateWildcard:        node.YANGPath != "/", // Do not generate wildcard for the fake root.
		WildcardSuffix:          WildcardSuffix,
		FakeRootName:            fakeRootName,
		CompressInfo:            node.CompressInfo,
		IsListContainer:         node.IsListContainer,
	}
}

// generateOneOff generates one-off free-form generated code.
func generateOneOff(b *strings.Builder, node *NodeData, tmplStruct gnmiStruct, compressPaths bool) error {
	if strings.TrimLeft(node.LocalGoTypeName, "*") == fakeRootName {
		tmplStruct.MethodName = "Query"
		if compressPaths {
			tmplStruct.MethodName = "State"
		}
		tmplStruct.IsState = true
		if err := batchStructTemplate.Execute(b, &tmplStruct); err != nil {
			return err
		}
		if err := batchTemplate.Execute(b, &tmplStruct); err != nil {
			return err
		}
		if compressPaths {
			tmplStruct.MethodName = "Config"
			tmplStruct.IsState = false
			if err := batchTemplate.Execute(b, &tmplStruct); err != nil {
				return err
			}
		}
	}
	if !packagesSeen[node.GoPathPackageName] {
		packagesSeen[node.GoPathPackageName] = true
		if err := oncePerPackageTmpl.Execute(b, struct{}{}); err != nil {
			return err
		}
	}

	return nil
}

// modifyQueryType modifies the default type of the query with custom types.
func modifyQueryType(node *NodeData, s *gnmiStruct) error {
	// use float32 instead of the binary type for ieeefloat32 since it is
	// more user friendly.
	if node.YANGTypeName == "ieeefloat32" {
		switch node.LocalGoTypeName {
		case "Binary":
			s.GoTypeName = "float32"
			s.SpecialConvertFunc = "ygot.BinaryToFloat32"
		case "[]Binary":
			s.GoTypeName = "[]float32"
			s.SpecialConvertFunc = "binarySliceToFloatSlice"
		default:
			return fmt.Errorf("ieeefloat32 is expected to be a binary, got %q", node.LocalGoTypeName)
		}
	}
	return nil
}

// GNMIFieldGenerator generates an embedded query field for a PathStruct.
//
// This is meant to be used for uncompressed generation, where there is no need
// to distinguish between config and state queries.
func GNMIFieldGenerator(pathStructName string, _ *ygen.ParsedDirectory, node *NodeData, wildcard bool) (string, error) {
	tmplStruct := gnmiStruct{
		GoTypeName: node.GoTypeName,
	}
	if err := modifyQueryType(node, &tmplStruct); err != nil {
		return "", err
	}

	var queryTypeName string
	switch {
	case wildcard:
		queryTypeName = wildcardQueryTypeName
	case node.ConfigFalse:
		queryTypeName = singletonQueryTypeName
	default:
		queryTypeName = configQueryTypeName
	}

	return fmt.Sprintf("\n\tygnmi.%s[%s]", queryTypeName, tmplStruct.GoTypeName), nil
}

// GNMIInitGenerator generates the initialization for an embedded query field
// for a PathStruct.
//
// This is meant to be used for uncompressed generation, where there is no need
// to distinguish between config and state queries.
func GNMIInitGenerator(pathStructName string, _ *ygen.ParsedDirectory, node *NodeData, wildcard bool) (string, error) {
	tmplStruct := defaultTmplStruct(pathStructName, node)
	tmplStruct.IsCompressedSchema = false
	if err := modifyQueryType(node, &tmplStruct); err != nil {
		return "", err
	}

	var queryTypeName string
	switch {
	case wildcard:
		queryTypeName = wildcardQueryTypeName
	case node.ConfigFalse:
		tmplStruct.SingletonTypeName = singletonQueryTypeName
		queryTypeName = singletonQueryTypeName
	default:
		tmplStruct.SingletonTypeName = configQueryTypeName
		queryTypeName = configQueryTypeName
	}

	var tmpl *template.Template
	switch {
	case wildcard && node.IsLeaf:
		tmpl = goGNMILeafTemplate.Lookup("leaf-gnmi-wildcard-query")
	case wildcard:
		tmpl = goGNMINonLeafTemplate.Lookup("nonleaf-gnmi-wildcard-query")
	case node.IsLeaf:
		tmpl = goGNMILeafTemplate.Lookup("leaf-gnmi-singleton-query")
	default:
		tmpl = goGNMINonLeafTemplate.Lookup("nonleaf-gnmi-singleton-query")
	}

	var b strings.Builder
	if err := tmpl.Execute(&b, &tmplStruct); err != nil {
		return "", err
	}
	return fmt.Sprintf("\n\tps.%s = %s", queryTypeName, b.String()), nil
}

// populateTmplForLeaf adds leaf specific fields to the gnmiStruct template.
func populateTmplForLeaf(dir *ygen.ParsedDirectory, fieldName string, shadow bool, tmplStruct *gnmiStruct) error {
	field, ok := dir.Fields[fieldName]
	if !ok {
		return fmt.Errorf("field %q does not exist in directory %s", fieldName, dir.Path)
	}
	// The longest path is the non-key path. This is the one we want to use
	// since the key is "compressed out".
	relPath := longestPath(field.MappedPaths)
	if shadow {
		relPath = longestPath(field.ShadowMappedPaths)
	}

	tmplStruct.RelPathList = `"` + strings.Join(relPath, `", "`) + `"`
	tmplStruct.AbsPath = field.YANGDetails.SchemaPath
	if shadow {
		tmplStruct.AbsPath = field.YANGDetails.ShadowSchemaPath
	}
	tmplStruct.RelPath = strings.Join(relPath, `/`)
	tmplStruct.InstantiatingModuleName = field.YANGDetails.BelongingModule
	tmplStruct.DefiningModuleName = field.YANGDetails.DefiningModule
	return nil
}

// generateConfigFunc determines if a node should have a .Config() method.
// For leaves, it checks if the directory has a shadow-path field.
// For non-leaves, it checks if the directory or any of its descendants are config nodes.
func generateConfigFunc(dir *ygen.ParsedDirectory, node *NodeData) bool {
	if node.IsLeaf {
		// Since we're generating with PreferOperationalState, we need
		// to check shadow-fields to know whether this has a config
		// sibling.
		field, ok := dir.Fields[node.YANGFieldName]
		return ok && len(field.ShadowMappedPaths) > 0
	}
	return !node.ConfigFalse
}

// mustParseDownstreamTemplate parses a new template.Template with the given
// tmpl as the associated template.
func mustParseDownstreamTemplate(tmpl *template.Template, src string) *template.Template {
	tmpl, err := template.Must(tmpl.Clone()).Parse(src)
	if err != nil {
		panic(err)
	}
	return tmpl
}

var (
	goGNMILeafQueryTemplate = mustTemplate("leaf-gnmi-query", `{{ define "leaf-gnmi-singleton-query" -}}
ygnmi.New{{ .SingletonTypeName }}[{{ .GoTypeName }}](
		"{{ .GoStructTypeName }}",
		{{ .IsState }},
		true,
		{{ .IsScalar }},
		{{ .IsCompressedSchema }},
		false,
		{{- if .IsCompressedSchema }}
		ygnmi.New{{ .PathBaseTypeName }}(
			[]string{ {{- .RelPathList -}} },
			nil,
			n.parent,
		),
		{{- else }}
		ps,
		{{- end }}
		func(gs ygot.ValidatedGoStruct) ({{ .GoTypeName }}, bool) { 
			ret := gs.(*{{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}).{{ .GoFieldName }}
			{{- if .IsScalar }}
			if ret == nil {
				var zero {{ .GoTypeName }}
				return zero, false
			}
			return *ret, true
			{{- else }}
			{{- if .SpecialConvertFunc }}
			return {{ .SpecialConvertFunc }}(ret), !reflect.ValueOf(ret).IsZero()
			{{- else}}
			return ret, !reflect.ValueOf(ret).IsZero()
			{{- end }}
			{{- end}}
		},
		func() ygot.ValidatedGoStruct { return new({{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &{{ .SchemaStructPkgAccessor }}{{ .FakeRootName }}{},
				SchemaTree: {{ .SchemaStructPkgAccessor }}SchemaTree,
				Unmarshal:  {{ .SchemaStructPkgAccessor }}Unmarshal,
			}
		},
		nil,
		nil,
	)
{{ end }}
{{ define "leaf-gnmi-wildcard-query" -}}
ygnmi.New{{ .WildcardTypeName }}[{{ .GoTypeName }}](
		"{{ .GoStructTypeName }}",
		{{ .IsState }},
		true,
		{{ .IsScalar }},
		{{ .IsCompressedSchema }},
		false,
		{{- if .IsCompressedSchema }}
		ygnmi.New{{ .PathBaseTypeName }}(
			[]string{ {{- .RelPathList -}} },
			nil,
			n.parent,
		),
		{{- else }}
		ps,
		{{- end }}
		func(gs ygot.ValidatedGoStruct) ({{ .GoTypeName }}, bool) { 
			ret := gs.(*{{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}).{{ .GoFieldName }}
			{{- if .IsScalar }}
			if ret == nil {
				var zero {{ .GoTypeName }}
				return zero, false
			}
			return *ret, true
			{{- else }}
			{{- if .SpecialConvertFunc }}
			return {{ .SpecialConvertFunc }}(ret), !reflect.ValueOf(ret).IsZero()
			{{- else}}
			return ret, !reflect.ValueOf(ret).IsZero()
			{{- end }}
			{{- end}}
		},
		func() ygot.ValidatedGoStruct { return new({{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &{{ .SchemaStructPkgAccessor }}{{ .FakeRootName }}{},
				SchemaTree: {{ .SchemaStructPkgAccessor }}SchemaTree,
				Unmarshal:  {{ .SchemaStructPkgAccessor }}Unmarshal,
			}
		},
		nil,
	)
{{ end }}`)

	goGNMILeafTemplate = mustParseDownstreamTemplate(goGNMILeafQueryTemplate, `
// {{ .MethodName }} returns a Query that can be used in gNMI operations.
// 	Defining module:      "{{ .DefiningModuleName }}"
// 	Instantiating module: "{{ .InstantiatingModuleName }}"
// 	Path from parent:     "{{ .RelPath }}"
// 	Path from root:       "{{ .AbsPath }}"
func (n *{{ .PathStructName }}) {{ .MethodName }}() ygnmi.{{ .SingletonTypeName }}[{{ .GoTypeName }}] {
	return {{ template "leaf-gnmi-singleton-query" . -}}
}

{{- if .GenerateWildcard }}

// {{ .MethodName }} returns a Query that can be used in gNMI operations.
// 	Defining module:      "{{ .DefiningModuleName }}"
// 	Instantiating module: "{{ .InstantiatingModuleName }}"
// 	Path from parent:     "{{ .RelPath }}"
// 	Path from root:       "{{ .AbsPath }}"
func (n *{{ .PathStructName }}{{ .WildcardSuffix }}) {{ .MethodName }}() ygnmi.{{ .WildcardTypeName }}[{{ .GoTypeName }}] {
	return {{ template "leaf-gnmi-wildcard-query" . -}}
}
{{- end }}
`)

	goGNMINonLeafQueryTemplate = mustTemplate("nonleaf-gnmi-query", `{{ define "nonleaf-gnmi-singleton-query" -}}
ygnmi.New{{ .SingletonTypeName }}[{{ .GoTypeName }}](
		"{{ .GoStructTypeName }}",
		{{ .IsState }},
		false,
		false,
		{{ .IsCompressedSchema }},
		{{ .IsListContainer }},
		{{- if .IsCompressedSchema }}
		n,
		{{- else }}
		ps,
		{{- end }}
		{{- if .IsListContainer }}
		func(gs ygot.ValidatedGoStruct) ({{ .GoTypeName }}, bool) { 
			ret := gs.(*{{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}).{{ .GoFieldName }}
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new({{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}) },
		{{- else }}
		nil,
		nil,
		{{- end }}
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &{{ .SchemaStructPkgAccessor }}{{ .FakeRootName }}{},
				SchemaTree: {{ .SchemaStructPkgAccessor }}SchemaTree,
				Unmarshal:  {{ .SchemaStructPkgAccessor }}Unmarshal,
			}
		},
		nil,
		{{- if .CompressInfo }}
		&ygnmi.CompressionInfo{
			PreRelPath: []string{ {{- .CompressInfo.PreRelPathList -}} },
			PostRelPath: []string{ {{- .CompressInfo.PostRelPathList -}} },
		},
		{{- else }}
		nil,
		{{- end }}
	)
{{ end }}
{{ define "nonleaf-gnmi-wildcard-query" -}}
ygnmi.New{{ .WildcardTypeName }}[{{ .GoTypeName }}](
		"{{ .GoStructTypeName }}",
		{{ .IsState }},
		false,
		false,
		{{ .IsCompressedSchema }},
		{{ .IsListContainer }},
		{{- if .IsCompressedSchema }}
		n,
		{{- else }}
		ps,
		{{- end }}
		{{- if .IsListContainer }}
		func(gs ygot.ValidatedGoStruct) ({{ .GoTypeName }}, bool) { 
			ret := gs.(*{{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}).{{ .GoFieldName }}
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new({{ .SchemaStructPkgAccessor }}{{ .GoStructTypeName }}) },
		{{- else }}
		nil,
		nil,
		{{- end }}
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &{{ .SchemaStructPkgAccessor }}{{ .FakeRootName }}{},
				SchemaTree: {{ .SchemaStructPkgAccessor }}SchemaTree,
				Unmarshal:  {{ .SchemaStructPkgAccessor }}Unmarshal,
			}
		},
		{{- if .CompressInfo }}
		&ygnmi.CompressionInfo{
			PreRelPath: []string{ {{- .CompressInfo.PreRelPathList -}} },
			PostRelPath: []string{ {{- .CompressInfo.PostRelPathList -}} },
		},
		{{- else }}
		nil,
		{{- end }}
	)
{{ end }}`)

	goGNMINonLeafTemplate = mustParseDownstreamTemplate(goGNMINonLeafQueryTemplate, `
// {{ .MethodName }} returns a Query that can be used in gNMI operations.
func (n *{{ .PathStructName }}) {{ .MethodName }}() ygnmi.{{ .SingletonTypeName }}[{{ .GoTypeName }}] {
	return {{ template "nonleaf-gnmi-singleton-query" . -}}
}

{{- if .GenerateWildcard }}

// {{ .MethodName }} returns a Query that can be used in gNMI operations.
func (n *{{ .PathStructName }}{{ .WildcardSuffix }}) {{ .MethodName }}() ygnmi.{{ .WildcardTypeName }}[{{ .GoTypeName }}] {
	return {{ template "nonleaf-gnmi-wildcard-query" . -}}
}
{{- end }}
`)

	batchStructTemplate = mustTemplate("batch-struct", `
// Batch contains a collection of paths.
// Use batch to call Lookup, Watch, etc. on multiple paths at once.
type Batch struct {
    paths []ygnmi.PathStruct
}

// AddPaths adds the paths to the batch.
func (b *Batch) AddPaths(paths ...ygnmi.PathStruct) *Batch {
    b.paths = append(b.paths, paths...)
    return b
}
`)

	batchTemplate = mustTemplate("batch", `
// {{ .MethodName }} returns a Query that can be used in gNMI operations.
// The returned query is immutable, adding paths does not modify existing queries.
func (b *Batch) {{ .MethodName }}() ygnmi.{{ .SingletonTypeName }}[{{ .GoTypeName }}] {
	queryPaths := make([]ygnmi.PathStruct, len(b.paths))
	copy(queryPaths, b.paths)
	return ygnmi.New{{ .SingletonTypeName }}[{{ .GoTypeName }}](
		"{{ .GoStructTypeName }}",
		{{ .IsState }},
		false,
		false,
		{{ .IsCompressedSchema }},
		false,
		ygnmi.NewDeviceRootBase(),
		nil,
		nil,
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &{{ .SchemaStructPkgAccessor }}{{ .FakeRootName }}{},
				SchemaTree: {{ .SchemaStructPkgAccessor }}SchemaTree,
				Unmarshal:  {{ .SchemaStructPkgAccessor }}Unmarshal,
			}
		},
		queryPaths,
		nil,
	)
}
`)

	oncePerPackageTmpl = mustTemplate("once-per-package", `
func binarySliceToFloatSlice(in []oc.Binary) []float32 {
	converted := make([]float32, 0, len(in))
	for _, binary := range in {
		converted = append(converted, ygot.BinaryToFloat32(binary))
	}
	return converted
}
`)
)

// Copyright 2019 Google Inc.
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

// Package pathgen contains a library to generate gNMI paths from a YANG model.
// The ygen library is used to parse YANG and obtain intermediate and some final
// information. The output always assumes the OpenConfig-specific conventions
// for a compressed schema.

package pathgen

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/openconfig/gnmi/errlist"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/genutil"
	"github.com/openconfig/ygot/gogen"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygen"
	"github.com/openconfig/ygot/ygot"
)

// Static default configuration values that differ from the zero value for their types.
const (
	// defaultPathPackageName specifies the default name that should be
	// used for the generated Go package.
	defaultPathPackageName = "ocpathstructs"
	// defaultFakeRootName is the default name for the root structure.
	defaultFakeRootName = "root"
	// defaultPathStructSuffix is the default suffix for generated
	// PathStructs to distinguish them from the generated GoStructs
	defaultPathStructSuffix = "Path"
	// defaultPackageSuffix is the default suffix for generated packages.
	defaultPackageSuffix = "path"
	// schemaStructPkgAlias is the package alias of the schema struct
	// package when the path struct package is to be generated in a
	// separate package.
	schemaStructPkgAlias = "oc"
	// WildcardSuffix is the suffix given to the wildcard versions of each
	// node as well as a list's wildcard child constructor methods that
	// distinguishes each from its non-wildcard counterpart.
	WildcardSuffix = "Any"
	// BuilderCtorSuffix is the suffix applied to the list builder
	// constructor method's name in order to indicate itself to the user.
	BuilderCtorSuffix = "Any"
	// WholeKeyedListSuffix is the suffix applied to a keyed list's
	// constructor method name in order to indicate itself to the user.
	WholeKeyedListSuffix = "Map"
	// BuilderKeyPrefix is the prefix applied to the key-modifying builder
	// method for a list PathStruct that uses the builder API.
	// NOTE: This cannot be "", as the builder method name would conflict
	// with the child constructor method for the keys.
	BuilderKeyPrefix = "With"
	// defaultYgnmiPath is the import path for the ygnmi library.
	defaultYgnmiPath = "github.com/openconfig/ygnmi/ygmni"
	// defaultYtypesPath is the import path for the ytypes library.
	defaultYtypesPath = "github.com/openconfig/ygot/ytypes"

	// yangTypeNameFlagKey is a custom flag for storing the YANG type's
	// name for a YANG node.
	yangTypeNameFlagKey = "YANG:typename"
)

// NewDefaultConfig creates a GenConfig with default configuration.
// schemaStructPkgPath is a required configuration parameter. It should be set
// to "" when the generated PathStruct package is to be the same package as the
// GoStructs package.
func NewDefaultConfig(schemaStructPkgPath string) *GenConfig {
	return &GenConfig{
		PackageName:   defaultPathPackageName,
		PackageSuffix: defaultPackageSuffix,
		GoImports: GoImports{
			SchemaStructPkgPath: schemaStructPkgPath,
			YgotImportPath:      genutil.GoDefaultYgotImportPath,
			YgnmiImportPath:     defaultYgnmiPath,
			YtypesImportPath:    defaultYtypesPath,
		},
		FakeRootName:     defaultFakeRootName,
		PathStructSuffix: defaultPathStructSuffix,
		GeneratingBinary: genutil.CallerName(),
	}
}

// Generator is func to returns extra generated code for a given node.
type Generator func(string, string, genutil.CompressBehaviour, *ygen.ParsedDirectory, *NodeData, bool) (string, error)

// ExtraGenerators is the set of all generators for a generation invocation.
type ExtraGenerators struct {
	// StructFields are extra PathStruct fields to be generated.
	StructFields []Generator
	// StructInits are extra PathStruct initializers to be generated when
	// the PathStruct is created.
	StructInits []Generator
	// Extras are other free-form extra generated code to accompany a
	// particular PathStruct.
	Extras []Generator
}

// GenConfig stores code generation configuration.
type GenConfig struct {
	// CompressBehaviour determines whether the direct children of some
	// entries are omitted in the generated code for brevity. It also uses
	// a shortened naming convention for type names, including
	// enumeration/identityref names.
	//
	// Compress options are recommended when using OpenConfig-styled YANG,
	// and Uncompress* options otherwise.
	CompressBehaviour genutil.CompressBehaviour
	// PackageName is the name that should be used for the generating package.
	PackageName string
	// GoImports contains package import options.
	GoImports GoImports
	// ParseOptions specifies the options for how the YANG schema is
	// produced.
	ParseOptions ygen.ParseOpts
	// FakeRootName specifies the name of the struct that should be generated
	// representing the root.
	FakeRootName string
	// PathStructSuffix is the suffix to be appended to generated
	// PathStructs to distinguish them from the generated GoStructs, which
	// assume a similar name.
	PathStructSuffix string
	// SkipEnumDeduplication specifies whether leaves of type 'enumeration' that
	// are used in multiple places in the schema should share a common type within
	// the generated code that is output by ygen. By default (false), a common type
	// is used.
	// This is the same flag used by ygen: they must match for pathgen's
	// generated code to be compatible with it.
	SkipEnumDeduplication bool
	// ShortenEnumLeafNames removes the module name from the name of
	// enumeration leaves.
	// This is the same flag used by ygen: they must match for pathgen's
	// generated code to be compatible with it.
	ShortenEnumLeafNames bool
	// EnumOrgPrefixesToTrim trims the organization name from the module
	// part of the name of enumeration leaves if there is a match.
	EnumOrgPrefixesToTrim []string
	// UseDefiningModuleForTypedefEnumNames uses the defining module name
	// to prefix typedef enumerated types instead of the module where the
	// typedef enumerated value is used.
	// This is the same flag used by ygen: they must match for pathgen's
	// generated code to be compatible with it.
	UseDefiningModuleForTypedefEnumNames bool
	// AppendEnumSuffixForSimpleUnionEnums appends an "Enum" suffix to the
	// enumeration name for simple (i.e. non-typedef) leaves which are
	// unions with an enumeration inside. This makes all inlined
	// enumerations within unions, whether typedef or not, have this
	// suffix, achieving consistency.  Since this flag is planned to be a
	// v1 compatibility flag along with
	// UseDefiningModuleForTypedefEnumNames, and will be removed in v1, it
	// only applies when useDefiningModuleForTypedefEnumNames is also set
	// to true.
	AppendEnumSuffixForSimpleUnionEnums bool
	// GeneratingBinary is the name of the binary calling the generator library, it is
	// included in the header of output files for debugging purposes. If a
	// string is not specified, the location of the library is utilised.
	GeneratingBinary string
	// GenerateWildcardPaths means to generate wildcard nodes and paths.
	GenerateWildcardPaths bool
	// SplitByModule controls whether to generate a go package for each yang module.
	SplitByModule bool
	// TrimPackageModulePrefix is the prefix to trim from generated go package names.
	TrimPackageModulePrefix string
	// BasePackagePath is used to create to full import path of the generated go packages.
	BasePackagePath string
	// PackageString is the string to apppend to the generated Go package names.
	PackageSuffix string
	// UnifyPathStructs controls whether to generate both config and states in the same package.
	UnifyPathStructs bool
	// ExtraGenerators are custom funcs that are used to extend the path struct generation.
	ExtraGenerators ExtraGenerators
	// IgnoreAtomicLists disables the following default behaviours:
	//  - All compressed lists will have a new accessor <ListName>Map() that
	//  retrieves the whole list.
	//  - Any child underneath lists whose compressed-out parent container
	//  is marked "telemetry-atomic" are no longer reachable.
	//  - lists marked `ordered-by user` will be represented using built-in Go maps
	//  instead of an ordered map Go structure.
	IgnoreAtomicLists bool
	// IgnoreAtomic avoids generating any descendants where a
	// non-compressed-out list or container is marked "telemetry-atomic".
	IgnoreAtomic bool
	// SplitPackagePaths specifies the prefix paths at or below which symbols will
	// reside in the provided package. The mapped string is the custom name for
	// the package (optional).
	SplitPackagePaths map[string]string
	//UseModuleNameAsPathOrigin uses the YANG module name to the origin for generated gNMI paths.
	UseModuleNameAsPathOrigin bool
	//PathOriginName specifies the origin name for generated gNMI paths.
	PathOriginName string
}

// GoImports contains package import options.
type GoImports struct {
	// SchemaStructPkgPath specifies the path to the ygen-generated structs, which
	// is used to get the enum and union type names used as the list key
	// for calling a list path accessor.
	SchemaStructPkgPath string
	// YgotImportPath specifies the path to the ygot library that should be used
	// in the generated code.
	YgotImportPath string
	// YtypesImportPath specifies the path to the ytypes library that should be used
	// in the generated code.
	YtypesImportPath string
	// YgnmiImportPath is the import path to the ygnmi library that should be used
	// in the generated code.
	YgnmiImportPath string
}

type goLangMapper struct {
	*gogen.GoLangMapper
}

// PopulateFieldFlags populates extra field information for pathgen.
func (goLangMapper) PopulateFieldFlags(nd ygen.NodeDetails, field *yang.Entry) map[string]string {
	flags := map[string]string{}
	if nd.Type == ygen.LeafNode || nd.Type == ygen.LeafListNode {
		// Only leaf or leaf-list nodes can have type statements.
		flags[yangTypeNameFlagKey] = field.Type.Name
	}
	return flags
}

// GeneratePathCode takes a slice of strings containing the path to a set of YANG
// files which contain YANG modules, and a second slice of strings which
// specifies the set of paths that are to be searched for associated models (e.g.,
// modules that are included by the specified set of modules, or submodules of those
// modules). It extracts the set of modules that are to be generated, and returns
// a map of package names to GeneratedPathCode structs. Each struct contains
// all the generated code of that package needed support the path-creation API.
// The important components of the generated code are listed below:
//  1. Struct definitions for each container, list, or leaf schema node,
//     as well as the fakeroot.
//  2. Next-level methods for the fakeroot and each non-leaf schema node,
//     which instantiate and return the next-level structs corresponding to
//     its child schema nodes.
//
// With these components, the generated API is able to support absolute path
// creation of any node of the input schema.
// Also returned is the NodeDataMap of the schema, i.e. information about each
// node in the generated code, which may help callers add customized
// augmentations to the basic generated path code.
// If errors are encountered during code generation, they are returned.
func (cg *GenConfig) GeneratePathCode(yangFiles, includePaths []string) (map[string]*GeneratedPathCode, NodeDataMap, util.Errors) {
	opts := ygen.IROptions{
		ParseOptions: cg.ParseOptions,
		TransformationOptions: ygen.TransformationOpts{
			CompressBehaviour:                    cg.CompressBehaviour,
			GenerateFakeRoot:                     true,
			FakeRootName:                         cg.FakeRootName,
			ShortenEnumLeafNames:                 cg.ShortenEnumLeafNames,
			EnumOrgPrefixesToTrim:                cg.EnumOrgPrefixesToTrim,
			UseDefiningModuleForTypedefEnumNames: cg.UseDefiningModuleForTypedefEnumNames,
			EnumerationsUseUnderscores:           true,
			SkipEnumDeduplication:                cg.SkipEnumDeduplication,
		},
		NestedDirectories:                   false,
		AbsoluteMapPaths:                    false,
		AppendEnumSuffixForSimpleUnionEnums: cg.AppendEnumSuffixForSimpleUnionEnums,
		UseModuleNameAsPathOrigin:           cg.UseModuleNameAsPathOrigin,
		PathOriginName:                      cg.PathOriginName,
	}

	var errs util.Errors
	ir, err := ygen.GenerateIR(yangFiles, includePaths, goLangMapper{GoLangMapper: gogen.NewGoLangMapper(true)}, opts)
	if err != nil {
		return nil, nil, util.AppendErr(errs, err)
	}

	var schemaStructPkgAccessor string
	if cg.GoImports.SchemaStructPkgPath != "" {
		schemaStructPkgAccessor = schemaStructPkgAlias + "."
	}

	// Get NodeDataMap for the schema.
	nodeDataMap, es := getNodeDataMap(ir, yang.CamelCase(cg.FakeRootName), schemaStructPkgAccessor, cg.PathStructSuffix, cg.PackageName, cg.PackageSuffix, cg.SplitByModule, cg.TrimPackageModulePrefix, !cg.IgnoreAtomicLists, cg.CompressBehaviour, cg.SplitPackagePaths)
	if es != nil {
		errs = util.AppendErrs(errs, es)
	}

	// Generate struct code.
	var structSnippets []GoPathStructCodeSnippet
	for _, directoryPath := range ir.OrderedDirectoryPathsByName() {
		directory := ir.Directories[directoryPath]

		structSnippet, es := generateDirectorySnippet(directory, ir.Directories, nodeDataMap, yang.CamelCase(cg.FakeRootName), cg.ExtraGenerators, schemaStructPkgAccessor, cg.PathStructSuffix, cg.GenerateWildcardPaths, cg.SplitByModule, cg.TrimPackageModulePrefix, cg.PackageName, cg.PackageSuffix, cg.SplitPackagePaths, cg.UnifyPathStructs, !cg.IgnoreAtomic, !cg.IgnoreAtomicLists, cg.CompressBehaviour)
		if es != nil {
			errs = util.AppendErrs(errs, es)
		}

		structSnippets = append(structSnippets, structSnippet...)
	}

	// Aggregate snippets by package and compute their deps.
	packages := map[string]*GeneratedPathCode{}
	for _, snippet := range structSnippets {
		if _, ok := packages[snippet.Package]; !ok {
			packages[snippet.Package] = &GeneratedPathCode{
				Deps: make(map[string]bool),
			}
		}
		packages[snippet.Package].Structs = append(packages[snippet.Package].Structs, snippet)
		for _, d := range snippet.Deps {
			packages[snippet.Package].Deps[d] = true
		}
	}
	for name, p := range packages {
		err := writeHeader(yangFiles, includePaths, name, cg, p)
		errs = util.AppendErr(errs, err)
	}

	if len(errs) == 0 {
		errs = nil
	}
	return packages, nodeDataMap, errs
}

// packageNameReplacePattern matches all characters allowed in yang modules, but
// not go packages, in addition to slash that is used to delimit path elements.
var packageNameReplacePattern = regexp.MustCompile("[._/-]")

// goPackageName returns the go package to use when generating code.
//
// - rootModuleName is the module in which the YANG tree the node is attached
// to was instantiated (rather than the module that has the same namespace as
// the node).
// - schemaPath is the path to the node without the module name and without
// choice/case statements.
// - isFakeRoot indicates whether the node belongs under the fake root.
// - Other parameters are directly passed from the generator's flags. See them
// for more information.
func goPackageName(rootModuleName string, schemaPath string, splitByModule, isFakeRoot bool, pkgName, trimPrefix, pkgSuffix string, splitPackagePaths map[string]string) string {
	if isFakeRoot {
		return pkgName
	}

	var name string
	var longestPrefix string
	for prefix, customPkgName := range splitPackagePaths {
		if len(prefix) <= len(longestPrefix) {
			continue
		}
		if strings.HasPrefix(schemaPath, prefix) {
			if customPkgName != "" {
				name = customPkgName
			} else {
				name = prefix
			}
			longestPrefix = prefix
		}
	}

	if name == "" {
		if !splitByModule {
			return pkgName
		}
		name = rootModuleName
	}

	if trimPrefix != "" {
		name = strings.TrimPrefix(name, trimPrefix)
	}

	name = packageNameReplacePattern.ReplaceAllString(name, "")
	return strings.ToLower(name) + pkgSuffix
}

// GeneratedPathCode contains generated code snippets that can be processed by the calling
// application. The generated code is divided into two types of objects - both represented
// as a slice of strings: Structs contains a set of Go structures that have been generated,
// and Enums contains the code for generated enumerated types (corresponding to identities,
// or enumerated values within the YANG models for which code is being generated). Additionally
// the header with package comment of the generated code is returned in Header, along with the
// a slice of strings containing the packages that are required for the generated Go code to
// be compiled is returned.
//
// For schemas that contain enumerated types (identities, or enumerations), a code snippet is
// returned as the EnumMap field that allows the string values from the YANG schema to be resolved.
// The keys of the map are strings corresponding to the name of the generated type, with the
// map values being maps of the int64 identifier for each value of the enumeration to the name of
// the element, as used in the YANG schema.
type GeneratedPathCode struct {
	Structs      []GoPathStructCodeSnippet // Structs is the generated set of structs representing containers or lists in the input YANG models.
	CommonHeader string                    // CommonHeader is the header that should be used for all output Go files.
	Deps         map[string]bool           // Deps is the list of packages that this package depends on.
}

// String method for GeneratedPathCode, which can be used to write all the
// generated code into a single file.
func (genCode GeneratedPathCode) String() string {
	var gotCode strings.Builder
	gotCode.WriteString(genCode.CommonHeader)
	for _, gotStruct := range genCode.Structs {
		gotCode.WriteString(gotStruct.String())
	}
	return gotCode.String()
}

// SplitFiles returns a slice of strings, each representing a file that
// together contains the entire generated code. fileN specifies the number of
// files to split the code into, and has to be between 1 and the total number
// of directory entries in the input schema. By splitting, the size of the
// output files can be roughly controlled.
func (genCode GeneratedPathCode) SplitFiles(fileN int) ([]string, error) {
	structN := len(genCode.Structs)
	if fileN < 1 || fileN > structN {
		return nil, fmt.Errorf("requested %d files, but must be between 1 and %d (number of structs)", fileN, structN)
	}

	files := make([]string, 0, fileN)
	structsPerFile := int(math.Ceil(float64(structN) / float64(fileN)))
	// Empty files could appear with certain structN/fileN combinations due
	// to the ceiling numbers being used for structsPerFile.
	// e.g. 4/3 gives two files of two structs.
	// This is a little more complex, but spreads out the structs more evenly.
	// If we instead use the floor number, and put all remainder structs in
	// the last file, we might double the last file's number of structs if we get unlucky.
	// e.g. 99/10 assigns 18 structs to the last file.
	emptyFiles := fileN - int(math.Ceil(float64(structN)/float64(structsPerFile)))
	var gotCode strings.Builder
	gotCode.WriteString(genCode.CommonHeader)
	for i, gotStruct := range genCode.Structs {
		gotCode.WriteString(gotStruct.String())
		// The last file contains the remainder of the structs.
		if i == structN-1 || (i+1)%structsPerFile == 0 {
			files = append(files, gotCode.String())
			gotCode.Reset()
			gotCode.WriteString(genCode.CommonHeader)
		}
	}
	for i := 0; i != emptyFiles; i++ {
		files = append(files, genCode.CommonHeader)
	}

	return files, nil
}

// GoPathStructCodeSnippet is used to store the generated code snippets associated with
// a particular Go struct entity (corresponding to a container, list, or leaf in the schema).
type GoPathStructCodeSnippet struct {
	// PathStructName is the name of the struct that is contained within the snippet.
	// It is stored such that callers can identify the struct to control where it
	// is output.
	PathStructName string
	// StructBase stores the basic code snippet that represents the struct that is
	// the input when code generation is performed, which includes its definition.
	StructBase string
	// ChildConstructors contains the method code snippets with the input struct as a
	// receiver, that is used to get the child path struct.
	ChildConstructors string
	// Package is the name of the package that this snippet belongs to.
	Package string
	// Deps are any packages that this snippet depends on.
	Deps []string
	// ExtraGeneration is the output of any extra generators passed into GenConfig.
	ExtraGeneration string
}

// String returns the contents of a GoPathStructCodeSnippet as a string by
// simply writing out all of its generated code.
func (g *GoPathStructCodeSnippet) String() string {
	var b strings.Builder
	for _, method := range []string{g.StructBase, g.ChildConstructors} {
		genutil.WriteIfNotEmpty(&b, method)
	}
	genutil.WriteIfNotEmpty(&b, g.ExtraGeneration)
	return b.String()
}

func genExtras(psName, fakeRootName string, compressBehaviour genutil.CompressBehaviour, nodeDataMap NodeDataMap, dir *ygen.ParsedDirectory, wildcard bool, extraGens []Generator, errs *errlist.List) string {
	node, ok := nodeDataMap[psName]
	if !ok {
		return ""
	}

	var b strings.Builder
	for _, gen := range extraGens {
		extra, err := gen(psName, fakeRootName, compressBehaviour, dir, node, wildcard)
		if err != nil {
			errs.Add(err)
			continue
		}
		if _, err := b.WriteString(extra); err != nil {
			errs.Add(err)
			continue
		}
	}
	return b.String()
}

// genExtras calls the extra generators to append to ExtraGeneration.
//
// If the PathStruct type is not present in nodeDataMap, then generation is
// skipped and no error is returned.
//
//   - nodeDataMap is the set of all nodes keyed by pathstruct name.
//   - dir is the directory to which this snippet belongs.
func (g *GoPathStructCodeSnippet) genExtras(nodeDataMap NodeDataMap, fakeRootName string, compressBehaviour genutil.CompressBehaviour, dir *ygen.ParsedDirectory, extraGens []Generator) error {
	var errs errlist.List
	g.ExtraGeneration += genExtras(g.PathStructName, fakeRootName, compressBehaviour, nodeDataMap, dir, false, extraGens, &errs)
	// Extra generation for keyed list types.
	g.ExtraGeneration += genExtras(g.PathStructName+WholeKeyedListSuffix, fakeRootName, compressBehaviour, nodeDataMap, dir, false, extraGens, &errs)

	return errs.Err()
}

// NodeDataMap is a map from the path struct type name of a schema node to its NodeData.
type NodeDataMap map[string]*NodeData

// NodeData contains information about the ygen-generated code of a YANG schema node.
type NodeData struct {
	// GoTypeName is the generated Go type name of a schema node. It is
	// qualified by the SchemaStructPkgAlias if necessary. It could be a
	// GoStruct or a leaf type.
	GoTypeName string
	// LocalGoTypeName is the generated Go type name of a schema node, but
	// always with the SchemaStructPkgAlias stripped. It could be a
	// GoStruct or a leaf type.
	LocalGoTypeName string
	// GoFieldName is the field name of the node under its parent struct.
	GoFieldName string
	// SubsumingGoStructName is the GoStruct type name corresponding to the node. If
	// the node is a leaf, then it is the parent GoStruct's name.
	SubsumingGoStructName string
	// IsLeaf indicates whether this child is a leaf node.
	IsLeaf bool
	// IsScalarField indicates a leaf that is stored as a pointer in its
	// parent struct.
	IsScalarField bool
	// IsListContainer indicates that the node represents the  container
	// surrounding a list.
	IsListContainer bool
	// HasDefault indicates whether this node has a default value
	// associated with it. This is only relevant to leaf or leaf-list
	// nodes.
	HasDefault bool
	// YANGTypeName is the type of the leaf given in the YANG file (without
	// the module prefix, if any, per goyang behaviour). If the node is not
	// a leaf this will be empty. Note that the current purpose for this is
	// to allow callers to handle certain types as special cases, but since
	// the name of the node is a very basic piece of information which
	// excludes the defining module, this is somewhat hacky, so it may be
	// removed or modified in the future.
	YANGTypeName string
	// YANGPath is the schema path of the YANG node.
	YANGPath string
	// GoPathPackageName is the Go package name containing the generated PathStruct for the schema node.
	GoPathPackageName string
	// YANGFieldName is the name of the field entry for this node, only set for leaves.
	YANGFieldName string
	// Directory is the name of the directory used to create this node. For leaves, it's the parent directory.
	DirectoryName string
	// CompressInfo contains information about a compressed path element for a node
	// which points to a path element that's compressed out.
	CompressInfo *CompressionInfo
	// ConfigFalse indicates whether the node is "config false" or "config
	// true" in its YANG schema definition.
	ConfigFalse bool
	// PathOriginName is the name of the origin for this node.
	PathOriginName string
}

// CompressionInfo contains information about a compressed path element for a
// node which points to a path element that's compressed out.
//
// e.g. for OpenConfig's /interfaces/interface, if a path points to
// /interfaces, then CompressionInfo will be populated with the following:
// - PreRelPath: []{"openconfig-interfaces:interfaces"}
// - PostRelPath: []{"openconfig-interfaces:interface"}
type CompressionInfo struct {
	// PreRelPath is the list of strings of qualified path elements prior
	// to the compressed-out node in Go list syntax.
	//
	// e.g. "openconfig-withlistval:atomic-lists", "openconfig-withlistval:atomic-lists-again"
	PreRelPathList string
	// PreRelPath is the list of strings of qualified path elements after
	// the compressed-out node in Go list syntax.
	PostRelPathList string
}

// GetOrderedNodeDataNames returns the alphabetically-sorted slice of keys
// (path struct names) for a given NodeDataMap.
func GetOrderedNodeDataNames(nodeDataMap NodeDataMap) []string {
	nodeDataNames := make([]string, 0, len(nodeDataMap))
	for name := range nodeDataMap {
		nodeDataNames = append(nodeDataNames, name)
	}
	sort.Slice(nodeDataNames, func(i, j int) bool {
		return nodeDataNames[i] < nodeDataNames[j]
	})
	return nodeDataNames
}

var (
	// goPathCommonHeaderTemplate is populated and output at the top of the generated code package
	goPathCommonHeaderTemplate = mustTemplate("commonHeader", `
{{- /**/ -}}
/*
Package {{ .PackageName }} is a generated package which contains definitions
of structs which generate gNMI paths for a YANG schema.

This package was generated by {{ .GeneratingBinary }}
using the following YANG input files:
{{- range $inputFile := .YANGFiles }}
	- {{ $inputFile }}
{{- end }}
Imported modules were sourced from:
{{- range $importPath := .IncludePaths }}
	- {{ $importPath }}
{{- end }}
*/
package {{ .PackageName }}

import (
	"reflect"
	{{- if .SchemaStructPkgPath }}
	{{ .SchemaStructPkgAlias }} "{{ .SchemaStructPkgPath }}"
	{{- end }}
	"{{ .YgotImportPath }}"
	"{{ .YgnmiImportPath }}"
	"{{ .YtypesImportPath }}"
{{- range $import := .ExtraImports }}
	"{{ $import }}"
{{- end }}
)
`)

	// goPathFakeRootTemplate defines a template for the type definition and
	// basic methods of the fakeroot object. The fakeroot object adheres to
	// the methods of PathStructInterfaceName and FakeRootBaseTypeName in
	// order to allow its path struct descendents to use the ygnmi.Resolve()
	// helper function for obtaining their absolute paths.
	goPathFakeRootTemplate = mustTemplate("fakeroot", `
// {{ .TypeName }} represents the {{ .YANGPath }} YANG schema element.
type {{ .TypeName }} struct {
	*ygnmi.{{ .FakeRootBaseTypeName }}
}

// Root returns a root path object from which YANG paths can be constructed.
func Root() *{{ .TypeName }} {
	return &{{ .TypeName }}{ygnmi.New{{- .FakeRootBaseTypeName }}()}
}
`)

	// goPathStructTemplate defines the template for the type definition of
	// a path node as well as its core method(s). A path struct/node is
	// either a container, list, or a leaf node in the openconfig schema
	// where the tree formed by the nodes mirrors the YANG
	// schema tree. The defined type stores the relative path to the
	// current node, as well as its parent node for obtaining its absolute
	// path. There are two versions of these, non-wildcard and wildcard.
	// The wildcard version is simply a type to indicate that the path it
	// holds contains a wildcard, but is otherwise the exact same.
	goPathStructTemplate = mustTemplate("struct", `
// {{ .TypeName }} represents the {{ .YANGPath }} YANG schema element.
type {{ .TypeName }} struct {
	*ygnmi.{{ .PathBaseTypeName }}
{{- if .GenerateParentField }}
	parent ygnmi.PathStruct
{{- end }}
	{{- .ExtraFields }}
}

{{- if .GenerateWildcardPaths }}

// {{ .TypeName }}{{ .WildcardSuffix }} represents the wildcard version of the {{ .YANGPath }} YANG schema element.
type {{ .TypeName }}{{ .WildcardSuffix }} struct {
	*ygnmi.{{ .PathBaseTypeName }}
{{- if .GenerateParentField }}
	parent ygnmi.PathStruct
{{- end }}
	{{- .ExtraWildcardFields }}
}
{{- end }}

// PathOrigin returns the name of the origin for the path object.
func (n *{{.TypeName}}) PathOriginName() string {
     return "{{ $.PathOriginName }}"
}
`)

	// goPathChildConstructorTemplate generates the child constructor method
	// for a generated struct by returning an instantiation of the child's
	// path struct object.
	goPathChildConstructorTemplate = mustTemplate("childConstructor", `
// {{ .MethodName }} ({{ .YANGNodeType }}): {{ .YANGDescription }}
// 	Defining module:      "{{ .DefiningModuleName }}"
// 	Instantiating module: "{{ .InstantiatingModuleName }}"
// 	Path from parent:     "{{ .RelPath }}"
// 	Path from root:       "{{ .AbsPath }}"
{{- if .KeyParamDocStrs }}
//
{{- end }}
{{- range $paramDocStr := .KeyParamDocStrs }}
// 	{{ $paramDocStr }}
{{- end }}
func (n *{{ .Struct.TypeName }}) {{ .MethodName -}} ({{ .KeyParamListStr }}) *{{ .ChildPkgAccessor }}{{ .TypeName }} {
	ps := &{{ .ChildPkgAccessor }}{{ .TypeName }}{
		{{ .Struct.PathBaseTypeName }}: ygnmi.New{{ .Struct.PathBaseTypeName }}(
			[]string{ {{- .RelPathList -}} },
			map[string]interface{}{ {{- .KeyEntriesStr -}} },
			n,
		),
{{- if .GenerateParentField }}
		parent: n,
{{- end }}
	}
{{- if .IsWildcard }}
	{{- .ExtraWildcardInits }}
{{- else }}
	{{- .ExtraInits }}
{{- end }}
	return ps
}
`)

	// goKeyBuilderTemplate generates a setter for a list key. This is used in the
	// builder style for the list API.
	goKeyBuilderTemplate = mustTemplate("goKeyBuilder", `
// {{ .MethodName }} sets {{ .TypeName }}'s key "{{ .KeySchemaName }}" to the specified value.
// {{ .KeyParamDocStr }}
func (n *{{ .TypeName }}) {{ .MethodName }}({{ .KeyParamName }} {{ .KeyParamType }}) *{{ .TypeName }} {
	ygnmi.ModifyKey(n.NodePath, "{{ .KeySchemaName }}", {{ .KeyParamName }})
	return n
}
`)
)

// mustTemplate generates a template.Template for a particular named source template
func mustTemplate(name, src string) *template.Template {
	return template.Must(template.New(name).Parse(src))
}

// getNodeDataMap returns the NodeDataMap for the provided schema given its
// parsed information.
// packageName, splitByModule, and trimOCPackage are used to determine
// the generated Go package name for the generated PathStructs.
//
// TODO(wenbli): Change this function to be callable while traversing the IR
// rather than traversing the IR itself again.
func getNodeDataMap(ir *ygen.IR, fakeRootName, schemaStructPkgAccessor, pathStructSuffix, packageName, packageSuffix string, splitByModule bool, trimPrefix string, generateAtomicLists bool, compressBehaviour genutil.CompressBehaviour, splitPackagePaths map[string]string) (NodeDataMap, util.Errors) {
	nodeDataMap := NodeDataMap{}
	var errs util.Errors
	for _, dir := range ir.Directories {
		if dir.IsFakeRoot {
			// Since we always generate the fake root, we add the
			// fake root GoStruct to the data map as well.
			nodeDataMap[dir.Name+pathStructSuffix] = &NodeData{
				GoTypeName:            "*" + schemaStructPkgAccessor + fakeRootName,
				LocalGoTypeName:       "*" + fakeRootName,
				GoFieldName:           "",
				SubsumingGoStructName: fakeRootName,
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "",
				YANGPath:              "/",
				GoPathPackageName:     goPackageName(dir.RootElementModule, dir.SchemaPath, splitByModule, true, packageName, trimPrefix, packageSuffix, splitPackagePaths),
				DirectoryName:         dir.Path,
			}
		}

		goFieldNameMap := ygen.GoFieldNameMap(dir)
		for _, fieldName := range dir.OrderedFieldNames() {
			field := dir.Fields[fieldName]
			pathStructName, err := getFieldTypeName(dir, fieldName, goFieldNameMap[fieldName], ir.Directories, pathStructSuffix)
			if err != nil {
				errs = util.AppendErr(errs, err)
				continue
			}

			mType := field.LangType
			isLeaf := mType != nil
			var fieldDir *ygen.ParsedDirectory
			if !isLeaf {
				var ok bool
				if fieldDir, ok = ir.Directories[field.YANGDetails.Path]; !ok {
					errs = util.AppendErr(errs, fmt.Errorf("field was a non-leaf but could not find it in IR: %s", field.YANGDetails.Path))
					continue
				}
			}

			subsumingGoStructName := dir.Name
			if !isLeaf {
				subsumingGoStructName = fieldDir.Name
			}

			var goTypeName, localGoTypeName string
			switch {
			case !isLeaf:
				goTypeName = "*" + schemaStructPkgAccessor + subsumingGoStructName
				localGoTypeName = "*" + subsumingGoStructName
			case field.Type == ygen.LeafListNode && ygen.IsYgenDefinedGoType(mType):
				goTypeName = "[]" + schemaStructPkgAccessor + mType.NativeType
				localGoTypeName = "[]" + mType.NativeType
			case ygen.IsYgenDefinedGoType(mType):
				goTypeName = schemaStructPkgAccessor + mType.NativeType
				localGoTypeName = mType.NativeType
			case field.Type == ygen.LeafListNode:
				goTypeName = "[]" + mType.NativeType
				localGoTypeName = goTypeName
			default:
				goTypeName = mType.NativeType
				if localGoTypeName == "" {
					localGoTypeName = goTypeName
				}
			}

			var yangTypeName string
			if field.Flags != nil {
				yangTypeName = field.Flags[yangTypeNameFlagKey]
			}
			nodeData := NodeData{
				GoTypeName:            goTypeName,
				LocalGoTypeName:       localGoTypeName,
				GoFieldName:           goFieldNameMap[fieldName],
				SubsumingGoStructName: subsumingGoStructName,
				IsLeaf:                isLeaf,
				IsScalarField:         gogen.IsScalarField(field),
				HasDefault:            isLeaf && (len(field.YANGDetails.Defaults) > 0 || mType.DefaultValue != nil),
				YANGTypeName:          yangTypeName,
				YANGPath:              field.YANGDetails.Path,
				GoPathPackageName:     goPackageName(field.YANGDetails.RootElementModule, field.YANGDetails.SchemaPath, splitByModule, false, packageName, trimPrefix, packageSuffix, splitPackagePaths),
				DirectoryName:         field.YANGDetails.Path,
				ConfigFalse:           field.YANGDetails.ConfigFalse,
			}
			// If NodeDetails.YANGDetails.Origin has a value, the value is set to the PathOriginName of the node.
			// If it is the empty string, we currently set it explicitly to "openconfig"
			// until the 'origin' extension of YANG model is parsed by ygot.
			origin := field.YANGDetails.Origin
			// TODO: remove this workaround when the origin attribute is returned in the ygot IR.
			if origin == "" {
				origin = "openconfig" // explicitly override the empty origin
			}
			nodeData.PathOriginName = origin

			switch {
			case !isLeaf && isKeyedList(fieldDir) && !(generateAtomicLists && isCompressedAtomicList(fieldDir)): // Non-atomic lists
				// Generate path to list element.
				nodeDataMap[pathStructName] = &nodeData
				fallthrough // Generate Map-suffixed PathStruct
			case !isLeaf && isKeyedList(fieldDir): // Atomic lists
				if !generateAtomicLists {
					break
				}
				// Generate path to container surrounding list element (i.e. Map-suffixed PathStruct).
				nodeData := nodeData
				nodeData.IsListContainer = true
				nodeData.SubsumingGoStructName = dir.Name
				if field.YANGDetails.OrderedByUser {
					typeName := gogen.OrderedMapTypeName(fieldDir.Name)
					nodeData.GoTypeName = "*" + schemaStructPkgAccessor + typeName
					nodeData.LocalGoTypeName = "*" + typeName
				} else {
					_, keyType, isDefined, err := gogen.UnorderedMapTypeName(field.YANGDetails.Path, goFieldNameMap[fieldName], dir.Name, ir.Directories)
					if err != nil {
						errs = util.AppendErr(errs, fmt.Errorf("could not determine type name of unordered keyed map: %v", err))
						continue
					}
					nonLocalKeyType := keyType
					if isDefined {
						nonLocalKeyType = schemaStructPkgAccessor + keyType
					}
					nodeData.GoTypeName = fmt.Sprintf("map[%s]%s", nonLocalKeyType, nodeData.GoTypeName)
					nodeData.LocalGoTypeName = fmt.Sprintf("map[%s]%s", keyType, nodeData.LocalGoTypeName)
				}

				if compressBehaviour.CompressEnabled() {
					relPath := longestPath(field.MappedPaths)
					relMods := longestPath(field.MappedPathModules)
					if gotLen := len(relPath); gotLen != 2 {
						errs = util.AppendErr(errs, fmt.Errorf("expected two path elements for the relative path of an atomic list or ordered list, got %d: %v", gotLen, relPath))
						continue
					}
					nodeData.CompressInfo = &CompressionInfo{
						PreRelPathList:  fmt.Sprintf(`"%s:%s"`, relMods[0], relPath[0]),
						PostRelPathList: fmt.Sprintf(`"%s:%s"`, relMods[1], relPath[1]),
					}
				}

				nodeDataMap[pathStructName+WholeKeyedListSuffix] = &nodeData
			case isLeaf:
				nodeData.YANGFieldName = fieldName
				nodeData.DirectoryName = dir.Path
				// A leaf resides in its parent's package since parents may need to instantiate leaves
				// with a reference to its unexported "parent" field (see .GenerateParentField).
				nodeData.GoPathPackageName = goPackageName(dir.RootElementModule, dir.SchemaPath, splitByModule, dir.IsFakeRoot, packageName, trimPrefix, packageSuffix, splitPackagePaths)
				fallthrough
			default:
				nodeDataMap[pathStructName] = &nodeData
			}
		}
	}

	if len(errs) != 0 {
		return nil, errs
	}
	return nodeDataMap, nil
}

// writeHeader parses the yangFiles from the includePaths, and fills the given
// *GeneratedPathCode with the header of the generated Go path code.
func writeHeader(yangFiles, includePaths []string, packageName string, cg *GenConfig, genCode *GeneratedPathCode) error {
	// Build input to the header template which stores parameters which are included
	// in the header of generated code.
	s := struct {
		GoImports                        // GoImports contains package import options.
		PackageName             string   // PackageName is the name that should be used for the generating package.
		GeneratingBinary        string   // GeneratingBinary is the name of the binary calling the generator library.
		YANGFiles               []string // YANGFiles contains the list of input YANG source files for code generation.
		IncludePaths            []string // IncludePaths contains the list of paths that included modules were searched for in.
		SchemaStructPkgAlias    string   // SchemaStructPkgAlias is the package alias for the imported ygen-generated file.
		PathBaseTypeName        string   // PathBaseTypeName is the type name of the common embedded path struct.
		PathStructInterfaceName string   // PathStructInterfaceName is the name of the interface which all path structs implement.
		FakeRootTypeName        string   // FakeRootTypeName is the type name of the fakeroot node in the generated code.
		ExtraImports            []string // ExtraImports for path structs that are in a different package.
	}{
		GoImports:               cg.GoImports,
		PackageName:             packageName,
		GeneratingBinary:        cg.GeneratingBinary,
		YANGFiles:               yangFiles,
		IncludePaths:            includePaths,
		SchemaStructPkgAlias:    schemaStructPkgAlias,
		PathBaseTypeName:        ygnmi.PathBaseTypeName,
		PathStructInterfaceName: ygnmi.PathStructInterfaceName,
		FakeRootTypeName:        yang.CamelCase(cg.FakeRootName),
	}
	// Create an ordered list of imports to include in the header.
	for dep := range genCode.Deps {
		s.ExtraImports = append(s.ExtraImports, fmt.Sprintf("%s/%s", cg.BasePackagePath, dep))
	}
	sort.Slice(s.ExtraImports, func(i, j int) bool { return s.ExtraImports[i] < s.ExtraImports[j] })

	var common strings.Builder
	if err := goPathCommonHeaderTemplate.Execute(&common, s); err != nil {
		return err
	}

	genCode.CommonHeader = common.String()
	return nil
}

// goPathStructData stores template information needed to generate a struct
// field's type definition.
type goPathStructData struct {
	// TypeName is the type name of the struct being output.
	TypeName string
	// YANGPath is the schema path of the struct being output.
	YANGPath string
	// PathBaseTypeName is the type name of the common embedded path struct.
	PathBaseTypeName string
	// PathStructInterfaceName is the name of the interface which all path structs implement.
	PathStructInterfaceName string
	// FakeRootBaseTypeName is the type name of the fake root struct which
	// should be embedded within the fake root path struct.
	FakeRootBaseTypeName string
	// WildcardSuffix is the suffix given to the wildcard versions of
	// each node that distinguishes each from its non-wildcard counterpart.
	WildcardSuffix string
	// GenerateWildcardPaths means to generate wildcard nodes and paths.
	GenerateWildcardPaths bool
	// GenerateParentField means to generate a parent PathStruct field in
	// the PathStruct so that it's accessible by code within the generated
	// file's package.
	GenerateParentField bool
	// ExtraFields are extra fields in the generated struct by extra
	// generators.
	ExtraFields string
	// ExtraFields are extra fields in the generated struct by extra
	// generators.
	ExtraWildcardFields string
	// PathOriginName is the name of the origin for the PathStruct.
	PathOriginName string
}

// genExtraFields calls the extra generators to append extra fields.
//
// If the PathStruct type is not present in nodeDataMap, then generation is
// skipped and no error is returned.
//
//   - nodeDataMap is the set of all nodes keyed by pathstruct name.
//   - dir is the directory to which this snippet belongs.
func (s *goPathStructData) genExtraFields(nodeDataMap NodeDataMap, fakeRootName string, compressBehaviour genutil.CompressBehaviour, dir *ygen.ParsedDirectory, extraFieldGens []Generator) error {
	var errs errlist.List
	s.ExtraFields += genExtras(s.TypeName, fakeRootName, compressBehaviour, nodeDataMap, dir, false, extraFieldGens, &errs)
	s.ExtraWildcardFields += genExtras(s.TypeName, fakeRootName, compressBehaviour, nodeDataMap, dir, true, extraFieldGens, &errs)

	return errs.Err()
}

func (s *goPathStructData) clearExtraFields() {
	s.ExtraFields = ""
	s.ExtraWildcardFields = ""
}

// getStructData returns the goPathStructData corresponding to a
// ParsedDirectory, which is used to store the attributes of the template for
// which code is being generated.
func getStructData(directory *ygen.ParsedDirectory, pathStructSuffix string, generateWildcardPaths bool) goPathStructData {
	return goPathStructData{
		TypeName:                directory.Name + pathStructSuffix,
		YANGPath:                directory.Path,
		PathBaseTypeName:        ygnmi.PathBaseTypeName,
		FakeRootBaseTypeName:    ygnmi.FakeRootBaseTypeName,
		PathStructInterfaceName: ygnmi.PathStructInterfaceName,
		WildcardSuffix:          WildcardSuffix,
		GenerateWildcardPaths:   generateWildcardPaths,
	}
}

// goPathFieldData stores template information needed to generate a struct
// field's child constructor method.
type goPathFieldData struct {
	MethodName              string           // MethodName is the name of the method that can be called to get to this field.
	SchemaName              string           // SchemaName is the field's original name in the schema.
	TypeName                string           // TypeName is the type name of the returned struct.
	YANGNodeType            string           // YANGNodeType is the type of YANG node for the node (e.g. "list", "container", "leaf").
	YANGDescription         string           // YANGDescription is the description for the node from its YANG definition.
	DefiningModuleName      string           // DefiningModuleName is the defining module for the node.
	InstantiatingModuleName string           // InstantiatingModuleName is the module that instantiated the node.
	AbsPath                 string           // AbsPath is the full path from root (not including keys).
	RelPath                 string           // RelPath is the relative path from its containing struct.
	RelPathList             string           // RelPathList is the list of strings that form the relative path from its containing struct.
	Struct                  goPathStructData // Struct stores template information for the field's containing struct.
	KeyParamListStr         string           // KeyParamListStr is the parameter list of the field's accessor method.
	KeyEntriesStr           string           // KeyEntriesStr is an ordered list of comma-separated ("schemaName": unique camel-case name) for a list's keys.
	KeyParamDocStrs         []string         // KeyParamDocStrs is an ordered slice of docstrings documenting the types of each list key parameter.
	ChildPkgAccessor        string           // ChildPkgAccessor is used if the child path struct exists in another package.
	// GenerateParentField means to generate a parent PathStruct field in
	// the PathStruct so that it's accessible by code within the generated
	// file's package.
	GenerateParentField bool
	// ExtraInits are extra field initializations in the generated struct
	// by extra generators.
	ExtraInits string
	// ExtraInits are extra field initializations in the generated struct
	// by extra generators.
	ExtraWildcardInits string
	// IsWildcard indicates whether to generate this as a wildcard or not a wildcard.
	IsWildcard bool
}

// genExtraInits calls the extra generators to append extra inits.
//
// If the PathStruct type is not present in nodeDataMap, then generation is
// skipped and no error is returned.
//
//   - nodeDataMap is the set of all nodes keyed by pathstruct name.
//   - dir is the directory to which this snippet belongs.
func (s *goPathFieldData) genExtraInits(nodeDataMap NodeDataMap, fakeRootName string, compressBehaviour genutil.CompressBehaviour, dir *ygen.ParsedDirectory, extraInitGens []Generator) error {
	var errs errlist.List
	s.ExtraInits += genExtras(s.TypeName, fakeRootName, compressBehaviour, nodeDataMap, dir, false, extraInitGens, &errs)
	s.ExtraWildcardInits += genExtras(s.TypeName, fakeRootName, compressBehaviour, nodeDataMap, dir, true, extraInitGens, &errs)

	return errs.Err()
}

func (s *goPathFieldData) clearExtraInit() {
	s.ExtraInits = ""
	s.ExtraWildcardInits = ""
}

func isKeyedList(directory *ygen.ParsedDirectory) bool {
	return (directory.Type == ygen.List || directory.Type == ygen.OrderedList) && len(directory.ListKeys) > 0
}

func isCompressedAtomicList(directory *ygen.ParsedDirectory) bool {
	return directory.Type == ygen.OrderedList || (directory.Type == ygen.List && directory.CompressedTelemetryAtomic)
}

// generateDirectorySnippet generates all Go code associated with a schema node
// (container, list, leaf, or fakeroot), all of which have a corresponding
// struct onto which to attach the necessary methods for path generation.
// When generating code for the fakeroot, several structs may be returned,
// one package containing the fake root struct and one package for each of the
// fake root's child lists that uses the builder API methods,
// since they must be defined in their respective child packages. Otherwise,
// the returned slice will only have a single struct containing the all the code.
// The code comprises of the type definition for the struct, and all accessors to
// the fields of the struct. directory is the parsed information of a schema
// node, and directories is a map from path to a parsed schema node for all
// directory nodes in the schema.
func generateDirectorySnippet(directory *ygen.ParsedDirectory, directories map[string]*ygen.ParsedDirectory, nodeDataMap NodeDataMap, fakeRootName string, extraGens ExtraGenerators, schemaStructPkgAccessor, pathStructSuffix string,
	generateWildcardPaths, splitByModule bool, trimPrefix, packageName, packageSuffix string, splitPackagePaths map[string]string, unified bool, generateAtomic, generateAtomicLists bool, compressBehaviour genutil.CompressBehaviour) ([]GoPathStructCodeSnippet, util.Errors) {

	var errs util.Errors
	// structBuf is used to store the code associated with the struct defined for
	// the target YANG entity.
	var structBuf strings.Builder
	var methodBuf strings.Builder

	// Output struct snippets.
	structData := getStructData(directory, pathStructSuffix, generateWildcardPaths)

	if directory.IsFakeRoot {
		// Fakeroot has its unique output.
		if err := goPathFakeRootTemplate.Execute(&structBuf, structData); err != nil {
			return nil, util.AppendErr(errs, err)
		}
	} else {
		if err := structData.genExtraFields(nodeDataMap, fakeRootName, compressBehaviour, directory, extraGens.StructFields); err != nil {
			return nil, util.AppendErr(errs, err)
		}
		// The PathOriginName from the nodeDataMap is set to that of the structData
		if nodeData, ok := nodeDataMap[structData.TypeName]; ok {
			structData.PathOriginName = nodeData.PathOriginName
		}

		if err := goPathStructTemplate.Execute(&structBuf, structData); err != nil {
			return nil, util.AppendErr(errs, err)
		}

		// Generate Map PathStructs for keyed list types.
		if isKeyedList(directory) {
			structData := structData
			structData.TypeName += WholeKeyedListSuffix
			structData.clearExtraFields()

			if err := structData.genExtraFields(nodeDataMap, fakeRootName, compressBehaviour, directory, extraGens.StructFields); err != nil {
				return nil, util.AppendErr(errs, err)
			}
			if err := goPathStructTemplate.Execute(&structBuf, structData); err != nil {
				return nil, util.AppendErr(errs, err)
			}
		}
	}

	parentPackage := goPackageName(directory.RootElementModule, directory.SchemaPath, splitByModule, directory.IsFakeRoot, packageName, trimPrefix, packageSuffix, splitPackagePaths)
	nonLeafSnippet := GoPathStructCodeSnippet{
		PathStructName: structData.TypeName,
		StructBase:     structBuf.String(),
		Package:        parentPackage,
	}
	errs = util.AppendErr(errs, nonLeafSnippet.genExtras(nodeDataMap, fakeRootName, compressBehaviour, directory, extraGens.Extras))

	// Since it is not possible for gNMI to refer to individual nodes
	// underneath an unkeyed list or a telemetry-atomic node (ordered lists
	// are always assumed to be telemetry-atomic), prevent constructing
	// paths that go below it by breaking the chain here.
	switch {
	case directory.Type == ygen.List && len(directory.ListKeys) == 0, generateAtomic && directory.TelemetryAtomic, generateAtomicLists && isCompressedAtomicList(directory):
		return []GoPathStructCodeSnippet{nonLeafSnippet}, nil
	}

	deps := map[string]bool{}
	listBuilderAPIBufs := map[string]*strings.Builder{}

	snippets := []GoPathStructCodeSnippet{}

	goFieldNameMap := ygen.GoFieldNameMap(directory)
	// Generate child constructor snippets for all fields of the node.
	// Alphabetically order fields to produce deterministic output.
	for _, fName := range directory.OrderedFieldNames() {
		field, ok := directory.Fields[fName]
		if !ok {
			errs = util.AppendErr(errs, fmt.Errorf("generateDirectorySnippet: field %s not found in directory %v", fName, directory))
			continue
		}
		goFieldName := goFieldNameMap[fName]
		var childPkgAccessor string

		// The most common case is that list builder API is written to same package as the rest of the child methods.
		buildBuf := &methodBuf

		// Only the fake root could be importing a child path struct from another package.
		// If it is, add that package as a dependency and set the accessor.
		if field.Type == ygen.ContainerNode || field.Type == ygen.ListNode {
			childPackage := goPackageName(field.YANGDetails.RootElementModule, field.YANGDetails.SchemaPath, splitByModule, false, packageName, trimPrefix, packageSuffix, splitPackagePaths)
			if parentPackage != childPackage {
				deps[childPackage] = true
				childPkgAccessor = childPackage + "."
				// The fake root could be generating a list builder API for one of its children which is in another package.
				// Write any list builders into the map, keyed by the child package name.
				if _, ok := listBuilderAPIBufs[childPackage]; !ok {
					listBuilderAPIBufs[childPackage] = &strings.Builder{}
				}
				buildBuf = listBuilderAPIBufs[childPackage]
			}
		}

		if es := generateChildConstructors(&methodBuf, buildBuf, directory, fName, goFieldName, directories, fakeRootName, schemaStructPkgAccessor, pathStructSuffix, generateWildcardPaths, childPkgAccessor, unified, generateAtomicLists, compressBehaviour, nodeDataMap, extraGens); es != nil {
			errs = util.AppendErrs(errs, es)
		}

		// Since leaves don't have their own Directory entries, we need
		// to output their struct snippets somewhere, and here is
		// convenient.
		if field.Type == ygen.LeafNode || field.Type == ygen.LeafListNode {
			var buf strings.Builder
			leafTypeName, err := getFieldTypeName(directory, fName, goFieldName, directories, pathStructSuffix)
			if err != nil {
				errs = util.AppendErr(errs, err)
			} else {
				structData := goPathStructData{
					TypeName:                leafTypeName,
					YANGPath:                field.YANGDetails.Path,
					PathBaseTypeName:        ygnmi.PathBaseTypeName,
					PathStructInterfaceName: ygnmi.PathStructInterfaceName,
					WildcardSuffix:          WildcardSuffix,
					GenerateWildcardPaths:   generateWildcardPaths,
					GenerateParentField:     unified,
				}
				errs = util.AppendErr(errs, structData.genExtraFields(nodeDataMap, fakeRootName, compressBehaviour, directory, extraGens.StructFields))
				errs = util.AppendErr(errs, goPathStructTemplate.Execute(&buf, structData))
			}
			leafSnippet := GoPathStructCodeSnippet{
				PathStructName: leafTypeName,
				StructBase:     buf.String(),
				Package:        goPackageName(directory.RootElementModule, directory.SchemaPath, splitByModule, directory.IsFakeRoot, packageName, trimPrefix, packageSuffix, splitPackagePaths),
			}
			errs = util.AppendErr(errs, leafSnippet.genExtras(nodeDataMap, fakeRootName, compressBehaviour, directory, extraGens.Extras))
			snippets = append(snippets, leafSnippet)
		}
	}

	if len(errs) == 0 {
		errs = nil
	}

	nonLeafSnippet.ChildConstructors = methodBuf.String()
	for dep := range deps {
		nonLeafSnippet.Deps = append(nonLeafSnippet.Deps, dep)
	}

	snippets = append(snippets, nonLeafSnippet)

	for pkg, build := range listBuilderAPIBufs {
		if code := build.String(); code != "" {
			snippets = append(snippets, GoPathStructCodeSnippet{
				PathStructName:    structData.TypeName,
				Package:           pkg,
				ChildConstructors: code,
			})
		}
	}

	return snippets, errs
}

func longestPath(ss [][]string) []string {
	var longest []string
	for _, s := range ss {
		if longest == nil {
			longest = s
			continue
		}
		if len(s) > len(longest) {
			longest = s
		}
	}
	return longest
}

func relPathListFn(path []string) string {
	return `"` + strings.Join(path, `", "`) + `"`
}

// generateChildConstructors generates and writes to methodBuf the Go methods
// that returns an instantiation of the child node's path struct object.
// When this is called on the fakeroot, the list builder API's methods
// need to be in the child package as opposed to the fakeroot's package.
// In all other cases, methodBuf and builderBuf can point to the same buffer.
// The func takes as input the buffers to store the method, a directory, the field name
// of the directory identifying the child yang.Entry, a directory-level unique
// field name to be used as the generated method's name and the incremental
// type name of of the child path struct, and a map of all directories of the
// whole schema keyed by their schema paths.
func generateChildConstructors(methodBuf *strings.Builder, builderBuf *strings.Builder, directory *ygen.ParsedDirectory, directoryFieldName string, goFieldName string, directories map[string]*ygen.ParsedDirectory, fakeRootName, schemaStructPkgAccessor,
	pathStructSuffix string, generateWildcardPaths bool, childPkgAccessor string, unified bool, generateAtomicLists bool, compressBehaviour genutil.CompressBehaviour, nodeDataMap NodeDataMap, extraGens ExtraGenerators) []error {

	field, ok := directory.Fields[directoryFieldName]
	if !ok {
		return []error{fmt.Errorf("generateChildConstructors: field %s not found in directory %v", directoryFieldName, directory)}
	}
	fieldTypeName, err := getFieldTypeName(directory, directoryFieldName, goFieldName, directories, pathStructSuffix)
	if err != nil {
		return []error{err}
	}

	structData := getStructData(directory, pathStructSuffix, generateWildcardPaths)
	// When there are multiple paths (this happens in the compressed
	// structs), the longest path is the non-key path. This is the one we
	// want to use since the key is "compressed out".
	relPath := longestPath(field.MappedPaths)
	// Be nil-tolerant for these two attributes. In real deployments (i.e.
	// not tests), these should be populated. Since these are just use for
	// documentation, it is not critical that they are populated.

	// Make copies of the path, to avoid modifying the directories.
	schemaPath := field.YANGDetails.SchemaPath
	path := make([]string, len(relPath))
	copy(path, relPath)

	// When generating unified path structs (one set for both state and config),
	// a leaf struct with a shadow path could be either state or config, so replace state with a wildcard.
	if unified && len(field.ShadowMappedPaths) > 0 && (field.Type == ygen.LeafNode || field.Type == ygen.LeafListNode) {
		path[0] = "*"
		parent := "/state/"
		if compressBehaviour == genutil.PreferIntendedConfig {
			parent = "/config/"
		}
		schemaPath = strings.ReplaceAll(schemaPath, parent, "/*/")
	}

	if err := structData.genExtraFields(nodeDataMap, fakeRootName, compressBehaviour, directory, extraGens.StructFields); err != nil {
		return []error{err}
	}

	fieldData := goPathFieldData{
		MethodName:              goFieldName,
		TypeName:                fieldTypeName,
		SchemaName:              field.Name,
		YANGNodeType:            field.Type.String(),
		YANGDescription:         strings.ReplaceAll(field.YANGDetails.Description, "\n", "\n// "),
		DefiningModuleName:      field.YANGDetails.DefiningModule,
		InstantiatingModuleName: field.YANGDetails.RootElementModule,
		AbsPath:                 schemaPath,
		Struct:                  structData,
		RelPath:                 strings.Join(path, `/`),
		RelPathList:             relPathListFn(path),
		ChildPkgAccessor:        childPkgAccessor,
	}

	isUnderFakeRoot := directory.IsFakeRoot

	// This is expected to be nil for leaf fields.
	fieldDirectory := directories[field.YANGDetails.Path]

	if err := fieldData.genExtraInits(nodeDataMap, fakeRootName, compressBehaviour, directory, extraGens.StructInits); err != nil {
		return []error{err}
	}

	switch {
	case field.Type == ygen.ListNode && len(fieldDirectory.ListKeys) == 0: // unkeyed lists
		// Generate a single wildcard constructor for keyless-lists.
		if errs := generateChildConstructorsForListBuilderFormat(methodBuf, builderBuf, nil, nil, fieldData, isUnderFakeRoot, schemaStructPkgAccessor); len(errs) > 0 {
			return errs
		}
		return nil
	case field.Type == ygen.ListNode && !(generateAtomicLists && isCompressedAtomicList(fieldDirectory)): // non-atomic keyed lists
		if generateWildcardPaths {
			if errs := generateChildConstructorsForListBuilderFormat(methodBuf, builderBuf, fieldDirectory.ListKeys, fieldDirectory.ListKeyYANGNames, fieldData, isUnderFakeRoot, schemaStructPkgAccessor); len(errs) > 0 {
				return errs
			}
		}
		if errs := generateChildConstructorsForList(methodBuf, fieldDirectory.ListKeys, fieldDirectory.ListKeyYANGNames, fieldData, isUnderFakeRoot, generateWildcardPaths, schemaStructPkgAccessor); len(errs) > 0 {
			return errs
		}
		fallthrough
	case field.Type == ygen.ListNode: // atomic keyed lists
		if !generateAtomicLists {
			return nil
		}
		if gotLen := len(path); compressBehaviour.CompressEnabled() && gotLen != 2 {
			return []error{fmt.Errorf("expected two path elements for the relative path of an atomic list or ordered list, got %d: %v", gotLen, path)}
		}
		fieldData.RelPathList = relPathListFn(path[:1])
		fieldData.TypeName += WholeKeyedListSuffix
		fieldData.MethodName += WholeKeyedListSuffix

		fieldData.clearExtraInit()
		if err := fieldData.genExtraInits(nodeDataMap, fakeRootName, compressBehaviour, directory, extraGens.StructInits); err != nil {
			return []error{err}
		}

		fallthrough
	default: // non-lists
		return generateChildConstructorsForLeafOrContainer(methodBuf, fieldData, isUnderFakeRoot, generateWildcardPaths, unified, field.Type == ygen.LeafNode || field.Type == ygen.LeafListNode)
	}
}

// generateChildConstructorsForLeafOrContainer writes into methodBuf the child
// constructor snippets for the container or leaf template output information
// contained in fieldData.
func generateChildConstructorsForLeafOrContainer(methodBuf *strings.Builder, fieldData goPathFieldData, isUnderFakeRoot, generateWildcardPaths, unified, isLeaf bool) []error {
	// Generate child constructor for the non-wildcard version of the parent struct.
	var errors []error
	if unified && isLeaf {
		fieldData.GenerateParentField = true
	}
	if err := goPathChildConstructorTemplate.Execute(methodBuf, fieldData); err != nil {
		errors = append(errors, err)
	}

	// The root node doesn't have a wildcard version of itself.
	if isUnderFakeRoot {
		return errors
	}

	if generateWildcardPaths {
		fieldData := fieldData
		// Generate child constructor for the wildcard version of the parent struct.
		fieldData.TypeName += WildcardSuffix
		fieldData.Struct.TypeName += WildcardSuffix
		fieldData.IsWildcard = true

		if err := goPathChildConstructorTemplate.Execute(methodBuf, fieldData); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// generateChildConstructorsForListBuilderFormat writes into methodBuf the
// child constructor method snippets for the list keys. The builder API helpers
// are written into the builderBuf, this allows the helpers to be written to
// package the child belongs to, not the parent if the child belongs in its own
// package. fieldData contains the childConstructor template output information
// for if the node were a container (which contains a subset of the basic
// information required for the list constructor methods).
func generateChildConstructorsForListBuilderFormat(methodBuf *strings.Builder, builderBuf *strings.Builder, keys map[string]*ygen.ListKey, keyNames []string, fieldData goPathFieldData, isUnderFakeRoot bool, schemaStructPkgAccessor string) []error {
	var errors []error
	// List of function parameters as would appear in the method definition.
	keyParams, err := makeKeyParams(keys, keyNames, schemaStructPkgAccessor)
	if err != nil {
		return append(errors, err)
	}
	keyN := len(keyParams)

	// Initialize ygnmi.NodePath's key list with wildcard values.
	var keyEntryStrs []string
	for i := 0; i != keyN; i++ {
		keyEntryStrs = append(keyEntryStrs, fmt.Sprintf(`"%s": "*"`, keyParams[i].name))
	}
	fieldData.KeyEntriesStr = strings.Join(keyEntryStrs, ", ")

	// There are no initial key parameters for the builder API.
	fieldData.KeyParamListStr = ""

	// Set the child type to be the wildcard version.
	fieldData.TypeName += WildcardSuffix

	// Add Builder suffix to the child constructor method name.
	fieldData.MethodName += BuilderCtorSuffix

	// Generate builder constructor method for non-wildcard version of parent struct.
	fieldData.IsWildcard = true // builder's child constructors are always wildcarded.
	if err := goPathChildConstructorTemplate.Execute(methodBuf, fieldData); err != nil {
		errors = append(errors, err)
	}

	// The root node doesn't have a wildcard version of itself.
	if !isUnderFakeRoot {
		// Generate builder constructor method for wildcard version of parent struct.
		fieldData.Struct.TypeName += WildcardSuffix
		if err := goPathChildConstructorTemplate.Execute(methodBuf, fieldData); err != nil {
			errors = append(errors, err)
		}
	}

	// A list with a single key doesn't need key builder methods because specifying
	// all keys or no keys are the only ways to construct its keys.
	if keyN == 1 {
		return errors
	}

	// Generate key-builder methods for the wildcard version of the PathStruct.
	// Although non-wildcard PathStruct is unnecessary, it is kept for generation simplicity.
	for i := 0; i != keyN; i++ {
		if err := goKeyBuilderTemplate.Execute(builderBuf,
			struct {
				MethodName     string
				TypeName       string
				KeySchemaName  string
				KeyParamType   string
				KeyParamName   string
				KeyParamDocStr string
			}{
				MethodName:     BuilderKeyPrefix + keyParams[i].varName,
				TypeName:       fieldData.TypeName,
				KeySchemaName:  keyParams[i].name,
				KeyParamName:   keyParams[i].varName,
				KeyParamType:   keyParams[i].typeName,
				KeyParamDocStr: keyParams[i].varName + ": " + keyParams[i].typeDocString,
			}); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// generateChildConstructorsForList writes into methodBuf the child constructor
// for a list field, where all input keys are parameters to the constructor.
func generateChildConstructorsForList(methodBuf *strings.Builder, keys map[string]*ygen.ListKey, keyNames []string, fieldData goPathFieldData, isUnderFakeRoot, generateWildcardPaths bool, schemaStructPkgAccessor string) []error {
	var errors []error
	// List of function parameters as would appear in the method definition.
	keyParams, err := makeKeyParams(keys, keyNames, schemaStructPkgAccessor)
	if err != nil {
		return append(errors, err)
	}

	parentTypeName := fieldData.Struct.TypeName
	wildcardParentTypeName := parentTypeName + WildcardSuffix
	fieldTypeName := fieldData.TypeName
	wildcardFieldTypeName := fieldTypeName + WildcardSuffix

	var paramListStrs, paramDocStrs, keyEntryStrs []string
	for _, param := range keyParams {
		paramDocStrs = append(paramDocStrs, param.varName+": "+param.typeDocString)
		paramListStrs = append(paramListStrs, fmt.Sprintf("%s %s", param.varName, param.typeName))
		keyEntryStrs = append(keyEntryStrs, fmt.Sprintf(`"%s": %s`, param.name, param.varName))
	}
	// Create the string for the method parameter list, docstrings, and ygnmi.NodePath's key list.
	fieldData.KeyParamListStr = strings.Join(paramListStrs, ", ")
	fieldData.KeyParamDocStrs = paramDocStrs
	fieldData.KeyEntriesStr = strings.Join(keyEntryStrs, ", ")

	if err := goPathChildConstructorTemplate.Execute(methodBuf, fieldData); err != nil {
		errors = append(errors, err)
	}

	// The root node doesn't have a wildcard version of itself.
	if isUnderFakeRoot {
		return errors
	}

	if generateWildcardPaths {
		// Generate child constructor method for wildcard version of parent struct.
		fieldData.Struct.TypeName = wildcardParentTypeName
		// Override the corner case for generating the non-wildcard child.
		fieldData.TypeName = wildcardFieldTypeName
		fieldData.IsWildcard = true
		if err := goPathChildConstructorTemplate.Execute(methodBuf, fieldData); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// getFieldTypeName returns the type name for a field node of a directory -
// handling the case where the field supplied is a leaf or directory. The input
// directories is a map from paths to directory entries, and goFieldName is the
// incremental type name to be used for the case that the directory field is a
// leaf. For non-leaves, their corresponding directories' "Name"s, which are the
// same names as their corresponding ygen Go struct type names, are re-used as
// their type names; for leaves, type names are synthesized.
func getFieldTypeName(directory *ygen.ParsedDirectory, directoryFieldName string, goFieldName string, directories map[string]*ygen.ParsedDirectory, pathStructSuffix string) (string, error) {
	field, ok := directory.Fields[directoryFieldName]
	if !ok {
		return "", fmt.Errorf("getFieldTypeName: field %s not found in directory %v", directoryFieldName, directory)
	}

	switch field.Type {
	case ygen.ContainerNode, ygen.ListNode:
		fieldDirectory, ok := directories[field.YANGDetails.Path]
		if !ok {
			return "", fmt.Errorf("getFieldTypeName: unexpected - field with path %q not found in parsed yang structs map: %v", field.YANGDetails.Path, directories)
		}
		return fieldDirectory.Name + pathStructSuffix, nil
	// Leaves do not have corresponding Directory entries, so their names need to be constructed.
	default:
		if directory.IsFakeRoot {
			// When a leaf resides at the root, its type name is its whole name -- we never want fakeroot's name as a prefix.
			return goFieldName + pathStructSuffix, nil
		}
		return directory.Name + "_" + goFieldName + pathStructSuffix, nil
	}
}

type keyParam struct {
	name          string
	varName       string
	typeName      string
	typeDocString string
}

// makeKeyParams generates the list of go parameter list components for a child
// list's constructor method given the list's key spec, as well as a
// list of each parameter's types as a comment string.
// It outputs the parameters in the same order as in the given keyNames.
// e.g.
//
//	  in: &map[string]*ygen.ListKey{
//			"fluorine": &ygen.ListKey{
//				Name: "Fluorine", LangType: &ygen.MappedType{NativeType: "string"}
//			},
//			"iodine-liquid": &ygen.ListKey{
//				Name: "IodineLiquid", LangType: &ygen.MappedType{NativeType: "A_Union", UnionTypes: {"Binary": 0, "uint64": 1}}
//			},
//	      }
//	      KeyNames: []string{"fluorine", "iodine-liquid"},
//
//	  {name, varName, typeName} out: [{"fluroine", "Fluorine", "string"}, {"iodine-liquid", "IodineLiquid", "oc.A_Union"}]
//	  docstring out: ["string", "[oc.Binary, oc.UnionUint64]"]
func makeKeyParams(keys map[string]*ygen.ListKey, keyNames []string, schemaStructPkgAccessor string) ([]keyParam, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	// Create parameter list *in order* of keys, which should be in schema order.
	var keyParams []keyParam
	// NOTE: Although the generated key names might not match their
	// corresponding ygen field names in case of a camelcase name
	// collision, we expect that the user is aware of the schema to know
	// the order of the keys, and not rely on the naming in that case.
	goKeyNameMap, err := getGoKeyNameMap(keys, keyNames)
	if err != nil {
		return nil, err
	}
	for _, keyName := range keyNames {
		listKey, ok := keys[keyName]
		switch {
		case !ok:
			return nil, fmt.Errorf("makeKeyParams: inconsistent IR: key doesn't exist: %s", keyName)
		case listKey.LangType == nil:
			return nil, fmt.Errorf("makeKeyParams: mappedType for key is nil: %s", keyName)
		}

		var typeName string
		switch {
		case listKey.LangType.NativeType == "interface{}": // ygen-unsupported types
			typeName = "string"
		case ygen.IsYgenDefinedGoType(listKey.LangType):
			typeName = schemaStructPkgAccessor + listKey.LangType.NativeType
		default:
			typeName = listKey.LangType.NativeType
		}
		varName := goKeyNameMap[keyName]

		typeDocString := typeName
		if len(listKey.LangType.UnionTypes) > 1 {
			var genTypes []string
			for _, name := range listKey.LangType.OrderedUnionTypes() {
				unionTypeName := name
				if simpleName, ok := ygot.SimpleUnionBuiltinGoTypes[name]; ok {
					unionTypeName = simpleName
				}
				// Add schemaStructPkgAccessor.
				if strings.HasPrefix(unionTypeName, "*") {
					unionTypeName = "*" + schemaStructPkgAccessor + unionTypeName[1:]
				} else {
					unionTypeName = schemaStructPkgAccessor + unionTypeName
				}
				genTypes = append(genTypes, unionTypeName)
			}
			// Create the subtype documentation string.
			typeDocString = "[" + strings.Join(genTypes, ", ") + "]"
		}

		keyParams = append(keyParams, keyParam{
			name:          keyName,
			varName:       varName,
			typeName:      typeName,
			typeDocString: typeDocString,
		})
	}
	return keyParams, nil
}

// getGoKeyNameMap returns a map of Go key names keyed by their schema names
// given a list of key entries and their order. Names are camelcased and
// uniquified to ensure compilation. Uniqification is done deterministically.
func getGoKeyNameMap(keys map[string]*ygen.ListKey, keyNames []string) (map[string]string, error) {
	goKeyNameMap := make(map[string]string, len(keyNames))

	usedKeyNames := map[string]bool{}
	for _, name := range keyNames {
		key, ok := keys[name]
		if !ok {
			return nil, fmt.Errorf("key %q doesn't exist in key map %v", name, keys)
		}
		goKeyNameMap[name] = genutil.MakeNameUnique(key.Name, usedKeyNames)
	}
	return goKeyNameMap, nil
}

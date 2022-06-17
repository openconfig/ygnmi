// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package generator contains the command for generating Go code from YANG modules.
package generator

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygnmi/pathgen"
	"github.com/openconfig/ygot/genutil"
	"github.com/openconfig/ygot/gogen"
	"github.com/openconfig/ygot/ygen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// New returns a new generator command.
//nolint:errcheck
func New() *cobra.Command {
	generator := &cobra.Command{
		Use:   "generator",
		RunE:  generate,
		Short: "Generates Go code for gNMI from a YANG schema.",
		Args:  cobra.MinimumNArgs(1),
	}

	generator.Flags().String("schema_struct_path", "", "The Go import path for the schema structs package. If struct generation is enabled, this defaults to base_package_path.")
	generator.Flags().String("ygot_path", "github.com/openconfig/ygot/ygot", "The import path to use for ygot.")
	generator.Flags().String("ygnmi_path", "github.com/openconfig/ygnmi/ygnmi", "The import path to use for ygnmi.")
	generator.Flags().String("ytypes_path", "github.com/openconfig/ygot/ytypes", "The import path to use for ytypes.")
	generator.Flags().String("goyang_path", "github.com/openconfig/goyang/pkg/yang", "The import path to use for goyang.")
	generator.Flags().String("base_package_path", "", "This needs to be set to the import path of the output_dir.")
	generator.Flags().String("trim_prefix", "", "A prefix (if any) to trim from generated package names and enums")
	generator.Flags().StringSlice("paths", nil, "Comma-separated list of paths to be recursively searched for included modules or submodules within the defined YANG modules.")
	generator.Flags().StringSlice("exclude_modules", nil, "Comma-separated YANG modules to exclude from code generation.")
	generator.Flags().String("output_dir", "", "The directory that the generated Go code should be written to. This directory is the base of the generated module packages. default (working dir)")
	generator.Flags().Int("structs_split_files_count", 1, "The number of files to split the generated schema structs into.")
	generator.Flags().Int("pathstructs_split_files_count", 1, "The number of files to split the generated path structs into.")

	generator.Flags().MarkHidden("schema_struct_path")

	return generator
}

// generate runs the ygnmi PathStruct and optionally the ygot GoStruct generation.
func generate(cmd *cobra.Command, args []string) error {
	if viper.Get("base_package_path") == "" {
		return fmt.Errorf("base_package_path must be set")
	}

	schemaStructPath := viper.GetString("schema_struct_path")
	if schemaStructPath == "" {
		schemaStructPath = viper.GetString("base_package_path")
	}

	version := "ygnmi version: " + cmd.Root().Version

	importPath := strings.Split(viper.GetString("base_package_path"), "/")
	rootPackageName := fmt.Sprintf("%spath", importPath[len(importPath)-1])

	pcg := pathgen.GenConfig{
		PackageName: rootPackageName,
		GoImports: pathgen.GoImports{
			SchemaStructPkgPath: schemaStructPath,
			YgotImportPath:      viper.GetString("ygot_path"),
			YgnmiImportPath:     viper.GetString("ygnmi_path"),
			YtypesImportPath:    viper.GetString("ytypes_path"),
		},
		PreferOperationalState:               true,
		ExcludeState:                         false,
		SkipEnumDeduplication:                false,
		ShortenEnumLeafNames:                 true,
		EnumOrgPrefixesToTrim:                []string{viper.GetString("trim_prefix")},
		UseDefiningModuleForTypedefEnumNames: false,
		AppendEnumSuffixForSimpleUnionEnums:  true,
		FakeRootName:                         "root",
		PathStructSuffix:                     "Path",
		ExcludeModules:                       viper.GetStringSlice("exclude_modules"),
		YANGParseOptions: yang.Options{
			IgnoreSubmoduleCircularDependencies: false,
		},
		GeneratingBinary:      version,
		GenerateWildcardPaths: true,
		TrimPackagePrefix:     viper.GetString("trim_prefix"),
		SplitByModule:         true,
		BasePackagePath:       viper.GetString("base_package_path"),
		PackageSuffix:         "",
		UnifyPathStructs:      true,
		ExtraGenerators:       []pathgen.Generator{pathgen.GNMIGenerator},
	}

	pathCode, _, errs := pcg.GeneratePathCode(args, viper.GetStringSlice("paths"))
	if errs != nil {
		return errs
	}

	for packageName, code := range pathCode {
		if packageName == rootPackageName {
			path := filepath.Join(viper.GetString("output_dir"), packageName, fmt.Sprintf("%s.go", packageName))
			if err := os.MkdirAll(filepath.Join(viper.GetString("output_dir"), packageName), 0755); err != nil {
				return fmt.Errorf("failed to create directory for package %q: %w", packageName, err)
			}
			if err := ioutil.WriteFile(path, []byte(code.String()), 0644); err != nil {
				return err
			}
			continue
		}
		files, err := code.SplitFiles(viper.GetInt("pathstructs_split_files_count"))
		if err != nil {
			return err
		}
		outFiles := map[string]string{}
		for i, f := range files {
			outFiles[fmt.Sprintf("%s-%d.go", packageName, i)] = f
		}
		if err := writeFiles(filepath.Join(viper.GetString("output_dir"), packageName), outFiles); err != nil {
			return err
		}
	}

	return generateStructs(args, schemaStructPath, version)
}

func generateStructs(modules []string, schemaPath, version string) error {
	cmp, err := genutil.TranslateToCompressBehaviour(true, false, true)
	if err != nil {
		return err
	}

	// Perform the code generation.
	cg := gogen.New(
		version,
		ygen.IROptions{
			ParseOptions: ygen.ParseOpts{
				ExcludeModules: viper.GetStringSlice("exclude_modules"),
				YANGParseOptions: yang.Options{
					IgnoreSubmoduleCircularDependencies: false,
				},
			},
			TransformationOptions: ygen.TransformationOpts{
				CompressBehaviour:                    cmp,
				SkipEnumDeduplication:                false,
				GenerateFakeRoot:                     true,
				FakeRootName:                         "root",
				ShortenEnumLeafNames:                 true,
				EnumOrgPrefixesToTrim:                []string{viper.GetString("trim_prefix")},
				UseDefiningModuleForTypedefEnumNames: false,
				EnumerationsUseUnderscores:           true,
			},
		},
		gogen.GoOpts{
			PackageName:                         path.Base(schemaPath),
			GenerateJSONSchema:                  true,
			IncludeDescriptions:                 false,
			IgnoreShadowSchemaPaths:             true,
			YgotImportPath:                      viper.GetString("ygot_path"),
			YtypesImportPath:                    viper.GetString("ytypes_path"),
			GoyangImportPath:                    viper.GetString("goyang_path"),
			GenerateRenameMethod:                false,
			AddAnnotationFields:                 false,
			AnnotationPrefix:                    "Λ",
			AddYangPresence:                     false,
			GenerateGetters:                     true,
			GenerateDeleteMethod:                true,
			GenerateAppendMethod:                true,
			GenerateLeafGetters:                 true,
			GeneratePopulateDefault:             true,
			ValidateFunctionName:                "Validate",
			GenerateSimpleUnions:                true,
			IncludeModelData:                    false,
			AppendEnumSuffixForSimpleUnionEnums: true,
		},
	)
	generatedGoCode, errs := cg.Generate(modules, viper.GetStringSlice("paths"))
	if errs != nil {
		return fmt.Errorf("error generating GoStruct Code: %v", errs)
	}
	out, err := splitCodeByFileN(generatedGoCode, viper.GetInt("structs_split_files_count"))
	if err != nil {
		return fmt.Errorf("error splitting GoStruct Code: %w", err)
	}
	if err := writeFiles(viper.GetString("output_dir"), out); err != nil {
		return fmt.Errorf("error while writing schema struct files: %w", err)
	}
	return nil
}

const (
	// enumMapFn is the filename to be used for the enum map when Go code is output to a directory.
	enumMapFn = "enum_map.go"
	// enumFn is the filename to be used for the enum code when Go code is output to a directory.
	enumFn = "enum.go"
	// schemaFn is the filename to be used for the schema code when outputting to a directory.
	schemaFn = "schema.go"
	// interfaceFn is the filename to be used for interface code when outputting to a directory.
	interfaceFn = "union.go"
	// structsFileFmt is the format string filename (missing index) to be
	// used for files containing structs when outputting to a directory.
	structsFileFmt = "structs-%d.go"
)

// writeFiles creates or truncates files in a given base directory and writes
// to them. Keys of the contents map are file names, and values are the
// contents to be written. An error is returned if the base directory does not
// exist. If a file cannot be written, the function aborts with the error,
// leaving an unspecified set of the other input files written with their given
// contents.
func writeFiles(dir string, out map[string]string) error {
	for filename, contents := range out {
		if len(contents) == 0 {
			continue
		}
		if dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}
		fh := genutil.OpenFile(filepath.Join(dir, filename))
		if fh == nil {
			return fmt.Errorf("could not open file %q", filename)
		}
		if _, err := fh.WriteString(contents); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		// flush & close written files before function finishes.
		defer genutil.SyncFile(fh)
	}

	return nil
}

func splitCodeByFileN(goCode *gogen.GeneratedCode, fileN int) (map[string]string, error) {
	structN := len(goCode.Structs)
	if fileN < 1 || fileN > structN {
		return nil, fmt.Errorf("requested %d files, but must be between 1 and %d (number of schema structs)", fileN, structN)
	}

	out := map[string]string{
		schemaFn: goCode.JSONSchemaCode,
		enumFn:   strings.Join(goCode.Enums, "\n"),
	}

	var structFiles []string
	var code, interfaceCode strings.Builder
	structsPerFile := int(math.Ceil(float64(structN) / float64(fileN)))
	// Empty files could appear with certain structN/fileN combinations due
	// to the ceiling numbers being used for structsPerFile.
	// e.g. 4/3 gives two files of two structs.
	// This is a little more complex, but spreads out the structs more evenly.
	// If we instead use the floor number, and put all remainder structs in
	// the last file, we might double the last file's number of structs if we get unlucky.
	// e.g. 99/10 assigns 18 structs to the last file.
	emptyFiles := fileN - int(math.Ceil(float64(structN)/float64(structsPerFile)))
	code.WriteString(goCode.OneOffHeader)
	for i, s := range goCode.Structs {
		code.WriteString(s.StructDef)
		code.WriteString(s.ListKeys)
		code.WriteString("\n")
		code.WriteString(s.Methods)
		if s.Methods != "" {
			code.WriteString("\n")
		}
		interfaceCode.WriteString(s.Interfaces)
		if s.Interfaces != "" {
			interfaceCode.WriteString("\n")
		}
		// The last file contains the remainder of the structs.
		if i == structN-1 || (i+1)%structsPerFile == 0 {
			structFiles = append(structFiles, code.String())
			code.Reset()
		}
	}
	for i := 0; i != emptyFiles; i++ {
		structFiles = append(structFiles, "")
	}

	for i, structFile := range structFiles {
		out[fmt.Sprintf(structsFileFmt, i)] = structFile
	}

	code.Reset()
	code.WriteString(goCode.EnumMap)
	if code.Len() != 0 {
		code.WriteString("\n")
	}
	code.WriteString(goCode.EnumTypeMap)

	out[enumMapFn] = code.String()
	out[interfaceFn] = interfaceCode.String()

	for name, code := range out {
		out[name] = goCode.CommonHeader + code
	}

	return out, nil
}

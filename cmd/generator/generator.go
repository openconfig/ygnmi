package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygnmi/pathgen"
	"github.com/openconfig/ygot/genutil"
	"github.com/spf13/cobra"
)

var (
	schemaStructPath string
	ygotImportPath   string
	ygnmiImportPath  string
	ytypesImportPath string
	baseImportPath   string
	paths            []string
	outputDir        string
)

// New returns a new generator command.
func New() *cobra.Command {
	generator := &cobra.Command{
		Use:   "generator",
		RunE:  generate,
		Short: "Generates Go code for gNMI from a YANG schema.",
	}

	generator.Flags().StringVar(&schemaStructPath, "schema_struct_path", "", "The Go import path for the schema structs package.")
	generator.Flags().StringVar(&ygotImportPath, "ygot_path", "github.com/openconfig/ygot/ygot", "The import path to use for ygot.")
	generator.Flags().StringVar(&ygnmiImportPath, "ygnmi_path", "github.com/openconfig/ygnmi/ygnmi", "The import path to use for ygot.")
	generator.Flags().StringVar(&ytypesImportPath, "ytypes_path", "github.com/openconfig/ygot/ytypes", "The import path to use for ytypes.")
	generator.Flags().StringVar(&baseImportPath, "base_import_path", "", "Base import path used to concatenate with module package relative paths for path struct imports.")
	generator.Flags().StringSliceVar(&paths, "path", nil, "Comma separated list of paths to be recursively searched for included modules or submodules within the defined YANG modules.")
	generator.Flags().StringVar(&outputDir, "output_dir", "", "The directory that the generated Go code should be written to. This directory is the base of the generated module packages.")

	return generator
}

func generate(cmd *cobra.Command, args []string) error {
	// Perform the code generation.
	pcg := pathgen.GenConfig{
		PackageName: "device",
		GoImports: pathgen.GoImports{
			SchemaStructPkgPath: schemaStructPath,
			YgotImportPath:      ygotImportPath,
			YgnmiImportPath:     ygnmiImportPath,
			YtypesImportPath:    ytypesImportPath,
		},
		PreferOperationalState:               true,
		ExcludeState:                         false,
		SkipEnumDeduplication:                false,
		ShortenEnumLeafNames:                 true,
		EnumOrgPrefixesToTrim:                []string{"openconfig"},
		UseDefiningModuleForTypedefEnumNames: true,
		AppendEnumSuffixForSimpleUnionEnums:  true,
		FakeRootName:                         "device",
		PathStructSuffix:                     "",
		ExcludeModules:                       nil,
		YANGParseOptions: yang.Options{
			IgnoreSubmoduleCircularDependencies: false,
		},
		GeneratingBinary:        genutil.CallerName(),
		ListBuilderKeyThreshold: 2,
		GenerateWildcardPaths:   true,
		SimplifyWildcardPaths:   false,
		TrimOCPackage:           true,
		SplitByModule:           true,
		BaseImportPath:          baseImportPath,
		PackageSuffix:           "",
		UnifyPathStructs:        true,
		ExtraGenerators:         []pathgen.Generator{pathgen.GNMIGenerator},
	}

	pathCode, _, errs := pcg.GeneratePathCode(args, paths)
	if errs != nil {
		return errs
	}

	for packageName, code := range pathCode {
		path := filepath.Join(outputDir, "device.go")
		if packageName != pcg.PackageName {
			if err := os.MkdirAll(filepath.Join(outputDir, packageName), 0755); err != nil {
				return fmt.Errorf("failed to create directory for package %q: %w", packageName, err)
			}
			path = filepath.Join(outputDir, packageName, fmt.Sprintf("%s.go", packageName))
		}
		err := ioutil.WriteFile(path, []byte(code.String()), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

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

package pathgen

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygot/genutil"
	"github.com/openconfig/ygot/testutil"
	"github.com/openconfig/ygot/ygen"
	"github.com/openconfig/ygot/ygot"
)

const (
	// TestRoot is the root of the test directory such that this is not
	// repeated when referencing files.
	TestRoot string = ""
	// deflakeRuns specifies the number of runs of code generation that
	// should be performed to check for flakes.
	deflakeRuns int = 10
	// datapath is the path to common YANG test modules.
	datapath = "testdata/yang/"
)

var updateGolden = flag.Bool("update_golden", false, "Update golden files")

func TestGeneratePathCode(t *testing.T) {
	tests := []struct {
		// Name is the identifier for the test.
		name string
		// inFiles is the set of inputFiles for the test.
		inFiles []string
		// inIncludePaths is the set of paths that should be searched for imports.
		inIncludePaths      []string
		inCompressBehaviour genutil.CompressBehaviour
		// inShortenEnumLeafNames says whether the enum leaf names are shortened (i.e. module name removed) in the generated Go code corresponding to the generated path library.
		inShortenEnumLeafNames bool
		// inUseDefiningModuleForTypedefEnumNames uses the defining module name to prefix typedef enumerated types instead of the module where the typedef enumerated value is used.
		inUseDefiningModuleForTypedefEnumNames bool
		// inGenerateWildcardPaths determines whether wildcard paths are generated.
		inGenerateWildcardPaths bool
		inSchemaStructPkgPath   string
		inPathStructSuffix      string
		inIgnoreUnsupported     bool
		inIgnoreAtomicLists     bool
		// inPathOriginName sets an arbitrary name to the path origin and inUseModuleNameAsPathOrigin uses the YANG module name there.
		inPathOriginName            string
		inUseModuleNameAsPathOrigin bool
		// checkYANGPath says whether to check for the YANG path in the NodeDataMap.
		checkYANGPath bool
		// wantStructsCodeFile is the path of the generated Go code that the output of the test should be compared to.
		wantStructsCodeFile string
		// wantNodeDataMap is the expected NodeDataMap to be produced to accompany the path struct outputs.
		wantNodeDataMap NodeDataMap
		// wantErr specifies whether the test should expect an error.
		wantErr bool
	}{{
		name:                    "simple openconfig test",
		inFiles:                 []string{filepath.Join(datapath, "openconfig-simple.yang")},
		wantStructsCodeFile:     filepath.Join(TestRoot, "testdata/structs/openconfig-simple.path-txt"),
		inCompressBehaviour:     genutil.PreferOperationalState,
		inShortenEnumLeafNames:  true,
		inGenerateWildcardPaths: true,
		inSchemaStructPkgPath:   "",
		inPathStructSuffix:      "Path",
		checkYANGPath:           true,
		wantNodeDataMap: NodeDataMap{
			"DevicePath": {
				GoTypeName:            "*Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"ParentPath": {
				GoTypeName:            "*Parent",
				LocalGoTypeName:       "*Parent",
				GoFieldName:           "Parent",
				SubsumingGoStructName: "Parent",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				YANGPath:              "/openconfig-simple/parent",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Parent_ChildPath": {
				GoTypeName:            "*Parent_Child",
				LocalGoTypeName:       "*Parent_Child",
				GoFieldName:           "Child",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				YANGPath:              "/openconfig-simple/parent/child",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_FourPath": {
				GoTypeName:            "Binary",
				LocalGoTypeName:       "Binary",
				GoFieldName:           "Four",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "binary",
				YANGPath:              "/openconfig-simple/parent/child/state/four",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "four",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Parent_Child_OnePath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "One",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				YANGPath:              "/openconfig-simple/parent/child/state/one",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "one",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Parent_Child_ThreePath": {
				GoTypeName:            "E_Child_Three",
				LocalGoTypeName:       "E_Child_Three",
				GoFieldName:           "Three",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "enumeration",
				YANGPath:              "/openconfig-simple/parent/child/state/three",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "three",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Parent_Child_TwoPath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "Two",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				YANGPath:              "/openconfig-simple/parent/child/state/two",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "two",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Parent_Child_FivePath": {
				GoTypeName:            "Binary",
				LocalGoTypeName:       "Binary",
				GoFieldName:           "Five",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				YANGTypeName:          "ieeefloat32",
				YANGPath:              "/openconfig-simple/parent/child/state/five",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "five",
				DirectoryName:         "/openconfig-simple/parent/child",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Parent_Child_SixPath": {
				GoTypeName:            "[]Binary",
				LocalGoTypeName:       "[]Binary",
				GoFieldName:           "Six",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				YANGTypeName:          "ieeefloat32",
				YANGPath:              "/openconfig-simple/parent/child/state/six",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "six",
				DirectoryName:         "/openconfig-simple/parent/child",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"RemoteContainerPath": {
				GoTypeName:            "*RemoteContainer",
				LocalGoTypeName:       "*RemoteContainer",
				GoFieldName:           "RemoteContainer",
				SubsumingGoStructName: "RemoteContainer",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				YANGPath:              "/openconfig-simple/remote-container",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/remote-container",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"RemoteContainer_ALeafPath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "ALeaf",
				SubsumingGoStructName: "RemoteContainer",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				YANGPath:              "/openconfig-simple/remote-container/state/a-leaf",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/remote-container",
				YANGFieldName:         "a-leaf",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			}},
	}, {
		name:                    "simple openconfig test with compression off",
		inFiles:                 []string{filepath.Join(datapath, "openconfig-simple.yang")},
		wantStructsCodeFile:     filepath.Join(TestRoot, "testdata/structs/openconfig-simple.uncompressed.path-txt"),
		inCompressBehaviour:     genutil.Uncompressed,
		inShortenEnumLeafNames:  true,
		inGenerateWildcardPaths: true,
		inSchemaStructPkgPath:   "",
		inPathStructSuffix:      "Path",
	}, {
		name:                                   "simple openconfig test with preferOperationalState=false",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-simple.yang")},
		inCompressBehaviour:                    genutil.PreferIntendedConfig,
		inShortenEnumLeafNames:                 true,
		inGenerateWildcardPaths:                true,
		inUseDefiningModuleForTypedefEnumNames: true,
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-simple.intendedconfig.path-txt"),
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		wantNodeDataMap: NodeDataMap{
			"DevicePath": {
				GoTypeName:            "*Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"ParentPath": {
				GoTypeName:            "*Parent",
				LocalGoTypeName:       "*Parent",
				GoFieldName:           "Parent",
				SubsumingGoStructName: "Parent",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Parent_ChildPath": {
				GoTypeName:            "*Parent_Child",
				LocalGoTypeName:       "*Parent_Child",
				GoFieldName:           "Child",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_FourPath": {
				GoTypeName:            "Binary",
				LocalGoTypeName:       "Binary",
				GoFieldName:           "Four",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "binary",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "four",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_OnePath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "One",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "one",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_ThreePath": {
				GoTypeName:            "E_Child_Three",
				LocalGoTypeName:       "E_Child_Three",
				GoFieldName:           "Three",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "enumeration",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "three",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_TwoPath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "Two",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "two",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Parent_Child_FivePath": {
				GoTypeName:            "Binary",
				LocalGoTypeName:       "Binary",
				GoFieldName:           "Five",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				YANGTypeName:          "ieeefloat32",
				YANGPath:              "/openconfig-simple/parent/child/config/five",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "five",
				DirectoryName:         "/openconfig-simple/parent/child",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_SixPath": {
				GoTypeName:            "[]Binary",
				LocalGoTypeName:       "[]Binary",
				GoFieldName:           "Six",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				YANGTypeName:          "ieeefloat32",
				YANGPath:              "/openconfig-simple/parent/child/config/six",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "six",
				DirectoryName:         "/openconfig-simple/parent/child",
				PathOriginName:        "openconfig",
			},
			"RemoteContainerPath": {
				GoTypeName:            "*RemoteContainer",
				LocalGoTypeName:       "*RemoteContainer",
				GoFieldName:           "RemoteContainer",
				SubsumingGoStructName: "RemoteContainer",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/remote-container",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"RemoteContainer_ALeafPath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "ALeaf",
				SubsumingGoStructName: "RemoteContainer",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/remote-container",
				YANGFieldName:         "a-leaf",
				PathOriginName:        "openconfig",
			}},
	}, {
		name:                                   "simple openconfig test with excludeState=true",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-simple.yang")},
		inCompressBehaviour:                    genutil.ExcludeDerivedState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-simple.excludestate.path-txt"),
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		wantNodeDataMap: NodeDataMap{
			"DevicePath": {
				GoTypeName:            "*Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"ParentPath": {
				GoTypeName:            "*Parent",
				LocalGoTypeName:       "*Parent",
				GoFieldName:           "Parent",
				SubsumingGoStructName: "Parent",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Parent_ChildPath": {
				GoTypeName:            "*Parent_Child",
				LocalGoTypeName:       "*Parent_Child",
				GoFieldName:           "Child",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_FourPath": {
				GoTypeName:            "Binary",
				LocalGoTypeName:       "Binary",
				GoFieldName:           "Four",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "binary",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "four",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_OnePath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "One",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "one",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_ThreePath": {
				GoTypeName:            "E_Child_Three",
				LocalGoTypeName:       "E_Child_Three",
				GoFieldName:           "Three",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "enumeration",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/parent/child",
				YANGFieldName:         "three",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_FivePath": {
				GoTypeName:            "Binary",
				LocalGoTypeName:       "Binary",
				GoFieldName:           "Five",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				YANGTypeName:          "ieeefloat32",
				YANGPath:              "/openconfig-simple/parent/child/config/five",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "five",
				DirectoryName:         "/openconfig-simple/parent/child",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_SixPath": {
				GoTypeName:            "[]Binary",
				LocalGoTypeName:       "[]Binary",
				GoFieldName:           "Six",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				YANGTypeName:          "ieeefloat32",
				YANGPath:              "/openconfig-simple/parent/child/config/six",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "six",
				DirectoryName:         "/openconfig-simple/parent/child",
				PathOriginName:        "openconfig",
			},
			"RemoteContainerPath": {
				GoTypeName:            "*RemoteContainer",
				LocalGoTypeName:       "*RemoteContainer",
				GoFieldName:           "RemoteContainer",
				SubsumingGoStructName: "RemoteContainer",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/remote-container",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"RemoteContainer_ALeafPath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "ALeaf",
				SubsumingGoStructName: "RemoteContainer",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple/remote-container",
				YANGFieldName:         "a-leaf",
				PathOriginName:        "openconfig",
			}},
	}, {
		name:                                   "simple openconfig test with ordered list",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-orderedlist.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		checkYANGPath:                          true,
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-orderedlist.path-txt"),
		wantNodeDataMap: NodeDataMap{
			"DevicePath": {
				GoTypeName:            "*Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"ModelPath": {
				GoTypeName:            "*Model",
				LocalGoTypeName:       "*Model",
				GoFieldName:           "Model",
				SubsumingGoStructName: "Model",
				YANGPath:              "/openconfig-orderedlist/model",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-orderedlist/model",
				PathOriginName:        "openconfig",
			},
			"Model_OrderedListPathMap": {
				GoTypeName:            "*Model_OrderedList_OrderedMap",
				LocalGoTypeName:       "*Model_OrderedList_OrderedMap",
				GoFieldName:           "OrderedList",
				SubsumingGoStructName: "Model",
				IsListContainer:       true,
				YANGPath:              "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				CompressInfo: &CompressionInfo{
					PreRelPathList:  `"openconfig-orderedlist:ordered-lists"`,
					PostRelPathList: `"openconfig-orderedlist:ordered-list"`,
				},
				PathOriginName: "openconfig",
			},
			"Model_OrderedList_KeyPath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "Key",
				SubsumingGoStructName: "Model_OrderedList",
				IsLeaf:                true,
				IsScalarField:         true,
				YANGTypeName:          "string",
				YANGPath:              "/openconfig-orderedlist/model/ordered-lists/ordered-list/state/key",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "key",
				DirectoryName:         "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			}},
	}, {
		name:                                   "simple openconfig test with ordered list but atomic generation turned off",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-orderedlist.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		inIgnoreAtomicLists:                    true,
		checkYANGPath:                          true,
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-orderedlist-unordered.path-txt"),
		wantNodeDataMap: NodeDataMap{
			"DevicePath": {
				GoTypeName:            "*Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"ModelPath": {
				GoTypeName:            "*Model",
				LocalGoTypeName:       "*Model",
				GoFieldName:           "Model",
				SubsumingGoStructName: "Model",
				YANGPath:              "/openconfig-orderedlist/model",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-orderedlist/model",
				PathOriginName:        "openconfig",
			},
			"Model_OrderedListPath": {
				GoTypeName:            "*Model_OrderedList",
				LocalGoTypeName:       "*Model_OrderedList",
				GoFieldName:           "OrderedList",
				SubsumingGoStructName: "Model_OrderedList",
				YANGPath:              "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				PathOriginName:        "openconfig",
			},
			"Model_OrderedList_KeyPath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "Key",
				SubsumingGoStructName: "Model_OrderedList",
				IsLeaf:                true,
				IsScalarField:         true,
				YANGTypeName:          "string",
				YANGPath:              "/openconfig-orderedlist/model/ordered-lists/ordered-list/state/key",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "key",
				DirectoryName:         "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			}},
	}, {
		name:                                   "simple openconfig test with list",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-withlist.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-withlist.path-txt"),
	}, {
		name:                                   "simple openconfig test with list without wildcard paths",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-withlist.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                false,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-withlist.nowildcard.path-txt"),
	}, {
		name:                                   "simple openconfig test with list in separate package",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-withlist.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "github.com/openconfig/ygot/ypathgen/testdata/exampleoc",
		inPathStructSuffix:                     "",
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-withlist-separate-package.path-txt"),
	}, {
		name:                                   "simple openconfig test with list in builder API",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-withlist.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-withlist.builder.path-txt"),
	}, {
		name:                                   "simple openconfig test with union & typedef & identity & enum",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-unione.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-unione.path-txt"),
		wantNodeDataMap: NodeDataMap{
			"DevicePath": {
				GoTypeName:            "*Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"DupEnumPath": {
				GoTypeName:            "*DupEnum",
				LocalGoTypeName:       "*DupEnum",
				GoFieldName:           "DupEnum",
				SubsumingGoStructName: "DupEnum",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-unione/dup-enum",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"DupEnum_APath": {
				GoTypeName:            "E_DupEnum_A",
				LocalGoTypeName:       "E_DupEnum_A",
				GoFieldName:           "A",
				SubsumingGoStructName: "DupEnum",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "enumeration",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-unione/dup-enum",
				YANGFieldName:         "A",
				PathOriginName:        "openconfig",
			},
			"DupEnum_BPath": {
				GoTypeName:            "E_DupEnum_B",
				LocalGoTypeName:       "E_DupEnum_B",
				GoFieldName:           "B",
				SubsumingGoStructName: "DupEnum",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "enumeration",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-unione/dup-enum",
				YANGFieldName:         "B",
				PathOriginName:        "openconfig",
			},
			"ComponentPath": {
				GoTypeName:            "*Component",
				LocalGoTypeName:       "*Component",
				GoFieldName:           "Component",
				SubsumingGoStructName: "Component",
				YANGPath:              "/openconfig-unione/platform/component",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			},
			"ComponentPathMap": {
				GoTypeName:            "map[string]*Component",
				LocalGoTypeName:       "map[string]*Component",
				GoFieldName:           "Component",
				SubsumingGoStructName: "Device",
				IsListContainer:       true,
				YANGPath:              "/openconfig-unione/platform/component",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-unione/platform/component",
				CompressInfo: &CompressionInfo{
					PreRelPathList:  `"openconfig-unione:platform"`,
					PostRelPathList: `"openconfig-unione:component"`,
				},
				PathOriginName: "openconfig",
			},
			"Component_E1Path": {
				GoTypeName:            "Component_E1_Union",
				LocalGoTypeName:       "Component_E1_Union",
				GoFieldName:           "E1",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				YANGTypeName:          "enumtypedef",
				YANGPath:              "/openconfig-unione/platform/component/state/e1",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "e1",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			},
			"Component_EnumeratedPath": {
				GoTypeName:            "Component_Enumerated_Union",
				LocalGoTypeName:       "Component_Enumerated_Union",
				GoFieldName:           "Enumerated",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				YANGTypeName:          "enumerated-union-type",
				YANGPath:              "/openconfig-unione/platform/component/state/enumerated",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "enumerated",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			},
			"Component_NamePath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "Name",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				IsScalarField:         true,
				YANGTypeName:          "string",
				YANGPath:              "/openconfig-unione/platform/component/state/name",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "name",
				DirectoryName:         "/openconfig-unione/platform/component",
				ConfigFalse:           false,
				PathOriginName:        "openconfig",
			},
			"Component_PowerPath": {
				GoTypeName:            "Component_Power_Union",
				LocalGoTypeName:       "Component_Power_Union",
				GoFieldName:           "Power",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				YANGTypeName:          "union",
				YANGPath:              "/openconfig-unione/platform/component/state/power",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "power",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			},
			"Component_R1Path": {
				GoTypeName:            "Component_E1_Union",
				LocalGoTypeName:       "Component_E1_Union",
				GoFieldName:           "R1",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				YANGTypeName:          "leafref",
				YANGPath:              "/openconfig-unione/platform/component/state/r1",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "r1",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			},
			"Component_TypePath": {
				GoTypeName:            "Component_Type_Union",
				LocalGoTypeName:       "Component_Type_Union",
				GoFieldName:           "Type",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				YANGTypeName:          "union",
				YANGPath:              "/openconfig-unione/platform/component/state/type",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "type",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			}},
	}, {
		name:                    "simple openconfig test with union & typedef & identity & enum, with enum names not shortened",
		inFiles:                 []string{filepath.Join(datapath, "openconfig-unione.yang")},
		inCompressBehaviour:     genutil.PreferOperationalState,
		inGenerateWildcardPaths: true,
		inSchemaStructPkgPath:   "",
		inPathStructSuffix:      "Path",
		wantStructsCodeFile:     filepath.Join(TestRoot, "testdata/structs/openconfig-unione.path-txt"),
		wantNodeDataMap: NodeDataMap{
			"DevicePath": {
				GoTypeName:            "*Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"DupEnumPath": {
				GoTypeName:            "*DupEnum",
				LocalGoTypeName:       "*DupEnum",
				GoFieldName:           "DupEnum",
				SubsumingGoStructName: "DupEnum",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-unione/dup-enum",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"DupEnum_APath": {
				GoTypeName:            "E_OpenconfigUnione_DupEnum_A",
				LocalGoTypeName:       "E_OpenconfigUnione_DupEnum_A",
				GoFieldName:           "A",
				SubsumingGoStructName: "DupEnum",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "enumeration",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-unione/dup-enum",
				YANGFieldName:         "A",
				PathOriginName:        "openconfig",
			},
			"DupEnum_BPath": {
				GoTypeName:            "E_OpenconfigUnione_DupEnum_B",
				LocalGoTypeName:       "E_OpenconfigUnione_DupEnum_B",
				GoFieldName:           "B",
				SubsumingGoStructName: "DupEnum",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "enumeration",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-unione/dup-enum",
				YANGFieldName:         "B",
				PathOriginName:        "openconfig",
			},
			"ComponentPath": {
				GoTypeName:            "*Component",
				LocalGoTypeName:       "*Component",
				GoFieldName:           "Component",
				SubsumingGoStructName: "Component",
				YANGPath:              "/openconfig-unione/platform/component",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			},
			"ComponentPathMap": {
				GoTypeName:            "map[string]*Component",
				LocalGoTypeName:       "map[string]*Component",
				GoFieldName:           "Component",
				SubsumingGoStructName: "Device",
				IsListContainer:       true,
				YANGPath:              "/openconfig-unione/platform/component",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-unione/platform/component",
				CompressInfo: &CompressionInfo{
					PreRelPathList:  `"openconfig-unione:platform"`,
					PostRelPathList: `"openconfig-unione:component"`,
				},
				PathOriginName: "openconfig",
			},
			"Component_E1Path": {
				GoTypeName:            "Component_E1_Union",
				LocalGoTypeName:       "Component_E1_Union",
				GoFieldName:           "E1",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				YANGTypeName:          "enumtypedef",
				YANGPath:              "/openconfig-unione/platform/component/state/e1",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "e1",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			},
			"Component_EnumeratedPath": {
				GoTypeName:            "Component_Enumerated_Union",
				LocalGoTypeName:       "Component_Enumerated_Union",
				GoFieldName:           "Enumerated",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				YANGTypeName:          "enumerated-union-type",
				YANGPath:              "/openconfig-unione/platform/component/state/enumerated",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "enumerated",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			},
			"Component_NamePath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "Name",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				IsScalarField:         true,
				YANGTypeName:          "string",
				YANGPath:              "/openconfig-unione/platform/component/state/name",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "name",
				DirectoryName:         "/openconfig-unione/platform/component",
				ConfigFalse:           false,
				PathOriginName:        "openconfig",
			},
			"Component_PowerPath": {
				GoTypeName:            "Component_Power_Union",
				LocalGoTypeName:       "Component_Power_Union",
				GoFieldName:           "Power",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				YANGTypeName:          "union",
				YANGPath:              "/openconfig-unione/platform/component/state/power",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "power",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			},
			"Component_R1Path": {
				GoTypeName:            "Component_E1_Union",
				LocalGoTypeName:       "Component_E1_Union",
				GoFieldName:           "R1",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				YANGTypeName:          "leafref",
				YANGPath:              "/openconfig-unione/platform/component/state/r1",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "r1",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			},
			"Component_TypePath": {
				GoTypeName:            "Component_Type_Union",
				LocalGoTypeName:       "Component_Type_Union",
				GoFieldName:           "Type",
				SubsumingGoStructName: "Component",
				IsLeaf:                true,
				YANGTypeName:          "union",
				YANGPath:              "/openconfig-unione/platform/component/state/type",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "type",
				DirectoryName:         "/openconfig-unione/platform/component",
				PathOriginName:        "openconfig",
			}},
	}, {
		name:                                   "simple openconfig test with submodule and union list key",
		inFiles:                                []string{filepath.Join(datapath, "enum-module.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/enum-module.path-txt"),
		wantNodeDataMap: NodeDataMap{
			"DevicePath": {
				GoTypeName:            "*Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"AListPath": {
				GoTypeName:            "*AList",
				LocalGoTypeName:       "*AList",
				GoFieldName:           "AList",
				SubsumingGoStructName: "AList",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/a-lists/a-list",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"AListPathMap": {
				GoTypeName:            "map[AList_Value_Union]*AList",
				LocalGoTypeName:       "map[AList_Value_Union]*AList",
				GoFieldName:           "AList",
				SubsumingGoStructName: "Device",
				IsListContainer:       true,
				YANGPath:              "/enum-module/a-lists/a-list",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/a-lists/a-list",
				CompressInfo:          &CompressionInfo{PreRelPathList: `"enum-module:a-lists"`, PostRelPathList: `"enum-module:a-list"`},
				PathOriginName:        "openconfig",
			},
			"AList_ValuePath": {
				GoTypeName:            "AList_Value_Union",
				LocalGoTypeName:       "AList_Value_Union",
				GoFieldName:           "Value",
				SubsumingGoStructName: "AList",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "td",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/a-lists/a-list",
				YANGFieldName:         "value",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"BListPath": {
				GoTypeName:            "*BList",
				LocalGoTypeName:       "*BList",
				GoFieldName:           "BList",
				SubsumingGoStructName: "BList",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/b-lists/b-list",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"BListPathMap": {
				GoTypeName:            "map[BList_Value_Union]*BList",
				LocalGoTypeName:       "map[BList_Value_Union]*BList",
				GoFieldName:           "BList",
				SubsumingGoStructName: "Device",
				IsListContainer:       true,
				YANGPath:              "/enum-module/b-lists/b-list",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/b-lists/b-list",
				CompressInfo:          &CompressionInfo{PreRelPathList: `"enum-module:b-lists"`, PostRelPathList: `"enum-module:b-list"`},
				PathOriginName:        "openconfig",
			},
			"BList_ValuePath": {
				GoTypeName:            "BList_Value_Union",
				LocalGoTypeName:       "BList_Value_Union",
				GoFieldName:           "Value",
				SubsumingGoStructName: "BList",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "td",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/b-lists/b-list",
				YANGFieldName:         "value",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"CPath": {
				GoTypeName:            "*C",
				LocalGoTypeName:       "*C",
				GoFieldName:           "C",
				SubsumingGoStructName: "C",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/c",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"C_ClPath": {
				GoTypeName:            "E_EnumModule_Cl",
				LocalGoTypeName:       "E_EnumModule_Cl",
				GoFieldName:           "Cl",
				SubsumingGoStructName: "C",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "enumeration",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/c",
				YANGFieldName:         "cl",
				ConfigFalse:           false,
				PathOriginName:        "openconfig",
			},
			"ParentPath": {
				GoTypeName:            "*Parent",
				LocalGoTypeName:       "*Parent",
				GoFieldName:           "Parent",
				SubsumingGoStructName: "Parent",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/parent",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Parent_ChildPath": {
				GoTypeName:            "*Parent_Child",
				LocalGoTypeName:       "*Parent_Child",
				GoFieldName:           "Child",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/parent/child",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Parent_Child_IdPath": {
				GoTypeName:            "E_EnumTypes_ID",
				LocalGoTypeName:       "E_EnumTypes_ID",
				GoFieldName:           "Id",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "identityref",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "id",
				DirectoryName:         "/enum-module/parent/child",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Parent_Child_Id2Path": {
				GoTypeName:            "E_EnumTypes_ID",
				LocalGoTypeName:       "E_EnumTypes_ID",
				GoFieldName:           "Id2",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            true,
				YANGTypeName:          "identityref",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/parent/child",
				YANGFieldName:         "id2",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Parent_Child_EnumPath": {
				GoTypeName:            "E_EnumTypes_TdEnum",
				LocalGoTypeName:       "E_EnumTypes_TdEnum",
				GoFieldName:           "Enum",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            true,
				YANGTypeName:          "td-enum",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/parent/child",
				YANGFieldName:         "enum",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Parent_Child_InlineEnumPath": {
				GoTypeName:            "E_Child_InlineEnum",
				LocalGoTypeName:       "E_Child_InlineEnum",
				GoFieldName:           "InlineEnum",
				SubsumingGoStructName: "Parent_Child",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            true,
				YANGTypeName:          "enumeration",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/enum-module/parent/child",
				YANGFieldName:         "inline-enum",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			}},
	}, {
		name:                                   "simple openconfig test with choice and cases",
		inFiles:                                []string{filepath.Join(datapath, "choice-case-example.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/choice-case-example.path-txt"),
	}, {
		name: "simple openconfig test with augmentations",
		inFiles: []string{
			filepath.Join(datapath, "openconfig-simple-target.yang"),
			filepath.Join(datapath, "openconfig-simple-augment.yang"),
		},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "github.com/openconfig/ygot/ypathgen/testdata/exampleoc",
		inPathStructSuffix:                     "",
		inIgnoreUnsupported:                    true,
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-augmented.path-txt"),
		wantNodeDataMap: NodeDataMap{
			"Device": {
				GoTypeName:            "*oc.Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"Native": {
				GoTypeName:            "*oc.Native",
				LocalGoTypeName:       "*Native",
				GoFieldName:           "Native",
				SubsumingGoStructName: "Native",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple-target/native",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Native_A": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "A",
				SubsumingGoStructName: "Native",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple-target/native",
				YANGFieldName:         "a",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Native_B": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "B",
				SubsumingGoStructName: "Native",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple-target/native",
				YANGFieldName:         "b",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"Target": {
				GoTypeName:            "*oc.Target",
				LocalGoTypeName:       "*Target",
				GoFieldName:           "Target",
				SubsumingGoStructName: "Target",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple-target/target",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Target_Foo": {
				GoTypeName:            "*oc.Target_Foo",
				LocalGoTypeName:       "*Target_Foo",
				GoFieldName:           "Foo",
				SubsumingGoStructName: "Target_Foo",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple-target/target/foo",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"Target_Foo_A": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "A",
				SubsumingGoStructName: "Target_Foo",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-simple-target/target/foo",
				YANGFieldName:         "a",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			}},
	}, {
		name:                                   "simple openconfig test with camelcase-name extension",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-enumcamelcase.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-enumcamelcase-compress.path-txt"),
	}, {
		name:                                   "simple openconfig test with camelcase-name extension in container and leaf",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-camelcase.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-camelcase-compress.path-txt"),
	}, {
		name:                                   "simple openconfig test with path origin name set",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-orderedlist.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		inPathOriginName:                       "test-origin",
		checkYANGPath:                          true,
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-setpathorigin.path-txt"),
		wantNodeDataMap: NodeDataMap{
			"DevicePath": {
				GoTypeName:            "*Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"ModelPath": {
				GoTypeName:            "*Model",
				LocalGoTypeName:       "*Model",
				GoFieldName:           "Model",
				SubsumingGoStructName: "Model",
				YANGPath:              "/openconfig-orderedlist/model",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-orderedlist/model",
				PathOriginName:        "test-origin",
			},
			"Model_OrderedListPathMap": {
				GoTypeName:            "*Model_OrderedList_OrderedMap",
				LocalGoTypeName:       "*Model_OrderedList_OrderedMap",
				GoFieldName:           "OrderedList",
				SubsumingGoStructName: "Model",
				IsListContainer:       true,
				YANGPath:              "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				CompressInfo: &CompressionInfo{
					PreRelPathList:  `"openconfig-orderedlist:ordered-lists"`,
					PostRelPathList: `"openconfig-orderedlist:ordered-list"`,
				},
				PathOriginName: "test-origin",
			},
			"Model_OrderedList_KeyPath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "Key",
				SubsumingGoStructName: "Model_OrderedList",
				IsLeaf:                true,
				IsScalarField:         true,
				YANGTypeName:          "string",
				YANGPath:              "/openconfig-orderedlist/model/ordered-lists/ordered-list/state/key",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "key",
				DirectoryName:         "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				ConfigFalse:           true,
				PathOriginName:        "test-origin",
			}},
	}, {
		name:                                   "simple openconfig test with path origin name as module name",
		inFiles:                                []string{filepath.Join(datapath, "openconfig-orderedlist.yang")},
		inCompressBehaviour:                    genutil.PreferOperationalState,
		inShortenEnumLeafNames:                 true,
		inUseDefiningModuleForTypedefEnumNames: true,
		inGenerateWildcardPaths:                true,
		inSchemaStructPkgPath:                  "",
		inPathStructSuffix:                     "Path",
		inUseModuleNameAsPathOrigin:            true,
		checkYANGPath:                          true,
		wantStructsCodeFile:                    filepath.Join(TestRoot, "testdata/structs/openconfig-usemodulenameaspathorigin.path-txt"),
		wantNodeDataMap: NodeDataMap{
			"DevicePath": {
				GoTypeName:            "*Device",
				LocalGoTypeName:       "*Device",
				SubsumingGoStructName: "Device",
				YANGPath:              "/",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/device",
			},
			"ModelPath": {
				GoTypeName:            "*Model",
				LocalGoTypeName:       "*Model",
				GoFieldName:           "Model",
				SubsumingGoStructName: "Model",
				YANGPath:              "/openconfig-orderedlist/model",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-orderedlist/model",
				PathOriginName:        "openconfig-orderedlist",
			},
			"Model_OrderedListPathMap": {
				GoTypeName:            "*Model_OrderedList_OrderedMap",
				LocalGoTypeName:       "*Model_OrderedList_OrderedMap",
				GoFieldName:           "OrderedList",
				SubsumingGoStructName: "Model",
				IsListContainer:       true,
				YANGPath:              "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				GoPathPackageName:     "ocstructs",
				DirectoryName:         "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				CompressInfo: &CompressionInfo{
					PreRelPathList:  `"openconfig-orderedlist:ordered-lists"`,
					PostRelPathList: `"openconfig-orderedlist:ordered-list"`,
				},
				PathOriginName: "openconfig-orderedlist",
			},
			"Model_OrderedList_KeyPath": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "Key",
				SubsumingGoStructName: "Model_OrderedList",
				IsLeaf:                true,
				IsScalarField:         true,
				YANGTypeName:          "string",
				YANGPath:              "/openconfig-orderedlist/model/ordered-lists/ordered-list/state/key",
				GoPathPackageName:     "ocstructs",
				YANGFieldName:         "key",
				DirectoryName:         "/openconfig-orderedlist/model/ordered-lists/ordered-list",
				ConfigFalse:           true,
				PathOriginName:        "openconfig-orderedlist",
			}},
	}}

	for _, tt := range tests {
		t.Run(tt.name+":"+strings.Join(tt.inFiles, ","), func(t *testing.T) {
			genCode := func() (string, NodeDataMap, *GenConfig) {
				cg := NewDefaultConfig(tt.inSchemaStructPkgPath)
				// Set the name of the caller explicitly to avoid issues when
				// the unit tests are called by external test entities.
				cg.GeneratingBinary = "pathgen-tests"
				cg.FakeRootName = "device"
				cg.PathStructSuffix = tt.inPathStructSuffix
				cg.CompressBehaviour = tt.inCompressBehaviour
				cg.ShortenEnumLeafNames = tt.inShortenEnumLeafNames
				cg.UseDefiningModuleForTypedefEnumNames = tt.inUseDefiningModuleForTypedefEnumNames
				cg.GenerateWildcardPaths = tt.inGenerateWildcardPaths
				cg.PackageName = "ocstructs"
				cg.IgnoreAtomicLists = tt.inIgnoreAtomicLists
				cg.ParseOptions.IgnoreUnsupportedStatements = tt.inIgnoreUnsupported
				cg.PathOriginName = tt.inPathOriginName
				cg.UseModuleNameAsPathOrigin = tt.inUseModuleNameAsPathOrigin

				gotCode, gotNodeDataMap, err := cg.GeneratePathCode(tt.inFiles, tt.inIncludePaths)
				if err != nil && !tt.wantErr {
					t.Fatalf("GeneratePathCode(%v, %v): Config: %v, got unexpected error: %v, want: nil", tt.inFiles, tt.inIncludePaths, cg, err)
				}

				return gotCode[cg.PackageName].String(), gotNodeDataMap, cg
			}

			gotCode, gotNodeDataMap, cg := genCode()

			if tt.wantNodeDataMap != nil {
				var cmpOpts []cmp.Option
				if !tt.checkYANGPath {
					cmpOpts = append(cmpOpts, cmpopts.IgnoreFields(NodeData{}, "YANGPath"))
				}
				if diff := cmp.Diff(tt.wantNodeDataMap, gotNodeDataMap, cmpOpts...); diff != "" {
					t.Errorf("(-wantNodeDataMap, +gotNodeDataMap):\n%s", diff)
				}
			}

			wantCodeBytes, rferr := os.ReadFile(tt.wantStructsCodeFile)
			if rferr != nil {
				t.Fatalf("os.ReadFile(%q) error: %v", tt.wantStructsCodeFile, rferr)
			}

			wantCode := string(wantCodeBytes)

			if gotCode != wantCode {
				if *updateGolden {
					if err := os.WriteFile(tt.wantStructsCodeFile, []byte(gotCode), 0644); err != nil {
						t.Fatal(err)
					}
				}
				// Use difflib to generate a unified diff between the
				// two code snippets such that this is simpler to debug
				// in the test output.
				diff, _ := testutil.GenerateUnifiedDiff(wantCode, gotCode)
				t.Errorf("GeneratePathCode(%v, %v), Config: %v, did not return correct code (file: %v), diff:\n%s",
					tt.inFiles, tt.inIncludePaths, cg, tt.wantStructsCodeFile, diff)
			}

			for i := 0; i < deflakeRuns; i++ {
				gotAttempt, _, _ := genCode()
				if gotAttempt != gotCode {
					diff, _ := testutil.GenerateUnifiedDiff(gotAttempt, gotCode)
					t.Fatalf("flaky code generation, diff:\n%s", diff)
				}
			}
		})
	}
}

func TestGeneratePathCodeSplitFiles(t *testing.T) {
	tests := []struct {
		name                  string   // Name is the identifier for the test.
		inFiles               []string // inFiles is the set of inputFiles for the test.
		inIncludePaths        []string // inIncludePaths is the set of paths that should be searched for imports.
		inFileNumber          int      // inFileNumber is the number of files into which to split the generated code.
		inSchemaStructPkgPath string
		wantStructsCodeFiles  []string // wantStructsCodeFiles is the paths of the generated Go code that the output of the test should be compared to.
		wantErr               bool     // whether an error is expected from the SplitFiles call
	}{{
		name:                  "fileNumber is higher than total number of structs",
		inFiles:               []string{filepath.Join(datapath, "openconfig-simple.yang")},
		inFileNumber:          1000,
		inSchemaStructPkgPath: "",
		wantErr:               true,
	}, {
		name:                  "fileNumber is exactly the total number of structs",
		inFiles:               []string{filepath.Join(datapath, "openconfig-simple.yang")},
		inFileNumber:          9,
		inSchemaStructPkgPath: "github.com/openconfig/ygot/ypathgen/testdata/exampleoc",
		wantStructsCodeFiles: []string{
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-90.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-91.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-92.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-93.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-94.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-95.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-96.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-97.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-98.path-txt"),
		},
	}, {
		name:                  "fileNumber is just under the total number of structs",
		inFiles:               []string{filepath.Join(datapath, "openconfig-simple.yang")},
		inFileNumber:          8,
		inSchemaStructPkgPath: "",
		wantStructsCodeFiles: []string{
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-80.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-81.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-82.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-83.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-84.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-85.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-86.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-87.path-txt"),
		},
	}, {
		name:                  "fileNumber is half the total number of structs",
		inFiles:               []string{filepath.Join(datapath, "openconfig-simple.yang")},
		inFileNumber:          4,
		inSchemaStructPkgPath: "github.com/openconfig/ygot/ypathgen/testdata/exampleoc",
		wantStructsCodeFiles: []string{
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-40.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-41.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-42.path-txt"),
			filepath.Join(TestRoot, "testdata/splitstructs/openconfig-simple-43.path-txt"),
		},
	}, {
		name:                  "single file",
		inFiles:               []string{filepath.Join(datapath, "openconfig-simple.yang")},
		inFileNumber:          1,
		inSchemaStructPkgPath: "",
		wantStructsCodeFiles:  []string{filepath.Join(TestRoot, "testdata/structs/openconfig-simple.path-txt")},
	}, {
		name:         "fileNumber is 0",
		inFiles:      []string{filepath.Join(datapath, "openconfig-simple.yang")},
		inFileNumber: 0,
		wantErr:      true,
	}}

	for _, tt := range tests {
		t.Run(tt.name+":"+strings.Join(tt.inFiles, ","), func(t *testing.T) {
			genCode := func() ([]string, *GenConfig) {
				cg := NewDefaultConfig(tt.inSchemaStructPkgPath)
				// Set the name of the caller explicitly to avoid issues when
				// the unit tests are called by external test entities.
				cg.GeneratingBinary = "pathgen-tests"
				cg.FakeRootName = "device"
				if tt.inSchemaStructPkgPath == "" {
					cg.PathStructSuffix = "Path"
				} else {
					cg.PathStructSuffix = ""
				}
				cg.CompressBehaviour = genutil.PreferOperationalState
				cg.GenerateWildcardPaths = true
				cg.PackageName = "ocstructs"

				gotCode, _, err := cg.GeneratePathCode(tt.inFiles, tt.inIncludePaths)
				if err != nil {
					t.Fatalf("GeneratePathCode(%v, %v): Config: %v, got unexpected error: %v", tt.inFiles, tt.inIncludePaths, cg, err)
				}

				files, e := gotCode[cg.PackageName].SplitFiles(tt.inFileNumber)
				if e != nil && !tt.wantErr {
					t.Fatalf("SplitFiles(%v): got unexpected error: %v", tt.inFileNumber, e)
				} else if e == nil && tt.wantErr {
					t.Fatalf("SplitFiles(%v): did not get expected error", tt.inFileNumber)
				}

				return files, cg
			}

			gotCode, cg := genCode()

			var wantCode []string
			for _, codeFile := range tt.wantStructsCodeFiles {
				wantCodeBytes, rferr := os.ReadFile(codeFile)
				if rferr != nil {
					t.Fatalf("os.ReadFile(%q) error: %v", tt.wantStructsCodeFiles, rferr)
				}
				wantCode = append(wantCode, string(wantCodeBytes))
			}

			if len(gotCode) != len(wantCode) {
				t.Errorf("GeneratePathCode(%v, %v), Config: %v, did not return correct code via SplitFiles function (files: %v), (gotfiles: %d, wantfiles: %d), diff (-want, +got):\n%s",
					tt.inFiles, tt.inIncludePaths, cg, tt.wantStructsCodeFiles, len(gotCode), len(wantCode), cmp.Diff(wantCode, gotCode))
			} else {
				for i := range gotCode {
					if gotCode[i] != wantCode[i] {
						if *updateGolden {
							if err := os.WriteFile(tt.wantStructsCodeFiles[i], []byte(gotCode[i]), 0644); err != nil {
								t.Fatal(err)
							}
						}
						// Use difflib to generate a unified diff between the
						// two code snippets such that this is simpler to debug
						// in the test output.
						diff, _ := testutil.GenerateUnifiedDiff(wantCode[i], gotCode[i])
						t.Errorf("GeneratePathCode(%v, %v), Config: %v, did not return correct code via SplitFiles function (file: %v), diff:\n%s",
							tt.inFiles, tt.inIncludePaths, cg, tt.wantStructsCodeFiles[i], diff)
					}
				}
			}
		})
	}
}

func TestGeneratePathCodeSplitModules(t *testing.T) {
	tests := []struct {
		// name is the identifier for the test.
		name string
		// inFiles is the set of inputFiles for the test.
		inFiles []string
		// inIncludePaths is the set of paths that should be searched for imports.
		inIncludePaths []string
		inTrimPrefix   string
		// wantStructsCodeFileDir map from package name to want source file.
		wantStructsCodeFiles map[string]string
	}{{
		name:    "oc simple",
		inFiles: []string{filepath.Join(datapath, "openconfig-simple.yang")},
		wantStructsCodeFiles: map[string]string{
			"openconfigsimplepath": "testdata/modules/oc-simple/simple.txt",
			"device":               "testdata/modules/oc-simple/device.txt",
		},
	}, {
		name:         "oc simple and trim",
		inFiles:      []string{filepath.Join(datapath, "openconfig-simple.yang")},
		inTrimPrefix: "openconfig",
		wantStructsCodeFiles: map[string]string{
			"simplepath": "testdata/modules/oc-simple-trim/simple.txt",
			"device":     "testdata/modules/oc-simple-trim/device.txt",
		},
	}, {
		name:    "oc list builder API",
		inFiles: []string{filepath.Join(datapath, "openconfig-withlist.yang")},
		wantStructsCodeFiles: map[string]string{
			"openconfigwithlistpath": "testdata/modules/oc-list/list.txt",
			"device":                 "testdata/modules/oc-list/device.txt",
		},
	}, {
		name:    "oc import",
		inFiles: []string{filepath.Join(datapath, "openconfig-import.yang")},
		wantStructsCodeFiles: map[string]string{
			"openconfigimportpath":       "testdata/modules/oc-import/import.txt",
			"openconfigsimpletargetpath": "testdata/modules/oc-import/simpletarget.txt",
			"device":                     "testdata/modules/oc-import/device.txt",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name+":"+strings.Join(tt.inFiles, ","), func(t *testing.T) {
			genCode := func() (map[string]string, *GenConfig) {
				cg := NewDefaultConfig("")
				// Set the name of the caller explicitly to avoid issues when
				// the unit tests are called by external test entities.
				cg.GeneratingBinary = "pathgen-tests"
				cg.FakeRootName = "device"
				cg.PackageName = "device"
				cg.CompressBehaviour = genutil.PreferOperationalState
				cg.GenerateWildcardPaths = true
				cg.SplitByModule = true
				cg.BasePackagePath = "example.com"
				cg.TrimPackageModulePrefix = tt.inTrimPrefix

				gotCode, _, err := cg.GeneratePathCode(tt.inFiles, tt.inIncludePaths)
				if err != nil {
					t.Fatalf("GeneratePathCode(%v, %v): Config: %v, got unexpected error: %v", tt.inFiles, tt.inIncludePaths, cg, err)
				}
				files := map[string]string{}
				for k, v := range gotCode {
					files[k] = v.String()
				}

				return files, cg
			}

			gotCode, cg := genCode()

			wantCode := map[string]string{}
			for pkg, codeFile := range tt.wantStructsCodeFiles {
				wantCodeBytes, rferr := os.ReadFile(codeFile)
				if rferr != nil {
					t.Fatalf("os.ReadFile(%q) error: %v", tt.wantStructsCodeFiles, rferr)
				}
				wantCode[pkg] = string(wantCodeBytes)
			}

			if len(gotCode) != len(wantCode) {
				t.Errorf("GeneratePathCode(%v, %v), Config: %v, did not return correct code via SplitFiles function (files: %v), (gotfiles: %d, wantfiles: %d), diff (-want, +got):\n%s",
					tt.inFiles, tt.inIncludePaths, cg, tt.wantStructsCodeFiles, len(gotCode), len(wantCode), cmp.Diff(wantCode, gotCode))
			} else {
				for pkg := range gotCode {
					if gotCode[pkg] != wantCode[pkg] {
						if *updateGolden {
							if err := os.WriteFile(tt.wantStructsCodeFiles[pkg], []byte(gotCode[pkg]), 0644); err != nil {
								t.Fatal(err)
							}
						}
						// Use difflib to generate a unified diff between the
						// two code snippets such that this is simpler to debug
						// in the test output.
						diff, _ := testutil.GenerateUnifiedDiff(wantCode[pkg], gotCode[pkg])
						t.Errorf("GeneratePathCode(%v, %v), Config: %v, did not return correct code via SplitFiles function (file: %v), diff:\n%s",
							tt.inFiles, tt.inIncludePaths, cg, tt.wantStructsCodeFiles[pkg], diff)
					}
				}
			}
		})
	}
}

// getIR is a helper returning an IR to be tested, and its corresponding
// Directory map with relevant fields filled out that would be returned from
// ygen.GenerateIR().
func getIR() *ygen.IR {
	ir := &ygen.IR{
		Directories: map[string]*ygen.ParsedDirectory{
			"/root": {
				Name:       "Root",
				Type:       ygen.Container,
				Path:       "/root",
				IsFakeRoot: true,
				Fields: map[string]*ygen.NodeDetails{
					"leaf": {
						Name: "Leaf",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "leaf",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/leaf",
							SchemaPath:        "/leaf",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Flags:                   map[string]string{yangTypeNameFlagKey: "ieeefloat32"},
						Type:                    ygen.LeafNode,
						LangType:                &ygen.MappedType{NativeType: "Binary"},
						MappedPaths:             [][]string{{"leaf"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
					"leaf-with-default": {
						Name: "LeafWithDefault",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "leaf-with-default",
							Defaults:          []string{"foo"},
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/leaf-with-default",
							SchemaPath:        "/leaf-with-default",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Flags:                   map[string]string{yangTypeNameFlagKey: "string"},
						Type:                    ygen.LeafNode,
						LangType:                &ygen.MappedType{NativeType: "string", DefaultValue: ygot.String(`"foo"`)},
						MappedPaths:             [][]string{{"leaf-with-default"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
					"container": {
						Name: "Container",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "container",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/container",
							SchemaPath:        "/container",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Type:                    ygen.ContainerNode,
						LangType:                nil,
						MappedPaths:             [][]string{{"container"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
					"container-with-config": {
						Name: "ContainerWithConfig",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "container-with-config",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/container-with-config",
							SchemaPath:        "/container-with-config",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Type:                    ygen.ContainerNode,
						LangType:                nil,
						MappedPaths:             [][]string{{"container-with-config"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
					"list": {
						Name: "List",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "list",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/list-container/list",
							SchemaPath:        "/list-container/list",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Type:                    ygen.ListNode,
						LangType:                nil,
						MappedPaths:             [][]string{{"list-container", "list"}},
						MappedPathModules:       [][]string{{"root-module", "root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
					// TODO(wenbli): Move this to a deeper level to test that the parent wildcard receivers are also generated.
					"list-with-state": {
						Name: "ListWithState",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "list-with-state",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/list-container-with-state/list-with-state",
							SchemaPath:        "/list-container-with-state/list-with-state",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Type:                    ygen.ListNode,
						LangType:                nil,
						MappedPaths:             [][]string{{"list-container-with-state", "list-with-state"}},
						MappedPathModules:       [][]string{{"root-module", "root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
					"keyless-list": {
						Name: "KeylessList",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "keyless-list",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/keyless-list-container/keyless-list",
							SchemaPath:        "/keyless-list-container/keyless-list",
							LeafrefTargetPath: "",
							Description:       "",
							ConfigFalse:       true,
						},
						Type:                    ygen.ListNode,
						LangType:                nil,
						MappedPaths:             [][]string{{"keyless-list-container", "keyless-list"}},
						MappedPathModules:       [][]string{{"root-module", "root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
				},
			},
			"/root-module/container": {
				Name: "Container",
				Type: ygen.List,
				Path: "/root-module/container",
				Fields: map[string]*ygen.NodeDetails{
					"leaf": {
						Name: "Leaf",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "leaf",
							Defaults:          []string{"foo"},
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/container/leaf",
							SchemaPath:        "/container/leaf",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Flags: map[string]string{yangTypeNameFlagKey: "int32"},
						Type:  ygen.LeafNode,
						LangType: &ygen.MappedType{
							NativeType: "int32",
						},
						MappedPaths:             [][]string{{"leaf"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
				},
				ListKeys:          nil,
				PackageName:       "",
				BelongingModule:   "root-module",
				RootElementModule: "root-module",
				DefiningModule:    "root-module",
			},
			"/root-module/container-with-config": {
				Name: "ContainerWithConfig",
				Type: ygen.Container,
				Path: "/root-module/container-with-config",
				Fields: map[string]*ygen.NodeDetails{
					"leaflist": {
						Name: "Leaflist",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "leaflist",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/container-with-config/state/leaflist",
							SchemaPath:        "/container-with-config/state/leaflist",
							ShadowSchemaPath:  "/container-with-config/config/leaflist",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Flags: map[string]string{yangTypeNameFlagKey: "uint32"},
						Type:  ygen.LeafListNode,
						LangType: &ygen.MappedType{
							NativeType: "uint32",
						},
						MappedPaths:             [][]string{{"state", "leaflist"}},
						MappedPathModules:       [][]string{{"root-module", "root-module"}},
						ShadowMappedPaths:       [][]string{{"config", "leaflist"}},
						ShadowMappedPathModules: [][]string{{"root-module", "root-module"}},
					},
				},
				ListKeys:          nil,
				PackageName:       "",
				BelongingModule:   "root-module",
				RootElementModule: "root-module",
				DefiningModule:    "root-module",
			},
			"/root-module/list-container-with-state/list-with-state": {
				Name:              "ListWithState",
				Type:              ygen.List,
				Path:              "/root-module/list-container-with-state/list-with-state",
				SchemaPath:        "/list-container-with-state/list-with-state",
				RootElementModule: "root-module",
				Fields: map[string]*ygen.NodeDetails{
					"key": {
						Name: "Key",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "key",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/list-container-with-state/list-with-state/state/key",
							SchemaPath:        "/list-container-with-state/list-with-state/state/key",
							LeafrefTargetPath: "",
							Description:       "",
							ConfigFalse:       true,
						},
						Flags: map[string]string{yangTypeNameFlagKey: "float64"},
						Type:  ygen.LeafNode,
						LangType: &ygen.MappedType{
							NativeType: "float64",
						},
						MappedPaths:             [][]string{{"state", "key"}, {"key"}},
						MappedPathModules:       [][]string{{"root-module", "root-module"}, {"root-module"}},
						ShadowMappedPaths:       [][]string{{"config", "key"}, {"key"}},
						ShadowMappedPathModules: [][]string{{"root-module", "root-module"}, {"root-module"}},
					},
				},
				ListKeys: map[string]*ygen.ListKey{
					"key": {
						Name: "Key",
						LangType: &ygen.MappedType{
							NativeType: "float64",
							ZeroValue:  "0",
						},
					},
				},
				ListKeyYANGNames: []string{"key"},
			},
			"/root-module/list-container/list": {
				Name: "List",
				Type: ygen.List,
				Path: "/root-module/list-container/list",
				Fields: map[string]*ygen.NodeDetails{
					"key1": {
						Name: "Key1",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "key1",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/list-container/list/key1",
							SchemaPath:        "/list-container/list/key1",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Flags: map[string]string{yangTypeNameFlagKey: "string"},
						Type:  ygen.LeafNode,
						LangType: &ygen.MappedType{
							NativeType: "string",
						},
						MappedPaths:             [][]string{{"key1"}},
						MappedPathModules:       [][]string{{"root-module", "root-module"}, {"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
					"key2": {
						Name: "Key2",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "key2",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/list-container/list/key2",
							SchemaPath:        "/list-container/list/key2",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Flags: map[string]string{yangTypeNameFlagKey: "binary"},
						Type:  ygen.LeafNode,
						LangType: &ygen.MappedType{
							NativeType: "Binary",
						},
						MappedPaths:             [][]string{{"key2"}},
						MappedPathModules:       [][]string{{"root-module", "root-module"}, {"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
					"union-key": {
						Name: "UnionKey",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "union-key",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/list-container/list/union-key",
							SchemaPath:        "/list-container/list/union-key",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Flags: map[string]string{yangTypeNameFlagKey: "union"},
						Type:  ygen.LeafNode,
						LangType: &ygen.MappedType{
							NativeType: "RootElementModule_List_UnionKey_Union",
							UnionTypes: map[string]ygen.MappedUnionSubtype{
								"string": {
									Index: 0,
								},
								"Binary": {
									Index: 1,
								},
							},
						},
						MappedPaths:             [][]string{{"union-key"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
				},
				ListKeys: map[string]*ygen.ListKey{
					"key1": {
						Name: "Key1",
						LangType: &ygen.MappedType{
							NativeType: "string",
							ZeroValue:  `""`,
						},
					},
					"key2": {
						Name: "Key2",
						LangType: &ygen.MappedType{
							NativeType: "Binary",
						},
					},
					"union-key": {
						Name: "UnionKey",
						LangType: &ygen.MappedType{
							NativeType: "RootElementModule_List_UnionKey_Union",
							UnionTypes: map[string]ygen.MappedUnionSubtype{
								"string": {
									Index: 0,
								},
								"Binary": {
									Index: 1,
								},
							},
							ZeroValue: "nil",
						},
					},
				},
				ListKeyYANGNames:  []string{"key1", "key2", "union-key"},
				PackageName:       "",
				IsFakeRoot:        false,
				BelongingModule:   "root-module",
				RootElementModule: "root-module",
				DefiningModule:    "root-module",
			},
			"/root-module/keyless-list-container/keyless-list": {
				Name: "KeylessList",
				Type: ygen.List,
				Path: "/root-module/keyless-list-container/keyless-list",
				Fields: map[string]*ygen.NodeDetails{
					"leaf": {
						Name: "Leaf",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "leaf",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/keyless-list-container/keyless-list/leaf",
							SchemaPath:        "/container/leaf",
							LeafrefTargetPath: "",
							Description:       "",
							ConfigFalse:       true,
						},
						Flags: map[string]string{yangTypeNameFlagKey: "int32"},
						Type:  ygen.LeafNode,
						LangType: &ygen.MappedType{
							NativeType: "int32",
						},
						MappedPaths:             [][]string{{"leaf"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
				},
				ListKeys:          nil,
				ListKeyYANGNames:  nil,
				PackageName:       "",
				IsFakeRoot:        false,
				BelongingModule:   "root-module",
				RootElementModule: "root-module",
				DefiningModule:    "root-module",
			},
		},
	}

	return ir
}

func TestGetNodeDataMap(t *testing.T) {
	ir := getIR()

	badIr := &ygen.IR{
		Directories: map[string]*ygen.ParsedDirectory{
			"/root": {
				Name:       "Root",
				Type:       ygen.Container,
				Path:       "/root",
				IsFakeRoot: true,
				Fields: map[string]*ygen.NodeDetails{
					"container": {
						Name: "Container",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "container",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/bad-path/container",
							SchemaPath:        "/container",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Type:                    ygen.ContainerNode,
						LangType:                nil,
						MappedPaths:             [][]string{{"container"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
				},
			},
			"/root-module/container": {
				Name: "Container",
				Type: ygen.List,
				Path: "/root-module/container",
				Fields: map[string]*ygen.NodeDetails{
					"leaf": {
						Name: "Leaf",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "leaf",
							Defaults:          []string{"foo"},
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/container/leaf",
							SchemaPath:        "/container/leaf",
							LeafrefTargetPath: "",
							Description:       "",
						},
						Flags: map[string]string{yangTypeNameFlagKey: "int32"},
						Type:  ygen.LeafNode,
						LangType: &ygen.MappedType{
							NativeType: "int32",
						},
						MappedPaths:             [][]string{{"leaf"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
				},
				ListKeys:          nil,
				PackageName:       "",
				BelongingModule:   "root-module",
				RootElementModule: "root-module",
				DefiningModule:    "root-module",
			},
		},
	}

	irWithSetOrigin := &ygen.IR{
		Directories: map[string]*ygen.ParsedDirectory{
			"/root": {
				Name:       "Root",
				Type:       ygen.Container,
				Path:       "/root",
				IsFakeRoot: true,
				Fields: map[string]*ygen.NodeDetails{
					"container": {
						Name: "Container",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "container",
							Defaults:          nil,
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/container",
							SchemaPath:        "/container",
							LeafrefTargetPath: "",
							Description:       "",
							Origin:            "test-origin",
						},
						Type:                    ygen.ContainerNode,
						LangType:                nil,
						MappedPaths:             [][]string{{"container"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
				},
			},
			"/root-module/container": {
				Name: "Container",
				Type: ygen.List,
				Path: "/root-module/container",
				Fields: map[string]*ygen.NodeDetails{
					"leaf": {
						Name: "Leaf",
						YANGDetails: ygen.YANGNodeDetails{
							Name:              "leaf",
							Defaults:          []string{"foo"},
							BelongingModule:   "root-module",
							RootElementModule: "root-module",
							DefiningModule:    "root-module",
							Path:              "/root-module/container/leaf",
							SchemaPath:        "/container/leaf",
							LeafrefTargetPath: "",
							Description:       "",
							Origin:            "test-origin",
						},
						Flags: map[string]string{yangTypeNameFlagKey: "int32"},
						Type:  ygen.LeafNode,
						LangType: &ygen.MappedType{
							NativeType: "int32",
						},
						MappedPaths:             [][]string{{"leaf"}},
						MappedPathModules:       [][]string{{"root-module"}},
						ShadowMappedPaths:       nil,
						ShadowMappedPathModules: nil,
					},
				},
				ListKeys:          nil,
				ListKeyYANGNames:  nil,
				PackageName:       "",
				IsFakeRoot:        false,
				BelongingModule:   "root-module",
				RootElementModule: "root-module",
				DefiningModule:    "root-module",
			},
		},
	}
	tests := []struct {
		name                      string
		inIR                      *ygen.IR
		inFakeRootName            string
		inSchemaStructPkgAccessor string
		inPathStructSuffix        string
		inPackageName             string
		inPackageSuffix           string
		inSplitByModule           bool
		wantNodeDataMap           NodeDataMap
		wantSorted                []string
		wantErrSubstrings         []string
	}{{
		name:                      "non-existent field path",
		inIR:                      badIr,
		inFakeRootName:            "device",
		inSchemaStructPkgAccessor: "oc.",
		inPathStructSuffix:        "Path",
		wantErrSubstrings:         []string{`field with path "/bad-path/container" not found`},
	}, {
		name:                      "big test with everything",
		inIR:                      ir,
		inFakeRootName:            "Root",
		inSchemaStructPkgAccessor: "struct.",
		inPathStructSuffix:        "_Path",
		inSplitByModule:           true,
		inPackageName:             "device",
		inPackageSuffix:           "path",
		wantNodeDataMap: NodeDataMap{
			"Container_Path": {
				GoTypeName:            "*struct.Container",
				LocalGoTypeName:       "*Container",
				GoFieldName:           "Container",
				SubsumingGoStructName: "Container",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/container",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"ContainerWithConfig_Path": {
				GoTypeName:            "*struct.ContainerWithConfig",
				LocalGoTypeName:       "*ContainerWithConfig",
				GoFieldName:           "ContainerWithConfig",
				SubsumingGoStructName: "ContainerWithConfig",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/container-with-config",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"ContainerWithConfig_Leaflist_Path": {
				GoTypeName:            "[]uint32",
				LocalGoTypeName:       "[]uint32",
				GoFieldName:           "Leaflist",
				SubsumingGoStructName: "ContainerWithConfig",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "uint32",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/container-with-config",
				YANGFieldName:         "leaflist",
				PathOriginName:        "openconfig",
			},
			"Container_Leaf_Path": {
				GoTypeName:            "int32",
				LocalGoTypeName:       "int32",
				GoFieldName:           "Leaf",
				SubsumingGoStructName: "Container",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            true,
				YANGTypeName:          "int32",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/container",
				YANGFieldName:         "leaf",
				PathOriginName:        "openconfig",
			},
			"Leaf_Path": {
				GoTypeName:            "struct.Binary",
				LocalGoTypeName:       "Binary",
				GoFieldName:           "Leaf",
				SubsumingGoStructName: "Root",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "ieeefloat32",
				GoPathPackageName:     "device",
				DirectoryName:         "/root",
				YANGFieldName:         "leaf",
				PathOriginName:        "openconfig",
			},
			"LeafWithDefault_Path": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "LeafWithDefault",
				SubsumingGoStructName: "Root",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            true,
				YANGTypeName:          "string",
				GoPathPackageName:     "device",
				DirectoryName:         "/root",
				YANGFieldName:         "leaf-with-default",
				PathOriginName:        "openconfig",
			},
			"List_Path": {
				GoTypeName:            "*struct.List",
				LocalGoTypeName:       "*List",
				GoFieldName:           "List",
				SubsumingGoStructName: "List",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/list-container/list",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"List_PathMap": {
				GoTypeName:            "map[struct.List_Key]*struct.List",
				LocalGoTypeName:       "map[List_Key]*List",
				GoFieldName:           "List",
				SubsumingGoStructName: "Root",
				IsListContainer:       true,
				YANGPath:              "/root-module/list-container/list",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/list-container/list",
				CompressInfo: &CompressionInfo{
					PreRelPathList:  `"root-module:list-container"`,
					PostRelPathList: `"root-module:list"`,
				},
				PathOriginName: "openconfig",
			},
			"ListWithState_Path": {
				GoTypeName:            "*struct.ListWithState",
				LocalGoTypeName:       "*ListWithState",
				GoFieldName:           "ListWithState",
				SubsumingGoStructName: "ListWithState",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/list-container-with-state/list-with-state",
				YANGFieldName:         "",
				PathOriginName:        "openconfig",
			},
			"ListWithState_PathMap": {
				GoTypeName:            "map[float64]*struct.ListWithState",
				LocalGoTypeName:       "map[float64]*ListWithState",
				GoFieldName:           "ListWithState",
				SubsumingGoStructName: "Root",
				IsListContainer:       true,
				YANGPath:              "/root-module/list-container-with-state/list-with-state",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/list-container-with-state/list-with-state",
				CompressInfo: &CompressionInfo{
					PreRelPathList:  `"root-module:list-container-with-state"`,
					PostRelPathList: `"root-module:list-with-state"`,
				},
				PathOriginName: "openconfig",
			},
			"ListWithState_Key_Path": {
				GoTypeName:            "float64",
				LocalGoTypeName:       "float64",
				GoFieldName:           "Key",
				SubsumingGoStructName: "ListWithState",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "float64",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/list-container-with-state/list-with-state",
				YANGFieldName:         "key",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"List_Key1_Path": {
				GoTypeName:            "string",
				LocalGoTypeName:       "string",
				GoFieldName:           "Key1",
				SubsumingGoStructName: "List",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            false,
				YANGTypeName:          "string",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/list-container/list",
				YANGFieldName:         "key1",
				PathOriginName:        "openconfig",
			},
			"List_Key2_Path": {
				GoTypeName:            "struct.Binary",
				LocalGoTypeName:       "Binary",
				GoFieldName:           "Key2",
				SubsumingGoStructName: "List",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "binary",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/list-container/list",
				YANGFieldName:         "key2",
				PathOriginName:        "openconfig",
			},
			"List_UnionKey_Path": {
				GoTypeName:            "struct.RootElementModule_List_UnionKey_Union",
				LocalGoTypeName:       "RootElementModule_List_UnionKey_Union",
				GoFieldName:           "UnionKey",
				SubsumingGoStructName: "List",
				IsLeaf:                true,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "union",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/list-container/list",
				YANGFieldName:         "union-key",
				PathOriginName:        "openconfig",
			},
			"Root_Path": {
				GoTypeName:            "*struct.Root",
				LocalGoTypeName:       "*Root",
				GoFieldName:           "",
				SubsumingGoStructName: "Root",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "device",
				DirectoryName:         "/root",
			},
			"KeylessList_Path": {
				GoTypeName:            "*struct.KeylessList",
				LocalGoTypeName:       "*KeylessList",
				GoFieldName:           "KeylessList",
				SubsumingGoStructName: "KeylessList",
				YANGPath:              "/root-module/keyless-list-container/keyless-list",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/keyless-list-container/keyless-list",
				YANGFieldName:         "",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			},
			"KeylessList_Leaf_Path": {
				GoTypeName:            "int32",
				LocalGoTypeName:       "int32",
				GoFieldName:           "Leaf",
				SubsumingGoStructName: "KeylessList",
				IsLeaf:                true,
				IsScalarField:         true,
				YANGTypeName:          "int32",
				YANGPath:              "/root-module/keyless-list-container/keyless-list/leaf",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/keyless-list-container/keyless-list",
				YANGFieldName:         "leaf",
				ConfigFalse:           true,
				PathOriginName:        "openconfig",
			}},
		wantSorted: []string{
			"ContainerWithConfig_Leaflist_Path",
			"ContainerWithConfig_Path",
			"Container_Leaf_Path",
			"Container_Path",
			"KeylessList_Leaf_Path",
			"KeylessList_Path",
			"LeafWithDefault_Path",
			"Leaf_Path",
			"ListWithState_Key_Path",
			"ListWithState_Path",
			"ListWithState_PathMap",
			"List_Key1_Path",
			"List_Key2_Path",
			"List_Path",
			"List_PathMap",
			"List_UnionKey_Path",
			"Root_Path",
		},
	}, {
		name:                      "smalle test with arbitrary PathOriginName",
		inIR:                      irWithSetOrigin,
		inFakeRootName:            "Root",
		inSchemaStructPkgAccessor: "struct.",
		inPathStructSuffix:        "_Path",
		inSplitByModule:           true,
		inPackageName:             "device",
		inPackageSuffix:           "path",
		wantNodeDataMap: NodeDataMap{
			"Container_Leaf_Path": {
				GoTypeName:            "int32",
				LocalGoTypeName:       "int32",
				GoFieldName:           "Leaf",
				SubsumingGoStructName: "Container",
				IsLeaf:                true,
				IsScalarField:         true,
				HasDefault:            true,
				YANGTypeName:          "int32",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/container",
				YANGFieldName:         "leaf",
				PathOriginName:        "test-origin",
			},
			"Container_Path": {
				GoTypeName:            "*struct.Container",
				LocalGoTypeName:       "*Container",
				GoFieldName:           "Container",
				SubsumingGoStructName: "Container",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/container",
				YANGFieldName:         "",
				PathOriginName:        "test-origin",
			},
			"Root_Path": {
				GoTypeName:            "*struct.Root",
				LocalGoTypeName:       "*Root",
				GoFieldName:           "",
				SubsumingGoStructName: "Root",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				GoPathPackageName:     "device",
				DirectoryName:         "/root",
			}},
		wantSorted: []string{
			"Container_Leaf_Path",
			"Container_Path",
			"Root_Path",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErrs := getNodeDataMap(tt.inIR, tt.inFakeRootName, tt.inSchemaStructPkgAccessor, tt.inPathStructSuffix, tt.inPackageName, tt.inPackageSuffix, tt.inSplitByModule, "", true, genutil.PreferOperationalState, nil)
			// TODO(wenbli): Enhance gNMI's errdiff with checking a slice of substrings and use here.
			var gotErrStrs []string
			for _, err := range gotErrs {
				gotErrStrs = append(gotErrStrs, err.Error())
			}
			if diff := cmp.Diff(tt.wantErrSubstrings, gotErrStrs, cmp.Comparer(func(x, y string) bool { return strings.Contains(x, y) || strings.Contains(y, x) })); diff != "" {
				t.Fatalf("Error substring check failed (-want, +got):\n%v", diff)
			}
			if diff := cmp.Diff(tt.wantNodeDataMap, got, cmpopts.IgnoreFields(NodeData{}, "YANGPath")); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantSorted, GetOrderedNodeDataNames(got), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("(-want sorted names, +got sorted names):\n%s", diff)
			}
		})
	}
}

// trimDocComments removes doc comments from the input code snippet string.
// Example:
//
//	// foo does bar
//	func foo() {
//	  // baz is need to do boo.
//	  baz()
//	}
//
//	// foo2 does bar2
//	func foo2() {
//	  // baz2 is need to do boo2.
//	  baz2()
//	}
//
// After:
//
//	func foo() {
//	  // baz is need to do boo.
//	  baz()
//	}
//
//	func foo2() {
//	  // baz2 is need to do boo2.
//	  baz2()
//	}
func trimDocComments(snippet string) string {
	var b strings.Builder
	for i, line := range strings.Split(snippet, "\n") {
		// i > 0 to prevent two newlines from being printed on an empty
		// line at the end of the string.
		if i > 0 && line == "" {
			b.WriteString("\n")
			continue
		}
		if !strings.HasPrefix(line, "//") {
			if i > 0 {
				b.WriteString("\n")
			}
			b.WriteString(line)
		}
	}
	return b.String()
}

const (
	// wantListMethodsNonWildcard is the expected non-wildcard child constructor
	// method for the test list node.
	wantListMethodsNonWildcard = `
// List (list): 
// 	Defining module:      "root-module"
// 	Instantiating module: "root-module"
// 	Path from parent:     "list-container/list"
// 	Path from root:       "/list-container/list"
//
// 	Key1: string
// 	Key2: oc.Binary
// 	UnionKey: [oc.UnionString, oc.Binary]
func (n *RootPath) List(Key1 string, Key2 oc.Binary, UnionKey oc.RootElementModule_List_UnionKey_Union) *ListPath {
	ps := &ListPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container", "list"},
			map[string]interface{}{"key1": Key1, "key2": Key2, "union-key": UnionKey},
			n,
		),
	}
	return ps
}

// ListMap (list): 
// 	Defining module:      "root-module"
// 	Instantiating module: "root-module"
// 	Path from parent:     "list-container/list"
// 	Path from root:       "/list-container/list"
//
// 	Key1: string
// 	Key2: oc.Binary
// 	UnionKey: [oc.UnionString, oc.Binary]
func (n *RootPath) ListMap() *ListPathMap {
	ps := &ListPathMap{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`

	// wantListMethods is the expected child constructor methods for the test list node.
	wantListMethods = `
// ListAny (list): 
// 	Defining module:      "root-module"
// 	Instantiating module: "root-module"
// 	Path from parent:     "list-container/list"
// 	Path from root:       "/list-container/list"
//
// 	Key1 (wildcarded): string
// 	Key2 (wildcarded): oc.Binary
// 	UnionKey (wildcarded): [oc.UnionString, oc.Binary]
func (n *RootPath) ListAny() *ListPathAny {
	ps := &ListPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container", "list"},
			map[string]interface{}{"key1": "*", "key2": "*", "union-key": "*"},
			n,
		),
	}
	return ps
}

func (n *ListPathAny) WithKey1(Key1 string) *ListPathAny {
	ygnmi.ModifyKey(n.NodePath, "key1", Key1)
	return n
}

func (n *ListPathAny) WithKey2(Key2 oc.Binary) *ListPathAny {
	ygnmi.ModifyKey(n.NodePath, "key2", Key2)
	return n
}

func (n *ListPathAny) WithUnionKey(UnionKey oc.RootElementModule_List_UnionKey_Union) *ListPathAny {
	ygnmi.ModifyKey(n.NodePath, "union-key", UnionKey)
	return n
}
` + wantListMethodsNonWildcard

	// wantNonListMethods is the expected child constructor methods for
	// non-list elements from the root.
	wantNonListMethods = `
// Container returns from RootPath the path struct for its child "container".
func (n *RootPath) Container() *ContainerPath {
	ps := &ContainerPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"container"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

// ContainerWithConfig returns from RootPath the path struct for its child "container-with-config".
func (n *RootPath) ContainerWithConfig() *ContainerWithConfigPath {
	ps := &ContainerWithConfigPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"container-with-config"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *RootPath) KeylessListAny() *KeylessListPathAny {
	ps := &KeylessListPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"keyless-list-container", "keyless-list"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

// Leaf returns from RootPath the path struct for its child "leaf".
func (n *RootPath) Leaf() *LeafPath {
	ps := &LeafPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"leaf"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

// LeafWithDefault returns from RootPath the path struct for its child "leaf-with-default".
func (n *RootPath) LeafWithDefault() *LeafWithDefaultPath {
	ps := &LeafWithDefaultPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"leaf-with-default"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`

	// wantNonListMethodsSplitModule is the expected child constructor
	// methods for non-list elements from the root with split modules.
	wantNonListMethodsSplitModule = `
// Container returns from RootPath the path struct for its child "container".
func (n *RootPath) Container() *rootmodulepath.ContainerPath {
	ps := &rootmodulepath.ContainerPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"container"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

// ContainerWithConfig returns from RootPath the path struct for its child "container-with-config".
func (n *RootPath) ContainerWithConfig() *rootmodulepath.ContainerWithConfigPath {
	ps := &rootmodulepath.ContainerWithConfigPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"container-with-config"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *RootPath) KeylessListAny() *rootmodulepath.KeylessListPathAny {
	ps := &rootmodulepath.KeylessListPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"keyless-list-container", "keyless-list"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

// Leaf returns from RootPath the path struct for its child "leaf".
func (n *RootPath) Leaf() *LeafPath {
	ps := &LeafPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"leaf"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

// LeafWithDefault returns from RootPath the path struct for its child "leaf-with-default".
func (n *RootPath) LeafWithDefault() *LeafWithDefaultPath {
	ps := &LeafWithDefaultPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"leaf-with-default"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`

	// wantStructBaseFakeRootNWC is the expected structs for the root device
	// when wildcards are disabled.
	wantFakeRootStructsNWC = `
// RootPath represents the /root YANG schema element.
type RootPath struct {
	*ygnmi.DeviceRootBase
}

// Root returns a root path object from which YANG paths can be constructed.
func Root() *RootPath {
	return &RootPath{ygnmi.NewDeviceRootBase()}
}
`

	// wantFakeRootStructsWC is the expected structs for the root device
	// when wildcards are enabled.
	wantFakeRootStructsWC = `
// RootPath represents the /root YANG schema element.
type RootPath struct {
	*ygnmi.DeviceRootBase
}

// Root returns a root path object from which YANG paths can be constructed.
func Root() *RootPath {
	return &RootPath{ygnmi.NewDeviceRootBase()}
}
`
)

func TestGenerateDirectorySnippet(t *testing.T) {
	directories := getIR().Directories

	tests := []struct {
		name               string
		inDirectory        *ygen.ParsedDirectory
		inPathStructSuffix string
		inSplitByModule    bool
		inPackageName      string
		inPackageSuffix    string
		inUnifiedPath      bool
		inNodeDataMap      NodeDataMap
		// want may be omitted to skip testing.
		want []GoPathStructCodeSnippet
		// wantNoWildcard may be omitted to skip testing.
		wantNoWildcard []GoPathStructCodeSnippet
	}{{
		name:            "container-with-config",
		inDirectory:     directories["/root-module/container-with-config"],
		inPackageName:   "device",
		inPackageSuffix: "path",
		want: []GoPathStructCodeSnippet{{
			PathStructName: "ContainerWithConfig_Leaflist",
			StructBase: `
// ContainerWithConfig_Leaflist represents the /root-module/container-with-config/state/leaflist YANG schema element.
type ContainerWithConfig_Leaflist struct {
	*ygnmi.NodePath
}

// ContainerWithConfig_LeaflistAny represents the wildcard version of the /root-module/container-with-config/state/leaflist YANG schema element.
type ContainerWithConfig_LeaflistAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *ContainerWithConfig_Leaflist) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "ContainerWithConfig",
			StructBase: `
// ContainerWithConfig represents the /root-module/container-with-config YANG schema element.
type ContainerWithConfig struct {
	*ygnmi.NodePath
}

// ContainerWithConfigAny represents the wildcard version of the /root-module/container-with-config YANG schema element.
type ContainerWithConfigAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *ContainerWithConfig) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: `
func (n *ContainerWithConfig) Leaflist() *ContainerWithConfig_Leaflist {
	ps := &ContainerWithConfig_Leaflist{
		NodePath: ygnmi.NewNodePath(
			[]string{"state", "leaflist"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *ContainerWithConfigAny) Leaflist() *ContainerWithConfig_LeaflistAny {
	ps := &ContainerWithConfig_LeaflistAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"state", "leaflist"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
			Package:         "device",
			ExtraGeneration: "",
		}},
		wantNoWildcard: []GoPathStructCodeSnippet{{
			PathStructName: "ContainerWithConfig_Leaflist",
			StructBase: `
// ContainerWithConfig_Leaflist represents the /root-module/container-with-config/state/leaflist YANG schema element.
type ContainerWithConfig_Leaflist struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *ContainerWithConfig_Leaflist) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "ContainerWithConfig",
			StructBase: `
// ContainerWithConfig represents the /root-module/container-with-config YANG schema element.
type ContainerWithConfig struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *ContainerWithConfig) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: `
func (n *ContainerWithConfig) Leaflist() *ContainerWithConfig_Leaflist {
	ps := &ContainerWithConfig_Leaflist{
		NodePath: ygnmi.NewNodePath(
			[]string{"state", "leaflist"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
			Package:         "device",
			ExtraGeneration: "",
		}},
	}, {
		name:            "unified-container-with-config",
		inDirectory:     directories["/root-module/container-with-config"],
		inPackageName:   "device",
		inPackageSuffix: "path",
		inUnifiedPath:   true,
		want: []GoPathStructCodeSnippet{{
			PathStructName: "ContainerWithConfig_Leaflist",
			StructBase: `
// ContainerWithConfig_Leaflist represents the /root-module/container-with-config/state/leaflist YANG schema element.
type ContainerWithConfig_Leaflist struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// ContainerWithConfig_LeaflistAny represents the wildcard version of the /root-module/container-with-config/state/leaflist YANG schema element.
type ContainerWithConfig_LeaflistAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// PathOrigin returns the name of the origin for the path object.
func (n *ContainerWithConfig_Leaflist) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "ContainerWithConfig",
			StructBase: `
// ContainerWithConfig represents the /root-module/container-with-config YANG schema element.
type ContainerWithConfig struct {
	*ygnmi.NodePath
}

// ContainerWithConfigAny represents the wildcard version of the /root-module/container-with-config YANG schema element.
type ContainerWithConfigAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *ContainerWithConfig) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: `
func (n *ContainerWithConfig) Leaflist() *ContainerWithConfig_Leaflist {
	ps := &ContainerWithConfig_Leaflist{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "leaflist"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

func (n *ContainerWithConfigAny) Leaflist() *ContainerWithConfig_LeaflistAny {
	ps := &ContainerWithConfig_LeaflistAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "leaflist"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}
`,
			Package:         "device",
			ExtraGeneration: "",
		}},
		wantNoWildcard: []GoPathStructCodeSnippet{{
			PathStructName: "ContainerWithConfig_Leaflist",
			StructBase: `
// ContainerWithConfig_Leaflist represents the /root-module/container-with-config/state/leaflist YANG schema element.
type ContainerWithConfig_Leaflist struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// PathOrigin returns the name of the origin for the path object.
func (n *ContainerWithConfig_Leaflist) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "ContainerWithConfig",
			StructBase: `
// ContainerWithConfig represents the /root-module/container-with-config YANG schema element.
type ContainerWithConfig struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *ContainerWithConfig) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: `
func (n *ContainerWithConfig) Leaflist() *ContainerWithConfig_Leaflist {
	ps := &ContainerWithConfig_Leaflist{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "leaflist"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}
`,
			Package:         "device",
			ExtraGeneration: "",
		}},
	}, {
		name:               "fakeroot",
		inDirectory:        directories["/root"],
		inPathStructSuffix: "Path",
		inPackageName:      "ocpathstructs",
		inPackageSuffix:    "path",
		want: []GoPathStructCodeSnippet{{
			PathStructName: "LeafPath",
			StructBase: `
// LeafPath represents the /root-module/leaf YANG schema element.
type LeafPath struct {
	*ygnmi.NodePath
}

// LeafPathAny represents the wildcard version of the /root-module/leaf YANG schema element.
type LeafPathAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *LeafPath) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "ocpathstructs",
			ExtraGeneration:   "",
		}, {
			PathStructName: "LeafWithDefaultPath",
			StructBase: `
// LeafWithDefaultPath represents the /root-module/leaf-with-default YANG schema element.
type LeafWithDefaultPath struct {
	*ygnmi.NodePath
}

// LeafWithDefaultPathAny represents the wildcard version of the /root-module/leaf-with-default YANG schema element.
type LeafWithDefaultPathAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *LeafWithDefaultPath) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "ocpathstructs",
			ExtraGeneration:   "",
		}, {
			PathStructName: "RootPath",
			Package:        "ocpathstructs",
			StructBase:     wantFakeRootStructsWC,
			ChildConstructors: trimDocComments(wantNonListMethods+wantListMethods) + `
func (n *RootPath) ListWithStateAny() *ListWithStatePathAny {
	ps := &ListWithStatePathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state", "list-with-state"},
			map[string]interface{}{"key": "*"},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListWithState(Key float64) *ListWithStatePath {
	ps := &ListWithStatePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state", "list-with-state"},
			map[string]interface{}{"key": Key},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListWithStateMap() *ListWithStatePathMap {
	ps := &ListWithStatePathMap{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
		}},
		wantNoWildcard: []GoPathStructCodeSnippet{{
			PathStructName: "LeafPath",
			StructBase: `
// LeafPath represents the /root-module/leaf YANG schema element.
type LeafPath struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *LeafPath) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "ocpathstructs",
			ExtraGeneration:   "",
		}, {
			PathStructName: "LeafWithDefaultPath",
			StructBase: `
// LeafWithDefaultPath represents the /root-module/leaf-with-default YANG schema element.
type LeafWithDefaultPath struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *LeafWithDefaultPath) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "ocpathstructs",
			ExtraGeneration:   "",
		}, {
			PathStructName: "RootPath",
			Package:        "ocpathstructs",
			StructBase:     wantFakeRootStructsWC,
			ChildConstructors: trimDocComments(wantNonListMethods+wantListMethodsNonWildcard) + `
func (n *RootPath) ListWithState(Key float64) *ListWithStatePath {
	ps := &ListWithStatePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state", "list-with-state"},
			map[string]interface{}{"key": Key},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListWithStateMap() *ListWithStatePathMap {
	ps := &ListWithStatePathMap{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
		}},
	}, {
		name:            "list",
		inDirectory:     directories["/root-module/list-container/list"],
		inPackageName:   "device",
		inPackageSuffix: "path",
		want: []GoPathStructCodeSnippet{{
			PathStructName: "List_Key1",
			StructBase: `
// List_Key1 represents the /root-module/list-container/list/key1 YANG schema element.
type List_Key1 struct {
	*ygnmi.NodePath
}

// List_Key1Any represents the wildcard version of the /root-module/list-container/list/key1 YANG schema element.
type List_Key1Any struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *List_Key1) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "List_Key2",
			StructBase: `
// List_Key2 represents the /root-module/list-container/list/key2 YANG schema element.
type List_Key2 struct {
	*ygnmi.NodePath
}

// List_Key2Any represents the wildcard version of the /root-module/list-container/list/key2 YANG schema element.
type List_Key2Any struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *List_Key2) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "List_UnionKey",
			StructBase: `
// List_UnionKey represents the /root-module/list-container/list/union-key YANG schema element.
type List_UnionKey struct {
	*ygnmi.NodePath
}

// List_UnionKeyAny represents the wildcard version of the /root-module/list-container/list/union-key YANG schema element.
type List_UnionKeyAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *List_UnionKey) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "List",
			StructBase: `
// List represents the /root-module/list-container/list YANG schema element.
type List struct {
	*ygnmi.NodePath
}

// ListAny represents the wildcard version of the /root-module/list-container/list YANG schema element.
type ListAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *List) PathOriginName() string {
     return ""
}

// ListMap represents the /root-module/list-container/list YANG schema element.
type ListMap struct {
	*ygnmi.NodePath
}

// ListMapAny represents the wildcard version of the /root-module/list-container/list YANG schema element.
type ListMapAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *ListMap) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: `
func (n *List) Key1() *List_Key1 {
	ps := &List_Key1{
		NodePath: ygnmi.NewNodePath(
			[]string{"key1"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *ListAny) Key1() *List_Key1Any {
	ps := &List_Key1Any{
		NodePath: ygnmi.NewNodePath(
			[]string{"key1"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *List) Key2() *List_Key2 {
	ps := &List_Key2{
		NodePath: ygnmi.NewNodePath(
			[]string{"key2"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *ListAny) Key2() *List_Key2Any {
	ps := &List_Key2Any{
		NodePath: ygnmi.NewNodePath(
			[]string{"key2"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *List) UnionKey() *List_UnionKey {
	ps := &List_UnionKey{
		NodePath: ygnmi.NewNodePath(
			[]string{"union-key"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *ListAny) UnionKey() *List_UnionKeyAny {
	ps := &List_UnionKeyAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"union-key"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
			Package:         "device",
			ExtraGeneration: "",
		}},
		wantNoWildcard: []GoPathStructCodeSnippet{{
			PathStructName: "List_Key1",
			StructBase: `
// List_Key1 represents the /root-module/list-container/list/key1 YANG schema element.
type List_Key1 struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *List_Key1) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "List_Key2",
			StructBase: `
// List_Key2 represents the /root-module/list-container/list/key2 YANG schema element.
type List_Key2 struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *List_Key2) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "List_UnionKey",
			StructBase: `
// List_UnionKey represents the /root-module/list-container/list/union-key YANG schema element.
type List_UnionKey struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *List_UnionKey) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "List",
			StructBase: `
// List represents the /root-module/list-container/list YANG schema element.
type List struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *List) PathOriginName() string {
     return ""
}

// ListMap represents the /root-module/list-container/list YANG schema element.
type ListMap struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *ListMap) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: `
func (n *List) Key1() *List_Key1 {
	ps := &List_Key1{
		NodePath: ygnmi.NewNodePath(
			[]string{"key1"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *List) Key2() *List_Key2 {
	ps := &List_Key2{
		NodePath: ygnmi.NewNodePath(
			[]string{"key2"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *List) UnionKey() *List_UnionKey {
	ps := &List_UnionKey{
		NodePath: ygnmi.NewNodePath(
			[]string{"union-key"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
			Package:         "device",
			ExtraGeneration: "",
		}},
	}, {
		name:               "fakeroot split by modules",
		inDirectory:        directories["/root"],
		inPathStructSuffix: "Path",
		inSplitByModule:    true,
		inPackageName:      "device",
		inPackageSuffix:    "path",
		wantNoWildcard: []GoPathStructCodeSnippet{{
			PathStructName: "LeafPath",
			StructBase: `
// LeafPath represents the /root-module/leaf YANG schema element.
type LeafPath struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *LeafPath) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "LeafWithDefaultPath",
			StructBase: `
// LeafWithDefaultPath represents the /root-module/leaf-with-default YANG schema element.
type LeafWithDefaultPath struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *LeafWithDefaultPath) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "RootPath",
			Package:        "device",
			Deps:           []string{"rootmodulepath"},
			StructBase:     wantFakeRootStructsNWC,
			ChildConstructors: trimDocComments(wantNonListMethodsSplitModule) + `
func (n *RootPath) List(Key1 string, Key2 oc.Binary, UnionKey oc.RootElementModule_List_UnionKey_Union) *rootmodulepath.ListPath {
	ps := &rootmodulepath.ListPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container", "list"},
			map[string]interface{}{"key1": Key1, "key2": Key2, "union-key": UnionKey},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListMap() *rootmodulepath.ListPathMap {
	ps := &rootmodulepath.ListPathMap{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListWithState(Key float64) *rootmodulepath.ListWithStatePath {
	ps := &rootmodulepath.ListWithStatePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state", "list-with-state"},
			map[string]interface{}{"key": Key},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListWithStateMap() *rootmodulepath.ListWithStatePathMap {
	ps := &rootmodulepath.ListWithStatePathMap{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
		}},
	}, {
		name:               "fakeroot split by modules and builder API",
		inDirectory:        directories["/root"],
		inPathStructSuffix: "Path",
		inSplitByModule:    true,
		inPackageName:      "device",
		inPackageSuffix:    "path",
		want: []GoPathStructCodeSnippet{{
			PathStructName: "LeafPath",
			StructBase: `
// LeafPath represents the /root-module/leaf YANG schema element.
type LeafPath struct {
	*ygnmi.NodePath
}

// LeafPathAny represents the wildcard version of the /root-module/leaf YANG schema element.
type LeafPathAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *LeafPath) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "LeafWithDefaultPath",
			StructBase: `
// LeafWithDefaultPath represents the /root-module/leaf-with-default YANG schema element.
type LeafWithDefaultPath struct {
	*ygnmi.NodePath
}

// LeafWithDefaultPathAny represents the wildcard version of the /root-module/leaf-with-default YANG schema element.
type LeafWithDefaultPathAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *LeafWithDefaultPath) PathOriginName() string {
     return ""
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		}, {
			PathStructName: "RootPath",
			StructBase:     wantFakeRootStructsWC,
			ChildConstructors: trimDocComments(wantNonListMethodsSplitModule) + `
func (n *RootPath) ListAny() *rootmodulepath.ListPathAny {
	ps := &rootmodulepath.ListPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container", "list"},
			map[string]interface{}{"key1": "*", "key2": "*", "union-key": "*"},
			n,
		),
	}
	return ps
}

func (n *RootPath) List(Key1 string, Key2 oc.Binary, UnionKey oc.RootElementModule_List_UnionKey_Union) *rootmodulepath.ListPath {
	ps := &rootmodulepath.ListPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container", "list"},
			map[string]interface{}{"key1": Key1, "key2": Key2, "union-key": UnionKey},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListMap() *rootmodulepath.ListPathMap {
	ps := &rootmodulepath.ListPathMap{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListWithStateAny() *rootmodulepath.ListWithStatePathAny {
	ps := &rootmodulepath.ListWithStatePathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state", "list-with-state"},
			map[string]interface{}{"key": "*"},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListWithState(Key float64) *rootmodulepath.ListWithStatePath {
	ps := &rootmodulepath.ListWithStatePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state", "list-with-state"},
			map[string]interface{}{"key": Key},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListWithStateMap() *rootmodulepath.ListWithStatePathMap {
	ps := &rootmodulepath.ListWithStatePathMap{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
			Package:         "device",
			Deps:            []string{"rootmodulepath"},
			ExtraGeneration: "",
		}, {
			PathStructName: "RootPath",
			StructBase:     ``,
			ChildConstructors: `
func (n *ListPathAny) WithKey1(Key1 string) *ListPathAny {
	ygnmi.ModifyKey(n.NodePath, "key1", Key1)
	return n
}

func (n *ListPathAny) WithKey2(Key2 oc.Binary) *ListPathAny {
	ygnmi.ModifyKey(n.NodePath, "key2", Key2)
	return n
}

func (n *ListPathAny) WithUnionKey(UnionKey oc.RootElementModule_List_UnionKey_Union) *ListPathAny {
	ygnmi.ModifyKey(n.NodePath, "union-key", UnionKey)
	return n
}
`,
			Package:         "rootmodulepath",
			ExtraGeneration: "",
		}},
	}, {
		name:            "container with PathOriginName",
		inDirectory:     directories["/root-module/container"],
		inPackageName:   "device",
		inPackageSuffix: "path",
		inNodeDataMap: NodeDataMap{
			"Container": {
				GoTypeName:            "*struct.Container",
				LocalGoTypeName:       "*Container",
				GoFieldName:           "Container",
				SubsumingGoStructName: "Container",
				IsLeaf:                false,
				IsScalarField:         false,
				HasDefault:            false,
				YANGTypeName:          "",
				GoPathPackageName:     "rootmodulepath",
				DirectoryName:         "/root-module/container",
				YANGFieldName:         "",
				PathOriginName:        "test-origin",
			},
		},
		want: []GoPathStructCodeSnippet{{
			PathStructName: "Container",
			StructBase: `
// Container represents the /root-module/container YANG schema element.
type Container struct {
	*ygnmi.NodePath
}

// ContainerAny represents the wildcard version of the /root-module/container YANG schema element.
type ContainerAny struct {
	*ygnmi.NodePath
}

// PathOrigin returns the name of the origin for the path object.
func (n *Container) PathOriginName() string {
     return "test-origin"
}
`,
			ChildConstructors: ``,
			Package:           "device",
			ExtraGeneration:   "",
		},
		},
	}}

	for _, tt := range tests {
		if tt.want != nil {
			t.Run(tt.name, func(t *testing.T) {
				got, gotErr := generateDirectorySnippet(tt.inDirectory, directories, tt.inNodeDataMap, "Root", ExtraGenerators{}, "oc.", tt.inPathStructSuffix, true, tt.inSplitByModule, "", tt.inPackageName, tt.inPackageSuffix, nil, tt.inUnifiedPath, true, true, genutil.PreferOperationalState)
				if gotErr != nil {
					t.Fatalf("func generateDirectorySnippet, unexpected error: %v", gotErr)
				}

				for i, s := range got {
					got[i].ChildConstructors = trimDocComments(s.ChildConstructors)
					t.Log(got[i])
				}

				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("func generateDirectorySnippet mismatch (-want, +got): %s\n", diff)
				}
			})
		}

		if tt.wantNoWildcard != nil {
			t.Run(tt.name+" no wildcard", func(t *testing.T) {
				got, gotErr := generateDirectorySnippet(tt.inDirectory, directories, nil, "Root", ExtraGenerators{}, "oc.", tt.inPathStructSuffix, false, tt.inSplitByModule, "", tt.inPackageName, tt.inPackageSuffix, nil, tt.inUnifiedPath, true, true, genutil.PreferOperationalState)
				if gotErr != nil {
					t.Fatalf("func generateDirectorySnippet, unexpected error: %v", gotErr)
				}

				for i, s := range got {
					got[i].ChildConstructors = trimDocComments(s.ChildConstructors)
					t.Log(got[i])
				}
				if diff := cmp.Diff(tt.wantNoWildcard, got); diff != "" {
					t.Errorf("func generateDirectorySnippet mismatch (-want, +got):%s\n", diff)
				}
			})
		}
	}
}

func TestGenerateChildConstructor(t *testing.T) {
	directories := getIR().Directories

	tests := []struct {
		name                    string
		inDirectory             *ygen.ParsedDirectory
		inDirectories           map[string]*ygen.ParsedDirectory
		inFieldName             string
		inUniqueFieldName       string
		inPathStructSuffix      string
		inGenerateWildcardPaths bool
		inUnifiedPaths          bool
		inChildAccessor         string
		testMethodDocComment    bool
		wantMethod              string
		// testMethodDocComment determines whether the doc comments for methods are tested.
		wantListBuilderAPI string
	}{{
		name:                    "container method",
		inDirectory:             directories["/root"],
		inDirectories:           directories,
		inFieldName:             "container",
		inUniqueFieldName:       "Container",
		inPathStructSuffix:      "Path",
		inGenerateWildcardPaths: true,
		wantMethod: `
func (n *RootPath) Container() *ContainerPath {
	ps := &ContainerPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"container"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
	}, {
		name:                    "container leaf method",
		inDirectory:             directories["/root-module/container"],
		inDirectories:           directories,
		inFieldName:             "leaf",
		inUniqueFieldName:       "Leaf",
		inPathStructSuffix:      "Path",
		inGenerateWildcardPaths: true,
		wantMethod: `
func (n *ContainerPath) Leaf() *Container_LeafPath {
	ps := &Container_LeafPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"leaf"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *ContainerPathAny) Leaf() *Container_LeafPathAny {
	ps := &Container_LeafPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"leaf"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
	}, {
		name:                    "unified container leaf method",
		inDirectory:             directories["/root-module/container"],
		inDirectories:           directories,
		inFieldName:             "leaf",
		inUniqueFieldName:       "Leaf",
		inPathStructSuffix:      "Path",
		inUnifiedPaths:          true,
		inGenerateWildcardPaths: true,
		wantMethod: `
func (n *ContainerPath) Leaf() *Container_LeafPath {
	ps := &Container_LeafPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"leaf"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

func (n *ContainerPathAny) Leaf() *Container_LeafPathAny {
	ps := &Container_LeafPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"leaf"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}
`,
	}, {
		name:                    "container leaf method without wildcard paths",
		inDirectory:             directories["/root-module/container"],
		inDirectories:           directories,
		inFieldName:             "leaf",
		inUniqueFieldName:       "Leaf",
		inPathStructSuffix:      "Path",
		inGenerateWildcardPaths: false,
		wantMethod: `
func (n *ContainerPath) Leaf() *Container_LeafPath {
	ps := &Container_LeafPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"leaf"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
	}, {
		name:                    "top-level leaf method",
		inDirectory:             directories["/root"],
		inDirectories:           directories,
		inFieldName:             "leaf",
		inUniqueFieldName:       "Leaf",
		inPathStructSuffix:      "Path",
		inGenerateWildcardPaths: true,
		wantMethod: `
func (n *RootPath) Leaf() *LeafPath {
	ps := &LeafPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"leaf"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
	}, {
		name:                    "container-with-config leaf-list method",
		inDirectory:             directories["/root-module/container-with-config"],
		inDirectories:           directories,
		inFieldName:             "leaflist",
		inUniqueFieldName:       "Leaflist",
		inPathStructSuffix:      "Path",
		inGenerateWildcardPaths: true,
		wantMethod: `
func (n *ContainerWithConfigPath) Leaflist() *ContainerWithConfig_LeaflistPath {
	ps := &ContainerWithConfig_LeaflistPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"state", "leaflist"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}

func (n *ContainerWithConfigPathAny) Leaflist() *ContainerWithConfig_LeaflistPathAny {
	ps := &ContainerWithConfig_LeaflistPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"state", "leaflist"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
	}, {
		name:                    "keyless list is skipped",
		inDirectory:             directories["/root"],
		inDirectories:           directories,
		inFieldName:             "keyless-list",
		inUniqueFieldName:       "KeylessList",
		inPathStructSuffix:      "Path",
		inGenerateWildcardPaths: true,
		testMethodDocComment:    true,
		wantMethod: `
// KeylessListAny (list): 
// 	Defining module:      "root-module"
// 	Instantiating module: "root-module"
// 	Path from parent:     "keyless-list-container/keyless-list"
// 	Path from root:       "/keyless-list-container/keyless-list"
func (n *RootPath) KeylessListAny() *KeylessListPathAny {
	ps := &KeylessListPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"keyless-list-container", "keyless-list"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
	}, {
		name:                    "list with state method",
		inDirectory:             directories["/root"],
		inDirectories:           directories,
		inFieldName:             "list-with-state",
		inUniqueFieldName:       "ListWithState",
		inPathStructSuffix:      "Path",
		inGenerateWildcardPaths: true,
		wantMethod: `
func (n *RootPath) ListWithStateAny() *ListWithStatePathAny {
	ps := &ListWithStatePathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state", "list-with-state"},
			map[string]interface{}{"key": "*"},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListWithState(Key float64) *ListWithStatePath {
	ps := &ListWithStatePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state", "list-with-state"},
			map[string]interface{}{"key": Key},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListWithStateMap() *ListWithStatePathMap {
	ps := &ListWithStatePathMap{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container-with-state"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
	}, {
		name:                    "root-level list methods over key threshold -- should use builder API",
		inDirectory:             directories["/root"],
		inDirectories:           directories,
		inFieldName:             "list",
		inUniqueFieldName:       "List",
		inPathStructSuffix:      "Path",
		inGenerateWildcardPaths: true,
		wantMethod: `
func (n *RootPath) ListAny() *ListPathAny {
	ps := &ListPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container", "list"},
			map[string]interface{}{"key1": "*", "key2": "*", "union-key": "*"},
			n,
		),
	}
	return ps
}

func (n *RootPath) List(Key1 string, Key2 oc.Binary, UnionKey oc.RootElementModule_List_UnionKey_Union) *ListPath {
	ps := &ListPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container", "list"},
			map[string]interface{}{"key1": Key1, "key2": Key2, "union-key": UnionKey},
			n,
		),
	}
	return ps
}

func (n *RootPath) ListMap() *ListPathMap {
	ps := &ListPathMap{
		NodePath: ygnmi.NewNodePath(
			[]string{"list-container"},
			map[string]interface{}{},
			n,
		),
	}
	return ps
}
`,
		wantListBuilderAPI: `
// WithKey1 sets ListPathAny's key "key1" to the specified value.
// Key1: string
func (n *ListPathAny) WithKey1(Key1 string) *ListPathAny {
	ygnmi.ModifyKey(n.NodePath, "key1", Key1)
	return n
}

// WithKey2 sets ListPathAny's key "key2" to the specified value.
// Key2: oc.Binary
func (n *ListPathAny) WithKey2(Key2 oc.Binary) *ListPathAny {
	ygnmi.ModifyKey(n.NodePath, "key2", Key2)
	return n
}

// WithUnionKey sets ListPathAny's key "union-key" to the specified value.
// UnionKey: [oc.UnionString, oc.Binary]
func (n *ListPathAny) WithUnionKey(UnionKey oc.RootElementModule_List_UnionKey_Union) *ListPathAny {
	ygnmi.ModifyKey(n.NodePath, "union-key", UnionKey)
	return n
}
`,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var methodBuf strings.Builder
			var builderBuf strings.Builder
			if errs := generateChildConstructors(&methodBuf, &builderBuf, tt.inDirectory, tt.inFieldName, tt.inUniqueFieldName, tt.inDirectories, "Root", "oc.", tt.inPathStructSuffix, tt.inGenerateWildcardPaths, tt.inChildAccessor, tt.inUnifiedPaths, true, genutil.PreferOperationalState, NodeDataMap{}, ExtraGenerators{}); errs != nil {
				t.Fatal(errs)
			}

			gotMethod := methodBuf.String()
			if !tt.testMethodDocComment {
				gotMethod = trimDocComments(gotMethod)
			}
			if got, want := gotMethod, tt.wantMethod; got != want {
				diff, _ := testutil.GenerateUnifiedDiff(want, got)
				t.Errorf("func generateChildConstructors methodBuf returned incorrect code, diff:\n%s", diff)
			}
			if got, want := builderBuf.String(), tt.wantListBuilderAPI; got != want {
				diff, _ := testutil.GenerateUnifiedDiff(want, got)
				t.Errorf("func generateChildConstructors builderBuf returned incorrect code, diff:\n%s", diff)
			}
		})
	}
}

func TestMakeKeyParams(t *testing.T) {
	tests := []struct {
		name             string
		inKeys           map[string]*ygen.ListKey
		inKeyNames       []string
		wantKeyParams    []keyParam
		wantErrSubstring string
	}{{
		name:       "empty listattr",
		inKeys:     nil,
		inKeyNames: nil,
	}, {
		name: "simple string param",
		inKeys: map[string]*ygen.ListKey{
			"fluorine": {
				Name: "Fluorine",
				LangType: &ygen.MappedType{
					NativeType: "string",
				},
			},
		},
		inKeyNames:    []string{"fluorine"},
		wantKeyParams: []keyParam{{name: "fluorine", varName: "Fluorine", typeName: "string", typeDocString: "string"}},
	}, {
		name: "simple int param, also testing camel-case",
		inKeys: map[string]*ygen.ListKey{
			"cl-cl": {
				Name: "ClCl",
				LangType: &ygen.MappedType{
					NativeType: "int",
				},
			},
		},
		inKeyNames:    []string{"cl-cl"},
		wantKeyParams: []keyParam{{name: "cl-cl", varName: "ClCl", typeName: "int", typeDocString: "int"}},
	}, {
		name: "name uniquification",
		inKeys: map[string]*ygen.ListKey{
			"cl-cl": {
				Name: "ClCl",
				LangType: &ygen.MappedType{
					NativeType: "int",
				},
			},
			"clCl": {
				Name: "ClCl",
				LangType: &ygen.MappedType{
					NativeType: "int",
				},
			},
		},
		inKeyNames: []string{"cl-cl", "clCl"},
		wantKeyParams: []keyParam{
			{name: "cl-cl", varName: "ClCl", typeName: "int", typeDocString: "int"},
			{name: "clCl", varName: "ClCl_", typeName: "int", typeDocString: "int"},
		},
	}, {
		name: "unsupported type",
		inKeys: map[string]*ygen.ListKey{
			"fluorine": {
				Name: "Fluorine",
				LangType: &ygen.MappedType{
					NativeType: "interface{}",
				},
			},
		},
		inKeyNames:    []string{"fluorine"},
		wantKeyParams: []keyParam{{name: "fluorine", varName: "Fluorine", typeName: "string", typeDocString: "string"}},
	}, {
		name: "keyElems doesn't match keys",
		inKeys: map[string]*ygen.ListKey{
			"neon": {
				Name: "Neon",
				LangType: &ygen.MappedType{
					NativeType: "light",
				},
			},
		},
		inKeyNames:       []string{"cl-cl"},
		wantErrSubstring: `key "cl-cl" doesn't exist in key map`,
	}, {
		name: "mappedType is nil",
		inKeys: map[string]*ygen.ListKey{
			"cl-cl": {
				Name:     "ClCl",
				LangType: nil,
			},
		},
		inKeyNames:       []string{"cl-cl"},
		wantErrSubstring: "mappedType for key is nil: cl-cl",
	}, {
		name: "multiple parameters",
		inKeys: map[string]*ygen.ListKey{
			"bromine": {
				Name: "Bromine",
				LangType: &ygen.MappedType{
					NativeType: "complex128",
				},
			},
			"cl-cl": {
				Name: "ClCl",
				LangType: &ygen.MappedType{
					NativeType: "int",
				},
			},
			"fluorine": {
				Name: "Fluorine",
				LangType: &ygen.MappedType{
					NativeType: "string",
				},
			},
			"iodine": {
				Name: "Iodine",
				LangType: &ygen.MappedType{
					NativeType: "float64",
				},
			},
		},
		inKeyNames: []string{"fluorine", "cl-cl", "bromine", "iodine"},
		wantKeyParams: []keyParam{
			{name: "fluorine", varName: "Fluorine", typeName: "string", typeDocString: "string"},
			{name: "cl-cl", varName: "ClCl", typeName: "int", typeDocString: "int"},
			{name: "bromine", varName: "Bromine", typeName: "complex128", typeDocString: "complex128"},
			{name: "iodine", varName: "Iodine", typeName: "float64", typeDocString: "float64"},
		},
	}, {
		name: "enumerated and union parameters",
		inKeys: map[string]*ygen.ListKey{
			"astatine": {
				Name: "Astatine",
				LangType: &ygen.MappedType{
					NativeType:        "Halogen",
					IsEnumeratedValue: true,
				},
			},
			"tennessine": {
				Name: "Tennessine",
				LangType: &ygen.MappedType{
					NativeType: "Ununseptium",
					UnionTypes: map[string]ygen.MappedUnionSubtype{
						"int32": {
							Index: 1,
						},
						"float64": {
							Index: 2,
						},
						"interface{}": {
							Index: 3,
						},
					},
				},
			},
		},
		inKeyNames: []string{"astatine", "tennessine"},
		wantKeyParams: []keyParam{
			{name: "astatine", varName: "Astatine", typeName: "oc.Halogen", typeDocString: "oc.Halogen"},
			{name: "tennessine", varName: "Tennessine", typeName: "oc.Ununseptium", typeDocString: "[oc.UnionInt32, oc.UnionFloat64, *oc.UnionUnsupported]"},
		},
	}, {
		name: "Binary and Empty",
		inKeys: map[string]*ygen.ListKey{
			"bromine": {
				Name: "Bromine",
				LangType: &ygen.MappedType{
					NativeType: "Binary",
				},
			},
			"cl-cl": {
				Name: "ClCl",
				LangType: &ygen.MappedType{
					NativeType: "YANGEmpty",
				},
			},
		},
		inKeyNames: []string{"cl-cl", "bromine"},
		wantKeyParams: []keyParam{
			{name: "cl-cl", varName: "ClCl", typeName: "oc.YANGEmpty", typeDocString: "oc.YANGEmpty"},
			{name: "bromine", varName: "Bromine", typeName: "oc.Binary", typeDocString: "oc.Binary"},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeyParams, err := makeKeyParams(tt.inKeys, tt.inKeyNames, "oc.")
			if diff := cmp.Diff(tt.wantKeyParams, gotKeyParams, cmp.AllowUnexported(keyParam{})); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}

			if diff := errdiff.Check(err, tt.wantErrSubstring); diff != "" {
				t.Errorf("func makeKeyParams, %v", diff)
			}
		})
	}
}

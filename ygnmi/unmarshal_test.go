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

package ygnmi

import (
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygot/ygot"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/testing/protocmp"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestUnmarshal(t *testing.T) {
	schemaStruct := testutil.GetSchemaStruct()
	superContainerSchema := schemaStruct().RootSchema().Dir["super-container"]

	passingTests := []struct {
		name                 string
		inData               []*DataPoint
		inQueryPath          *gpb.Path
		inStructSchema       *yang.Entry
		inStruct             ygot.ValidatedGoStruct
		inLeaf               bool
		inPreferShadowPath   bool
		wantUnmarshalledData []*DataPoint
		wantStruct           ygot.ValidatedGoStruct
	}{{
		name: "retrieve uint64",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 43}},
			Timestamp: time.Unix(1, 1),
		}},
		wantStruct: &testutil.LeafContainerStruct{Uint64Leaf: ygot.Uint64(43)},
	}, {
		name: "retrieve uint64 into fake root",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
		inStructSchema: schemaStruct().RootSchema(),
		inStruct:       &testutil.Device{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 43}},
			Timestamp: time.Unix(1, 1),
		}},
		wantStruct: &testutil.Device{SuperContainer: &testutil.SuperContainer{LeafContainerStruct: &testutil.LeafContainerStruct{Uint64Leaf: ygot.Uint64(43)}}},
	}, {
		name: "successfully retrieve uint64 with positive int",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 42}},
			Timestamp: time.Unix(1, 1),
		}},
		wantStruct: &testutil.LeafContainerStruct{Uint64Leaf: ygot.Uint64(42)},
	}, {
		name: "successfully retrieve uint64 with zero",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 0}},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 0}},
			Timestamp: time.Unix(1, 1),
		}},
		wantStruct: &testutil.LeafContainerStruct{Uint64Leaf: ygot.Uint64(0)},
	}, {
		name: "delete uint64",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{Uint64Leaf: ygot.Uint64(0)},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Timestamp: time.Unix(1, 1),
		}},
		wantStruct: &testutil.LeafContainerStruct{},
	}, {
		name: "retrieve union",
		inData: []*DataPoint{{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "aaaa",
				},
			},
			Timestamp: time.Unix(2, 2),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "aaaa",
				},
			},
			Timestamp: time.Unix(2, 2),
		}},
		wantStruct: &testutil.LeafContainerStruct{UnionLeaf: &testutil.UnionLeafType_String{String: "aaaa"}},
	}, {
		name: "delete union",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Timestamp: time.Unix(2, 2),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{UnionLeaf: &testutil.UnionLeafType_String{String: "forty two"}},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Timestamp: time.Unix(2, 2),
		}},
		wantStruct: &testutil.LeafContainerStruct{},
	}, {
		name: "delete union that's already deleted",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Timestamp: time.Unix(2, 2),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Timestamp: time.Unix(2, 2),
		}},
		wantStruct: &testutil.LeafContainerStruct{},
	}, {
		name: "retrieve union with a single enum inside",
		inData: []*DataPoint{{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf2"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "E_VALUE_FORTY_FOUR",
				},
			},
			Timestamp: time.Unix(2, 2),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf2"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf2"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "E_VALUE_FORTY_FOUR",
				},
			},
			Timestamp: time.Unix(2, 2),
		}},
		wantStruct: &testutil.LeafContainerStruct{UnionLeaf2: testutil.EnumType(44)},
	}, {
		name: "retrieve leaflist",
		inData: []*DataPoint{{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_LeaflistVal{
					LeaflistVal: &gpb.ScalarArray{
						Element: []*gpb.TypedValue{
							{Value: &gpb.TypedValue_StringVal{StringVal: "aaaaa"}},
							{Value: &gpb.TypedValue_StringVal{StringVal: "b"}},
						},
					},
				},
			},
			Timestamp: time.Unix(1, 3),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{UnionLeaf2: testutil.EnumType(44)},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_LeaflistVal{
					LeaflistVal: &gpb.ScalarArray{
						Element: []*gpb.TypedValue{
							{Value: &gpb.TypedValue_StringVal{StringVal: "aaaaa"}},
							{Value: &gpb.TypedValue_StringVal{StringVal: "b"}},
						},
					},
				},
			},
			Timestamp: time.Unix(1, 3),
		}},
		wantStruct: &testutil.LeafContainerStruct{
			UnionLeaf2:          testutil.EnumType(44),
			UnionLeafSingleType: []string{"aaaaa", "b"},
		},
	}, {
		name: "delete leaflist",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Timestamp: time.Unix(1, 3),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct: &testutil.LeafContainerStruct{
			UnionLeaf2:          testutil.EnumType(44),
			UnionLeafSingleType: []string{"forty two", "forty three"},
		},
		inLeaf: true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Timestamp: time.Unix(1, 3),
		}},
		wantStruct: &testutil.LeafContainerStruct{
			UnionLeaf2: testutil.EnumType(44),
		},
	}, {
		name: "retrieve leaf container, setting uint64 and enum",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(1, 1),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "E_VALUE_FORTY_THREE",
				},
			},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 42}},
			Timestamp: time.Unix(1, 1),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "E_VALUE_FORTY_THREE",
				},
			},
			Timestamp: time.Unix(1, 1),
		}},
		wantStruct: &testutil.LeafContainerStruct{
			Uint64Leaf: ygot.Uint64(42),
			EnumLeaf:   43,
		},
	}, {
		name: "retrieve leaf container, setting uint64, enum, union, union with single enum, and leaflist",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
			Timestamp: time.Unix(10, 10),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "E_VALUE_ONE_HUNDRED",
				},
			},
			Timestamp: time.Unix(10, 10),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_UintVal{
					UintVal: 100,
				},
			},
			Timestamp: time.Unix(10, 10),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf2"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "E_VALUE_ONE_HUNDRED",
				},
			},
			Timestamp: time.Unix(10, 10),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_LeaflistVal{
					LeaflistVal: &gpb.ScalarArray{
						Element: []*gpb.TypedValue{
							{Value: &gpb.TypedValue_StringVal{StringVal: "aa"}},
							{Value: &gpb.TypedValue_StringVal{StringVal: "bb"}},
						},
					},
				},
			},
			Timestamp: time.Unix(10, 10),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct: &testutil.LeafContainerStruct{
			Uint64Leaf: ygot.Uint64(42),
			EnumLeaf:   43,
		},
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 100}},
			Timestamp: time.Unix(10, 10),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "E_VALUE_ONE_HUNDRED",
				},
			},
			Timestamp: time.Unix(10, 10),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_UintVal{
					UintVal: 100,
				},
			},
			Timestamp: time.Unix(10, 10),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf2"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "E_VALUE_ONE_HUNDRED",
				},
			},
			Timestamp: time.Unix(10, 10),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_LeaflistVal{
					LeaflistVal: &gpb.ScalarArray{
						Element: []*gpb.TypedValue{
							{Value: &gpb.TypedValue_StringVal{StringVal: "aa"}},
							{Value: &gpb.TypedValue_StringVal{StringVal: "bb"}},
						},
					},
				},
			},
			Timestamp: time.Unix(10, 10),
		}},
		wantStruct: &testutil.LeafContainerStruct{
			Uint64Leaf:          ygot.Uint64(100),
			EnumLeaf:            testutil.EnumType(100),
			UnionLeaf:           &testutil.UnionLeafType_Uint32{Uint32: 100},
			UnionLeaf2:          testutil.EnumType(100),
			UnionLeafSingleType: []string{"aa", "bb"},
		},
	}, {
		name: "set union uint leaf with positive int value",
		inData: []*DataPoint{{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_IntVal{
					IntVal: 100,
				},
			},
			Timestamp: time.Unix(10, 10),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_UintVal{
					UintVal: 100,
				},
			},
			Timestamp: time.Unix(10, 10),
		}},
		wantStruct: &testutil.LeafContainerStruct{
			UnionLeaf: &testutil.UnionLeafType_Uint32{Uint32: 100},
		},
	}, {
		name:           "empty datapoint slice",
		inData:         []*DataPoint{},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantStruct:     &testutil.LeafContainerStruct{},
	}, {
		name:           "nil datapoint slice",
		inData:         nil,
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantStruct:     &testutil.LeafContainerStruct{},
	}, {
		name: "not all timestamps are the same -- the values should apply in order they are in the slice",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
			Timestamp: time.Unix(2, 1),
		}, {
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 200}},
			Timestamp: time.Unix(1, 0),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_LeaflistVal{
					LeaflistVal: &gpb.ScalarArray{
						Element: []*gpb.TypedValue{
							{Value: &gpb.TypedValue_StringVal{StringVal: "aa"}},
							{Value: &gpb.TypedValue_StringVal{StringVal: "bb"}},
						},
					},
				},
			},
			Timestamp: time.Unix(10, 10),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_LeaflistVal{
					LeaflistVal: &gpb.ScalarArray{
						Element: []*gpb.TypedValue{
							{Value: &gpb.TypedValue_StringVal{StringVal: "a"}},
							{Value: &gpb.TypedValue_StringVal{StringVal: "bbb"}},
						},
					},
				},
			},
			Timestamp: time.Unix(20, 20),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct"),
		inStructSchema: superContainerSchema,
		inStruct:       &testutil.SuperContainer{},
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 100}},
			Timestamp: time.Unix(2, 1),
		}, {
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 200}},
			Timestamp: time.Unix(1, 0),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_LeaflistVal{
					LeaflistVal: &gpb.ScalarArray{
						Element: []*gpb.TypedValue{
							{Value: &gpb.TypedValue_StringVal{StringVal: "aa"}},
							{Value: &gpb.TypedValue_StringVal{StringVal: "bb"}},
						},
					},
				},
			},
			Timestamp: time.Unix(10, 10),
		}, {
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_LeaflistVal{
					LeaflistVal: &gpb.ScalarArray{
						Element: []*gpb.TypedValue{
							{Value: &gpb.TypedValue_StringVal{StringVal: "a"}},
							{Value: &gpb.TypedValue_StringVal{StringVal: "bbb"}},
						},
					},
				},
			},
			Timestamp: time.Unix(20, 20),
		}},
		wantStruct: &testutil.SuperContainer{
			LeafContainerStruct: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(200),
				// TODO: If Collect calls are to be
				// implemented, then need to add tests for adds
				// and deletes to same and different children,
				// whether leaf or non-leaf, under a non-leaf.
				UnionLeafSingleType: []string{"a", "bbb"},
			},
		},
	}, {
		name: "retrieve single list key, setting the key value",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/key"),
		inStructSchema: superContainerSchema.Dir["model"].Dir["a"].Dir["single-key"],
		inStruct:       &testutil.Model_SingleKey{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(4, 4),
		}},
		wantStruct: &testutil.Model_SingleKey{Key: ygot.Int32(42)},
	}, {
		name: "retrieve single shadow-path list key, not setting the key value",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/key"),
		inStructSchema: superContainerSchema.Dir["model"].Dir["a"].Dir["single-key"],
		inStruct:       &testutil.Model_SingleKey{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(4, 4),
		}},
		wantStruct: &testutil.Model_SingleKey{},
	}, {
		name: "retrieve single shadow-path list key, setting the key value with preferShadowPath=true",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:        testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/key"),
		inStructSchema:     superContainerSchema.Dir["model"].Dir["a"].Dir["single-key"],
		inStruct:           &testutil.Model_SingleKey{},
		inLeaf:             true,
		inPreferShadowPath: true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(4, 4),
		}},
		wantStruct: &testutil.Model_SingleKey{Key: ygot.Int32(42)},
	}, {
		name: "retrieve single non-shadow-path list value, setting the value",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 4242}},
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
		inStructSchema: superContainerSchema.Dir["model"].Dir["a"].Dir["single-key"],
		inStruct:       &testutil.Model_SingleKey{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 4242}},
			Timestamp: time.Unix(4, 4),
		}},
		wantStruct: &testutil.Model_SingleKey{Value: ygot.Int64(4242)},
	}, {
		name: "retrieve single shadow-path list value, not setting the value",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/value"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 4242}},
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/value"),
		inStructSchema: superContainerSchema.Dir["model"].Dir["a"].Dir["single-key"],
		inStruct:       &testutil.Model_SingleKey{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/value"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 4242}},
			Timestamp: time.Unix(4, 4),
		}},
		wantStruct: &testutil.Model_SingleKey{},
	}, {
		name: "retrieve single non-shadow-path list value, not setting the value with preferShadowPath=true",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 4242}},
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:        testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
		inStructSchema:     superContainerSchema.Dir["model"].Dir["a"].Dir["single-key"],
		inStruct:           &testutil.Model_SingleKey{},
		inLeaf:             true,
		inPreferShadowPath: true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 4242}},
			Timestamp: time.Unix(4, 4),
		}},
		wantStruct: &testutil.Model_SingleKey{},
	}, {
		name: "retrieve single shadow-path list value, setting the value with preferShadowPath=true",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/value"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 4242}},
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:        testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/value"),
		inStructSchema:     superContainerSchema.Dir["model"].Dir["a"].Dir["single-key"],
		inStruct:           &testutil.Model_SingleKey{},
		inLeaf:             true,
		inPreferShadowPath: true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/state/value"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 4242}},
			Timestamp: time.Unix(4, 4),
		}},
		wantStruct: &testutil.Model_SingleKey{Value: ygot.Int64(4242)},
	}, {
		name: "retrieve entire list, setting multiple list keys",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(4, 4),
		}, {
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=43]/config/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/model/a/single-key"),
		inStructSchema: superContainerSchema.Dir["model"],
		inStruct:       &testutil.Model{},
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(4, 4),
		}, {
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=43]/config/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			Timestamp: time.Unix(4, 4),
		}},
		wantStruct: &testutil.Model{SingleKey: map[int32]*testutil.Model_SingleKey{
			42: {Key: ygot.Int32(42)},
			43: {Key: ygot.Int32(43)},
		}},
	}, {
		name: "delete a list key value, success",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
		inStructSchema: superContainerSchema.Dir["model"],
		inStruct: &testutil.Model{SingleKey: map[int32]*testutil.Model_SingleKey{
			42: {Key: ygot.Int32(42), Value: ygot.Int64(4242)},
		}},
		inLeaf: true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
			Timestamp: time.Unix(4, 4),
		}},
		wantStruct: &testutil.Model{SingleKey: map[int32]*testutil.Model_SingleKey{
			42: {Key: ygot.Int32(42)},
		}},
	}, {
		name: "delete a list key, no-op since preferShadowPath=true",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
		inStructSchema: superContainerSchema.Dir["model"],
		inStruct: &testutil.Model{SingleKey: map[int32]*testutil.Model_SingleKey{
			42: {Key: ygot.Int32(42), Value: ygot.Int64(4242)},
		}},
		inLeaf:             true,
		inPreferShadowPath: true,
		wantUnmarshalledData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=42]/config/value"),
			Timestamp: time.Unix(4, 4),
		}},
		wantStruct: &testutil.Model{SingleKey: map[int32]*testutil.Model_SingleKey{
			42: {Key: ygot.Int32(42), Value: ygot.Int64(4242)},
		}},
	}}

	for _, tt := range passingTests {
		t.Run(tt.name, func(t *testing.T) {
			unmarshalledData, complianceErrs, err := unmarshal(tt.inData, tt.inStructSchema, tt.inStruct, tt.inQueryPath, schemaStruct(), tt.inLeaf, tt.inPreferShadowPath, &opt{})
			if err != nil {
				t.Fatalf("unmarshal: got error, want none: %v", err)
			}
			if complianceErrs != nil {
				t.Fatalf("unmarshal: got compliance errors, want none: %v", complianceErrs)
			}

			if diff := cmp.Diff(tt.wantUnmarshalledData, unmarshalledData, protocmp.Transform()); diff != "" {
				t.Errorf("unmarshal: successfully unmarshalled datapoints do not match (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantStruct, tt.inStruct); diff != "" {
				t.Errorf("unmarshal: struct after unmarshalling does not match (-want +got):\n%s", diff)
			}
		})
	}

	failingTests := []struct {
		name                  string
		inData                []*DataPoint
		inQueryPath           *gpb.Path
		inStructSchema        *yang.Entry
		inStruct              ygot.ValidatedGoStruct
		inLeaf                bool
		inPreferShadowPath    bool
		wantUnmarshalledData  []*DataPoint
		wantStruct            ygot.ValidatedGoStruct
		wantErrSubstr         string
		wantPathErrSubstr     *TelemetryError
		wantTypeErrSubstr     *TelemetryError
		wantValidateErrSubstr string
	}{{
		name: "fail to retrieve uint64 due to wrong type",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantTypeErrSubstr: &TelemetryError{
			Path:  testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			Err:   errors.New("failed to unmarshal"),
		},
	}, {
		name: "multiple datapoints for leaf node",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 42}},
			Timestamp: time.Unix(1, 1),
		}, {
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 43}},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantPathErrSubstr: &TelemetryError{
			Err: errors.New("got multiple"),
		},
	}, {
		name: "failed to retrieve uint64 with negative int",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: -42}},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantTypeErrSubstr: &TelemetryError{
			Path:  testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: -42}},
			Err:   errors.New("failed to unmarshal"),
		},
	}, {
		name: "fail to retrieve uint64 due to wrong path",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/xxxxxxxxx/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantPathErrSubstr: &TelemetryError{
			Path:  testutil.GNMIPath(t, "super-container/xxxxxxxxx/uint64-leaf"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			Err:   errors.New(`does not match the query path "/super-container/leaf-container-struct/uint64-leaf"`),
		},
	}, {
		name: "retrieve union with field that doesn't match regex",
		inData: []*DataPoint{{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "forty two",
				},
			},
			Timestamp: time.Unix(2, 2),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantUnmarshalledData: []*DataPoint{{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
			Value: &gpb.TypedValue{
				Value: &gpb.TypedValue_StringVal{
					StringVal: "forty two",
				},
			},
			Timestamp: time.Unix(2, 2),
		}},
		wantValidateErrSubstr: "does not match regular expression pattern",
	}, {
		name: "delete at a non-existent path",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/dne"),
			Timestamp: time.Unix(2, 2),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantPathErrSubstr: &TelemetryError{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/dne"),
			Err:  errors.New("does not match the query path"),
		},
	}, {
		name: "fail to delete union with a single enum inside due to wrong path prefix",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "not-valid-prefix/leaf-container-struct/union-leaf2"),
			Timestamp: time.Unix(2, 2),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf2"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{UnionLeaf2: testutil.EnumType(44)},
		inLeaf:         true,
		wantPathErrSubstr: &TelemetryError{
			Path: testutil.GNMIPath(t, "not-valid-prefix/leaf-container-struct/union-leaf2"),
			Err:  errors.New(`does not match the query path "/super-container/leaf-container-struct/union-leaf2"`),
		},
	}, {
		name: "fail to delete union with a single enum inside due to wrong path suffix",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/leaf-container-struct/union-needle2"),
			Timestamp: time.Unix(2, 2),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/union-leaf2"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{UnionLeaf2: testutil.EnumType(44)},
		inLeaf:         true,
		wantPathErrSubstr: &TelemetryError{
			Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/union-needle2"),
			Err:  errors.New(`does not match the query path "/super-container/leaf-container-struct/union-leaf2"`),
		},
	}, {
		name: "retrieve a list value with an invalid list key",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=forty-four]/config/key"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Timestamp: time.Unix(4, 4),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/model/a/single-key[key=44]/config/key"),
		inStructSchema: superContainerSchema.Dir["model"].Dir["a"].Dir["single-key"],
		inStruct:       &testutil.Model_SingleKey{},
		inLeaf:         true,
		wantPathErrSubstr: &TelemetryError{
			Path:  testutil.GNMIPath(t, "super-container/model/a/single-key[key=forty-four]/config/key"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			Err:   errors.New(`does not match the query path "/super-container/model/a/single-key[key=44]/config/key"`),
		},
	}, {
		name: "invalid input: parent schema is not parent of input data's path.",
		inData: []*DataPoint{{
			Path:      testutil.GNMIPath(t, "different-container/leaf-container-struct/uint64-leaf"),
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantPathErrSubstr: &TelemetryError{
			Path:  testutil.GNMIPath(t, "different-container/leaf-container-struct/uint64-leaf"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			Err:   errors.New(`does not match the query path "/super-container/leaf-container-struct/uint64-leaf"`),
		},
	}, {
		name: "invalid input: deprecated elements field in path",
		inData: []*DataPoint{{
			Path: &gpb.Path{
				Element: []string{"super-container", "model", "a", "single-key", "forty-four", "config", "key"},
			},
			Value:     &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			Timestamp: time.Unix(1, 1),
		}},
		inQueryPath:    testutil.GNMIPath(t, "super-container/model/a/single-key[key=44]/config/key"),
		inStructSchema: superContainerSchema.Dir["leaf-container-struct"],
		inStruct:       &testutil.LeafContainerStruct{},
		inLeaf:         true,
		wantPathErrSubstr: &TelemetryError{
			Path: &gpb.Path{
				Element: []string{"super-container", "model", "a", "single-key", "forty-four", "config", "key"},
			},
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			Err:   errors.New(`path uses deprecated and unsupported Element field`),
		},
	}}

	for _, tt := range failingTests {
		t.Run(tt.name, func(t *testing.T) {
			unmarshalledData, complianceErrs, errs := unmarshal(tt.inData, tt.inStructSchema, tt.inStruct, tt.inQueryPath, schemaStruct(), tt.inLeaf, tt.inPreferShadowPath, &opt{})
			if errs != nil {
				t.Fatalf("unmarshal: got more than one error: %v", errs)
			}
			if diff := cmp.Diff(tt.wantUnmarshalledData, unmarshalledData, protocmp.Transform()); diff != "" {
				t.Errorf("unmarshal: successfully unmarshalled datapoints do not match (-want +got):\n%s", diff)
			}

			var pathErrs, typeErrs []*TelemetryError
			var validateErrs []error
			if complianceErrs != nil {
				pathErrs = complianceErrs.PathErrors
				typeErrs = complianceErrs.TypeErrors
				validateErrs = complianceErrs.ValidateErrors
				// Validate documentation on error
				if !strings.Contains(complianceErrs.String(), "https://github.com/openconfig/ygnmi#noncompliance-errors") {
					t.Errorf("ComplianceError String() didn't contain expected reference to documentation.")
				}
			}
			if len(pathErrs) > 1 {
				t.Fatalf("unmarshal: got more than one path unmarshal error: %v", pathErrs)
			}
			if len(typeErrs) > 1 {
				t.Fatalf("unmarshal: got more than one type unmarshal error: %v", typeErrs)
			}
			if len(validateErrs) > 1 {
				t.Fatalf("unmarshal: got more than one validate error: %v", validateErrs)
			}

			// Populate errors for validation.
			var err, validateErr error
			var pathErr, typeErr *TelemetryError
			if len(pathErrs) == 1 {
				pathErr = pathErrs[0]
			}
			if len(typeErrs) == 1 {
				typeErr = typeErrs[0]
			}
			if len(validateErrs) == 1 {
				validateErr = validateErrs[0]
			}

			// Validate expected errors
			if diff := errdiff.Substring(err, tt.wantErrSubstr); diff != "" {
				t.Fatalf("unmarshal: did not get expected error substring:\n%s", diff)
			}

			verifyTelemetryError := func(t *testing.T, gotErr, wantErrSubstr *TelemetryError) {
				t.Helper()
				// Only do exact verification on the Path and Value fields of the Telemetry errors.
				if diff := cmp.Diff(wantErrSubstr, gotErr, protocmp.Transform(), cmp.FilterPath(
					func(p cmp.Path) bool {
						return p.String() == "Err"
					},
					cmp.Ignore(),
				)); diff != "" {
					t.Fatalf("unmarshal: did not get expected path compliance error (-want, +got):\n%s", diff)
				}
				if gotErr != nil && wantErrSubstr != nil {
					if diff := errdiff.Substring(gotErr.Err, wantErrSubstr.Err.Error()); diff != "" {
						t.Fatalf("unmarshal: did not get expected compliance error substring:\n%s", diff)
					}
				}
			}
			verifyTelemetryError(t, pathErr, tt.wantPathErrSubstr)
			verifyTelemetryError(t, typeErr, tt.wantTypeErrSubstr)

			if diff := errdiff.Substring(validateErr, tt.wantValidateErrSubstr); diff != "" {
				t.Fatalf("unmarshal: did not get expected validateErr substring:\n%s", diff)
			}
		})
	}
}

func TestLatestTimestamp(t *testing.T) {
	tests := []struct {
		desc     string
		in       []*DataPoint
		want     time.Time
		wantRecv time.Time
	}{{
		desc: "basic",
		in: []*DataPoint{{
			Path:          testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:         &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
			Timestamp:     time.Unix(0, 1),
			RecvTimestamp: time.Unix(5, 5),
		}, {
			Path:          testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			Value:         &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 200}},
			Timestamp:     time.Unix(3, 3),
			RecvTimestamp: time.Unix(4, 4),
		}, {
			Path:          testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Value:         &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 300}},
			Timestamp:     time.Unix(2, 2),
			RecvTimestamp: time.Unix(3, 3),
		}, {
			Path:          testutil.GNMIPath(t, "super-container/leaf-container-struct/union-stleaflist"),
			Value:         &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 400}},
			Timestamp:     time.Unix(1, 1),
			RecvTimestamp: time.Unix(2, 2),
		}},
		want:     time.Unix(3, 3),
		wantRecv: time.Unix(5, 5),
	}, {
		desc:     "empty list",
		want:     time.Time{},
		wantRecv: time.Time{},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got := LatestTimestamp(tt.in); !got.Equal(tt.want) {
				t.Errorf("LatestTimestamp: got %v, want %v", got, tt.want)
			}
			if got := LatestRecvTimestamp(tt.in); !got.Equal(tt.wantRecv) {
				t.Errorf("LatestRecvTimestamp: got %v, want %v", got, tt.wantRecv)
			}
		})
	}
}

func TestBundleDatapoints(t *testing.T) {
	tests := []struct {
		desc         string
		inDatapoints []*DataPoint
		inPrefixLen  int
		want         map[string][]*DataPoint
		wantPrefixes []string
		wantErr      bool
	}{{
		desc: "leaf-paths",
		inDatapoints: []*DataPoint{{
			Path:  testutil.GNMIPath(t, "alpha[key=un]/bravo/leaf0"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha[key=deux]/bravo/leaf0"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 200}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha[key=un]/bravo/leaf0/leaf-under-leaf"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 300}},
		}},
		inPrefixLen: 3,
		want: map[string][]*DataPoint{
			"/alpha[key=un]/bravo/leaf0": {{
				Path:  testutil.GNMIPath(t, "alpha[key=un]/bravo/leaf0"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
			}, {
				Path:  testutil.GNMIPath(t, "alpha[key=un]/bravo/leaf0/leaf-under-leaf"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 300}},
			}},
			"/alpha[key=deux]/bravo/leaf0": {{
				Path:  testutil.GNMIPath(t, "alpha[key=deux]/bravo/leaf0"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 200}},
			}},
		},
		wantPrefixes: []string{
			"/alpha[key=deux]/bravo/leaf0",
			"/alpha[key=un]/bravo/leaf0",
		},
	}, {
		desc: "non-leaf-paths",
		inDatapoints: []*DataPoint{{
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=un]/leaf0"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=deux]/leaf2"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 300}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=trois]/leaf3"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 400}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=trois]/leaf4"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 500}},
		}, {
			// duplicate path
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=trois]/leaf3"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 1000}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=deux]/foo/leaf5"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 600}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=un]/foo/bar/leaf6"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 700}},
		}},
		inPrefixLen: 2,
		want: map[string][]*DataPoint{
			"/alpha/bravo[key=un]": {{
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=un]/leaf0"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
			}, {
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=un]/foo/bar/leaf6"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 700}},
			}},
			"/alpha/bravo[key=deux]": {{
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=deux]/leaf2"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 300}},
			}, {
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=deux]/foo/leaf5"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 600}},
			}},
			"/alpha/bravo[key=trois]": {{
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=trois]/leaf3"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 400}},
			}, {
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=trois]/leaf4"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 500}},
			}, {
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=trois]/leaf3"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 1000}},
			}},
		},
		wantPrefixes: []string{
			"/alpha/bravo[key=deux]",
			"/alpha/bravo[key=trois]",
			"/alpha/bravo[key=un]",
		},
	}, {
		desc: "path-shorter-than-prefixLen",
		inDatapoints: []*DataPoint{{
			Path:  testutil.GNMIPath(t, "alpha/"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=un]/leaf0"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=deux]/leaf2"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 300}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=trois]/leaf3"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 400}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=trois]/leaf4"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 500}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=deux]/foo/leaf5"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 600}},
		}, {
			Path:  testutil.GNMIPath(t, "alpha/bravo[key=un]/foo/bar/leaf6"),
			Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 700}},
		}},
		inPrefixLen: 2,
		want: map[string][]*DataPoint{
			"/": {{
				Path:  testutil.GNMIPath(t, "alpha/"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
			}},
			"/alpha/bravo[key=un]": {{
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=un]/leaf0"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
			}, {
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=un]/foo/bar/leaf6"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 700}},
			}},
			"/alpha/bravo[key=deux]": {{
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=deux]/leaf2"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 300}},
			}, {
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=deux]/foo/leaf5"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 600}},
			}},
			"/alpha/bravo[key=trois]": {{
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=trois]/leaf3"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 400}},
			}, {
				Path:  testutil.GNMIPath(t, "alpha/bravo[key=trois]/leaf4"),
				Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 500}},
			}},
		},
		wantPrefixes: []string{
			"/",
			"/alpha/bravo[key=deux]",
			"/alpha/bravo[key=trois]",
			"/alpha/bravo[key=un]",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, gotPrefixes, err := bundleDatapoints(tt.inDatapoints, tt.inPrefixLen)
			if gotErr := err != nil; gotErr != tt.wantErr {
				t.Errorf("Got error: %v, want error: %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Errorf("Datapoint groups (-got, +want):\n%s", diff)
			}
			if diff := cmp.Diff(gotPrefixes, tt.wantPrefixes); diff != "" {
				t.Errorf("Prefixes (-got, +want):\n%s", diff)
			}
		})
	}
}

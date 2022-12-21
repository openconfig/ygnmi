// Copyright 2022 Google LLC
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

package testutil

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
	"github.com/pkg/errors"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

var globalEnumMap = map[string]map[int64]ygot.EnumDefinition{
	"EnumType": {
		43:  {Name: "E_VALUE_FORTY_THREE"},
		44:  {Name: "E_VALUE_FORTY_FOUR"},
		100: {Name: "E_VALUE_ONE_HUNDRED"},
	},
	"EnumType2": {
		42: {Name: "E_VALUE_FORTY_TWO"},
	},
}

// EnumType is used as an enum type in various tests in the ytypes package.
type EnumType int64

func (EnumType) ΛMap() map[string]map[int64]ygot.EnumDefinition {
	return globalEnumMap
}

func (EnumType) IsYANGGoEnum() {}

func (e EnumType) String() string {
	return ygot.EnumLogString(e, int64(e), "EnumType")
}

// EnumType2 is used as an enum type in various tests in the ytypes package.
type EnumType2 int64

func (EnumType2) ΛMap() map[string]map[int64]ygot.EnumDefinition {
	return globalEnumMap
}

func (EnumType2) IsYANGGoEnum() {}

type UnionLeafType interface {
	Is_UnionLeafType()
}

type UnionLeafType_String struct {
	String string
}

func (*UnionLeafType_String) Is_UnionLeafType() {}

type UnionLeafType_Uint32 struct {
	Uint32 uint32
}

func (*UnionLeafType_Uint32) Is_UnionLeafType() {}

type UnionLeafType_EnumType struct {
	EnumType EnumType
}

func (*UnionLeafType_EnumType) Is_UnionLeafType() {}

func (*UnionLeafType_EnumType) ΛMap() map[string]map[int64]ygot.EnumDefinition {
	return globalEnumMap
}

type UnionLeafType_EnumType2 struct {
	EnumType2 EnumType2
}

func (*UnionLeafType_EnumType2) Is_UnionLeafType() {}

func (*UnionLeafType_EnumType2) ΛMap() map[string]map[int64]ygot.EnumDefinition {
	return globalEnumMap
}

// Device is the fake root.
type Device struct {
	SuperContainer *SuperContainer `path:"super-container" module:"yang-module"`
}

func (*Device) IsYANGGoStruct()                              {}
func (*Device) Validate(opts ...ygot.ValidationOption) error { return nil }
func (*Device) ΛBelongingModule() string                     { return "" }
func (*Device) ΛEnumTypeMap() map[string][]reflect.Type      { return nil }

func unmarshalFunc([]byte, ygot.GoStruct, ...ytypes.UnmarshalOpt) error { return nil }

type SuperContainer struct {
	LeafContainerStruct *LeafContainerStruct `path:"leaf-container-struct"`
	Model               *Model               `path:"model"`
}

func (*SuperContainer) IsYANGGoStruct()                              {}
func (*SuperContainer) Validate(opts ...ygot.ValidationOption) error { return nil }
func (*SuperContainer) ΛBelongingModule() string                     { return "" }
func (*SuperContainer) ΛEnumTypeMap() map[string][]reflect.Type      { return nil }

type LeafContainerStruct struct {
	Uint64Leaf          *uint64       `path:"uint64-leaf" shadow-path:"state/uint64-leaf"`
	EnumLeaf            EnumType      `path:"enum-leaf"`
	UnionLeaf           UnionLeafType `path:"union-leaf"`
	UnionLeaf2          EnumType      `path:"union-leaf2"`
	UnionLeafSingleType []string      `path:"union-stleaflist"`
}

func (*LeafContainerStruct) IsYANGGoStruct()                              {}
func (*LeafContainerStruct) Validate(opts ...ygot.ValidationOption) error { return nil }
func (*LeafContainerStruct) ΛBelongingModule() string                     { return "" }

type Model struct {
	SingleKey map[int32]*Model_SingleKey `path:"a/single-key"`
}

func (*Model) IsYANGGoStruct()                              {}
func (*Model) Validate(opts ...ygot.ValidationOption) error { return nil }
func (*Model) ΛBelongingModule() string                     { return "" }
func (*Model) ΛEnumTypeMap() map[string][]reflect.Type      { return nil }

type Model_SingleKey struct {
	Key   *int32 `path:"config/key|key" shadow-path:"state/key"`
	Value *int64 `path:"config/value" shadow-path:"state/value"`
}

func (*Model_SingleKey) IsYANGGoStruct()                              {}
func (*Model_SingleKey) Validate(opts ...ygot.ValidationOption) error { return nil }
func (*Model_SingleKey) ΛBelongingModule() string                     { return "" }
func (*Model_SingleKey) ΛEnumTypeMap() map[string][]reflect.Type      { return nil }

// NewSingleKey is a *generated* method for Model which may be used by an
// unmarshal function in ytype's reflect library, and is kept here in case.
func (t *Model) NewSingleKey(Key int32) (*Model_SingleKey, error) {
	// Initialise the list within the receiver struct if it has not already been
	// created.
	if t.SingleKey == nil {
		t.SingleKey = make(map[int32]*Model_SingleKey)
	}

	key := Key

	// Ensure that this key has not already been used in the
	// list. Keyed YANG lists do not allow duplicate keys to
	// be created.
	if _, ok := t.SingleKey[key]; ok {
		return nil, fmt.Errorf("duplicate key %q for list SingleKey", key)
	}

	t.SingleKey[key] = &Model_SingleKey{
		Key: &Key,
	}

	return t.SingleKey[key], nil
}

// RenameSingleKey is a *generated* method for Model which may be used by an
// unmarshal function in ytype's reflect library, and is kept here in case.
// RenameSingleKey renames an entry in the list SingleKey within
// the Model struct. The entry with key oldK is renamed to newK updating
// the key within the value.
func (t *Model) RenameSingleKey(oldK, newK int32) error {
	if _, ok := t.SingleKey[newK]; ok {
		return fmt.Errorf("key %q already exists in SingleKey", newK)
	}

	e, ok := t.SingleKey[oldK]
	if !ok {
		return fmt.Errorf("key %q not found in SingleKey", oldK)
	}
	e.Key = &newK

	t.SingleKey[newK] = e
	delete(t.SingleKey, oldK)
	return nil
}

// ΛListKeyMap returns the keys of the Model_SingleKey struct, which is a YANG list entry.
func (t *Model_SingleKey) ΛListKeyMap() (map[string]interface{}, error) {
	if t.Key == nil {
		return nil, fmt.Errorf("nil value for key Key")
	}

	return map[string]interface{}{
		"key": *t.Key,
	}, nil
}

func (*LeafContainerStruct) ΛEnumTypeMap() map[string][]reflect.Type {
	return map[string][]reflect.Type{
		"/super-container/leaf-container-struct/enum-leaf":   {reflect.TypeOf(EnumType(0))},
		"/super-container/leaf-container-struct/union-leaf":  {reflect.TypeOf(EnumType(0)), reflect.TypeOf(EnumType2(0))},
		"/super-container/leaf-container-struct/union-leaf2": {reflect.TypeOf(EnumType(0))},
	}
}

func (*LeafContainerStruct) To_UnionLeafType(i interface{}) (UnionLeafType, error) {
	switch v := i.(type) {
	case string:
		return &UnionLeafType_String{v}, nil
	case uint32:
		return &UnionLeafType_Uint32{v}, nil
	case EnumType:
		return &UnionLeafType_EnumType{v}, nil
	case EnumType2:
		return &UnionLeafType_EnumType2{v}, nil
	default:
		return nil, errors.Errorf("cannot convert %v to To_UnionLeafType, unknown union type, got: %T, want any of [string, uint32]", i, i)
	}
}

// addParents adds parent pointers for a schema tree.
func addParents(e *yang.Entry) {
	for _, c := range e.Dir {
		c.Parent = e
		addParents(c)
	}
}

// mustPath converts a string to its path proto.
func GNMIPath(t testing.TB, s string) *gpb.Path {
	p, err := ygot.StringToStructuredPath(s)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: remove when fixed https://github.com/openconfig/ygot/issues/615
	if p.Origin == "" && (len(p.Elem) == 0 || p.Elem[0].Name != "meta") {
		p.Origin = "openconfig"
	}
	return p
}

func GetSchemaStruct() func() *ytypes.Schema {
	rootSchema := &yang.Entry{
		Name:       "device",
		Kind:       yang.DirectoryEntry,
		Annotation: map[string]interface{}{"isFakeRoot": true},
		Dir: map[string]*yang.Entry{
			"super-container": {
				Name: "super-container",
				Kind: yang.DirectoryEntry,
				Dir: map[string]*yang.Entry{
					"model": {
						Name: "model",
						Kind: yang.DirectoryEntry,
						Dir: map[string]*yang.Entry{
							"a": {
								Name: "a",
								Kind: yang.DirectoryEntry,
								Dir: map[string]*yang.Entry{
									"single-key": {
										Name:     "single-key",
										Kind:     yang.DirectoryEntry,
										ListAttr: yang.NewDefaultListAttr(),
										Key:      "key",
										Dir: map[string]*yang.Entry{
											"key": {
												Name: "key",
												Kind: yang.LeafEntry,
												Type: &yang.YangType{Kind: yang.Yleafref, Path: "../config/key"},
											},
											"config": {
												Name: "config",
												Kind: yang.DirectoryEntry,
												Dir: map[string]*yang.Entry{
													"key": {
														Name: "key",
														Kind: yang.LeafEntry,
														Type: &yang.YangType{Kind: yang.Yint32},
													},
													"value": {
														Name: "value",
														Kind: yang.LeafEntry,
														Type: &yang.YangType{Kind: yang.Yint64},
													},
												},
											},
											"state": {
												Name: "state",
												Kind: yang.DirectoryEntry,
												Dir: map[string]*yang.Entry{
													"key": {
														Name: "key",
														Kind: yang.LeafEntry,
														Type: &yang.YangType{Kind: yang.Yint32},
													},
													"value": {
														Name: "value",
														Kind: yang.LeafEntry,
														Type: &yang.YangType{Kind: yang.Yint64},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"leaf-container-struct": {
						Name: "leaf-container-struct",
						Kind: yang.DirectoryEntry,
						Dir: map[string]*yang.Entry{
							"state": {
								Name: "state",
								Kind: yang.DirectoryEntry,
								Dir: map[string]*yang.Entry{
									"uint64-leaf": {
										Kind: yang.LeafEntry,
										Type: &yang.YangType{
											Kind: yang.Yuint64,
										},
									},
								},
							},
							"uint64-leaf": {
								Name: "uint64-leaf",
								Kind: yang.LeafEntry,
								Type: &yang.YangType{
									Kind: yang.Yuint64,
								},
							},
							"enum-leaf": {
								Name: "enum-leaf",
								Kind: yang.LeafEntry,
								Type: &yang.YangType{
									Kind: yang.Yenum,
								},
							},
							"union-leaf": {
								Name: "union-leaf",
								Kind: yang.LeafEntry,
								Type: &yang.YangType{
									Kind: yang.Yunion,
									Type: []*yang.YangType{
										{
											Kind:    yang.Ystring,
											Pattern: []string{"a+"},
										},
										{
											Kind: yang.Yuint32,
										},
										{
											Kind: yang.Yenum,
										},
										{
											Kind: yang.Yenum,
										},
										{
											Kind: yang.Yleafref,
											Path: "../enum-leaf",
										},
									},
								},
							},
							"union-leaf2": {
								Name: "union-leaf2",
								Kind: yang.LeafEntry,
								Type: &yang.YangType{
									Kind: yang.Yunion,
									Type: []*yang.YangType{
										{
											Kind: yang.Yenum,
										},
									},
								},
							},
							"union-stleaflist": {
								Name:     "union-stleaflist",
								Kind:     yang.LeafEntry,
								ListAttr: yang.NewDefaultListAttr(),
								Type: &yang.YangType{
									Kind: yang.Yunion,
									Type: []*yang.YangType{
										{
											// Note that Validate is not called as part of Unmarshal,
											// therefore any string pattern will actually match.
											Kind:    yang.Ystring,
											Pattern: []string{"a+"},
										},
										{
											Kind:    yang.Ystring,
											Pattern: []string{"b+"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	addParents(rootSchema)
	schemaStruct := func() *ytypes.Schema {
		return &ytypes.Schema{
			Root: &Device{},
			SchemaTree: map[string]*yang.Entry{
				"Device":                rootSchema,
				"super-container":       rootSchema.Dir["super-container"],
				"leaf-container-struct": rootSchema.Dir["super-container"].Dir["leaf-container-struct"],
				"Model_SingleKey":       rootSchema.Dir["super-container"].Dir["model"].Dir["a"].Dir["single-key"],
			},
			Unmarshal: unmarshalFunc,
		}
	}
	return schemaStruct
}

// Code generated from ygnmi_test.go. DO NOT EDIT.

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

package ygnmi_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygnmi/internal/exampleocconfig"
	"github.com/openconfig/ygnmi/internal/exampleocconfig/exampleocconfigpath"
	"github.com/openconfig/ygnmi/internal/gnmitestutil"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygnmi/schemaless"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

func getSamplePreferConfigOrderedMap(t *testing.T) *exampleocconfig.Model_SingleKey_OrderedList_OrderedMap {
	om := &exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{}
	ol, err := om.AppendNew("foo")
	if err != nil {
		t.Fatal(err)
	}
	ol.SetValue(42)
	ol, err = om.AppendNew("bar")
	if err != nil {
		t.Fatal(err)
	}
	ol.SetValue(43)
	ol, err = om.AppendNew("baz")
	if err != nil {
		t.Fatal(err)
	}
	ol.SetValue(44)
	return om
}

func getSamplePreferConfigOrderedMapIncomplete(t *testing.T) *exampleocconfig.Model_SingleKey_OrderedList_OrderedMap {
	om := &exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{}
	ol, err := om.AppendNew("foo")
	if err != nil {
		t.Fatal(err)
	}
	ol.SetValue(42)
	ol, err = om.AppendNew("bar")
	if err != nil {
		t.Fatal(err)
	}
	ol.SetValue(43)
	return om
}

func getSamplePreferConfigSingleKeyedMap(t *testing.T) map[string]*exampleocconfig.Model_SingleKey {
	model := &exampleocconfig.Model{}
	model.GetOrCreateSingleKey("foo").SetValue(42)
	model.GetOrCreateSingleKey("bar").SetValue(43)
	model.GetOrCreateSingleKey("baz").SetValue(44)
	return model.SingleKey
}

func getSamplePreferConfigSingleKeyedMapIncomplete(t *testing.T) map[string]*exampleocconfig.Model_SingleKey {
	model := &exampleocconfig.Model{}
	model.GetOrCreateSingleKey("foo").SetValue(42)
	model.GetOrCreateSingleKey("bar").SetValue(43)
	return model.SingleKey
}

func getSamplePreferConfigInnerSingleKeyedMap(t *testing.T) map[string]*exampleocconfig.Model_SingleKey_SingleKey {
	sk := &exampleocconfig.Model_SingleKey{}
	sk.GetOrCreateSingleKey("foo").SetValue(42)
	sk.GetOrCreateSingleKey("bar").SetValue(43)
	sk.GetOrCreateSingleKey("baz").SetValue(44)
	return sk.SingleKey
}

func getSamplePreferConfigInnerSingleKeyedMapIncomplete(t *testing.T) map[string]*exampleocconfig.Model_SingleKey_SingleKey {
	sk := &exampleocconfig.Model_SingleKey{}
	sk.GetOrCreateSingleKey("foo").SetValue(42)
	sk.GetOrCreateSingleKey("bar").SetValue(43)
	return sk.SingleKey
}

func TestPreferConfigLookup(t *testing.T) {
	fakeGNMI, c := newClient(t)
	leafPath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := exampleocconfigpath.Root().RemoteContainer().ALeaf().State()

	leafTests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		inQuery              ygnmi.SingletonQuery[string]
		wantRequestValues    *ygnmi.RequestValues
		wantSubscriptionPath *gpb.Path
		wantVal              *ygnmi.Value[string]
		wantErr              string
	}{{
		desc:    "success update and sync",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  false,
			ConfigFiltered: true,
		},
		wantSubscriptionPath: leafPath,
		wantVal: (&ygnmi.Value[string]{
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal("foo"),
	}, {
		desc:    "success update and no sync",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			})
		},
		wantSubscriptionPath: leafPath,
		wantVal: (&ygnmi.Value[string]{
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal("foo"),
	}, {
		desc:    "success with prefix",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Prefix:    testutil.GNMIPath(t, "remote-container"),
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "state/a-leaf"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: leafPath,
		wantVal: (&ygnmi.Value[string]{
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal("foo"),
	}, {
		desc:    "success multiple notifs and first no value",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Update: []*gpb.Update{},
			}).Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: leafPath,
		wantVal: (&ygnmi.Value[string]{
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal("foo"),
	}, {
		desc:    "success no value",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		wantSubscriptionPath: leafPath,
		wantVal: (&ygnmi.Value[string]{
			Path: leafPath,
		}),
	}, {
		desc:    "error multiple values",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Sync()
		},
		wantErr: "noncompliant data encountered while unmarshalling leaf",
	}, {
		desc:    "error deprecated path",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: &gpb.Path{
						Element: []string{"super-container", "leaf-container-struct", "uint64-leaf"},
					},
					Val: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Sync()
		},
		wantErr: "noncompliant data encountered while unmarshalling leaf",
	}, {
		desc:    "error last path element wrong",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "leaf-container-struct/enum-leaf"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "E_VALUE_FORTY_THREE"}},
				}},
			}).Sync()
		},
		wantErr: "noncompliant data encountered while unmarshalling leaf",
	}, {
		desc:    "error non existant path",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "leaf-container-struct/does-not-exist"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		wantErr: "does-not-exist",
	}, {
		desc:    "error nil update",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  nil,
				}},
			}).Sync()
		},
		wantErr: "invalid nil Val",
	}, {
		desc:    "error wrong type",
		inQuery: lq,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantErr: "failed to unmarshal",
	}}
	for _, tt := range leafTests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			lookupCheckFn(t, fakeGNMI, c, tt.inQuery, tt.wantErr, tt.wantRequestValues, tt.wantSubscriptionPath, tt.wantVal)
		})
	}

	rootPath := testutil.GNMIPath(t, "parent/child")
	strPath := testutil.GNMIPath(t, "parent/child/state/one")
	enumPath := testutil.GNMIPath(t, "parent/child/state/three")
	strCfgPath := testutil.GNMIPath(t, "parent/child/config/one")

	configQuery := exampleocconfigpath.Root().Parent().Child().Config()
	stateQuery := exampleocconfigpath.Root().Parent().Child().State()

	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		inQuery              ygnmi.SingletonQuery[*exampleocconfig.Parent_Child]
		wantRequestValues    *ygnmi.RequestValues
		wantSubscriptionPath *gpb.Path
		wantVal              *ygnmi.Value[*exampleocconfig.Parent_Child]
		wantErr              string
	}{{
		desc: "success one update and state false",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: strCfgPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery: configQuery,
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  true,
			ConfigFiltered: false,
		},
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(&exampleocconfig.Parent_Child{
			One: ygot.String("foo"),
		}),
	}, {
		desc: "success one update and state true",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery: stateQuery,
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  false,
			ConfigFiltered: true,
		},
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(&exampleocconfig.Parent_Child{
			One: ygot.String("foo"),
		}),
	}, {
		desc: "success one update with prefix",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Prefix:    testutil.GNMIPath(t, "parent"),
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "child/state/one"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery: stateQuery,
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  false,
			ConfigFiltered: true,
		},
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(&exampleocconfig.Parent_Child{
			One: ygot.String("foo"),
		}),
	}, {
		desc: "success ignore state update when state false",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery: configQuery,
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  true,
			ConfigFiltered: false,
		},
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}),
	}, {
		desc: "success ignore non-state update when state true",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: strCfgPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery:              stateQuery,
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}),
	}, {
		desc: "success multiple updates in single notification",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "ONE"}},
				}, {
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery:              stateQuery,
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(&exampleocconfig.Parent_Child{
			One:   ygot.String("foo"),
			Three: exampleocconfig.Child_Three_ONE,
		}),
	}, {
		desc: "success multiple notifications",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "ONE"}},
				}},
			}).Notification(&gpb.Notification{
				Timestamp: 102,
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery:              stateQuery,
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 102),
		}).SetVal(&exampleocconfig.Parent_Child{
			One:   ygot.String("foo"),
			Three: exampleocconfig.Child_Three_ONE,
		}),
	}, {
		desc: "success no values",
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		inQuery: stateQuery,
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  false,
			ConfigFiltered: true,
		},
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path: rootPath,
		}),
	}}

	for _, tt := range nonLeafTests {
		t.Run("nonleaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			lookupCheckFn(t, fakeGNMI, c, tt.inQuery, tt.wantErr, tt.wantRequestValues, tt.wantSubscriptionPath, tt.wantVal)
		})
	}

	t.Run("success with ieeefloat32", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/model/a/single-key[key=foo]/state/counter"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_BytesVal{BytesVal: []byte{0xc0, 0x00, 0x00, 0x00}}},
			}},
		}).Sync()

		lookupCheckFn(
			t, fakeGNMI, c,
			exampleocconfigpath.Root().Model().SingleKey("foo").Counter().State(),
			"",
			nil,
			testutil.GNMIPath(t, "/model/a/single-key[key=foo]/state/counter"),
			(&ygnmi.Value[float32]{
				Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/state/counter"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(-2),
		)
	})

	t.Run("success with leaf-list ieeefloat32", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/model/a/single-key[key=foo]/state/counters"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_LeaflistVal{LeaflistVal: &gpb.ScalarArray{Element: []*gpb.TypedValue{{Value: &gpb.TypedValue_BytesVal{BytesVal: []byte{0xc0, 0x00, 0x00, 0x00}}}}}}},
			}},
		}).Sync()

		lookupCheckFn(
			t, fakeGNMI, c,
			exampleocconfigpath.Root().Model().SingleKey("foo").Counters().State(),
			"",
			nil,
			testutil.GNMIPath(t, "/model/a/single-key[key=foo]/state/counters"),
			(&ygnmi.Value[[]float32]{
				Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/state/counters"),
				Timestamp: time.Unix(0, 100),
			}).SetVal([]float32{-2}),
		)
	})

	t.Run("success ordered map", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Atomic:    true,
			Prefix:    testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/config/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/config/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/config/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/config/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/config/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/config/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 44}},
			}},
		}).Sync()

		lookupCheckFn(
			t, fakeGNMI, c,
			ygnmi.SingletonQuery[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap](exampleocconfigpath.Root().Model().SingleKey("foo").OrderedListMap().Config()),
			"",
			nil,
			testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap]{
				Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(getSamplePreferConfigOrderedMap(t)),
		)
	})

	t.Run("success whole single-keyed map", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Prefix:    testutil.GNMIPath(t, "/model/a"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `single-key[key=foo]/config/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/config/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/config/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/config/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/config/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/config/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 44}},
			}},
		}).Sync()

		lookupCheckFn(
			t, fakeGNMI, c,
			ygnmi.SingletonQuery[map[string]*exampleocconfig.Model_SingleKey](exampleocconfigpath.Root().Model().SingleKeyMap().Config()),
			"",
			nil,
			testutil.GNMIPath(t, "/model/a"),
			(&ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "/model/a"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(getSamplePreferConfigSingleKeyedMap(t)),
		)
	})
}

func TestPreferConfigLookupWithGet(t *testing.T) {
	fakeGNMI, c := newClient(t)
	leafPath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")

	tests := []struct {
		desc        string
		stub        func(s *gnmitestutil.Stubber)
		wantVal     *ygnmi.Value[string]
		wantRequest *gpb.GetRequest
		wantErr     string
	}{{
		desc: "success",
		stub: func(s *gnmitestutil.Stubber) {
			s.GetResponse(&gpb.GetResponse{
				Notification: []*gpb.Notification{{
					Timestamp: 100,
					Update: []*gpb.Update{{
						Path: leafPath,
						Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`"foo"`)}},
					}},
				}},
			}, nil)
		},
		wantVal: (&ygnmi.Value[string]{
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal("foo"),
		wantRequest: &gpb.GetRequest{
			Encoding: gpb.Encoding_JSON_IETF,
			Type:     gpb.GetRequest_STATE,
			Prefix:   &gpb.Path{},
			Path:     []*gpb.Path{leafPath},
		},
	}, {
		desc: "not found error",
		stub: func(s *gnmitestutil.Stubber) {
			s.GetResponse(nil, status.Error(codes.NotFound, "test"))
		},
		wantVal: (&ygnmi.Value[string]{
			Path: leafPath,
		}),
		wantRequest: &gpb.GetRequest{
			Encoding: gpb.Encoding_JSON_IETF,
			Type:     gpb.GetRequest_STATE,
			Prefix:   &gpb.Path{},
			Path:     []*gpb.Path{leafPath},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			lookupWithGetCheckFn(
				t, fakeGNMI, c,
				exampleocconfigpath.Root().RemoteContainer().ALeaf().State(),
				"",
				&ygnmi.RequestValues{
					StateFiltered:  false,
					ConfigFiltered: true,
				},
				tt.wantRequest,
				tt.wantVal,
			)
		})
	}

	nonLeafPath := testutil.GNMIPath(t, "/parent/child")
	nonLeafTests := []struct {
		desc    string
		stub    func(s *gnmitestutil.Stubber)
		wantVal *ygnmi.Value[*exampleocconfig.Parent_Child]
		wantErr string
	}{{
		desc: "single leaf",
		stub: func(s *gnmitestutil.Stubber) {
			s.GetResponse(&gpb.GetResponse{
				Notification: []*gpb.Notification{{
					Timestamp: 100,
					Update: []*gpb.Update{{
						Path: nonLeafPath,
						Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`{"config": {"three": "ONE" }}`)}},
					}},
				}},
			}, nil)
		},
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path:      nonLeafPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(&exampleocconfig.Parent_Child{Three: exampleocconfig.Child_Three_ONE}),
	}, {
		desc: "with extra value",
		stub: func(s *gnmitestutil.Stubber) {
			s.GetResponse(&gpb.GetResponse{
				Notification: []*gpb.Notification{{
					Timestamp: 100,
					Update: []*gpb.Update{{
						Path: nonLeafPath,
						Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`{"config": {"three": "ONE", "ten": "ten" }}`)}},
					}},
				}},
			}, nil)
		},
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path:      nonLeafPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(&exampleocconfig.Parent_Child{Three: exampleocconfig.Child_Three_ONE}),
	}, {
		desc: "with invalid type", // TODO: When partial unmarshaling of JSON is supports, this test case should have a value.
		stub: func(s *gnmitestutil.Stubber) {
			s.GetResponse(&gpb.GetResponse{
				Notification: []*gpb.Notification{{
					Timestamp: 100,
					Update: []*gpb.Update{{
						Path: nonLeafPath,
						Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`{"config": {"three": "ONE", "one": 10 }}`)}},
					}},
				}},
			}, nil)
		},
		wantVal: &ygnmi.Value[*exampleocconfig.Parent_Child]{
			Path: nonLeafPath,
			ComplianceErrors: &ygnmi.ComplianceErrors{
				TypeErrors: []*ygnmi.TelemetryError{{
					Path:  nonLeafPath,
					Value: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`{"config": {"three": "ONE", "one": 10 }}`)}},
				}},
			},
		},
	}}
	for _, tt := range nonLeafTests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			path := exampleocconfigpath.Root().Parent().Child().Config()
			lookupWithGetCheckFn(
				t, fakeGNMI, c,
				ygnmi.SingletonQuery[*exampleocconfig.Parent_Child](path),
				"",
				&ygnmi.RequestValues{
					StateFiltered:  true,
					ConfigFiltered: false,
				},
				&gpb.GetRequest{
					Encoding: gpb.Encoding_JSON_IETF,
					Type:     gpb.GetRequest_CONFIG,
					Prefix:   &gpb.Path{},
					Path:     []*gpb.Path{nonLeafPath},
				},
				tt.wantVal,
				cmpopts.IgnoreFields(ygnmi.TelemetryError{}, "Err"),
			)
		})
	}

	t.Run("success ordered map", func(t *testing.T) {
		fakeGNMI.Stub().GetResponse(&gpb.GetResponse{
			Notification: []*gpb.Notification{{
				Timestamp: 100,
				Atomic:    true,
				Prefix:    testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, ""),
					Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`{
  "openconfig-withlistval:ordered-list": [
    {
      "config": {
        "key": "foo",
        "value": "42"
      },
      "key": "foo"
    },
    {
      "config": {
        "key": "bar",
        "value": "43"
      },
      "key": "bar"
    },
    {
      "config": {
        "key": "baz",
        "value": "44"
      },
      "key": "baz"
    }
  ]
}`)}},
				}},
			}},
		}, nil)

		lookupWithGetCheckFn(
			t, fakeGNMI, c,
			ygnmi.SingletonQuery[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap](exampleocconfigpath.Root().Model().SingleKey("foo").OrderedListMap().Config()),
			"",
			nil,
			&gpb.GetRequest{
				Encoding: gpb.Encoding_JSON_IETF,
				Type:     gpb.GetRequest_CONFIG,
				Prefix:   &gpb.Path{},
				Path:     []*gpb.Path{testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists")},
			},
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap]{
				Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(getSamplePreferConfigOrderedMap(t)),
		)
	})

	t.Run("success whole single-keyed map", func(t *testing.T) {
		fakeGNMI.Stub().GetResponse(&gpb.GetResponse{
			Notification: []*gpb.Notification{{
				Timestamp: 100,
				Prefix:    testutil.GNMIPath(t, "/model/a"),
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, ""),
					Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`{
  "openconfig-withlistval:single-key": [
    {
      "config": {
        "key": "foo",
        "value": "42"
      },
      "key": "foo"
    },
    {
      "config": {
        "key": "bar",
        "value": "43"
      },
      "key": "bar"
    },
    {
      "config": {
        "key": "baz",
        "value": "44"
      },
      "key": "baz"
    }
  ]
}`)}},
				}},
			}},
		}, nil)

		lookupWithGetCheckFn(
			t, fakeGNMI, c,
			ygnmi.SingletonQuery[map[string]*exampleocconfig.Model_SingleKey](exampleocconfigpath.Root().Model().SingleKeyMap().Config()),
			"",
			nil,
			&gpb.GetRequest{
				Encoding: gpb.Encoding_JSON_IETF,
				Type:     gpb.GetRequest_CONFIG,
				Prefix:   &gpb.Path{},
				Path:     []*gpb.Path{testutil.GNMIPath(t, "/model/a")},
			},
			(&ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "/model/a"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(getSamplePreferConfigSingleKeyedMap(t)),
		)
	})
}

func TestPreferConfigGet(t *testing.T) {
	fakeGNMI, c := newClient(t)
	leafPath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := exampleocconfigpath.Root().RemoteContainer().ALeaf().State()

	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		wantSubscriptionPath *gpb.Path
		want                 string
		wantVal              string
		wantErr              string
	}{{
		desc: "value present",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: leafPath,
		wantVal:              "foo",
	}, {
		desc: "value not present",
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		wantSubscriptionPath: leafPath,
		wantErr:              "value not present",
	}, {
		desc: "error nil update",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  nil,
				}},
			}).Sync()
		},
		wantErr: "invalid nil Val",
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			getCheckFn(
				t, fakeGNMI, c, lq, tt.wantErr,
				&ygnmi.RequestValues{
					StateFiltered:  false,
					ConfigFiltered: true,
				},
				tt.wantSubscriptionPath, tt.wantVal)
		})
	}

	t.Run("use get", func(t *testing.T) {
		fakeGNMI.Stub().GetResponse(&gpb.GetResponse{
			Notification: []*gpb.Notification{{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "/remote-container/config/a-leaf"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`"foo"`)}},
				}},
			}},
		}, nil)
		wantGetRequest := &gpb.GetRequest{
			Encoding: gpb.Encoding_JSON_IETF,
			Type:     gpb.GetRequest_CONFIG,
			Prefix:   &gpb.Path{},
			Path:     []*gpb.Path{testutil.GNMIPath(t, "/remote-container/config/a-leaf")},
		}
		wantVal := "foo"

		got, err := ygnmi.Get[string](context.Background(), c, exampleocconfigpath.Root().RemoteContainer().ALeaf().Config(), ygnmi.WithUseGet())
		if err != nil {
			t.Fatalf("Get() returned unexpected error: %v", err)
		}
		if diff := cmp.Diff(wantVal, got, cmp.AllowUnexported(ygnmi.Value[string]{}), cmpopts.IgnoreFields(ygnmi.Value[string]{}, "RecvTimestamp"), protocmp.Transform()); diff != "" {
			t.Errorf("Get() returned unexpected diff: %s", diff)
		}
		if diff := cmp.Diff(wantGetRequest, fakeGNMI.GetRequests()[0], protocmp.Transform()); diff != "" {
			t.Errorf("Get() GetRequest different from expected: %s", diff)
		}
	})

	t.Run("success ordered map", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Atomic:    true,
			Prefix:    testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 44}},
			}},
		}).Sync()

		getCheckFn(t, fakeGNMI, c,
			exampleocconfigpath.Root().Model().SingleKey("foo").OrderedListMap().State(),
			"",
			nil,
			testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
			getSamplePreferConfigOrderedMap(t),
		)
	})
}

func TestPreferConfigWatch(t *testing.T) {
	fakeGNMI, client := newClient(t)
	path := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := exampleocconfigpath.Root().RemoteContainer().ALeaf().State()

	startTime := time.Now()
	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[string]
		wantVals             []*ygnmi.Value[string]
		wantErr              string
		wantMode             gpb.SubscriptionMode
		wantInterval         uint64
		opts                 []ygnmi.Option
	}{{
		desc: "single notif and pred true",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		dur: time.Second,
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      path,
			}).SetVal("foo")},
		wantSubscriptionPath: path,
		wantLastVal: (&ygnmi.Value[string]{
			Timestamp: startTime,
			Path:      path,
		}).SetVal("foo"),
	}, {
		desc: "single notif and pred true with custom mode",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		dur:      time.Second,
		opts:     []ygnmi.Option{ygnmi.WithSubscriptionMode(gpb.SubscriptionMode_ON_CHANGE)},
		wantMode: gpb.SubscriptionMode_ON_CHANGE,
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      path,
			}).SetVal("foo")},
		wantSubscriptionPath: path,
		wantLastVal: (&ygnmi.Value[string]{
			Timestamp: startTime,
			Path:      path,
		}).SetVal("foo"),
	}, {
		desc: "single notif and pred true with custom interval",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		dur:          time.Second,
		opts:         []ygnmi.Option{ygnmi.WithSampleInterval(time.Millisecond)},
		wantInterval: 1000000,
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      path,
			}).SetVal("foo")},
		wantSubscriptionPath: path,
		wantLastVal: (&ygnmi.Value[string]{
			Timestamp: startTime,
			Path:      path,
		}).SetVal("foo"),
	}, {
		desc: "single notif and pred false error EOF",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		dur: time.Second,
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      path,
			}).SetVal("bar"),
		},
		wantSubscriptionPath: path,
		wantLastVal: (&ygnmi.Value[string]{
			Timestamp: startTime,
			Path:      path,
		}).SetVal("bar"),
		wantErr: "error receiving gNMI response: EOF",
	}, {
		desc: "multiple notif and pred true",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			})
		},
		dur: time.Second,
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      path,
			}).SetVal("bar"),
			(&ygnmi.Value[string]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      path,
			}).SetVal("foo"),
		},
		wantSubscriptionPath: path,
		wantLastVal: (&ygnmi.Value[string]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      path,
		}).SetVal("foo"),
	}, {
		desc: "multiple notif with deletes",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Delete:    []*gpb.Path{path},
			})
		},
		dur: time.Second,
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      path,
			}).SetVal("bar"),
			(&ygnmi.Value[string]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      path,
			}),
		},
		wantSubscriptionPath: path,
		wantLastVal: (&ygnmi.Value[string]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      path,
		}),
		wantErr: "EOF",
	}, {
		desc: "negative duration",
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		dur:     -1 * time.Second,
		wantErr: "context deadline exceeded",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			watchCheckFn(t, fakeGNMI, tt.dur, client,
				lq,
				tt.opts,
				func(val string) bool { return val == "foo" },
				tt.wantErr,
				&ygnmi.RequestValues{
					StateFiltered:  false,
					ConfigFiltered: true,
				},
				[]*gpb.Path{tt.wantSubscriptionPath},
				[]gpb.SubscriptionMode{tt.wantMode},
				[]uint64{tt.wantInterval},
				tt.wantVals,
				tt.wantLastVal,
			)
		})
	}

	t.Run("multiple awaits", func(t *testing.T) {
		fakeGNMI.Stub().Sync()
		w := ygnmi.Watch(context.Background(), client, exampleocconfigpath.Root().RemoteContainer().ALeaf().State(), func(v *ygnmi.Value[string]) error { return nil })
		want := &ygnmi.Value[string]{
			Path: path,
		}
		val, err := w.Await()
		if err != nil {
			t.Fatalf("Await() got unexpected error: %v", err)
		}
		if diff := cmp.Diff(want, val, cmp.AllowUnexported(ygnmi.Value[string]{}), protocmp.Transform()); diff != "" {
			t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
		}
		_, err = w.Await()
		if d := errdiff.Check(err, "Await already called and Watcher is closed"); d != "" {
			t.Fatalf("Await() returned unexpected diff: %s", d)
		}
	})

	rootPath := testutil.GNMIPath(t, "parent/child")
	strPath := testutil.GNMIPath(t, "parent/child/state/one")
	enumPath := testutil.GNMIPath(t, "parent/child/state/three")
	startTime = time.Now()
	nonLeafQuery := exampleocconfigpath.Root().Parent().Child().State()

	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		opts                 []ygnmi.Option
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[*exampleocconfig.Parent_Child]
		wantVals             []*ygnmi.Value[*exampleocconfig.Parent_Child]
		wantErr              string
	}{{
		desc: "single notif and pred false",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[*exampleocconfig.Parent_Child]{
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleocconfig.Parent_Child{
				One: ygot.String("bar"),
			}),
		},
		wantErr:              "EOF",
		wantSubscriptionPath: rootPath,
		wantLastVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Timestamp: startTime,
			Path:      rootPath,
		}).SetVal(&exampleocconfig.Parent_Child{
			One: ygot.String("bar"),
		}),
	}, {
		desc: "multiple notif and pred true",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "ONE"}},
				}},
			})
		},
		wantVals: []*ygnmi.Value[*exampleocconfig.Parent_Child]{
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleocconfig.Parent_Child{
				One: ygot.String("foo"),
			}),
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      rootPath,
			}).SetVal(&exampleocconfig.Parent_Child{
				Three: exampleocconfig.Child_Three_ONE,
				One:   ygot.String("foo"),
			}),
		},
		wantSubscriptionPath: rootPath,
		wantLastVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}).SetVal(&exampleocconfig.Parent_Child{
			Three: exampleocconfig.Child_Three_ONE,
			One:   ygot.String("foo"),
		}),
	}, {
		desc: "multiple notif before sync",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "ONE"}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[*exampleocconfig.Parent_Child]{
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      rootPath,
			}).SetVal(&exampleocconfig.Parent_Child{
				Three: exampleocconfig.Child_Three_ONE,
				One:   ygot.String("foo"),
			})},
		wantSubscriptionPath: rootPath,
		wantLastVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}).SetVal(&exampleocconfig.Parent_Child{
			Three: exampleocconfig.Child_Three_ONE,
			One:   ygot.String("foo"),
		}),
	}, {
		desc: "delete leaf in container",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}, {
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "ONE"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Delete:    []*gpb.Path{strPath},
			})
		},
		wantVals: []*ygnmi.Value[*exampleocconfig.Parent_Child]{
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleocconfig.Parent_Child{
				Three: exampleocconfig.Child_Three_ONE,
				One:   ygot.String("bar"),
			}),
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      rootPath,
			}).SetVal(&exampleocconfig.Parent_Child{
				Three: exampleocconfig.Child_Three_ONE,
			}),
		},
		wantSubscriptionPath: rootPath,
		wantErr:              "EOF",
		wantLastVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}).SetVal(&exampleocconfig.Parent_Child{
			Three: exampleocconfig.Child_Three_ONE,
		}),
	}, {
		desc: "delete at container level",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}, {
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "ONE"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Delete:    []*gpb.Path{testutil.GNMIPath(t, "parent/child")},
			})
		},
		wantVals: []*ygnmi.Value[*exampleocconfig.Parent_Child]{
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleocconfig.Parent_Child{
				Three: exampleocconfig.Child_Three_ONE,
				One:   ygot.String("bar"),
			}),
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      rootPath,
			}),
		},
		wantSubscriptionPath: rootPath,
		wantErr:              "EOF",
		wantLastVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}).SetVal(&exampleocconfig.Parent_Child{
			Three: exampleocconfig.Child_Three_ONE,
		}),
	}}

	for _, tt := range nonLeafTests {
		t.Run("nonleaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			watchCheckFn(t, fakeGNMI, 2*time.Second, client,
				nonLeafQuery,
				tt.opts,
				func(val *exampleocconfig.Parent_Child) bool {
					return val.One != nil && *val.One == "foo" && val.Three == exampleocconfig.Child_Three_ONE
				},
				tt.wantErr,
				&ygnmi.RequestValues{
					StateFiltered:  false,
					ConfigFiltered: true,
				},
				[]*gpb.Path{tt.wantSubscriptionPath},
				[]gpb.SubscriptionMode{gpb.SubscriptionMode_TARGET_DEFINED},
				[]uint64{0},
				tt.wantVals,
				tt.wantLastVal,
			)
		})
	}

	t.Run("success ordered map", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: startTime.UnixNano(),
			Atomic:    true,
			Prefix:    testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}},
		}).Sync().Notification(&gpb.Notification{
			Timestamp: startTime.Add(time.Millisecond).UnixNano(),
			Atomic:    true,
			Prefix:    testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 44}},
			}},
		})

		want := getSamplePreferConfigOrderedMap(t)
		watchCheckFn(t, fakeGNMI, 2*time.Second, client,
			exampleocconfigpath.Root().Model().SingleKey("foo").OrderedListMap().State(),
			nil,
			func(val *exampleocconfig.Model_SingleKey_OrderedList_OrderedMap) bool {
				return cmp.Equal(val, want, cmp.AllowUnexported(exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{}))
			},
			"",
			nil,
			[]*gpb.Path{testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists")},
			[]gpb.SubscriptionMode{gpb.SubscriptionMode_TARGET_DEFINED},
			[]uint64{0},
			[]*ygnmi.Value[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap]{
				(&ygnmi.Value[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap]{
					Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
					Timestamp: startTime,
				}).SetVal(getSamplePreferConfigOrderedMapIncomplete(t)),
				(&ygnmi.Value[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap]{
					Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
					Timestamp: startTime.Add(time.Millisecond),
				}).SetVal(getSamplePreferConfigOrderedMap(t)),
			},
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap]{
				Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
				Timestamp: startTime.Add(time.Millisecond),
			}).SetVal(getSamplePreferConfigOrderedMap(t)),
		)
	})

	t.Run("success whole single-keyed map", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: startTime.UnixNano(),
			Prefix:    testutil.GNMIPath(t, "/model/a"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}},
		}).Sync().Notification(&gpb.Notification{
			Timestamp: startTime.Add(time.Millisecond).UnixNano(),
			Prefix:    testutil.GNMIPath(t, "/model/a"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 44}},
			}},
		})

		want := getSamplePreferConfigSingleKeyedMap(t)
		watchCheckFn(t, fakeGNMI, 2*time.Second, client,
			exampleocconfigpath.Root().Model().SingleKeyMap().State(),
			nil,
			func(val map[string]*exampleocconfig.Model_SingleKey) bool {
				return cmp.Equal(val, want)
			},
			"",
			nil,
			[]*gpb.Path{testutil.GNMIPath(t, "/model/a")},
			[]gpb.SubscriptionMode{gpb.SubscriptionMode_TARGET_DEFINED},
			[]uint64{0},
			[]*ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey]{
				(&ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey]{
					Path:      testutil.GNMIPath(t, "/model/a"),
					Timestamp: startTime,
				}).SetVal(getSamplePreferConfigSingleKeyedMapIncomplete(t)),
				(&ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey]{
					Path:      testutil.GNMIPath(t, "/model/a"),
					Timestamp: startTime.Add(time.Millisecond),
				}).SetVal(getSamplePreferConfigSingleKeyedMap(t)),
			},
			(&ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "/model/a"),
				Timestamp: startTime.Add(time.Millisecond),
			}).SetVal(getSamplePreferConfigSingleKeyedMap(t)),
		)
	})
}

func TestPreferConfigAwait(t *testing.T) {
	fakeGNMI, client := newClient(t)
	path := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := exampleocconfigpath.Root().RemoteContainer().ALeaf().State()

	startTime := time.Now()
	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantVal              *ygnmi.Value[string]
		wantErr              string
		wantMode             gpb.SubscriptionMode
		opts                 []ygnmi.Option
	}{{
		desc: "value never equal",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		dur:                  time.Second,
		wantSubscriptionPath: path,
		wantErr:              "EOF",
	}, {
		desc: "success",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		dur:                  time.Second,
		wantSubscriptionPath: path,
		wantVal: (&ygnmi.Value[string]{
			Timestamp: startTime,
			Path:      path,
		}).SetVal("foo"),
	}, {
		desc: "success with custom mode",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		dur:                  time.Second,
		opts:                 []ygnmi.Option{ygnmi.WithSubscriptionMode(gpb.SubscriptionMode_ON_CHANGE)},
		wantMode:             gpb.SubscriptionMode_ON_CHANGE,
		wantSubscriptionPath: path,
		wantVal: (&ygnmi.Value[string]{
			Timestamp: startTime,
			Path:      path,
		}).SetVal("foo"),
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			ctx, cancel := context.WithTimeout(context.Background(), tt.dur)
			defer cancel()
			val, err := ygnmi.Await(ctx, client, lq, "foo", tt.opts...)
			verifySubscriptionModesSent(t, fakeGNMI, tt.wantMode)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Await() returned unexpected diff: %s", diff)
			}
			if err != nil {
				return
			}
			if val != nil {
				checkJustReceived(t, val.RecvTimestamp)
				tt.wantVal.RecvTimestamp = val.RecvTimestamp
			}
			if diff := cmp.Diff(tt.wantVal, val, cmp.AllowUnexported(ygnmi.Value[string]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}

	rootPath := testutil.GNMIPath(t, "parent/child")
	strPath := testutil.GNMIPath(t, "parent/child/state/one")
	enumPath := testutil.GNMIPath(t, "parent/child/state/three")
	startTime = time.Now()
	nonLeafQuery := exampleocconfigpath.Root().Parent().Child().State()

	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[*exampleocconfig.Parent_Child]
		wantErr              string
	}{{
		desc: "value never equal",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		wantErr:              "EOF",
		wantSubscriptionPath: rootPath,
		wantLastVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Timestamp: startTime,
			Path:      rootPath,
		}).SetVal(&exampleocconfig.Parent_Child{
			One: ygot.String("bar"),
		}),
	}, {
		desc: "success",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "ONE"}},
				}},
			})
		},
		wantSubscriptionPath: rootPath,
		wantLastVal: (&ygnmi.Value[*exampleocconfig.Parent_Child]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}).SetVal(&exampleocconfig.Parent_Child{
			Three: exampleocconfig.Child_Three_ONE,
			One:   ygot.String("foo"),
		}),
	}}

	for _, tt := range nonLeafTests {
		t.Run("nonleaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			val, err := ygnmi.Await(context.Background(), client, nonLeafQuery, &exampleocconfig.Parent_Child{One: ygot.String("foo"), Three: exampleocconfig.Child_Three_ONE})
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Await() returned unexpected diff: %s", diff)
			}
			if err != nil {
				return
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath)
			if val != nil {
				checkJustReceived(t, val.RecvTimestamp)
				tt.wantLastVal.RecvTimestamp = val.RecvTimestamp
			}
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(ygnmi.Value[*exampleocconfig.Parent_Child]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestPreferConfigCollect(t *testing.T) {
	fakeGNMI, client := newClient(t)
	path := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := exampleocconfigpath.Root().RemoteContainer().ALeaf().State()

	startTime := time.Now()
	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		dur                  time.Duration
		wantRequestValues    *ygnmi.RequestValues
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[string]
		wantErr              string
		wantMode             gpb.SubscriptionMode
		opts                 []ygnmi.Option
	}{{
		desc: "no values",
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		dur:                  time.Second,
		wantSubscriptionPath: path,
		wantErr:              "EOF",
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Path: path,
			}),
		},
	}, {
		desc: "multiple values",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			})
		},
		dur:                  100 * time.Millisecond,
		wantSubscriptionPath: path,
		wantErr:              "EOF",
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      path,
			}).SetVal("foo"),
			(&ygnmi.Value[string]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      path,
			}).SetVal("bar"),
		},
	}, {
		desc: "multiple values and custom mode",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			})
		},
		dur:                  100 * time.Millisecond,
		wantSubscriptionPath: path,
		opts:                 []ygnmi.Option{ygnmi.WithSubscriptionMode(gpb.SubscriptionMode_ON_CHANGE)},
		wantMode:             gpb.SubscriptionMode_ON_CHANGE,
		wantErr:              "EOF",
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      path,
			}).SetVal("foo"),
			(&ygnmi.Value[string]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      path,
			}).SetVal("bar"),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			collectCheckFn(
				t, fakeGNMI, client, lq, tt.wantErr,
				&ygnmi.RequestValues{
					StateFiltered:  false,
					ConfigFiltered: true,
				},
				tt.wantSubscriptionPath, tt.wantVals,
			)
		})
	}

	rootPath := testutil.GNMIPath(t, "parent/child")
	strPath := testutil.GNMIPath(t, "parent/child/state/one")
	enumPath := testutil.GNMIPath(t, "parent/child/state/three")
	startTime = time.Now()
	nonLeafQuery := exampleocconfigpath.Root().Parent().Child().State()

	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[*exampleocconfig.Parent_Child]
		wantErr              string
	}{{
		desc: "one val",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		wantErr:              "EOF",
		wantSubscriptionPath: rootPath,
		wantVals: []*ygnmi.Value[*exampleocconfig.Parent_Child]{
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleocconfig.Parent_Child{
				One: ygot.String("bar"),
			}),
		},
	}, {
		desc: "multiple values",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "ONE"}},
				}},
			})
		},
		wantErr:              "EOF",
		wantSubscriptionPath: rootPath,
		wantVals: []*ygnmi.Value[*exampleocconfig.Parent_Child]{
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleocconfig.Parent_Child{
				One: ygot.String("foo"),
			}),
			(&ygnmi.Value[*exampleocconfig.Parent_Child]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      rootPath,
			}).SetVal(&exampleocconfig.Parent_Child{
				Three: exampleocconfig.Child_Three_ONE,
				One:   ygot.String("foo"),
			}),
		},
	}}

	for _, tt := range nonLeafTests {
		t.Run("nonleaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			collectCheckFn(
				t, fakeGNMI, client, nonLeafQuery, tt.wantErr,
				&ygnmi.RequestValues{
					StateFiltered:  false,
					ConfigFiltered: true,
				},
				tt.wantSubscriptionPath, tt.wantVals)
		})
	}
}

func TestPreferConfigLookupAll(t *testing.T) {
	fakeGNMI, c := newClient(t)
	leafPath := testutil.GNMIPath(t, "model/a/single-key[key=*]/state/value")
	lq := exampleocconfigpath.Root().Model().SingleKeyAny().Value().State()

	leafTests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[int64]
		wantErr              string
	}{{
		desc: "success one value",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[int64]{
			(&ygnmi.Value[int64]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(10),
		},
		wantSubscriptionPath: leafPath,
	}, {
		desc: "success no values",
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		wantVals:             nil,
		wantSubscriptionPath: leafPath,
	}, {
		desc: "non compliant value",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}, {
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/state/fake-val"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[int64]{
			(&ygnmi.Value[int64]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(10),
		},
		wantSubscriptionPath: leafPath,
	}, {
		desc: "success multiples value in same notification",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}, {
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[int64]{
			(&ygnmi.Value[int64]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(10),
			(&ygnmi.Value[int64]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(11)},
		wantSubscriptionPath: leafPath,
	}, {
		desc: "success multiples value in different notifications",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[int64]{
			(&ygnmi.Value[int64]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(10),
			(&ygnmi.Value[int64]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value"),
				Timestamp: time.Unix(0, 101),
			}).SetVal(11)},
		wantSubscriptionPath: leafPath,
	}, {
		desc: "success ignore mismatched paths",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/config/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Sync()
		},
		wantVals:             nil,
		wantSubscriptionPath: leafPath,
	}, {
		desc: "success ignore mismatched types",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: ""}},
				}},
			}).Sync()
		},
		wantVals:             nil,
		wantSubscriptionPath: leafPath,
	}, {
		desc: "error nil update val",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value"),
					Val:  nil,
				}},
			}).Sync()
		},
		wantErr:              "failed to receive to data",
		wantSubscriptionPath: leafPath,
	}}

	for _, tt := range leafTests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			lookupAllCheckFn(t, fakeGNMI, c, lq, tt.wantErr,
				&ygnmi.RequestValues{
					StateFiltered:  false,
					ConfigFiltered: true,
				},
				tt.wantSubscriptionPath, tt.wantVals, false,
			)
		})
	}

	nonLeafPath := testutil.GNMIPath(t, "model/a/single-key[key=*]")
	nonLeafQ := exampleocconfigpath.Root().Model().SingleKeyAny().State()
	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[*exampleocconfig.Model_SingleKey]
		wantErr              string
	}{{
		desc: "one value",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=10]"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(10),
			}),
		},
		wantSubscriptionPath: nonLeafPath,
	}, {
		desc: "multiple values",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}, {
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 101}},
				}, {
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/key"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "10"}},
				}},
			}).Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/state/key"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "11"}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=10]"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(100),
				Key:   ygot.String("10"),
			}),
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=11]"),
				Timestamp: time.Unix(0, 101),
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(101),
				Key:   ygot.String("11"),
			}),
		},
		wantSubscriptionPath: nonLeafPath,
	}, {
		desc: "non compliant values",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value-fake"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}, {
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/key"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "10"}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=10]"),
				Timestamp: time.Unix(0, 100),
				ComplianceErrors: &ygnmi.ComplianceErrors{
					PathErrors: []*ygnmi.TelemetryError{{
						Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
						Path:  testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value-fake"),
					}},
				},
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Key: ygot.String("10"),
			}),
		},
		wantSubscriptionPath: nonLeafPath,
	}, {
		desc: "no values",
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		wantVals:             nil,
		wantSubscriptionPath: nonLeafPath,
	}}

	for _, tt := range nonLeafTests {
		t.Run("nonLeaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			lookupAllCheckFn(
				t, fakeGNMI, c, nonLeafQ, tt.wantErr,
				&ygnmi.RequestValues{
					StateFiltered:  false,
					ConfigFiltered: true,
				},
				tt.wantSubscriptionPath, tt.wantVals, true,
			)
		})
	}

	t.Run("success ordered map", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Atomic:    true,
			Prefix:    testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=baz]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 44}},
			}},
		}).Notification(&gpb.Notification{
			Timestamp: 101,
			Atomic:    true,
			Prefix:    testutil.GNMIPath(t, "/model/a/single-key[key=bar]/ordered-lists"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}, {
				Path: testutil.GNMIPath(t, `ordered-list[key=bar]/state/dne-value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}},
		}).Sync()

		lookupAllCheckFn(
			t, fakeGNMI, c,
			exampleocconfigpath.Root().Model().SingleKeyAny().OrderedListMap().State(),
			"",
			nil,
			testutil.GNMIPath(t, "/model/a/single-key[key=*]/ordered-lists"),
			[]*ygnmi.Value[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap]{
				// In alphabetical order.
				(&ygnmi.Value[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap]{
					Path:      testutil.GNMIPath(t, "/model/a/single-key[key=bar]/ordered-lists"),
					Timestamp: time.Unix(0, 101),
					ComplianceErrors: &ygnmi.ComplianceErrors{
						PathErrors: []*ygnmi.TelemetryError{{
							Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
							Path:  testutil.GNMIPath(t, "model/a/single-key[key=bar]/ordered-lists/ordered-list[key=bar]/state/dne-value"),
						}},
					},
				}).SetVal(getSamplePreferConfigOrderedMapIncomplete(t)),
				(&ygnmi.Value[*exampleocconfig.Model_SingleKey_OrderedList_OrderedMap]{
					Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
					Timestamp: time.Unix(0, 100),
				}).SetVal(getSamplePreferConfigOrderedMap(t)),
			},
			true,
		)
	})

	t.Run("success whole single-keyed map", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Atomic:    true,
			Prefix:    testutil.GNMIPath(t, "/model/a/single-key[key=foo]/inner-a"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 44}},
			}},
		}).Notification(&gpb.Notification{
			Timestamp: 101,
			Atomic:    true,
			Prefix:    testutil.GNMIPath(t, "/model/a/single-key[key=bar]/inner-a"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/dne-value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}},
		}).Sync()

		lookupAllCheckFn(
			t, fakeGNMI, c,
			exampleocconfigpath.Root().Model().SingleKeyAny().SingleKeyMap().State(),
			"",
			nil,
			testutil.GNMIPath(t, "/model/a/single-key[key=*]/inner-a"),
			[]*ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey_SingleKey]{
				// In alphabetical order.
				(&ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey_SingleKey]{
					Path:      testutil.GNMIPath(t, "/model/a/single-key[key=bar]/inner-a"),
					Timestamp: time.Unix(0, 101),
					ComplianceErrors: &ygnmi.ComplianceErrors{
						PathErrors: []*ygnmi.TelemetryError{{
							Value: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
							Path:  testutil.GNMIPath(t, "model/a/single-key[key=bar]/inner-a/single-key[key=bar]/state/dne-value"),
						}},
					},
				}).SetVal(getSamplePreferConfigInnerSingleKeyedMapIncomplete(t)),
				(&ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey_SingleKey]{
					Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/inner-a"),
					Timestamp: time.Unix(0, 100),
				}).SetVal(getSamplePreferConfigInnerSingleKeyedMap(t)),
			},
			true,
		)
	})

	t.Run("use get", func(t *testing.T) {
		fakeGNMI.Stub().GetResponse(&gpb.GetResponse{
			Notification: []*gpb.Notification{{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`"1"`)}},
				}},
			}},
		}, nil)
		wantGetRequest := &gpb.GetRequest{
			Encoding: gpb.Encoding_JSON_IETF,
			Type:     gpb.GetRequest_STATE,
			Prefix:   &gpb.Path{},
			Path:     []*gpb.Path{leafPath},
		}
		wantVal := []*ygnmi.Value[int64]{(&ygnmi.Value[int64]{
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(1)}

		got, err := ygnmi.LookupAll(context.Background(), c, lq, ygnmi.WithUseGet())
		if err != nil {
			t.Fatalf("LookupAll() returned unexpected error: %v", err)
		}
		if diff := cmp.Diff(wantVal, got, cmp.AllowUnexported(ygnmi.Value[int64]{}), cmpopts.IgnoreFields(ygnmi.Value[int64]{}, "RecvTimestamp"), protocmp.Transform()); diff != "" {
			t.Errorf("LookupAll() returned unexpected diff: %s", diff)
		}
		if diff := cmp.Diff(wantGetRequest, fakeGNMI.GetRequests()[0], protocmp.Transform()); diff != "" {
			t.Errorf("LookupAll() GetRequest different from expected: %s", diff)
		}
	})
}

func TestPreferConfigGetAll(t *testing.T) {
	fakeGNMI, c := newClient(t)
	leafPath := testutil.GNMIPath(t, "model/a/single-key[key=*]/state/value")
	lq := exampleocconfigpath.Root().Model().SingleKeyAny().Value().State()

	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []int64
		wantErr              string
	}{{
		desc: "success",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantVals:             []int64{10},
		wantSubscriptionPath: leafPath,
	}, {
		desc: "success no values",
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		wantErr:              ygnmi.ErrNotPresent.Error(),
		wantSubscriptionPath: leafPath,
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			got, err := ygnmi.GetAll(context.Background(), c, lq)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("LookupAll(ctx, c, %v) returned unexpected diff: %s", lq, diff)
			}
			if err != nil {
				return
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath)
			if diff := cmp.Diff(tt.wantVals, got); diff != "" {
				t.Errorf("LookupAll() returned unexpected diff (-want,+got):\n%s", diff)
			}
		})
	}
	t.Run("use get", func(t *testing.T) {
		fakeGNMI.Stub().GetResponse(&gpb.GetResponse{
			Notification: []*gpb.Notification{{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`"1"`)}},
				}},
			}},
		}, nil)
		wantGetRequest := &gpb.GetRequest{
			Encoding: gpb.Encoding_JSON_IETF,
			Type:     gpb.GetRequest_STATE,
			Prefix:   &gpb.Path{},
			Path:     []*gpb.Path{leafPath},
		}
		wantVal := []int64{1}

		got, err := ygnmi.GetAll(context.Background(), c, exampleocconfigpath.Root().Model().SingleKeyAny().Value().State(), ygnmi.WithUseGet())
		if err != nil {
			t.Fatalf("Get() returned unexpected error: %v", err)
		}
		if diff := cmp.Diff(wantVal, got, cmp.AllowUnexported(ygnmi.Value[string]{}), cmpopts.IgnoreFields(ygnmi.Value[string]{}, "RecvTimestamp"), protocmp.Transform()); diff != "" {
			t.Errorf("Get() returned unexpected diff: %s", diff)
		}
		if diff := cmp.Diff(wantGetRequest, fakeGNMI.GetRequests()[0], protocmp.Transform()); diff != "" {
			t.Errorf("Get() GetRequest different from expected: %s", diff)
		}
	})
}

func TestPreferConfigWatchAll(t *testing.T) {
	fakeGNMI, client := newClient(t)
	leafQueryPath := testutil.GNMIPath(t, "model/a/single-key[key=*]/state/value")
	key10Path := testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value")
	key11Path := testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value")

	startTime := time.Now()
	lq := exampleocconfigpath.Root().Model().SingleKeyAny().Value().State()
	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[int64]
		wantVals             []*ygnmi.Value[int64]
		wantErr              string
		wantMode             gpb.SubscriptionMode
		opts                 []ygnmi.Option
	}{{
		desc: "predicate not true",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: leafQueryPath,
		wantVals: []*ygnmi.Value[int64]{
			(&ygnmi.Value[int64]{
				Timestamp: startTime,
				Path:      key10Path,
			}).SetVal(100),
		},
		wantLastVal: (&ygnmi.Value[int64]{
			Timestamp: startTime,
			Path:      key10Path,
		}).SetVal(100),
		wantErr: "EOF",
	}, {
		desc: "predicate becomes true",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: key11Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 101}},
				}},
			})
		},
		wantSubscriptionPath: leafQueryPath,
		wantVals: []*ygnmi.Value[int64]{
			(&ygnmi.Value[int64]{
				Timestamp: startTime,
				Path:      key10Path,
			}).SetVal(100),
			(&ygnmi.Value[int64]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      key11Path,
			}).SetVal(101),
		},
		wantLastVal: (&ygnmi.Value[int64]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      key11Path,
		}).SetVal(101),
	}, {
		desc: "predicate becomes true with custom mode",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: key11Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 101}},
				}},
			})
		},
		opts:                 []ygnmi.Option{ygnmi.WithSubscriptionMode(gpb.SubscriptionMode_ON_CHANGE)},
		wantMode:             gpb.SubscriptionMode_ON_CHANGE,
		wantSubscriptionPath: leafQueryPath,
		wantVals: []*ygnmi.Value[int64]{
			(&ygnmi.Value[int64]{
				Timestamp: startTime,
				Path:      key10Path,
			}).SetVal(100),
			(&ygnmi.Value[int64]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      key11Path,
			}).SetVal(101),
		},
		wantLastVal: (&ygnmi.Value[int64]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      key11Path,
		}).SetVal(101),
	}, {
		desc: "multiple values in notification",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}, {
					Path: key11Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 101}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: leafQueryPath,
		wantVals: []*ygnmi.Value[int64]{
			(&ygnmi.Value[int64]{
				Timestamp: startTime,
				Path:      key10Path,
			}).SetVal(100),
			(&ygnmi.Value[int64]{
				Timestamp: startTime,
				Path:      key11Path,
			}).SetVal(101),
		},
		wantLastVal: (&ygnmi.Value[int64]{
			Timestamp: startTime,
			Path:      key11Path,
		}).SetVal(101),
	}, {
		desc: "error nil value",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  nil,
				}},
			}).Sync()
		},
		wantSubscriptionPath: leafQueryPath,
		wantLastVal:          nil,
		wantErr:              "invalid nil Val in update",
	}, {
		desc: "subscribe fails",
		dur:  -1 * time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  nil,
				}},
			}).Sync()
		},
		wantSubscriptionPath: leafQueryPath,
		wantLastVal:          nil,
		wantErr:              "gNMI failed to Subscribe",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			i := 0
			ctx, cancel := context.WithTimeout(context.Background(), tt.dur)
			defer cancel()
			var key10Cond, key11Cond bool

			w := ygnmi.WatchAll(ctx, client, lq, func(v *ygnmi.Value[int64]) error {
				if i > len(tt.wantVals) {
					t.Fatalf("Predicate(%d) expected no more values but got: %+v", i, v)
				}
				if diff := cmp.Diff(tt.wantVals[i], v, cmpopts.IgnoreFields(ygnmi.Value[int64]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[int64]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate(%d) got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", i, diff, v.ComplianceErrors)
				}
				val, present := v.Val()
				key10Cond = key10Cond || (present && proto.Equal(v.Path, key10Path) && val == 100)
				key11Cond = key11Cond || (present && proto.Equal(v.Path, key11Path) && val == 101)
				i++
				if key10Cond && key11Cond {
					return nil
				}
				return ygnmi.Continue
			}, tt.opts...)
			val, err := w.Await()
			if i < len(tt.wantVals) {
				t.Errorf("Predicate received too few values: got %d, want %d", i, len(tt.wantVals))
			}
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Await() returned unexpected diff: %s", diff)
			}
			verifySubscriptionModesSent(t, fakeGNMI, tt.wantMode)
			if val != nil {
				checkJustReceived(t, val.RecvTimestamp)
				tt.wantLastVal.RecvTimestamp = val.RecvTimestamp
			}
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(ygnmi.Value[int64]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}

	nonLeafPath := testutil.GNMIPath(t, "model/a/single-key[key=*]")
	nonLeafKey10Path := testutil.GNMIPath(t, "model/a/single-key[key=10]")
	nonLeafKey11Path := testutil.GNMIPath(t, "model/a/single-key[key=11]")

	nonLeafQ := exampleocconfigpath.Root().Model().SingleKeyAny().State()
	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[*exampleocconfig.Model_SingleKey]
		wantVals             []*ygnmi.Value[*exampleocconfig.Model_SingleKey]
		wantErr              string
	}{{
		desc: "predicate not true",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: nonLeafPath,
		wantVals: []*ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime,
				Path:      nonLeafKey10Path,
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(100),
			}),
		},
		wantLastVal: (&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			Timestamp: startTime,
			Path:      nonLeafKey10Path,
		}).SetVal(&exampleocconfig.Model_SingleKey{
			Value: ygot.Int64(100),
		}),
		wantErr: "EOF",
	}, {
		desc: "predicate becomes true",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}, {
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/state/key"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "test"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: key11Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 101}},
				}},
			})
		},
		wantSubscriptionPath: nonLeafPath,
		wantVals: []*ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime,
				Path:      nonLeafKey10Path,
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(100),
			}),
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime,
				Path:      nonLeafKey11Path,
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Key: ygot.String("test"),
			}),
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      nonLeafKey11Path,
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(101),
				Key:   ygot.String("test"),
			}),
		},
		wantLastVal: (&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      nonLeafKey11Path,
		}).SetVal(&exampleocconfig.Model_SingleKey{
			Value: ygot.Int64(101),
			Key:   ygot.String("test"),
		}),
	}, {
		desc: "predicate becomes true after some deletions",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}, {
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/state/key"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "test"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Delete:    []*gpb.Path{testutil.GNMIPath(t, "model/a/single-key[key=11]/state/key")},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(2 * time.Millisecond).UnixNano(),
				Delete:    []*gpb.Path{testutil.GNMIPath(t, "model/a/single-key[key=10]")},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(3 * time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}, {
					Path: testutil.GNMIPath(t, "model/a/single-key[key=11]/state/key"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "test"}},
				}, {
					Path: key11Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 101}},
				}},
			})
		},
		wantSubscriptionPath: nonLeafPath,
		wantVals: []*ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime,
				Path:      nonLeafKey10Path,
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(100),
			}),
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime,
				Path:      nonLeafKey11Path,
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Key: ygot.String("test"),
			}),
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      nonLeafKey11Path,
			}),
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime.Add(2 * time.Millisecond),
				Path:      nonLeafKey10Path,
			}),
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime.Add(3 * time.Millisecond),
				Path:      nonLeafKey10Path,
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(100),
			}),
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime.Add(3 * time.Millisecond),
				Path:      nonLeafKey11Path,
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(101),
				Key:   ygot.String("test"),
			}),
		},
		wantLastVal: (&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			Timestamp: startTime.Add(3 * time.Millisecond),
			Path:      nonLeafKey11Path,
		}).SetVal(&exampleocconfig.Model_SingleKey{
			Value: ygot.Int64(101),
			Key:   ygot.String("test"),
		}),
	}}
	for _, tt := range nonLeafTests {
		t.Run("nonLeaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			i := 0
			ctx, cancel := context.WithTimeout(context.Background(), tt.dur)
			defer cancel()
			var key10Cond, key11Cond bool

			w := ygnmi.WatchAll(ctx, client, nonLeafQ, func(v *ygnmi.Value[*exampleocconfig.Model_SingleKey]) error {
				if i > len(tt.wantVals) {
					t.Fatalf("Predicate(%d) expected no more values but got: %+v", i, v)
				}
				if diff := cmp.Diff(tt.wantVals[i], v, cmpopts.IgnoreFields(ygnmi.Value[*exampleocconfig.Model_SingleKey]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[*exampleocconfig.Model_SingleKey]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate(%d) got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", i, diff, v.ComplianceErrors)
				}
				val, present := v.Val()
				key10Cond = key10Cond || (present && proto.Equal(v.Path, nonLeafKey10Path) && val.Value != nil && *val.Value == 100)
				key11Cond = key11Cond || (present && proto.Equal(v.Path, nonLeafKey11Path) && val.Value != nil && *val.Value == 101)
				i++
				if key10Cond && key11Cond {
					return nil
				}
				return ygnmi.Continue
			})
			val, err := w.Await()
			if i < len(tt.wantVals) {
				t.Errorf("Predicate received too few values: got %d, want %d", i, len(tt.wantVals))
			}
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Await() returned unexpected diff: %s", diff)
			}
			if val != nil {
				checkJustReceived(t, val.RecvTimestamp)
				tt.wantLastVal.RecvTimestamp = val.RecvTimestamp
			}
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(ygnmi.Value[*exampleocconfig.Model_SingleKey]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestPreferConfigCollectAll(t *testing.T) {
	fakeGNMI, client := newClient(t)
	leafQueryPath := testutil.GNMIPath(t, "model/a/single-key[key=*]/state/value")
	key10Path := testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value")
	key11Path := testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value")

	startTime := time.Now()
	lq := exampleocconfigpath.Root().Model().SingleKeyAny().Value().State()
	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[int64]
		wantErr              string
		wantMode             gpb.SubscriptionMode
		opts                 []ygnmi.Option
	}{{
		desc: "no values",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		wantErr:              "EOF",
		wantSubscriptionPath: leafQueryPath,
		wantVals:             nil,
	}, {
		desc: "no values with custom mode",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		opts:                 []ygnmi.Option{ygnmi.WithSubscriptionMode(gpb.SubscriptionMode_ON_CHANGE)},
		wantMode:             gpb.SubscriptionMode_ON_CHANGE,
		wantErr:              "EOF",
		wantSubscriptionPath: leafQueryPath,
		wantVals:             nil,
	}, {
		desc: "multiple values",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: key11Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 101}},
				}},
			})
		},
		wantErr:              "EOF",
		wantSubscriptionPath: leafQueryPath,
		wantVals: []*ygnmi.Value[int64]{
			(&ygnmi.Value[int64]{
				Timestamp: startTime,
				Path:      key10Path,
			}).SetVal(100),
			(&ygnmi.Value[int64]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      key11Path,
			}).SetVal(101),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			ctx, cancel := context.WithTimeout(context.Background(), tt.dur)
			defer cancel()

			vals, err := ygnmi.CollectAll(ctx, client, lq, tt.opts...).Await()
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Await() returned unexpected diff: %s", diff)
			}
			verifySubscriptionModesSent(t, fakeGNMI, tt.wantMode)
			for _, val := range vals {
				checkJustReceived(t, val.RecvTimestamp)
			}
			if diff := cmp.Diff(tt.wantVals, vals, cmpopts.IgnoreFields(ygnmi.Value[int64]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[int64]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}

	nonLeafPath := testutil.GNMIPath(t, "model/a/single-key[key=*]")
	nonLeafKey10Path := testutil.GNMIPath(t, "model/a/single-key[key=10]")
	nonLeafKey11Path := testutil.GNMIPath(t, "model/a/single-key[key=11]")

	nonLeafQ := exampleocconfigpath.Root().Model().SingleKeyAny().State()
	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[*exampleocconfig.Model_SingleKey]
		wantErr              string
	}{{
		desc: "no values",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Sync()
		},
		wantSubscriptionPath: nonLeafPath,
		wantVals:             nil,
		wantErr:              "EOF",
	}, {
		desc: "multiple values",
		dur:  time.Second,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: key11Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 101}},
				}},
			})
		},
		wantSubscriptionPath: nonLeafPath,
		wantErr:              "EOF",
		wantVals: []*ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime,
				Path:      nonLeafKey10Path,
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(100),
			}),
			(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      nonLeafKey11Path,
			}).SetVal(&exampleocconfig.Model_SingleKey{
				Value: ygot.Int64(101),
			}),
		},
	}}
	for _, tt := range nonLeafTests {
		t.Run("nonLeaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			ctx, cancel := context.WithTimeout(context.Background(), tt.dur)
			defer cancel()

			vals, err := ygnmi.CollectAll(ctx, client, nonLeafQ).Await()
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Await() returned unexpected diff: %s", diff)
			}
			for _, val := range vals {
				checkJustReceived(t, val.RecvTimestamp)
			}
			if diff := cmp.Diff(tt.wantVals, vals, cmpopts.IgnoreFields(ygnmi.Value[*exampleocconfig.Model_SingleKey]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[*exampleocconfig.Model_SingleKey]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestPreferConfigUpdate(t *testing.T) {
	setClient := &gnmitestutil.SetClient{}
	client, err := ygnmi.NewClient(setClient, ygnmi.WithTarget("dut"))
	if err != nil {
		t.Fatalf("Unexpected error creating client: %v", err)
	}

	tests := []struct {
		desc         string
		op           func(*ygnmi.Client) (*ygnmi.Result, error)
		wantErr      string
		wantRequest  *gpb.SetRequest
		stubResponse *gpb.SetResponse
		stubErr      error
	}{{
		desc: "scalar leaf",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, exampleocconfigpath.Root().Parent().Child().One().Config(), "10")
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/config/one"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"10\"")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "non scalar leaf",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, exampleocconfigpath.Root().Parent().Child().Three().Config(), exampleocconfig.Child_Three_ONE)
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/config/three"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"ONE\"")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "non leaf",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, exampleocconfigpath.Root().Parent().Child().Config(), &exampleocconfig.Parent_Child{One: ygot.String("10")})
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("{\n  \"openconfig-simple:config\": {\n    \"one\": \"10\"\n  }\n}")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "server error",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, exampleocconfigpath.Root().Parent().Child().One().Config(), "10")
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/config/one"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`"10"`)}},
			}},
		},
		stubErr: fmt.Errorf("fake"),
		wantErr: "fake",
	}, {
		desc: "YANG ordered list",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			om := &exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{}
			ol, err := om.AppendNew("foo")
			if err != nil {
				t.Fatal(err)
			}
			ol.SetValue(42)
			ol, err = om.AppendNew("bar")
			if err != nil {
				t.Fatal(err)
			}
			ol.SetValue(43)
			ol, err = om.AppendNew("baz")
			if err != nil {
				t.Fatal(err)
			}
			ol.SetValue(44)
			return ygnmi.Update(context.Background(), c, exampleocconfigpath.Root().Model().SingleKey("foo").OrderedListMap().Config(), om)
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{
  "openconfig-withlistval:ordered-list": [
    {
      "config": {
        "key": "foo",
        "value": "42"
      },
      "key": "foo"
    },
    {
      "config": {
        "key": "bar",
        "value": "43"
      },
      "key": "bar"
    },
    {
      "config": {
        "key": "baz",
        "value": "44"
      },
      "key": "baz"
    }
  ]
}`))}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "whole single-keyed list",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, exampleocconfigpath.Root().Model().SingleKeyMap().Config(), getSamplePreferConfigSingleKeyedMap(t))
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/model/a"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{
  "openconfig-withlistval:single-key": [
    {
      "config": {
        "key": "bar",
        "value": "43"
      },
      "key": "bar"
    },
    {
      "config": {
        "key": "baz",
        "value": "44"
      },
      "key": "baz"
    },
    {
      "config": {
        "key": "foo",
        "value": "42"
      },
      "key": "foo"
    }
  ]
}`))}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "leaf and prefer proto",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, exampleocconfigpath.Root().Parent().Child().One().Config(), "10", ygnmi.WithSetPreferProtoEncoding())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/config/one"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "10"}},
			}},
		},
	}, {
		desc: "non leaf and prefer proto",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, exampleocconfigpath.Root().Parent().Child().Config(), &exampleocconfig.Parent_Child{One: ygot.String("10")}, ygnmi.WithSetPreferProtoEncoding())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("{\n  \"openconfig-simple:config\": {\n    \"one\": \"10\"\n  }\n}")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "fallback openconfig origin",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, mustSchemaless[*gpb.CapabilityResponse](t, "/foo", "openconfig"), &gpb.CapabilityResponse{GNMIVersion: "1"}, ygnmi.WithSetFallbackEncoding())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "foo"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_AnyVal{AnyVal: mustAnyNew(t, &gpb.CapabilityResponse{GNMIVersion: "1"})}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
		wantErr: "failed to encode set request",
	}, {
		desc: "fallback empty origin",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, mustSchemaless[*gpb.CapabilityResponse](t, "/foo", "openconfig"), &gpb.CapabilityResponse{GNMIVersion: "1"}, ygnmi.WithSetFallbackEncoding())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "foo"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_AnyVal{AnyVal: mustAnyNew(t, &gpb.CapabilityResponse{GNMIVersion: "1"})}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
		wantErr: "failed to encode set request",
	}, {
		desc: "fallback proto",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, mustSchemaless[*gpb.CapabilityResponse](t, "/foo", "test"), &gpb.CapabilityResponse{GNMIVersion: "1"}, ygnmi.WithSetFallbackEncoding())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: &gpb.Path{
					Elem:   []*gpb.PathElem{{Name: "foo"}},
					Origin: "test",
				},
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_AnyVal{AnyVal: mustAnyNew(t, &gpb.CapabilityResponse{GNMIVersion: "1"})}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "fallback json",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(context.Background(), c, mustSchemaless[*testStruct](t, "/foo", "test"), &testStruct{Val: "test"}, ygnmi.WithSetFallbackEncoding())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: &gpb.Path{
					Origin: "test",
					Elem:   []*gpb.PathElem{{Name: "foo"}},
				},
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonVal{JsonVal: []byte(`{"Val":"test"}`)}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			setClient.Reset()
			setClient.AddResponse(tt.stubResponse, tt.stubErr)

			got, err := tt.op(client)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Update() returned unexpected diff: %s", diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.wantRequest, setClient.Requests[0], protocmp.Transform()); diff != "" {
				t.Errorf("Update() sent unexpected request (-want,+got):\n%s", diff)
			}
			want := &ygnmi.Result{
				RawResponse: tt.stubResponse,
				Timestamp:   time.Unix(0, tt.stubResponse.GetTimestamp()),
			}
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Update() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestPreferConfigReplace(t *testing.T) {
	setClient := &gnmitestutil.SetClient{}
	client, err := ygnmi.NewClient(setClient, ygnmi.WithTarget("dut"))
	if err != nil {
		t.Fatalf("Unexpected error creating client: %v", err)
	}
	tests := []struct {
		desc         string
		op           func(*ygnmi.Client) (*ygnmi.Result, error)
		wantErr      string
		wantRequest  *gpb.SetRequest
		stubResponse *gpb.SetResponse
		stubErr      error
	}{{
		desc: "scalar leaf",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Replace(context.Background(), c, exampleocconfigpath.Root().Parent().Child().One().Config(), "10")
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/config/one"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"10\"")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "non scalar leaf",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Replace(context.Background(), c, exampleocconfigpath.Root().Parent().Child().Three().Config(), exampleocconfig.Child_Three_ONE)

		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/config/three"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"ONE\"")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "non leaf",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Replace(context.Background(), c, exampleocconfigpath.Root().Parent().Child().Config(), &exampleocconfig.Parent_Child{One: ygot.String("10")})
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("{\n  \"openconfig-simple:config\": {\n    \"one\": \"10\"\n  }\n}")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "YANG ordered list",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			om := &exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{}
			ol, err := om.AppendNew("foo")
			if err != nil {
				t.Fatal(err)
			}
			ol.SetValue(42)
			ol, err = om.AppendNew("bar")
			if err != nil {
				t.Fatal(err)
			}
			ol.SetValue(43)
			ol, err = om.AppendNew("baz")
			if err != nil {
				t.Fatal(err)
			}
			ol.SetValue(44)
			return ygnmi.Replace(context.Background(), c, exampleocconfigpath.Root().Model().SingleKey("foo").OrderedListMap().Config(), om)
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{
  "openconfig-withlistval:ordered-list": [
    {
      "config": {
        "key": "foo",
        "value": "42"
      },
      "key": "foo"
    },
    {
      "config": {
        "key": "bar",
        "value": "43"
      },
      "key": "bar"
    },
    {
      "config": {
        "key": "baz",
        "value": "44"
      },
      "key": "baz"
    }
  ]
}`))}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "whole single-keyed list",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Replace(context.Background(), c, exampleocconfigpath.Root().Model().SingleKeyMap().Config(), getSamplePreferConfigSingleKeyedMap(t))
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/model/a"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{
  "openconfig-withlistval:single-key": [
    {
      "config": {
        "key": "bar",
        "value": "43"
      },
      "key": "bar"
    },
    {
      "config": {
        "key": "baz",
        "value": "44"
      },
      "key": "baz"
    },
    {
      "config": {
        "key": "foo",
        "value": "42"
      },
      "key": "foo"
    }
  ]
}`))}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "server error",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Replace(context.Background(), c, exampleocconfigpath.Root().Parent().Child().One().Config(), "10")
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/config/one"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`"10"`)}},
			}},
		},
		stubErr: fmt.Errorf("fake"),
		wantErr: "fake",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			configCheckFn(t, setClient, client, tt.op, tt.wantRequest, tt.stubResponse, tt.wantErr, tt.stubErr)
		})
	}
}

func TestPreferConfigDelete(t *testing.T) {
	setClient := &gnmitestutil.SetClient{}
	client, err := ygnmi.NewClient(setClient, ygnmi.WithTarget("dut"))
	if err != nil {
		t.Fatalf("Unexpected error creating client: %v", err)
	}
	tests := []struct {
		desc         string
		op           func(*ygnmi.Client) (*ygnmi.Result, error)
		wantErr      string
		wantRequest  *gpb.SetRequest
		stubResponse *gpb.SetResponse
		stubErr      error
	}{{
		desc: "success",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Delete(context.Background(), c, exampleocconfigpath.Root().Parent().Child().One().Config())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Delete: []*gpb.Path{
				testutil.GNMIPath(t, "parent/child/config/one"),
			},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "YANG ordered list",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Delete(context.Background(), c, exampleocconfigpath.Root().Model().SingleKey("foo").OrderedListMap().Config())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Delete: []*gpb.Path{
				testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
			},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "whole single-keyed list",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Delete(context.Background(), c, exampleocconfigpath.Root().Model().SingleKeyMap().Config())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Delete: []*gpb.Path{
				testutil.GNMIPath(t, "/model/a"),
			},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "server error",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Delete(context.Background(), c, exampleocconfigpath.Root().Parent().Child().One().Config())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Delete: []*gpb.Path{
				testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			},
		},
		stubErr: fmt.Errorf("fake"),
		wantErr: "fake",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			configCheckFn(t, setClient, client, tt.op, tt.wantRequest, tt.stubResponse, tt.wantErr, tt.stubErr)
		})
	}
}

func TestPreferConfigBatchGet(t *testing.T) {
	fakeGNMI, c := newClient(t)
	aLeafStatePath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	aLeafConfigPath := testutil.GNMIPath(t, "/remote-container/config/a-leaf")
	twoPath := testutil.GNMIPath(t, "/parent/child/state/two")
	aLeafSubPath := testutil.GNMIPath(t, "/remote-container/*/a-leaf")

	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		config               bool
		paths                []ygnmi.PathStruct
		wantSubscriptionPath []*gpb.Path
		wantVal              *ygnmi.Value[*exampleocconfig.Root]
		wantErr              string
	}{{
		desc: "state leaves",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: aLeafStatePath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}, {
					Path: aLeafConfigPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "config"}},
				}, {
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.PathStruct{
			exampleocconfigpath.Root().RemoteContainer().ALeaf(),
			exampleocconfigpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantVal: (&ygnmi.Value[*exampleocconfig.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleocconfig.Root{
			RemoteContainer: &exampleocconfig.RemoteContainer{ALeaf: ygot.String("foo")},
			Parent:          &exampleocconfig.Parent{Child: &exampleocconfig.Parent_Child{Two: ygot.String("bar")}},
		}),
	}, {
		desc:   "config ignore state leaves",
		config: true,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: aLeafStatePath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}, {
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.PathStruct{
			exampleocconfigpath.Root().RemoteContainer().ALeaf(),
			exampleocconfigpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantVal: (&ygnmi.Value[*exampleocconfig.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleocconfig.Root{
			RemoteContainer: &exampleocconfig.RemoteContainer{},
			Parent:          &exampleocconfig.Parent{Child: &exampleocconfig.Parent_Child{}},
		}),
	}, {
		desc: "non leaves",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: aLeafStatePath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}, {
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.PathStruct{
			exampleocconfigpath.Root().RemoteContainer(),
			exampleocconfigpath.Root().Parent(),
		},
		wantSubscriptionPath: []*gpb.Path{
			testutil.GNMIPath(t, "/remote-container"),
			testutil.GNMIPath(t, "/parent"),
		},
		wantVal: (&ygnmi.Value[*exampleocconfig.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleocconfig.Root{
			RemoteContainer: &exampleocconfig.RemoteContainer{ALeaf: ygot.String("foo")},
			Parent:          &exampleocconfig.Parent{Child: &exampleocconfig.Parent_Child{Two: ygot.String("bar")}},
		}),
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			b := &exampleocconfigpath.Batch{}
			b.AddPaths(tt.paths...)
			query := b.State()
			if tt.config {
				query = b.Config()
			}
			got, err := ygnmi.Lookup(context.Background(), c, query)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Lookup() returned unexpected diff: %s", diff)
			}
			if err != nil {
				return
			}
			checkJustReceived(t, got.RecvTimestamp)
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath...)
			tt.wantVal.RecvTimestamp = got.RecvTimestamp

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[*exampleocconfig.Root]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Lookup() returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", diff, got.ComplianceErrors)
			}
		})
	}
	t.Run("immutable query", func(t *testing.T) {
		fakeGNMI.Stub().Sync()
		b := &exampleocconfigpath.Batch{}
		b.AddPaths(exampleocconfigpath.Root().Model())
		q := b.State()
		if _, err := ygnmi.Lookup(context.Background(), c, q); err != nil {
			t.Fatal(err)
		}
		verifySubscriptionPathsSent(t, fakeGNMI, testutil.GNMIPath(t, "/model"))
		b.AddPaths(exampleocconfigpath.Root().A(), exampleocconfigpath.Root().A().B())
		if _, err := ygnmi.Lookup(context.Background(), c, q); err != nil {
			t.Fatal(err)
		}
		verifySubscriptionPathsSent(t, fakeGNMI, testutil.GNMIPath(t, "/model"))
	})
}

func TestPreferConfigBatchWatch(t *testing.T) {
	fakeGNMI, c := newClient(t)
	aLeafStatePath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	twoPath := testutil.GNMIPath(t, "/parent/child/state/two")
	aLeafSubPath := testutil.GNMIPath(t, "/remote-container/*/a-leaf")

	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		config               bool
		paths                []ygnmi.PathStruct
		wantSubscriptionPath []*gpb.Path
		wantVal              *ygnmi.Value[*exampleocconfig.Root]
		wantErr              string
	}{{
		desc: "predicate true",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: aLeafStatePath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}, {
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.PathStruct{
			exampleocconfigpath.Root().RemoteContainer().ALeaf(),
			exampleocconfigpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantVal: (&ygnmi.Value[*exampleocconfig.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleocconfig.Root{
			RemoteContainer: &exampleocconfig.RemoteContainer{ALeaf: ygot.String("foo")},
			Parent:          &exampleocconfig.Parent{Child: &exampleocconfig.Parent_Child{Two: ygot.String("bar")}},
		}),
	}, {
		desc: "predicate false true false",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: aLeafStatePath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Notification(&gpb.Notification{
				Timestamp: 102,
				Update: []*gpb.Update{{
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "sample"}},
				}},
			})
		},
		paths: []ygnmi.PathStruct{
			exampleocconfigpath.Root().RemoteContainer().ALeaf(),
			exampleocconfigpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantVal: (&ygnmi.Value[*exampleocconfig.Root]{
			Timestamp: time.Unix(0, 101),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleocconfig.Root{
			RemoteContainer: &exampleocconfig.RemoteContainer{ALeaf: ygot.String("foo")},
			Parent:          &exampleocconfig.Parent{Child: &exampleocconfig.Parent_Child{Two: ygot.String("bar")}},
		}),
	}, {
		desc:   "predicate false",
		config: true,
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.PathStruct{
			exampleocconfigpath.Root().RemoteContainer().ALeaf(),
			exampleocconfigpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantErr: "EOF",
	}, {
		desc: "non leaves",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: aLeafStatePath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}, {
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.PathStruct{
			exampleocconfigpath.Root().RemoteContainer(),
			exampleocconfigpath.Root().Parent(),
		},
		wantSubscriptionPath: []*gpb.Path{
			testutil.GNMIPath(t, "/remote-container"),
			testutil.GNMIPath(t, "/parent"),
		},
		wantVal: (&ygnmi.Value[*exampleocconfig.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleocconfig.Root{
			RemoteContainer: &exampleocconfig.RemoteContainer{ALeaf: ygot.String("foo")},
			Parent:          &exampleocconfig.Parent{Child: &exampleocconfig.Parent_Child{Two: ygot.String("bar")}},
		}),
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			b := &exampleocconfigpath.Batch{}
			b.AddPaths(tt.paths...)
			query := b.State()
			if tt.config {
				query = b.Config()
			}

			got, err := ygnmi.Watch(context.Background(), c, query, func(v *ygnmi.Value[*exampleocconfig.Root]) error {
				if v, ok := v.Val(); ok && v.GetRemoteContainer().GetALeaf() == "foo" && v.GetParent().GetChild().GetTwo() == "bar" {
					return nil
				}
				return ygnmi.Continue
			}).Await()
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Watch() returned unexpected diff: %s", diff)
			}
			if err != nil {
				return
			}
			checkJustReceived(t, got.RecvTimestamp)
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath...)
			tt.wantVal.RecvTimestamp = got.RecvTimestamp

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[*exampleocconfig.Root]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Watch() returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", diff, got.ComplianceErrors)
			}
		})
	}
}

func TestPreferConfigCustomRootBatch(t *testing.T) {
	fakeGNMI, c := newClient(t)
	twoPath := testutil.GNMIPath(t, "/parent/child/state/two")

	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		paths                []ygnmi.UntypedQuery
		wantSubscriptionPath []*gpb.Path
		wantVal              *ygnmi.Value[*exampleocconfig.Parent]
		wantAddErr           string
		wantLookupErr        string
	}{{
		desc: "not prefix",
		stub: func(s *gnmitestutil.Stubber) {},
		paths: []ygnmi.UntypedQuery{
			exampleocconfigpath.Root().Model().Config(),
		},
		wantAddErr: "is not a prefix",
	}, {
		desc: "success",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.UntypedQuery{
			exampleocconfigpath.Root().Parent().Child().Two().State(),
		},
		wantSubscriptionPath: []*gpb.Path{
			twoPath,
		},
		wantVal: (&ygnmi.Value[*exampleocconfig.Parent]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/parent"),
		}).SetVal(&exampleocconfig.Parent{
			Child: &exampleocconfig.Parent_Child{
				Two: ygot.String("foo"),
			},
		}),
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			b := ygnmi.NewBatch(exampleocconfigpath.Root().Parent().State())
			err := b.AddPaths(tt.paths...)
			if diff := errdiff.Substring(err, tt.wantAddErr); diff != "" {
				t.Fatalf("AddPaths returned unexpected diff: %s", diff)
			}
			if err != nil {
				return
			}
			got, gotErr := ygnmi.Lookup(context.Background(), c, b.Query())
			if diff := errdiff.Substring(gotErr, tt.wantLookupErr); diff != "" {
				t.Fatalf("Watch() returned unexpected diff: %s", diff)
			}
			if gotErr != nil {
				return
			}
			checkJustReceived(t, got.RecvTimestamp)
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath...)
			tt.wantVal.RecvTimestamp = got.RecvTimestamp

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[*exampleocconfig.Parent]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Watch() returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", diff, got.ComplianceErrors)
			}
		})
	}

	fakeGNMI, client := newClient(t)
	startTime := time.Now()
	t.Run("success whole single-keyed map", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: startTime.UnixNano(),
			Prefix:    testutil.GNMIPath(t, "/model/a"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}},
		}).Sync().Notification(&gpb.Notification{
			Timestamp: startTime.Add(time.Millisecond).UnixNano(),
			Prefix:    testutil.GNMIPath(t, "/model/a"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=bar]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 43}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "baz"}},
			}, {
				Path: testutil.GNMIPath(t, `single-key[key=baz]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 44}},
			}},
		})

		modelPath := exampleocconfigpath.Root().Model()
		b := ygnmi.NewBatch(modelPath.SingleKeyMap().State())
		if err := b.AddPaths(
			modelPath.SingleKeyAny().Key().State(),
			modelPath.SingleKeyAny().Value().State(),
		); err != nil {
			t.Fatal(err)
		}

		want := getSamplePreferConfigSingleKeyedMap(t)
		watchCheckFn(t, fakeGNMI, 2*time.Second, client,
			b.Query(),
			nil,
			func(val map[string]*exampleocconfig.Model_SingleKey) bool {
				return cmp.Equal(val, want)
			},
			"",
			nil,
			[]*gpb.Path{
				testutil.GNMIPath(t, "/model/a/single-key[key=*]/state/key"),
				testutil.GNMIPath(t, "/model/a/single-key[key=*]/state/value"),
			},
			[]gpb.SubscriptionMode{
				gpb.SubscriptionMode_TARGET_DEFINED,
				gpb.SubscriptionMode_TARGET_DEFINED,
			},
			[]uint64{0, 0},
			[]*ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey]{
				(&ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey]{
					Path:      testutil.GNMIPath(t, "/model/a"),
					Timestamp: startTime,
				}).SetVal(getSamplePreferConfigSingleKeyedMapIncomplete(t)),
				(&ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey]{
					Path:      testutil.GNMIPath(t, "/model/a"),
					Timestamp: startTime.Add(time.Millisecond),
				}).SetVal(getSamplePreferConfigSingleKeyedMap(t)),
			},
			(&ygnmi.Value[map[string]*exampleocconfig.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "/model/a"),
				Timestamp: startTime.Add(time.Millisecond),
			}).SetVal(getSamplePreferConfigSingleKeyedMap(t)),
		)
	})
}

func TestPreferConfigCustomRootWildcardBatch(t *testing.T) {
	fakeGNMI, c := newClient(t)
	valuePathWild := testutil.GNMIPath(t, "/model/a/single-key[key=*]/state/value")
	valuePath := testutil.GNMIPath(t, "/model/a/single-key[key=foo]/state/value")
	keyPathWild := testutil.GNMIPath(t, "/model/a/single-key[key=*]/state/key")
	keyPath := testutil.GNMIPath(t, "/model/a/single-key[key=foo]/state/key")

	tests := []struct {
		desc                 string
		stub                 func(s *gnmitestutil.Stubber)
		paths                []ygnmi.UntypedQuery
		wantSubscriptionPath []*gpb.Path
		wantVal              []*ygnmi.Value[*exampleocconfig.Model_SingleKey]
		wantAddErr           string
		wantLookupErr        string
	}{{
		desc: "not prefix",
		stub: func(s *gnmitestutil.Stubber) {},
		paths: []ygnmi.UntypedQuery{
			exampleocconfigpath.Root().Model().Config(),
		},
		wantAddErr: "is not a prefix",
	}, {
		desc: "success",
		stub: func(s *gnmitestutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: valuePath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 42}},
				}, {
					Path: keyPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.UntypedQuery{
			exampleocconfigpath.Root().Model().SingleKeyAny().Value().State(),
			exampleocconfigpath.Root().Model().SingleKeyAny().Key().State(),
		},
		wantSubscriptionPath: []*gpb.Path{
			keyPathWild,
			valuePathWild,
		},
		wantVal: []*ygnmi.Value[*exampleocconfig.Model_SingleKey]{(&ygnmi.Value[*exampleocconfig.Model_SingleKey]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]"),
		}).SetVal(&exampleocconfig.Model_SingleKey{
			Key:   ygot.String("foo"),
			Value: ygot.Int64(42),
		})},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			b := ygnmi.NewWildcardBatch(exampleocconfigpath.Root().Model().SingleKeyAny().State())
			err := b.AddPaths(tt.paths...)
			if diff := errdiff.Substring(err, tt.wantAddErr); diff != "" {
				t.Fatalf("AddPaths returned unexpected diff: %s", diff)
			}
			if err != nil {
				return
			}
			got, gotErr := ygnmi.LookupAll(context.Background(), c, b.Query())
			if diff := errdiff.Substring(gotErr, tt.wantLookupErr); diff != "" {
				t.Fatalf("Watch() returned unexpected diff: %s", diff)
			}
			if gotErr != nil {
				return
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath...)

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[*exampleocconfig.Model_SingleKey]{}), protocmp.Transform(), cmpopts.IgnoreFields(ygnmi.Value[*exampleocconfig.Model_SingleKey]{}, "RecvTimestamp")); diff != "" {
				t.Errorf("Watch() returned unexpected diff (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestPreferConfigSetBatch(t *testing.T) {
	setClient := &gnmitestutil.SetClient{}
	client, err := ygnmi.NewClient(setClient, ygnmi.WithTarget("dut"))
	if err != nil {
		t.Fatalf("Unexpected error creating client: %v", err)
	}
	tests := []struct {
		desc         string
		addPaths     func(*ygnmi.SetBatch)
		wantErr      string
		wantRequest  *gpb.SetRequest
		stubResponse *gpb.SetResponse
		stubErr      error
	}{{
		desc: "leaf update replace delete unionreplace",
		addPaths: func(sb *ygnmi.SetBatch) {
			cliPath, err := schemaless.NewConfig[string]("", "cli")
			if err != nil {
				t.Fatalf("Failed to create CLI ygnmi query: %v", err)
			}
			ygnmi.BatchUpdate(sb, cliPath, "hello, mercury")
			ygnmi.BatchUpdate(sb, exampleocconfigpath.Root().Parent().Child().One().Config(), "foo")
			ygnmi.BatchReplace(sb, cliPath, "hello, venus")
			ygnmi.BatchReplace(sb, exampleocconfigpath.Root().Parent().Child().One().Config(), "bar")
			ygnmi.BatchDelete(sb, cliPath)
			ygnmi.BatchDelete(sb, exampleocconfigpath.Root().Parent().Child().One().Config())
			ygnmi.BatchUnionReplace(sb, exampleocconfigpath.Root().Parent().Child().One().Config(), "baz")
			ygnmi.BatchUnionReplaceCLI(sb, "openos", "open sesame")
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: &gpb.Path{Origin: "cli"},
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_AsciiVal{AsciiVal: "hello, mercury"}},
			}, {
				Path: testutil.GNMIPath(t, "parent/child/config/one"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"foo\"")}},
			}},
			Replace: []*gpb.Update{{
				Path: &gpb.Path{Origin: "cli"},
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_AsciiVal{AsciiVal: "hello, venus"}},
			}, {
				Path: testutil.GNMIPath(t, "parent/child/config/one"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"bar\"")}},
			}},
			Delete: []*gpb.Path{
				{Origin: "cli"},
				testutil.GNMIPath(t, "parent/child/config/one"),
			},
			UnionReplace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/config/one"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"baz\"")}},
			}, {
				Path: &gpb.Path{Origin: "openos_cli"},
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_AsciiVal{AsciiVal: "open sesame"}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "non leaf update delete replace",
		addPaths: func(sb *ygnmi.SetBatch) {
			ygnmi.BatchUpdate(sb, exampleocconfigpath.Root().Parent().Child().Config(), &exampleocconfig.Parent_Child{One: ygot.String("foo")})
			ygnmi.BatchDelete(sb, exampleocconfigpath.Root().Parent().Child().One().Config())

			ygnmi.BatchReplace(sb, exampleocconfigpath.Root().Model().SingleKeyMap().Config(), getSamplePreferConfigSingleKeyedMap(t))
			ygnmi.BatchUpdate(sb, exampleocconfigpath.Root().Model().SingleKeyMap().Config(), getSamplePreferConfigSingleKeyedMap(t))
			ygnmi.BatchReplace(sb, exampleocconfigpath.Root().Model().SingleKey("foo").OrderedListMap().Config(), getSamplePreferConfigOrderedMap(t))
			ygnmi.BatchUpdate(sb, exampleocconfigpath.Root().Model().SingleKey("bar").OrderedListMap().Config(), getSamplePreferConfigOrderedMap(t))
			ygnmi.BatchDelete(sb, exampleocconfigpath.Root().Model().SingleKeyMap().Config())
			ygnmi.BatchDelete(sb, exampleocconfigpath.Root().Model().SingleKey("baz").OrderedListMap().Config())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/model/a"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{
  "openconfig-withlistval:single-key": [
    {
      "config": {
        "key": "bar",
        "value": "43"
      },
      "key": "bar"
    },
    {
      "config": {
        "key": "baz",
        "value": "44"
      },
      "key": "baz"
    },
    {
      "config": {
        "key": "foo",
        "value": "42"
      },
      "key": "foo"
    }
  ]
}`))}},
			}, {
				Path: testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{
  "openconfig-withlistval:ordered-list": [
    {
      "config": {
        "key": "foo",
        "value": "42"
      },
      "key": "foo"
    },
    {
      "config": {
        "key": "bar",
        "value": "43"
      },
      "key": "bar"
    },
    {
      "config": {
        "key": "baz",
        "value": "44"
      },
      "key": "baz"
    }
  ]
}`))}},
			}},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("{\n  \"openconfig-simple:config\": {\n    \"one\": \"foo\"\n  }\n}")}},
			}, {
				Path: testutil.GNMIPath(t, "/model/a"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{
  "openconfig-withlistval:single-key": [
    {
      "config": {
        "key": "bar",
        "value": "43"
      },
      "key": "bar"
    },
    {
      "config": {
        "key": "baz",
        "value": "44"
      },
      "key": "baz"
    },
    {
      "config": {
        "key": "foo",
        "value": "42"
      },
      "key": "foo"
    }
  ]
}`))}},
			}, {
				Path: testutil.GNMIPath(t, "/model/a/single-key[key=bar]/ordered-lists"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{
  "openconfig-withlistval:ordered-list": [
    {
      "config": {
        "key": "foo",
        "value": "42"
      },
      "key": "foo"
    },
    {
      "config": {
        "key": "bar",
        "value": "43"
      },
      "key": "bar"
    },
    {
      "config": {
        "key": "baz",
        "value": "44"
      },
      "key": "baz"
    }
  ]
}`))}},
			}},
			Delete: []*gpb.Path{
				testutil.GNMIPath(t, "parent/child/config/one"),
				testutil.GNMIPath(t, "/model/a"),
				testutil.GNMIPath(t, "/model/a/single-key[key=baz]/ordered-lists"),
			},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "non leaf update delete replace nil objects",
		addPaths: func(sb *ygnmi.SetBatch) {
			ygnmi.BatchReplace(sb, exampleocconfigpath.Root().Model().SingleKeyMap().Config(), map[string]*exampleocconfig.Model_SingleKey{})
			ygnmi.BatchUpdate(sb, exampleocconfigpath.Root().Model().SingleKeyMap().Config(), nil)
			ygnmi.BatchReplace(sb, exampleocconfigpath.Root().Model().SingleKey("foo").OrderedListMap().Config(), &exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{})
			ygnmi.BatchUpdate(sb, exampleocconfigpath.Root().Model().SingleKey("bar").OrderedListMap().Config(), nil)
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/model/a"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{
  "openconfig-withlistval:single-key": []
}`))}},
			}, {
				Path: testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{
  "openconfig-withlistval:ordered-list": []
}`))}},
			}},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/model/a"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{}`))}},
			}, {
				Path: testutil.GNMIPath(t, "/model/a/single-key[key=bar]/ordered-lists"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(removeWhitespace(`{}`))}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "server error",
		addPaths: func(sb *ygnmi.SetBatch) {
			ygnmi.BatchDelete(sb, exampleocconfigpath.Root().Parent().Child().One().Config())
		},
		stubErr: fmt.Errorf("fake"),
		wantErr: "fake",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			setClient.Reset()
			setClient.AddResponse(tt.stubResponse, tt.stubErr)
			b := &ygnmi.SetBatch{}
			tt.addPaths(b)

			got, err := b.Set(context.Background(), client)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Set() returned unexpected diff: %s", diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.wantRequest, setClient.Requests[0], protocmp.Transform()); diff != "" {
				t.Errorf("Set() sent unexpected request (-want,+got):\n%s", diff)
			}
			want := &ygnmi.Result{
				RawResponse: tt.stubResponse,
				Timestamp:   time.Unix(0, tt.stubResponse.GetTimestamp()),
			}
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Set() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestPreferConfigWatchCancel(t *testing.T) {
	srv := &gnmiS{
		errCh: make(chan error, 1),
	}
	s := grpc.NewServer(grpc.Creds(local.NewCredentials()))
	gpb.RegisterGNMIServer(s, srv)
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		//nolint:errcheck // Don't care about this error.
		s.Serve(l)
	}()
	conn, err := grpc.Dial(l.Addr().String(), grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	c, _ := ygnmi.NewClient(gpb.NewGNMIClient(conn))

	w := ygnmi.Watch(context.Background(), c, exampleocconfigpath.Root().RemoteContainer().ALeaf().State(), func(v *ygnmi.Value[string]) error {
		return nil
	})
	if _, err := w.Await(); err != nil {
		t.Fatal(err)
	}
	if err := <-srv.errCh; err == nil {
		t.Fatalf("Watch() unexpected error: got %v, want context.Cancel", err)
	}
}

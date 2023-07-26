// Copyright 2023 Google LLC
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

//go:build go1.21

package ygnmi_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygnmi/internal/uexampleoc"
	"github.com/openconfig/ygnmi/internal/uexampleoc/uexampleocpath"
	"github.com/openconfig/ygnmi/schemaless"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/testing/protocmp"
)

func getSampleSingleKeyedMapUncompressed(t *testing.T) map[string]*uexampleoc.OpenconfigWithlistval_Model_A_SingleKey {
	model := &uexampleoc.OpenconfigWithlistval_Model{}
	model.GetOrCreateA().GetOrCreateSingleKey("foo").GetOrCreateConfig().SetKey("foo")
	model.GetOrCreateA().GetOrCreateSingleKey("foo").GetOrCreateConfig().SetValue(42)
	model.GetOrCreateA().GetOrCreateSingleKey("foo").GetOrCreateState().SetValue(84)
	model.GetOrCreateA().GetOrCreateSingleKey("bar").GetOrCreateConfig().SetValue(43)
	model.GetOrCreateA().GetOrCreateSingleKey("bar").GetOrCreateConfig().SetKey("bar")
	model.GetOrCreateA().GetOrCreateSingleKey("baz").GetOrCreateConfig().SetKey("baz")
	model.GetOrCreateA().GetOrCreateSingleKey("baz").GetOrCreateConfig().SetValue(44)
	return model.A.SingleKey
}

func TestUncompressedTelemetry(t *testing.T) {
	fakeGNMI, c := newClient(t)

	t.Run("lookup-config-leaf", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `/model/a/single-key[key=foo]/config/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}},
		}).Sync()

		lookupCheckFn(
			t, fakeGNMI, c,
			ygnmi.SingletonQuery[string](uexampleocpath.Root().Model().A().SingleKey("foo").Config().Key()),
			"",
			testutil.GNMIPath(t, `/model/a/single-key[key=foo]/config/key`),
			(&ygnmi.Value[string]{
				Path:      testutil.GNMIPath(t, `/model/a/single-key[key=foo]/config/key`),
				Timestamp: time.Unix(0, 100),
			}).SetVal("foo"),
		)
	})

	t.Run("lookup-state-leaf", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `/model/a/single-key[key=foo]/state/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}},
		}).Sync()

		lookupCheckFn(
			t, fakeGNMI, c,
			uexampleocpath.Root().Model().A().SingleKey("foo").State().Key(),
			"",
			testutil.GNMIPath(t, `/model/a/single-key[key=foo]/state/key`),
			(&ygnmi.Value[string]{
				Path:      testutil.GNMIPath(t, `/model/a/single-key[key=foo]/state/key`),
				Timestamp: time.Unix(0, 100),
			}).SetVal("foo"),
		)
	})

	t.Run("lookup-container", func(t *testing.T) {
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
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 84}},
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
			uexampleocpath.Root().Model().A(),
			"",
			testutil.GNMIPath(t, "/model/a"),
			(&ygnmi.Value[*uexampleoc.OpenconfigWithlistval_Model_A]{
				Path:      testutil.GNMIPath(t, "/model/a"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(&uexampleoc.OpenconfigWithlistval_Model_A{SingleKey: getSampleSingleKeyedMapUncompressed(t)}),
		)
	})

	t.Run("lookup-whole-single-keyed-map", func(t *testing.T) {
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
				Path: testutil.GNMIPath(t, `single-key[key=foo]/state/value`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 84}},
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
			uexampleocpath.Root().Model().A().SingleKeyMap(),
			"",
			testutil.GNMIPath(t, "/model/a/single-key"),
			(&ygnmi.Value[map[string]*uexampleoc.OpenconfigWithlistval_Model_A_SingleKey]{
				Path:      testutil.GNMIPath(t, "/model/a/single-key"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(getSampleSingleKeyedMapUncompressed(t)),
		)
	})

	startTime := time.Now()
	t.Run("collect-leaf", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: startTime.UnixNano(),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `/model/a/single-key[key=foo]/config/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}},
		}).Sync().Notification(&gpb.Notification{
			Timestamp: startTime.Add(time.Millisecond).UnixNano(),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `/model/a/single-key[key=foo]/config/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
			}},
		}).Sync()

		collectCheckFn(
			t, fakeGNMI, c,
			uexampleocpath.Root().Model().A().SingleKey("foo").Config().Key(),
			"EOF",
			testutil.GNMIPath(t, `/model/a/single-key[key=foo]/config/key`),
			[]*ygnmi.Value[string]{
				(&ygnmi.Value[string]{
					Path:      testutil.GNMIPath(t, `/model/a/single-key[key=foo]/config/key`),
					Timestamp: startTime,
				}).SetVal("foo"),
				(&ygnmi.Value[string]{
					Path:      testutil.GNMIPath(t, `/model/a/single-key[key=foo]/config/key`),
					Timestamp: startTime.Add(time.Millisecond),
				}).SetVal("bar"),
			},
		)
	})

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
			uexampleocpath.Root().Model().A().SingleKey("foo").State().Counter(),
			"",
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
			uexampleocpath.Root().Model().A().SingleKey("foo").State().Counters(),
			"",
			testutil.GNMIPath(t, "/model/a/single-key[key=foo]/state/counters"),
			(&ygnmi.Value[[]float32]{
				Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/state/counters"),
				Timestamp: time.Unix(0, 100),
			}).SetVal([]float32{-2}),
		)
	})
}

func TestUncompressedConfig(t *testing.T) {
	setClient := &testutil.SetClient{}
	c, err := ygnmi.NewClient(setClient, ygnmi.WithTarget("dut"))
	if err != nil {
		t.Fatalf("Unexpected error creating client: %v", err)
	}

	t.Run("YANG ordered list for uncompressed schema", func(t *testing.T) {
		configCheckFn(t, setClient, c,
			func(c *ygnmi.Client) (*ygnmi.Result, error) {
				om := &uexampleoc.OpenconfigWithlistval_Model_A_SingleKey_OrderedLists_OrderedList_OrderedMap{}
				ol, err := om.AppendNew("foo")
				if err != nil {
					t.Fatal(err)
				}
				ol.GetOrCreateConfig().SetKey("foo")
				ol.GetOrCreateConfig().SetValue(42)
				ol, err = om.AppendNew("bar")
				if err != nil {
					t.Fatal(err)
				}
				ol.GetOrCreateConfig().SetKey("bar")
				ol.GetOrCreateConfig().SetValue(43)
				ol, err = om.AppendNew("baz")
				if err != nil {
					t.Fatal(err)
				}
				ol.GetOrCreateConfig().SetKey("baz")
				ol.GetOrCreateConfig().SetValue(44)
				return ygnmi.Replace(context.Background(), c, uexampleocpath.Root().Model().A().SingleKey("foo").OrderedLists().OrderedListMap(), om)
			},
			&gpb.SetRequest{
				Prefix: &gpb.Path{
					Target: "dut",
				},
				Replace: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists/ordered-list"),
					Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`[
  {
    "openconfig-withlistval:config": {
      "key": "foo",
      "value": "42"
    },
    "openconfig-withlistval:key": "foo"
  },
  {
    "openconfig-withlistval:config": {
      "key": "bar",
      "value": "43"
    },
    "openconfig-withlistval:key": "bar"
  },
  {
    "openconfig-withlistval:config": {
      "key": "baz",
      "value": "44"
    },
    "openconfig-withlistval:key": "baz"
  }
]`)}},
				}},
			},
			&gpb.SetResponse{
				Prefix: &gpb.Path{
					Target: "dut",
				},
			},
			"",
			nil,
		)
	})

	t.Run("whole multi-keyed list for uncompressed schema", func(t *testing.T) {
		configCheckFn(t, setClient, c,
			func(c *ygnmi.Client) (*ygnmi.Result, error) {
				return ygnmi.Delete(context.Background(), c, uexampleocpath.Root().Model().B().MultiKeyMap())
			},
			&gpb.SetRequest{
				Prefix: &gpb.Path{
					Target: "dut",
				},
				Delete: []*gpb.Path{
					testutil.GNMIPath(t, "/model/b/multi-key"),
				},
			},
			&gpb.SetResponse{
				Prefix: &gpb.Path{
					Target: "dut",
				},
			},
			"",
			nil,
		)
	})
}

func TestUncompressedBatchGet(t *testing.T) {
	fakeGNMI, c := newClient(t)
	twoPath := testutil.GNMIPath(t, "/parent/child/state/two")
	threePath := testutil.GNMIPath(t, "/parent/child/state/three")
	aLeafSubPath := testutil.GNMIPath(t, "/remote-container/config/a-leaf")

	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		paths                []ygnmi.PathStruct
		wantSubscriptionPath []*gpb.Path
		wantVal              *ygnmi.Value[*uexampleoc.Root]
		wantErr              string
	}{{
		desc: "config and state leaves",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: aLeafSubPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}, {
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}, {
					Path: threePath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "ONE"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.PathStruct{
			uexampleocpath.Root().RemoteContainer().Config().ALeaf(),
			uexampleocpath.Root().Parent().Child().State().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantVal: (&ygnmi.Value[*uexampleoc.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&uexampleoc.Root{
			RemoteContainer: &uexampleoc.OpenconfigSimple_RemoteContainer{
				Config: &uexampleoc.OpenconfigSimple_RemoteContainer_Config{ALeaf: ygot.String("foo")},
			},
			Parent: &uexampleoc.OpenconfigSimple_Parent{Child: &uexampleoc.OpenconfigSimple_Parent_Child{
				State: &uexampleoc.OpenconfigSimple_Parent_Child_State{
					Two:   ygot.String("bar"),
					Three: uexampleoc.Simple_Parent_Child_Config_Three_ONE,
				}},
			},
		}),
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			b := &uexampleocpath.Batch{}
			b.AddPaths(tt.paths...)
			query := b.Query()
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

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[*uexampleoc.Root]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Lookup() returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", diff, got.ComplianceErrors)
			}
		})
	}
	t.Run("immutable query", func(t *testing.T) {
		fakeGNMI.Stub().Sync()
		b := &uexampleocpath.Batch{}
		b.AddPaths(uexampleocpath.Root().Model())
		q := b.Query()
		if _, err := ygnmi.Lookup(context.Background(), c, q); err != nil {
			t.Fatal(err)
		}
		verifySubscriptionPathsSent(t, fakeGNMI, testutil.GNMIPath(t, "/model"))
		b.AddPaths(uexampleocpath.Root().A(), uexampleocpath.Root().A().B())
		if _, err := ygnmi.Lookup(context.Background(), c, q); err != nil {
			t.Fatal(err)
		}
		verifySubscriptionPathsSent(t, fakeGNMI, testutil.GNMIPath(t, "/model"))
	})
}

func TestUncompressedCustomRootBatch(t *testing.T) {
	fakeGNMI, c := newClient(t)
	twoPath := testutil.GNMIPath(t, "/parent/child/state/two")

	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		paths                []ygnmi.PathStruct
		wantSubscriptionPath []*gpb.Path
		wantVal              *ygnmi.Value[*uexampleoc.OpenconfigSimple_Parent]
		wantAddErr           string
		wantLookupErr        string
	}{{
		desc: "not prefix",
		stub: func(s *testutil.Stubber) {},
		paths: []ygnmi.PathStruct{
			uexampleocpath.Root().Model(),
		},
		wantAddErr: "is not a prefix",
	}, {
		desc: "success",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.PathStruct{
			uexampleocpath.Root().Parent().Child().State().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			twoPath,
		},
		wantVal: (&ygnmi.Value[*uexampleoc.OpenconfigSimple_Parent]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/parent"),
		}).SetVal(&uexampleoc.OpenconfigSimple_Parent{
			Child: &uexampleoc.OpenconfigSimple_Parent_Child{
				State: &uexampleoc.OpenconfigSimple_Parent_Child_State{
					Two: ygot.String("foo"),
				},
			},
		}),
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			b := ygnmi.NewBatch(uexampleocpath.Root().Parent())
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

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[*uexampleoc.OpenconfigSimple_Parent]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Watch() returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", diff, got.ComplianceErrors)
			}
		})
	}
}

func TestUncompressedSetBatch(t *testing.T) {
	setClient := &testutil.SetClient{}
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
		desc: "leaf update replace delete",
		addPaths: func(sb *ygnmi.SetBatch) {
			cliPath, err := schemaless.NewConfig[string]("", "cli")
			if err != nil {
				t.Fatalf("Failed to create CLI ygnmi query: %v", err)
			}
			ygnmi.BatchUpdate(sb, cliPath, "hello, mercury")
			ygnmi.BatchUpdate(sb, uexampleocpath.Root().Parent().Child().Config().One(), "foo")
			ygnmi.BatchReplace(sb, cliPath, "hello, venus")
			ygnmi.BatchReplace(sb, uexampleocpath.Root().Parent().Child().Config().One(), "bar")
			ygnmi.BatchDelete(sb, cliPath)
			ygnmi.BatchDelete(sb, uexampleocpath.Root().Parent().Child().Config().One())
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

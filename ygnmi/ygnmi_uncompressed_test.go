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

package ygnmi_test

import (
	"context"
	"testing"
	"time"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygnmi/internal/uexampleoc"
	"github.com/openconfig/ygnmi/internal/uexampleoc/uexampleocpath"
	"github.com/openconfig/ygnmi/ygnmi"
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

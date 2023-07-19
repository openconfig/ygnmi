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

func TestUncompressed(t *testing.T) {
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
			ygnmi.SingletonQuery[string](uexampleocpath.Root().Model().A().SingleKey("foo").Config().Key().Query()),
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
			ygnmi.SingletonQuery[string](uexampleocpath.Root().Model().A().SingleKey("foo").State().Key().Query()),
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
			ygnmi.SingletonQuery[*uexampleoc.OpenconfigWithlistval_Model_A](uexampleocpath.Root().Model().A().Query()),
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
			ygnmi.SingletonQuery[map[string]*uexampleoc.OpenconfigWithlistval_Model_A_SingleKey](uexampleocpath.Root().Model().A().SingleKeyMap().Query()),
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
			ygnmi.SingletonQuery[string](uexampleocpath.Root().Model().A().SingleKey("foo").Config().Key().Query()),
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

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

	"github.com/openconfig/ygnmi/internal/exampleocunordered/exampleocunorderedpath"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygnmi/ygnmi"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestUnorderedOrderedMap(t *testing.T) {
	fakeGNMI, c := newClient(t)
	t.Run("success unordered ordered map leaf", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Prefix:    testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists"),
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, `ordered-list[key=foo]/config/key`),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
			}},
		}).Sync()

		lookupCheckFn(
			t, fakeGNMI, c,
			ygnmi.SingletonQuery[string](exampleocunorderedpath.Root().Model().SingleKey("foo").OrderedList("foo").Key().Config()),
			"",
			testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists/ordered-list[key=foo]/config/key"),
			(&ygnmi.Value[string]{
				Path:      testutil.GNMIPath(t, "/model/a/single-key[key=foo]/ordered-lists/ordered-list[key=foo]/config/key"),
				Timestamp: time.Unix(0, 100),
			}).SetVal("foo"),
		)
	})
}

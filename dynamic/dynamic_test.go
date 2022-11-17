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

package dynamic

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygnmi/ygnmi"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

type dynamicData struct {
	Name string
}

func TestGet(t *testing.T) {
	fakeGNMI, c := newClient(t)

	t.Run("primitive", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/foo/bar"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{
					StringVal: "test",
				}},
			}},
		}).Sync()
		query, err := NewConfig[string]("/foo/bar", "")
		if err != nil {
			t.Fatal(err)
		}

		want := "test"
		got, err := ygnmi.Get[string](context.Background(), c, query)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Lookup() returned unexpected diff: %s", diff)
		}
	})

	t.Run("struct", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/foo/bar"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonVal{
					JsonVal: []byte(`{ "name": "test" }`),
				}},
			}},
		}).Sync()
		query, err := NewConfig[*dynamicData]("/foo/bar", "")
		if err != nil {
			t.Fatal(err)
		}

		want := &dynamicData{
			Name: "test",
		}

		got, err := ygnmi.Get[*dynamicData](context.Background(), c, query)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Lookup() returned unexpected diff: %s", diff)
		}
	})
}

func TestGetAll(t *testing.T) {
	fakeGNMI, c := newClient(t)

	t.Run("primitive", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/foo/bar[name=test]"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{
					IntVal: 10,
				}},
			}, {
				Path: testutil.GNMIPath(t, "/foo/bar[name=test2]"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{
					IntVal: 15,
				}},
			}},
		}).Sync()
		query, err := NewWildcard[int]("/foo/bar[name=*]", "")
		if err != nil {
			t.Fatal(err)
		}

		want := []int{15, 10}
		got, err := ygnmi.GetAll(context.Background(), c, query)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Lookup() returned unexpected diff: %s", diff)
		}
	})

	t.Run("struct", func(t *testing.T) {
		fakeGNMI.Stub().Notification(&gpb.Notification{
			Timestamp: 100,
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "/foo/bar[name=test]"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{
					JsonIetfVal: []byte(`{ "name": "test" }`),
				}},
			}, {
				Path: testutil.GNMIPath(t, "/foo/bar[name=test2]"),
				Val: &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{
					JsonIetfVal: []byte(`{ "name": "test2" }`),
				}},
			}},
		}).Sync()
		query, err := NewWildcard[*dynamicData]("/foo/bar[name=*]", "")
		if err != nil {
			t.Fatal(err)
		}

		want := []*dynamicData{{
			Name: "test2",
		}, {
			Name: "test",
		}}

		got, err := ygnmi.GetAll(context.Background(), c, query)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Lookup() returned unexpected diff: %s", diff)
		}
	})
}

func newClient(t testing.TB) (*testutil.FakeGNMI, *ygnmi.Client) {
	fakeGNMI, err := testutil.StartGNMI(0)
	if err != nil {
		t.Fatal(err)
	}
	gnmiClient, err := fakeGNMI.Dial(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	c, err := ygnmi.NewClient(gnmiClient)
	if err != nil {
		t.Fatal(err)
	}
	return fakeGNMI, c
}

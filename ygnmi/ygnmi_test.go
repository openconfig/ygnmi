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
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygnmi/internal/exampleoc"
	"github.com/openconfig/ygnmi/internal/exampleoc/exampleocpath"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	ygottestutil "github.com/openconfig/ygot/testutil"
)

func TestLookup(t *testing.T) {
	fakeGNMI, c := newClient(t)
	leafPath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := exampleocpath.Root().RemoteContainer().ALeaf().State()

	leafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		inQuery              ygnmi.SingletonQuery[string]
		wantSubscriptionPath *gpb.Path
		wantVal              *ygnmi.Value[string]
		wantErr              string
	}{{
		desc:    "success update and sync",
		inQuery: lq,
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
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
		desc:    "success update and no sync",
		inQuery: lq,
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		wantSubscriptionPath: leafPath,
		wantVal: (&ygnmi.Value[string]{
			Path: leafPath,
		}),
	}, {
		desc:    "error multiple values",
		inQuery: lq,
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
			got, err := ygnmi.Lookup(context.Background(), c, tt.inQuery)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Lookup(ctx, c, %v) returned unexpected diff: %s", tt.inQuery, diff)
			}
			if err != nil {
				return
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath)
			checkJustReceived(t, got.RecvTimestamp)
			tt.wantVal.RecvTimestamp = got.RecvTimestamp

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[string]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Lookup(ctx, c, %v) returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", tt.inQuery, diff, got.ComplianceErrors)
			}
		})
	}

	rootPath := testutil.GNMIPath(t, "parent/child")
	strPath := testutil.GNMIPath(t, "parent/child/state/one")
	enumPath := testutil.GNMIPath(t, "parent/child/state/three")
	strCfgPath := testutil.GNMIPath(t, "parent/child/config/one")

	configQuery := exampleocpath.Root().Parent().Child().Config()
	stateQuery := exampleocpath.Root().Parent().Child().State()

	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		inQuery              ygnmi.SingletonQuery[*exampleoc.Parent_Child]
		wantSubscriptionPath *gpb.Path
		wantVal              *ygnmi.Value[*exampleoc.Parent_Child]
		wantErr              string
	}{{
		desc: "success one update and state false",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: strCfgPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery:              configQuery,
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(&exampleoc.Parent_Child{
			One: ygot.String("foo"),
		}),
	}, {
		desc: "success one update and state true",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery:              stateQuery,
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(&exampleoc.Parent_Child{
			One: ygot.String("foo"),
		}),
	}, {
		desc: "success one update with prefix",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Prefix:    testutil.GNMIPath(t, "parent"),
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "child/state/one"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery:              stateQuery,
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(&exampleoc.Parent_Child{
			One: ygot.String("foo"),
		}),
	}, {
		desc: "success ignore state update when state false",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		inQuery:              configQuery,
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}),
	}, {
		desc: "success ignore non-state update when state true",
		stub: func(s *testutil.Stubber) {
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
		wantVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}),
	}, {
		desc: "success multiple updates in single notification",
		stub: func(s *testutil.Stubber) {
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
		wantVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal(&exampleoc.Parent_Child{
			One:   ygot.String("foo"),
			Three: exampleoc.Child_Three_ONE,
		}),
	}, {
		desc: "success multiple notifications",
		stub: func(s *testutil.Stubber) {
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
		wantVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Path:      rootPath,
			Timestamp: time.Unix(0, 102),
		}).SetVal(&exampleoc.Parent_Child{
			One:   ygot.String("foo"),
			Three: exampleoc.Child_Three_ONE,
		}),
	}, {
		desc: "success no values",
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		inQuery:              stateQuery,
		wantSubscriptionPath: rootPath,
		wantVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Path: rootPath,
		}),
	}}

	for _, tt := range tests {
		t.Run("nonleaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			got, err := ygnmi.Lookup(context.Background(), c, tt.inQuery)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Lookup(ctx, c, %v) returned unexpected diff: %s", tt.inQuery, diff)
			}
			if err != nil {
				return
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath)
			checkJustReceived(t, got.RecvTimestamp)
			tt.wantVal.RecvTimestamp = got.RecvTimestamp

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[*exampleoc.Parent_Child]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Lookup(ctx, c, %v) returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", tt.inQuery, diff, got.ComplianceErrors)
			}
		})
	}
	t.Run("use get", func(t *testing.T) {
		fakeGNMI.Stub().GetResponse(&gpb.GetResponse{
			Notification: []*gpb.Notification{{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`"foo"`)}},
				}},
			}},
		})
		wantGetRequest := &gpb.GetRequest{
			Encoding: gpb.Encoding_JSON_IETF,
			Type:     gpb.GetRequest_STATE,
			Prefix:   &gpb.Path{},
			Path:     []*gpb.Path{leafPath},
		}
		wantVal := (&ygnmi.Value[string]{
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		}).SetVal("foo")

		got, err := ygnmi.Lookup(context.Background(), c, exampleocpath.Root().RemoteContainer().ALeaf().State(), ygnmi.WithUseGet())
		if err != nil {
			t.Fatalf("Lookup() returned unexpected error: %v", err)
		}
		if diff := cmp.Diff(wantVal, got, cmp.AllowUnexported(ygnmi.Value[string]{}), cmpopts.IgnoreFields(ygnmi.Value[string]{}, "RecvTimestamp"), protocmp.Transform()); diff != "" {
			t.Errorf("Lookup() returned unexpected diff: %s", diff)
		}
		if diff := cmp.Diff(wantGetRequest, fakeGNMI.GetRequests()[0], protocmp.Transform()); diff != "" {
			t.Errorf("Lookup() GetRequest different from expected: %s", diff)
		}
	})
}

func TestGet(t *testing.T) {
	fakeGNMI, c := newClient(t)
	leafPath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := exampleocpath.Root().RemoteContainer().ALeaf().State()

	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		want                 string
		wantVal              string
		wantErr              string
	}{{
		desc: "value present",
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		wantSubscriptionPath: leafPath,
		wantErr:              "value not present",
	}, {
		desc: "error nil update",
		stub: func(s *testutil.Stubber) {
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
			got, err := ygnmi.Get(context.Background(), c, lq)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Get(ctx, c, %v) returned unexpected diff: %s", lq, diff)
			}
			if err != nil {
				return
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath)

			if diff := cmp.Diff(tt.wantVal, got); diff != "" {
				t.Errorf("Get(ctx, c, %v) returned unexpected diff (-want,+got):\n %s", lq, diff)
			}
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
		})
		wantGetRequest := &gpb.GetRequest{
			Encoding: gpb.Encoding_JSON_IETF,
			Type:     gpb.GetRequest_CONFIG,
			Prefix:   &gpb.Path{},
			Path:     []*gpb.Path{testutil.GNMIPath(t, "/remote-container/config/a-leaf")},
		}
		wantVal := "foo"

		got, err := ygnmi.Get[string](context.Background(), c, exampleocpath.Root().RemoteContainer().ALeaf().Config(), ygnmi.WithUseGet())
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

func TestWatch(t *testing.T) {
	fakeGNMI, client := newClient(t)
	path := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := exampleocpath.Root().RemoteContainer().ALeaf().State()

	startTime := time.Now()
	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[string]
		wantVals             []*ygnmi.Value[string]
		wantErr              string
		wantMode             gpb.SubscriptionMode
		opts                 []ygnmi.Option
	}{{
		desc: "single notif and pred true",
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		desc: "single notif and pred false error EOF",
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		dur:     -1 * time.Second,
		wantErr: "context deadline exceeded",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			i := 0
			ctx, cancel := context.WithTimeout(context.Background(), tt.dur)
			defer cancel()
			w := ygnmi.Watch(ctx, client, lq, func(v *ygnmi.Value[string]) error {
				if i > len(tt.wantVals) {
					t.Fatalf("Predicate(%d) expected no more values but got: %+v", i, v)
				}
				if diff := cmp.Diff(tt.wantVals[i], v, cmpopts.IgnoreFields(ygnmi.Value[string]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[string]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate(%d) got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", i, diff, v.ComplianceErrors)
				}
				val, present := v.Val()
				i++
				if present && val == "foo" {
					return nil
				}
				return ygnmi.Continue
			}, tt.opts...)
			val, err := w.Await()
			if i < len(tt.wantVals) {
				t.Errorf("Predicate received too few values: got %d, want %d", i, len(tt.wantVals))
			}
			verifySubscriptionModesSent(t, fakeGNMI, tt.wantMode)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Await() returned unexpected diff: %s", diff)
			}
			if val != nil {
				checkJustReceived(t, val.RecvTimestamp)
				tt.wantLastVal.RecvTimestamp = val.RecvTimestamp
			}
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(ygnmi.Value[string]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}

	t.Run("multiple awaits", func(t *testing.T) {
		fakeGNMI.Stub().Sync()
		w := ygnmi.Watch(context.Background(), client, exampleocpath.Root().RemoteContainer().ALeaf().State(), func(v *ygnmi.Value[string]) error { return nil })
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
	nonLeafQuery := exampleocpath.Root().Parent().Child().State()

	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[*exampleoc.Parent_Child]
		wantVals             []*ygnmi.Value[*exampleoc.Parent_Child]
		wantErr              string
	}{{
		desc: "single notif and pred false",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: strPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[*exampleoc.Parent_Child]{
			(&ygnmi.Value[*exampleoc.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleoc.Parent_Child{
				One: ygot.String("bar"),
			}),
		},
		wantErr:              "EOF",
		wantSubscriptionPath: rootPath,
		wantLastVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Timestamp: startTime,
			Path:      rootPath,
		}).SetVal(&exampleoc.Parent_Child{
			One: ygot.String("bar"),
		}),
	}, {
		desc: "multiple notif and pred true",
		stub: func(s *testutil.Stubber) {
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
		wantVals: []*ygnmi.Value[*exampleoc.Parent_Child]{
			(&ygnmi.Value[*exampleoc.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleoc.Parent_Child{
				One: ygot.String("foo"),
			}),
			(&ygnmi.Value[*exampleoc.Parent_Child]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      rootPath,
			}).SetVal(&exampleoc.Parent_Child{
				Three: exampleoc.Child_Three_ONE,
				One:   ygot.String("foo"),
			}),
		},
		wantSubscriptionPath: rootPath,
		wantLastVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}).SetVal(&exampleoc.Parent_Child{
			Three: exampleoc.Child_Three_ONE,
			One:   ygot.String("foo"),
		}),
	}, {
		desc: "multiple notif before sync",
		stub: func(s *testutil.Stubber) {
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
		wantVals: []*ygnmi.Value[*exampleoc.Parent_Child]{
			(&ygnmi.Value[*exampleoc.Parent_Child]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      rootPath,
			}).SetVal(&exampleoc.Parent_Child{
				Three: exampleoc.Child_Three_ONE,
				One:   ygot.String("foo"),
			})},
		wantSubscriptionPath: rootPath,
		wantLastVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}).SetVal(&exampleoc.Parent_Child{
			Three: exampleoc.Child_Three_ONE,
			One:   ygot.String("foo"),
		}),
	}, {
		desc: "delete leaf in container",
		stub: func(s *testutil.Stubber) {
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
		wantVals: []*ygnmi.Value[*exampleoc.Parent_Child]{
			(&ygnmi.Value[*exampleoc.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleoc.Parent_Child{
				Three: exampleoc.Child_Three_ONE,
				One:   ygot.String("bar"),
			}),
			(&ygnmi.Value[*exampleoc.Parent_Child]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      rootPath,
			}).SetVal(&exampleoc.Parent_Child{
				Three: exampleoc.Child_Three_ONE,
			}),
		},
		wantSubscriptionPath: rootPath,
		wantErr:              "EOF",
		wantLastVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}).SetVal(&exampleoc.Parent_Child{
			Three: exampleoc.Child_Three_ONE,
		}),
	}}

	for _, tt := range nonLeafTests {
		t.Run("nonleaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			w := ygnmi.Watch(context.Background(), client, nonLeafQuery, func(v *ygnmi.Value[*exampleoc.Parent_Child]) error {
				if len(tt.wantVals) == 0 {
					t.Fatalf("Predicate expected no more values but got: %+v", v)
				}
				if diff := cmp.Diff(tt.wantVals[0], v, cmpopts.IgnoreFields(ygnmi.Value[*exampleoc.Parent_Child]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[*exampleoc.Parent_Child]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", diff, v.ComplianceErrors)
				}
				tt.wantVals = tt.wantVals[1:]
				val, present := v.Val()
				if present && val.One != nil && *val.One == "foo" && val.Three == exampleoc.Child_Three_ONE {
					return nil
				}
				return ygnmi.Continue
			})
			val, err := w.Await()
			if len(tt.wantVals) > 0 {
				t.Errorf("Predicate received too few values, remaining: %+v", tt.wantVals)
			}
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
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(ygnmi.Value[*exampleoc.Parent_Child]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestAwait(t *testing.T) {
	fakeGNMI, client := newClient(t)
	path := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := exampleocpath.Root().RemoteContainer().ALeaf().State()

	startTime := time.Now()
	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantVal              *ygnmi.Value[string]
		wantErr              string
		wantMode             gpb.SubscriptionMode
		opts                 []ygnmi.Option
	}{{
		desc: "value never equal",
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
	nonLeafQuery := exampleocpath.Root().Parent().Child().State()

	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[*exampleoc.Parent_Child]
		wantErr              string
	}{{
		desc: "value never equal",
		stub: func(s *testutil.Stubber) {
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
		wantLastVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Timestamp: startTime,
			Path:      rootPath,
		}).SetVal(&exampleoc.Parent_Child{
			One: ygot.String("bar"),
		}),
	}, {
		desc: "success",
		stub: func(s *testutil.Stubber) {
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
		wantLastVal: (&ygnmi.Value[*exampleoc.Parent_Child]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}).SetVal(&exampleoc.Parent_Child{
			Three: exampleoc.Child_Three_ONE,
			One:   ygot.String("foo"),
		}),
	}}

	for _, tt := range nonLeafTests {
		t.Run("nonleaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			val, err := ygnmi.Await(context.Background(), client, nonLeafQuery, &exampleoc.Parent_Child{One: ygot.String("foo"), Three: exampleoc.Child_Three_ONE})
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
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(ygnmi.Value[*exampleoc.Parent_Child]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestCollect(t *testing.T) {
	fakeGNMI, client := newClient(t)
	path := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := exampleocpath.Root().RemoteContainer().ALeaf().State()

	startTime := time.Now()
	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[string]
		wantErr              string
		wantMode             gpb.SubscriptionMode
		opts                 []ygnmi.Option
	}{{
		desc: "no values",
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
			ctx, cancel := context.WithTimeout(context.Background(), tt.dur)
			defer cancel()
			vals, err := ygnmi.Collect(ctx, client, lq, tt.opts...).Await()
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Await() returned unexpected diff: %s", diff)
			}
			for _, val := range vals {
				checkJustReceived(t, val.RecvTimestamp)
			}
			if diff := cmp.Diff(tt.wantVals, vals, cmpopts.IgnoreFields(ygnmi.Value[string]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[string]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}

	rootPath := testutil.GNMIPath(t, "parent/child")
	strPath := testutil.GNMIPath(t, "parent/child/state/one")
	enumPath := testutil.GNMIPath(t, "parent/child/state/three")
	startTime = time.Now()
	nonLeafQuery := exampleocpath.Root().Parent().Child().State()

	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[*exampleoc.Parent_Child]
		wantErr              string
	}{{
		desc: "one val",
		stub: func(s *testutil.Stubber) {
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
		wantVals: []*ygnmi.Value[*exampleoc.Parent_Child]{
			(&ygnmi.Value[*exampleoc.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleoc.Parent_Child{
				One: ygot.String("bar"),
			}),
		},
	}, {
		desc: "multiple values",
		stub: func(s *testutil.Stubber) {
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
		wantVals: []*ygnmi.Value[*exampleoc.Parent_Child]{
			(&ygnmi.Value[*exampleoc.Parent_Child]{
				Timestamp: startTime,
				Path:      rootPath,
			}).SetVal(&exampleoc.Parent_Child{
				One: ygot.String("foo"),
			}),
			(&ygnmi.Value[*exampleoc.Parent_Child]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      rootPath,
			}).SetVal(&exampleoc.Parent_Child{
				Three: exampleoc.Child_Three_ONE,
				One:   ygot.String("foo"),
			}),
		},
	}}

	for _, tt := range nonLeafTests {
		t.Run("nonleaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			vals, err := ygnmi.Collect(context.Background(), client, nonLeafQuery).Await()
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Errorf("Await() returned unexpected diff: %s", diff)
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath)
			for _, val := range vals {
				checkJustReceived(t, val.RecvTimestamp)
			}
			if diff := cmp.Diff(tt.wantVals, vals, cmpopts.IgnoreFields(ygnmi.Value[*exampleoc.Parent_Child]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[*exampleoc.Parent_Child]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestLookupAll(t *testing.T) {
	fakeGNMI, c := newClient(t)
	leafPath := testutil.GNMIPath(t, "model/a/single-key[key=*]/state/value")
	lq := exampleocpath.Root().Model().SingleKeyAny().Value().State()

	leafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[int64]
		wantErr              string
	}{{
		desc: "success one value",
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		wantVals:             nil,
		wantSubscriptionPath: leafPath,
	}, {
		desc: "success multiples value in same notification",
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
			got, err := ygnmi.LookupAll(context.Background(), c, lq)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("LookupAll(ctx, c, %v) returned unexpected diff: %s", lq, diff)
			}
			if err != nil {
				return
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath)
			for _, val := range got {
				checkJustReceived(t, val.RecvTimestamp)
			}
			if diff := cmp.Diff(tt.wantVals, got, cmp.AllowUnexported(ygnmi.Value[int64]{}), cmpopts.IgnoreFields(ygnmi.Value[int64]{}, "RecvTimestamp"), protocmp.Transform()); diff != "" {
				t.Errorf("LookupAll() returned unexpected diff (-want,+got):\n%s", diff)
			}
		})
	}

	nonLeafPath := testutil.GNMIPath(t, "model/a/single-key[key=*]")
	nonLeafQ := exampleocpath.Root().Model().SingleKeyAny().State()
	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[*exampleoc.Model_SingleKey]
		wantErr              string
	}{{
		desc: "one value",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[*exampleoc.Model_SingleKey]{
			(&ygnmi.Value[*exampleoc.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=10]"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(&exampleoc.Model_SingleKey{
				Value: ygot.Int64(10),
			}),
		},
		wantSubscriptionPath: nonLeafPath,
	}, {
		desc: "multiple values",
		stub: func(s *testutil.Stubber) {
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
		wantVals: []*ygnmi.Value[*exampleoc.Model_SingleKey]{
			(&ygnmi.Value[*exampleoc.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=10]"),
				Timestamp: time.Unix(0, 100),
			}).SetVal(&exampleoc.Model_SingleKey{
				Value: ygot.Int64(100),
				Key:   ygot.String("10"),
			}),
			(&ygnmi.Value[*exampleoc.Model_SingleKey]{
				Path:      testutil.GNMIPath(t, "model/a/single-key[key=11]"),
				Timestamp: time.Unix(0, 101),
			}).SetVal(&exampleoc.Model_SingleKey{
				Value: ygot.Int64(101),
				Key:   ygot.String("11"),
			}),
		},
		wantSubscriptionPath: nonLeafPath,
	}, {
		desc: "no values",
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		wantVals:             nil,
		wantSubscriptionPath: nonLeafPath,
	}}
	for _, tt := range nonLeafTests {
		t.Run("nonLeaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			got, err := ygnmi.LookupAll(context.Background(), c, nonLeafQ)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("LookupAll(ctx, c, %v) returned unexpected diff: %s", nonLeafQ, diff)
			}
			if err != nil {
				return
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath)
			for _, val := range got {
				checkJustReceived(t, val.RecvTimestamp)
			}
			if diff := cmp.Diff(tt.wantVals, got, cmp.AllowUnexported(ygnmi.Value[*exampleoc.Model_SingleKey]{}), cmpopts.IgnoreFields(ygnmi.Value[*exampleoc.Model_SingleKey]{}, "RecvTimestamp"), protocmp.Transform()); diff != "" {
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
		})
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

func TestGetAll(t *testing.T) {
	fakeGNMI, c := newClient(t)
	leafPath := testutil.GNMIPath(t, "model/a/single-key[key=*]/state/value")
	lq := exampleocpath.Root().Model().SingleKeyAny().Value().State()

	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []int64
		wantErr              string
	}{{
		desc: "success",
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		})
		wantGetRequest := &gpb.GetRequest{
			Encoding: gpb.Encoding_JSON_IETF,
			Type:     gpb.GetRequest_STATE,
			Prefix:   &gpb.Path{},
			Path:     []*gpb.Path{leafPath},
		}
		wantVal := []int64{1}

		got, err := ygnmi.GetAll(context.Background(), c, exampleocpath.Root().Model().SingleKeyAny().Value().State(), ygnmi.WithUseGet())
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

func TestWatchAll(t *testing.T) {
	fakeGNMI, client := newClient(t)
	leafQueryPath := testutil.GNMIPath(t, "model/a/single-key[key=*]/state/value")
	key10Path := testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value")
	key11Path := testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value")

	startTime := time.Now()
	lq := exampleocpath.Root().Model().SingleKeyAny().Value().State()
	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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

	nonLeafQ := exampleocpath.Root().Model().SingleKeyAny().State()
	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[*exampleoc.Model_SingleKey]
		wantVals             []*ygnmi.Value[*exampleoc.Model_SingleKey]
		wantErr              string
	}{{
		desc: "predicate not true",
		dur:  time.Second,
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: key10Path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: nonLeafPath,
		wantVals: []*ygnmi.Value[*exampleoc.Model_SingleKey]{
			(&ygnmi.Value[*exampleoc.Model_SingleKey]{
				Timestamp: startTime,
				Path:      nonLeafKey10Path,
			}).SetVal(&exampleoc.Model_SingleKey{
				Value: ygot.Int64(100),
			}),
		},
		wantLastVal: (&ygnmi.Value[*exampleoc.Model_SingleKey]{
			Timestamp: startTime,
			Path:      nonLeafKey10Path,
		}).SetVal(&exampleoc.Model_SingleKey{
			Value: ygot.Int64(100),
		}),
		wantErr: "EOF",
	}, {
		desc: "predicate becomes true",
		dur:  time.Second,
		stub: func(s *testutil.Stubber) {
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
		wantVals: []*ygnmi.Value[*exampleoc.Model_SingleKey]{
			(&ygnmi.Value[*exampleoc.Model_SingleKey]{
				Timestamp: startTime,
				Path:      nonLeafKey10Path,
			}).SetVal(&exampleoc.Model_SingleKey{
				Value: ygot.Int64(100),
			}),
			(&ygnmi.Value[*exampleoc.Model_SingleKey]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      nonLeafKey11Path,
			}).SetVal(&exampleoc.Model_SingleKey{
				Value: ygot.Int64(101),
			}),
		},
		wantLastVal: (&ygnmi.Value[*exampleoc.Model_SingleKey]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      nonLeafKey11Path,
		}).SetVal(&exampleoc.Model_SingleKey{
			Value: ygot.Int64(101),
		}),
	}}
	for _, tt := range nonLeafTests {
		t.Run("nonLeaf "+tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			i := 0
			ctx, cancel := context.WithTimeout(context.Background(), tt.dur)
			defer cancel()
			var key10Cond, key11Cond bool

			w := ygnmi.WatchAll(ctx, client, nonLeafQ, func(v *ygnmi.Value[*exampleoc.Model_SingleKey]) error {
				if i > len(tt.wantVals) {
					t.Fatalf("Predicate(%d) expected no more values but got: %+v", i, v)
				}
				if diff := cmp.Diff(tt.wantVals[i], v, cmpopts.IgnoreFields(ygnmi.Value[*exampleoc.Model_SingleKey]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[*exampleoc.Model_SingleKey]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate(%d) got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", i, diff, v.ComplianceErrors)
				}
				val, present := v.Val()
				key10Cond = key10Cond || (present && proto.Equal(v.Path, nonLeafKey10Path) && *val.Value == 100)
				key11Cond = key11Cond || (present && proto.Equal(v.Path, nonLeafKey11Path) && *val.Value == 101)
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
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(ygnmi.Value[*exampleoc.Model_SingleKey]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestCollectAll(t *testing.T) {
	fakeGNMI, client := newClient(t)
	leafQueryPath := testutil.GNMIPath(t, "model/a/single-key[key=*]/state/value")
	key10Path := testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value")
	key11Path := testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value")

	startTime := time.Now()
	lq := exampleocpath.Root().Model().SingleKeyAny().Value().State()
	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[int64]
		wantErr              string
		wantMode             gpb.SubscriptionMode
		opts                 []ygnmi.Option
	}{{
		desc: "no values",
		dur:  time.Second,
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		wantErr:              "EOF",
		wantSubscriptionPath: leafQueryPath,
		wantVals:             nil,
	}, {
		desc: "no values with custom mode",
		dur:  time.Second,
		stub: func(s *testutil.Stubber) {
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
		stub: func(s *testutil.Stubber) {
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

	nonLeafQ := exampleocpath.Root().Model().SingleKeyAny().State()
	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[*exampleoc.Model_SingleKey]
		wantErr              string
	}{{
		desc: "no values",
		dur:  time.Second,
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		wantSubscriptionPath: nonLeafPath,
		wantVals:             nil,
		wantErr:              "EOF",
	}, {
		desc: "multiple values",
		dur:  time.Second,
		stub: func(s *testutil.Stubber) {
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
		wantVals: []*ygnmi.Value[*exampleoc.Model_SingleKey]{
			(&ygnmi.Value[*exampleoc.Model_SingleKey]{
				Timestamp: startTime,
				Path:      nonLeafKey10Path,
			}).SetVal(&exampleoc.Model_SingleKey{
				Value: ygot.Int64(100),
			}),
			(&ygnmi.Value[*exampleoc.Model_SingleKey]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      nonLeafKey11Path,
			}).SetVal(&exampleoc.Model_SingleKey{
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
			if diff := cmp.Diff(tt.wantVals, vals, cmpopts.IgnoreFields(ygnmi.Value[*exampleoc.Model_SingleKey]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[*exampleoc.Model_SingleKey]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	setClient := &fakeGNMISetClient{}
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
			return ygnmi.Update(context.Background(), c, exampleocpath.Root().Parent().Child().One().Config(), "10")
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
			return ygnmi.Update(context.Background(), c, exampleocpath.Root().Parent().Child().Three().Config(), exampleoc.Child_Three_ONE)
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
			return ygnmi.Update(context.Background(), c, exampleocpath.Root().Parent().Child().Config(), &exampleoc.Parent_Child{One: ygot.String("10")})
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
			return ygnmi.Update(context.Background(), c, exampleocpath.Root().Parent().Child().One().Config(), "10")
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
			if diff := cmp.Diff(tt.wantRequest, setClient.requests[0], protocmp.Transform()); diff != "" {
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

func TestReplace(t *testing.T) {
	setClient := &fakeGNMISetClient{}
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
			return ygnmi.Replace(context.Background(), c, exampleocpath.Root().Parent().Child().One().Config(), "10")
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
			return ygnmi.Replace(context.Background(), c, exampleocpath.Root().Parent().Child().Three().Config(), exampleoc.Child_Three_ONE)

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
			return ygnmi.Replace(context.Background(), c, exampleocpath.Root().Parent().Child().Config(), &exampleoc.Parent_Child{One: ygot.String("10")})
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
		desc: "server error",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Replace(context.Background(), c, exampleocpath.Root().Parent().Child().One().Config(), "10")
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
			setClient.Reset()
			setClient.AddResponse(tt.stubResponse, tt.stubErr)

			got, err := tt.op(client)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Replace() returned unexpected diff: %s", diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.wantRequest, setClient.requests[0], protocmp.Transform()); diff != "" {
				t.Errorf("Replace() sent unexpected request (-want,+got):\n%s", diff)
			}
			want := &ygnmi.Result{
				RawResponse: tt.stubResponse,
				Timestamp:   time.Unix(0, tt.stubResponse.GetTimestamp()),
			}
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Replace() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	setClient := &fakeGNMISetClient{}
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
			return ygnmi.Delete(context.Background(), c, exampleocpath.Root().Parent().Child().One().Config())
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
		desc: "server error",
		op: func(c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Delete(context.Background(), c, exampleocpath.Root().Parent().Child().One().Config())
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
			setClient.Reset()
			setClient.AddResponse(tt.stubResponse, tt.stubErr)

			got, err := tt.op(client)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Delete() returned unexpected diff: %s", diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.wantRequest, setClient.requests[0], protocmp.Transform()); diff != "" {
				t.Errorf("Delete() sent unexpected request (-want,+got):\n%s", diff)
			}
			want := &ygnmi.Result{
				RawResponse: tt.stubResponse,
				Timestamp:   time.Unix(0, tt.stubResponse.GetTimestamp()),
			}
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Delete() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestBatchGet(t *testing.T) {
	fakeGNMI, c := newClient(t)
	aLeafStatePath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	aLeafConfigPath := testutil.GNMIPath(t, "/remote-container/config/a-leaf")
	twoPath := testutil.GNMIPath(t, "/parent/child/state/two")
	aLeafSubPath := testutil.GNMIPath(t, "/remote-container/*/a-leaf")

	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		config               bool
		paths                []ygnmi.PathStruct
		wantSubscriptionPath []*gpb.Path
		wantVal              *ygnmi.Value[*exampleoc.Root]
		wantErr              string
	}{{
		desc: "state leaves",
		stub: func(s *testutil.Stubber) {
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
			exampleocpath.Root().RemoteContainer().ALeaf(),
			exampleocpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantVal: (&ygnmi.Value[*exampleoc.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleoc.Root{
			RemoteContainer: &exampleoc.RemoteContainer{ALeaf: ygot.String("foo")},
			Parent:          &exampleoc.Parent{Child: &exampleoc.Parent_Child{Two: ygot.String("bar")}},
		}),
	}, {
		desc:   "config ignore state leaves",
		config: true,
		stub: func(s *testutil.Stubber) {
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
			exampleocpath.Root().RemoteContainer().ALeaf(),
			exampleocpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantVal: (&ygnmi.Value[*exampleoc.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleoc.Root{
			RemoteContainer: &exampleoc.RemoteContainer{},
			Parent:          &exampleoc.Parent{Child: &exampleoc.Parent_Child{}},
		}),
	}, {
		desc: "non leaves",
		stub: func(s *testutil.Stubber) {
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
			exampleocpath.Root().RemoteContainer(),
			exampleocpath.Root().Parent(),
		},
		wantSubscriptionPath: []*gpb.Path{
			testutil.GNMIPath(t, "/remote-container"),
			testutil.GNMIPath(t, "/parent"),
		},
		wantVal: (&ygnmi.Value[*exampleoc.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleoc.Root{
			RemoteContainer: &exampleoc.RemoteContainer{ALeaf: ygot.String("foo")},
			Parent:          &exampleoc.Parent{Child: &exampleoc.Parent_Child{Two: ygot.String("bar")}},
		}),
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			b := &exampleocpath.Batch{}
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

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[*exampleoc.Root]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Lookup() returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", diff, got.ComplianceErrors)
			}
		})
	}
	t.Run("immutable query", func(t *testing.T) {
		fakeGNMI.Stub().Sync()
		b := &exampleocpath.Batch{}
		b.AddPaths(exampleocpath.Root().Model())
		q := b.State()
		if _, err := ygnmi.Lookup(context.Background(), c, q); err != nil {
			t.Fatal(err)
		}
		verifySubscriptionPathsSent(t, fakeGNMI, testutil.GNMIPath(t, "/model"))
		b.AddPaths(exampleocpath.Root().A(), exampleocpath.Root().A().B())
		if _, err := ygnmi.Lookup(context.Background(), c, q); err != nil {
			t.Fatal(err)
		}
		verifySubscriptionPathsSent(t, fakeGNMI, testutil.GNMIPath(t, "/model"))
	})
}

func TestBatchWatch(t *testing.T) {
	fakeGNMI, c := newClient(t)
	aLeafStatePath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	twoPath := testutil.GNMIPath(t, "/parent/child/state/two")
	aLeafSubPath := testutil.GNMIPath(t, "/remote-container/*/a-leaf")

	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		config               bool
		paths                []ygnmi.PathStruct
		wantSubscriptionPath []*gpb.Path
		wantVal              *ygnmi.Value[*exampleoc.Root]
		wantErr              string
	}{{
		desc: "predicate true",
		stub: func(s *testutil.Stubber) {
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
			exampleocpath.Root().RemoteContainer().ALeaf(),
			exampleocpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantVal: (&ygnmi.Value[*exampleoc.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleoc.Root{
			RemoteContainer: &exampleoc.RemoteContainer{ALeaf: ygot.String("foo")},
			Parent:          &exampleoc.Parent{Child: &exampleoc.Parent_Child{Two: ygot.String("bar")}},
		}),
	}, {
		desc: "predicate false true false",
		stub: func(s *testutil.Stubber) {
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
			exampleocpath.Root().RemoteContainer().ALeaf(),
			exampleocpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantVal: (&ygnmi.Value[*exampleoc.Root]{
			Timestamp: time.Unix(0, 101),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleoc.Root{
			RemoteContainer: &exampleoc.RemoteContainer{ALeaf: ygot.String("foo")},
			Parent:          &exampleoc.Parent{Child: &exampleoc.Parent_Child{Two: ygot.String("bar")}},
		}),
	}, {
		desc:   "predicate false",
		config: true,
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: twoPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		paths: []ygnmi.PathStruct{
			exampleocpath.Root().RemoteContainer().ALeaf(),
			exampleocpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			aLeafSubPath,
			twoPath,
		},
		wantErr: "EOF",
	}, {
		desc: "non leaves",
		stub: func(s *testutil.Stubber) {
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
			exampleocpath.Root().RemoteContainer(),
			exampleocpath.Root().Parent(),
		},
		wantSubscriptionPath: []*gpb.Path{
			testutil.GNMIPath(t, "/remote-container"),
			testutil.GNMIPath(t, "/parent"),
		},
		wantVal: (&ygnmi.Value[*exampleoc.Root]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/"),
		}).SetVal(&exampleoc.Root{
			RemoteContainer: &exampleoc.RemoteContainer{ALeaf: ygot.String("foo")},
			Parent:          &exampleoc.Parent{Child: &exampleoc.Parent_Child{Two: ygot.String("bar")}},
		}),
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			b := &exampleocpath.Batch{}
			b.AddPaths(tt.paths...)
			query := b.State()
			if tt.config {
				query = b.Config()
			}

			got, err := ygnmi.Watch(context.Background(), c, query, func(v *ygnmi.Value[*exampleoc.Root]) error {
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

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[*exampleoc.Root]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Watch() returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", diff, got.ComplianceErrors)
			}
		})
	}
}

func TestCustomRootBatch(t *testing.T) {
	fakeGNMI, c := newClient(t)
	twoPath := testutil.GNMIPath(t, "/parent/child/state/two")

	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		paths                []ygnmi.PathStruct
		wantSubscriptionPath []*gpb.Path
		wantVal              *ygnmi.Value[*exampleoc.Parent]
		wantAddErr           string
		wantLookupErr        string
	}{{
		desc: "not prefix",
		stub: func(s *testutil.Stubber) {},
		paths: []ygnmi.PathStruct{
			exampleocpath.Root().Model(),
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
			exampleocpath.Root().Parent().Child().Two(),
		},
		wantSubscriptionPath: []*gpb.Path{
			twoPath,
		},
		wantVal: (&ygnmi.Value[*exampleoc.Parent]{
			Timestamp: time.Unix(0, 100),
			Path:      testutil.GNMIPath(t, "/parent"),
		}).SetVal(&exampleoc.Parent{
			Child: &exampleoc.Parent_Child{
				Two: ygot.String("foo"),
			},
		}),
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			b := ygnmi.NewBatch(exampleocpath.Root().Parent().State())
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

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(ygnmi.Value[*exampleoc.Parent]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Watch() returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", diff, got.ComplianceErrors)
			}
		})
	}

}

func TestSetBatch(t *testing.T) {
	setClient := &fakeGNMISetClient{}
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
			ygnmi.BatchUpdate(sb, exampleocpath.Root().Parent().Child().One().Config(), "foo")
			ygnmi.BatchReplace(sb, exampleocpath.Root().Parent().Child().One().Config(), "bar")
			ygnmi.BatchDelete(sb, exampleocpath.Root().Parent().Child().One().Config())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/config/one"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"foo\"")}},
			}},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/config/one"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"bar\"")}},
			}},
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
		desc: "non leaf update delete",
		addPaths: func(sb *ygnmi.SetBatch) {
			ygnmi.BatchUpdate(sb, exampleocpath.Root().Parent().Child().Config(), &exampleoc.Parent_Child{One: ygot.String("foo")})
			ygnmi.BatchDelete(sb, exampleocpath.Root().Parent().Child().One().Config())
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "parent/child/"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("{\n  \"openconfig-simple:config\": {\n    \"one\": \"foo\"\n  }\n}")}},
			}},
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
		desc: "server error",
		addPaths: func(sb *ygnmi.SetBatch) {
			ygnmi.BatchDelete(sb, exampleocpath.Root().Parent().Child().One().Config())
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
			if diff := cmp.Diff(tt.wantRequest, setClient.requests[0], protocmp.Transform()); diff != "" {
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

type fakeGNMISetClient struct {
	gpb.GNMIClient
	// responses are the gNMI responses to return from calls to Set.
	responses []*gpb.SetResponse
	// requests received by the client are stored in the slice.
	requests []*gpb.SetRequest
	// responseErrs are the errors to return from calls to Set.
	responseErrs []error
	// i is index current index of the response and error to return.
	i int
}

func (f *fakeGNMISetClient) Reset() {
	f.requests = nil
	f.responses = nil
	f.responseErrs = nil
	f.i = 0
}

func (f *fakeGNMISetClient) AddResponse(resp *gpb.SetResponse, err error) *fakeGNMISetClient {
	f.responses = append(f.responses, resp)
	f.responseErrs = append(f.responseErrs, err)
	return f
}

func (f *fakeGNMISetClient) Set(_ context.Context, req *gpb.SetRequest, opts ...grpc.CallOption) (*gpb.SetResponse, error) {
	defer func() { f.i++ }()
	f.requests = append(f.requests, req)
	return f.responses[f.i], f.responseErrs[f.i]
}

// checkJustReceived checks that the received time is just before now.
func checkJustReceived(t *testing.T, recvTime time.Time) {
	if diffSecs := time.Since(recvTime).Seconds(); diffSecs <= 0 && diffSecs > 1 {
		t.Errorf("Received time is too far (%v seconds) away from now", diffSecs)
	}
}

// verifySubscriptionPathsSent verifies the paths of the sent subscription requests is the same as wantPaths.
func verifySubscriptionPathsSent(t *testing.T, fakeGNMI *testutil.FakeGNMI, wantPaths ...*gpb.Path) {
	t.Helper()
	requests := fakeGNMI.Requests()
	if len(requests) != 1 {
		t.Errorf("Number of subscription requests sent is not 1: %v", requests)
		return
	}

	var gotPaths []*gpb.Path
	req := requests[0].GetSubscribe()
	for _, sub := range req.GetSubscription() {
		got, err := util.JoinPaths(req.GetPrefix(), sub.GetPath())
		if err != nil {
			t.Fatal(err)
		}
		got.Target = ""
		gotPaths = append(gotPaths, got)
	}
	if diff := cmp.Diff(wantPaths, gotPaths, protocmp.Transform(), cmpopts.SortSlices(ygottestutil.PathLess)); diff != "" {
		t.Errorf("Subscription paths (-want, +got):\n%s", diff)
	}
}

// verifySubscriptionModesSent verifies the modes of the sent subscription requests is the same as wantModes.
func verifySubscriptionModesSent(t *testing.T, fakeGNMI *testutil.FakeGNMI, wantModes ...gpb.SubscriptionMode) {
	t.Helper()
	requests := fakeGNMI.Requests()
	if len(requests) != 1 {
		t.Errorf("Number of subscription requests sent is not 1: %v", requests)
		return
	}

	var gotModes []gpb.SubscriptionMode
	req := requests[0].GetSubscribe()
	for _, sub := range req.GetSubscription() {
		gotModes = append(gotModes, sub.Mode)
	}
	if diff := cmp.Diff(wantModes, gotModes, protocmp.Transform()); diff != "" {
		t.Errorf("Subscription modes (-want, +got):\n%s", diff)
	}
}

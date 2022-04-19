package ygnmi

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/testing/protocmp"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	ygottestutil "github.com/openconfig/ygot/testutil"
)

func getClient(t testing.TB) (*testutil.FakeGNMI, *Client) {
	fakeGNMI, err := testutil.Start(0)
	if err != nil {
		t.Fatal(err)
	}
	gnmiClient, err := fakeGNMI.Dial(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewClient(gnmiClient)
	if err != nil {
		t.Fatal(err)
	}
	return fakeGNMI, c
}

func TestLookup(t *testing.T) {
	fakeGNMI, c := getClient(t)
	leafPath := testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf")
	lq := &LeafSingletonQuery[uint64]{
		parentDir:  "leaf-container-struct",
		state:      false,
		ps:         ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "uint64-leaf"}, nil, ygot.NewDeviceRootBase("")),
		extractFn:  func(vgs ygot.ValidatedGoStruct) uint64 { return *(vgs.(*testutil.LeafContainerStruct)).Uint64Leaf },
		goStructFn: func() ygot.ValidatedGoStruct { return new(testutil.LeafContainerStruct) },
		yschema:    testutil.GetSchemaStruct()(),
	}

	leaftests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		inQuery              SingletonQuery[uint64]
		wantSubscriptionPath *gpb.Path
		wantVal              *Value[uint64]
		wantErr              string
	}{{
		desc:    "success update and sync",
		inQuery: lq,
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: leafPath,
		wantVal: &Value[uint64]{
			val:       10,
			present:   true,
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		},
	}, {
		desc:    "success update and no sync",
		inQuery: lq,
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			})
		},
		wantSubscriptionPath: leafPath,
		wantVal: &Value[uint64]{
			val:       10,
			present:   true,
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		},
	}, {
		desc:    "success with prefix",
		inQuery: lq,
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Prefix:    testutil.GNMIPath(t, "super-container"),
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "leaf-container-struct/uint64-leaf"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: leafPath,
		wantVal: &Value[uint64]{
			val:       10,
			present:   true,
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		},
	}, {
		desc:    "success multiple notifs and first no value",
		inQuery: lq,
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Update: []*gpb.Update{},
			}).Notification(&gpb.Notification{
				Timestamp: 100,
				Prefix:    testutil.GNMIPath(t, "super-container"),
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "leaf-container-struct/uint64-leaf"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: leafPath,
		wantVal: &Value[uint64]{
			val:       10,
			present:   true,
			Path:      leafPath,
			Timestamp: time.Unix(0, 100),
		},
	}, {
		desc:    "success no value",
		inQuery: lq,
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		wantSubscriptionPath: leafPath,
		wantVal: &Value[uint64]{
			present: false,
			Path:    leafPath,
		},
	}, {
		desc:    "error multiple values",
		inQuery: lq,
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
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
					Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf"),
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
					Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/does-not-exist"),
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
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		wantErr: "failed to unmarshal",
	}}
	for _, tt := range leaftests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			got, err := Lookup(context.Background(), c, tt.inQuery)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Lookup(ctx, c, %v) returned unexpected diff: %s", tt.inQuery, diff)
			}
			if err != nil {
				return
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath)
			checkJustReceived(t, got.RecvTimestamp)
			tt.wantVal.RecvTimestamp = got.RecvTimestamp

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(Value[uint64]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Lookup(ctx, c, %v) returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", tt.inQuery, diff, got.ComplianceErrors)
			}
		})
	}
}

func TestLookupNonLeaf(t *testing.T) {
	fakeGNMI, c := getClient(t)
	rootPath := testutil.GNMIPath(t, "super-container/leaf-container-struct")
	intPath := testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf")
	enumPath := testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf")
	intStatePath := testutil.GNMIPath(t, "super-container/leaf-container-struct/state/uint64-leaf")

	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		inState              bool
		wantSubscriptionPath *gpb.Path
		wantVal              *Value[*testutil.LeafContainerStruct]
		wantErr              string
	}{{
		desc: "success one update and state false",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: intPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 10}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: rootPath,
		wantVal: &Value[*testutil.LeafContainerStruct]{
			val: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(10),
			},
			present:   true,
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		},
	}, {
		desc: "success one update and state true",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: intStatePath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 12}},
				}},
			}).Sync()
		},
		inState:              true,
		wantSubscriptionPath: rootPath,
		wantVal: &Value[*testutil.LeafContainerStruct]{
			val: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(12),
			},
			present:   true,
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		},
	}, {
		desc: "success one update with prefix",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Prefix:    testutil.GNMIPath(t, "super-container"),
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "leaf-container-struct/uint64-leaf"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 12}},
				}},
			}).Sync()
		},
		inState:              false,
		wantSubscriptionPath: rootPath,
		wantVal: &Value[*testutil.LeafContainerStruct]{
			val: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(12),
			},
			present:   true,
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		},
	}, {
		desc: "success ignore state update when state false",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: intStatePath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 10}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: rootPath,
		wantVal: &Value[*testutil.LeafContainerStruct]{
			// TODO(DanG100): fix the check to correctly mark this as not present.
			present:   true,
			val:       &testutil.LeafContainerStruct{},
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		},
	}, {
		desc: "success ignore non-state update when state true",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: intPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 10}},
				}},
			}).Sync()
		},
		inState:              true,
		wantSubscriptionPath: rootPath,
		wantVal: &Value[*testutil.LeafContainerStruct]{
			// TODO(DanG100): fix the check to correctly mark this as not present.
			present:   true,
			val:       &testutil.LeafContainerStruct{},
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		},
	}, {
		desc: "success multiple updates in single notification",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "E_VALUE_FORTY_THREE"}},
				}, {
					Path: intPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 10}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: rootPath,
		wantVal: &Value[*testutil.LeafContainerStruct]{
			val: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(10),
				EnumLeaf:   43,
			},
			present:   true,
			Path:      rootPath,
			Timestamp: time.Unix(0, 100),
		},
	}, {
		desc: "success multiple notifications",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "E_VALUE_FORTY_THREE"}},
				}},
			}).Notification(&gpb.Notification{
				Timestamp: 102,
				Update: []*gpb.Update{{
					Path: intPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_UintVal{UintVal: 10}},
				}},
			}).Sync()
		},
		wantSubscriptionPath: rootPath,
		wantVal: &Value[*testutil.LeafContainerStruct]{
			val: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(10),
				EnumLeaf:   43,
			},
			present:   true,
			Path:      rootPath,
			Timestamp: time.Unix(0, 102),
		},
	}, {
		desc: "success no values",
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		wantSubscriptionPath: rootPath,
		wantVal: &Value[*testutil.LeafContainerStruct]{
			Path: rootPath,
		},
	}}
	// Failing tests cases are covered by the TestLookup() as these tests test the same func.

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			q := &NonLeafSingletonQuery[*testutil.LeafContainerStruct]{
				dir:     "leaf-container-struct",
				ps:      ygot.NewNodePath([]string{"super-container", "leaf-container-struct"}, nil, ygot.NewDeviceRootBase("")),
				yschema: testutil.GetSchemaStruct()(),
				state:   tt.inState,
			}

			got, err := Lookup[*testutil.LeafContainerStruct](context.Background(), c, q)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("Lookup(ctx, c, %v) returned unexpected diff: %s", q, diff)
			}
			if err != nil {
				return
			}
			verifySubscriptionPathsSent(t, fakeGNMI, tt.wantSubscriptionPath)
			checkJustReceived(t, got.RecvTimestamp)
			tt.wantVal.RecvTimestamp = got.RecvTimestamp

			if diff := cmp.Diff(tt.wantVal, got, cmp.AllowUnexported(Value[*testutil.LeafContainerStruct]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Lookup(ctx, c, %v) returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", q, diff, got.ComplianceErrors)
			}
		})
	}
func TestWatch(t *testing.T) {
	fakeGNMI, c := getClient(t)
	path := testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf")
	q := &LeafSingletonQuery[uint64]{
		parentDir: "leaf-container-struct",
		state:     false,
		ps:        ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "uint64-leaf"}, nil, ygot.NewDeviceRootBase("")),
		extractFn: func(vgs ygot.ValidatedGoStruct) uint64 {
			lcs := vgs.(*testutil.LeafContainerStruct)
			if lcs.Uint64Leaf == nil {
				return 0
			}
			return *lcs.Uint64Leaf
		},
		goStructFn: func() ygot.ValidatedGoStruct { return new(testutil.LeafContainerStruct) },
		yschema:    testutil.GetSchemaStruct()(),
	}

	startTime := time.Now()
	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantLastVal          *Value[uint64]
		wantVals             []*Value[uint64]
		wantStatus           bool
		wantErr              string
	}{{
		desc: "single notif and pred true",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 12}},
				}},
			}).Sync()
		},
		wantVals: []*Value[uint64]{{
			present:   true,
			val:       12,
			Timestamp: startTime,
			Path:      path,
		}},
		wantSubscriptionPath: path,
		wantStatus:           true,
		wantLastVal: &Value[uint64]{
			present:   true,
			val:       12,
			Timestamp: startTime,
			Path:      path,
		},
	}, {
		desc: "single notif and pred false",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 9}},
				}},
			}).Sync()
		},
		wantVals: []*Value[uint64]{{
			present:   true,
			val:       9,
			Timestamp: startTime,
			Path:      path,
		}},
		wantSubscriptionPath: path,
		wantStatus:           false,
		wantLastVal: &Value[uint64]{
			present:   true,
			val:       9,
			Timestamp: startTime,
			Path:      path,
		},
	}, {
		desc: "multiple notif and pred true",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 8}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			})
		},
		wantVals: []*Value[uint64]{{
			present:   true,
			val:       8,
			Timestamp: startTime,
			Path:      path,
		}, {
			present:   true,
			val:       11,
			Timestamp: startTime.Add(time.Millisecond),
			Path:      path,
		}},
		wantSubscriptionPath: path,
		wantStatus:           true,
		wantLastVal: &Value[uint64]{
			present:   true,
			val:       11,
			Timestamp: startTime.Add(time.Millisecond),
			Path:      path,
		},
	}, {
		desc: "multiple notif with deletes",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 8}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Delete:    []*gpb.Path{path},
			})
		},
		wantVals: []*Value[uint64]{{
			present:   true,
			val:       8,
			Timestamp: startTime,
			Path:      path,
		}, {
			present:   false,
			Timestamp: startTime.Add(time.Millisecond),
			Path:      path,
		}},
		wantSubscriptionPath: path,
		wantStatus:           false,
		wantLastVal: &Value[uint64]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      path,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			w, err := Watch[uint64](context.Background(), c, q, 1*time.Second, func(v *Value[uint64]) bool {
				if len(tt.wantVals) == 0 {
					t.Fatalf("Predicate expected no more values but got: %+v", v)
				}
				if diff := cmp.Diff(tt.wantVals[0], v, cmpopts.IgnoreFields(Value[uint64]{}, "RecvTimestamp"), cmp.AllowUnexported(Value[uint64]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", diff, v.ComplianceErrors)
				}
				tt.wantVals = tt.wantVals[1:]
				val, present := v.Val()
				return present && val > 10
			})
			if err != nil {
				t.Fatalf("Watch() returned unexpected error: %v", err)
			}
			val, status, err := w.Await()
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

			if tt.wantStatus != status {
				t.Errorf("Await() returned unexpected status got: %v, want %v", status, tt.wantStatus)
			}
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(Value[uint64]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

// checks that the received time is just before now
func checkJustReceived(t *testing.T, recvTime time.Time) {
	if diffSecs := time.Now().Sub(recvTime).Seconds(); diffSecs <= 0 && diffSecs > 1 {
		t.Errorf("received time is too far (%v seconds) away from now", diffSecs)
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
		t.Errorf("subscription paths (-want, +got):\n%s", diff)
	}
}

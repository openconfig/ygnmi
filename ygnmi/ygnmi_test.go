package ygnmi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
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
		leafBaseQuery: leafBaseQuery[uint64]{
			parentDir:  "leaf-container-struct",
			state:      false,
			ps:         ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "uint64-leaf"}, nil, ygot.NewDeviceRootBase("")),
			extractFn:  func(vgs ygot.ValidatedGoStruct) uint64 { return *(vgs.(*testutil.LeafContainerStruct)).Uint64Leaf },
			goStructFn: func() ygot.ValidatedGoStruct { return new(testutil.LeafContainerStruct) },
			yschema:    testutil.GetSchemaStruct()(),
		},
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
				nonLeafBaseQuery[*testutil.LeafContainerStruct]{
					dir:     "leaf-container-struct",
					ps:      ygot.NewNodePath([]string{"super-container", "leaf-container-struct"}, nil, ygot.NewDeviceRootBase("")),
					yschema: testutil.GetSchemaStruct()(),
					state:   tt.inState,
				},
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
}

func TestWatch(t *testing.T) {
	fakeGNMI, client := getClient(t)
	path := testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf")
	q := &LeafSingletonQuery[uint64]{
		leafBaseQuery: leafBaseQuery[uint64]{
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
		},
	}

	startTime := time.Now()
	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *Value[uint64]
		wantVals             []*Value[uint64]
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
		dur: time.Second,
		wantVals: []*Value[uint64]{{
			present:   true,
			val:       12,
			Timestamp: startTime,
			Path:      path,
		}},
		wantSubscriptionPath: path,
		wantLastVal: &Value[uint64]{
			present:   true,
			val:       12,
			Timestamp: startTime,
			Path:      path,
		},
	}, {
		desc: "single notif and pred false error EOF",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: path,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 9}},
				}},
			}).Sync()
		},
		dur: time.Second,
		wantVals: []*Value[uint64]{{
			present:   true,
			val:       9,
			Timestamp: startTime,
			Path:      path,
		}},
		wantSubscriptionPath: path,
		wantLastVal: &Value[uint64]{
			present:   true,
			val:       9,
			Timestamp: startTime,
			Path:      path,
		},
		wantErr: "error receiving gNMI response: EOF",
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
		dur: time.Second,
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
		dur: time.Second,
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
		wantLastVal: &Value[uint64]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      path,
		},
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
			w := Watch[uint64](ctx, client, q, func(v *Value[uint64]) bool {
				if i > len(tt.wantVals) {
					t.Fatalf("Predicate(%d) expected no more values but got: %+v", i, v)
				}
				if diff := cmp.Diff(tt.wantVals[i], v, cmpopts.IgnoreFields(Value[uint64]{}, "RecvTimestamp"), cmp.AllowUnexported(Value[uint64]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate(%d) got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", i, diff, v.ComplianceErrors)
				}
				val, present := v.Val()
				i++
				return present && val > 10
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
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(Value[uint64]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}

	t.Run("multiple awaits", func(t *testing.T) {
		fakeGNMI.Stub().Sync()
		w := Watch[uint64](context.Background(), client, q, func(v *Value[uint64]) bool { return true })
		want := &Value[uint64]{
			Path: path,
		}
		val, err := w.Await()
		if err != nil {
			t.Fatalf("Await() got unexpected error: %v", err)
		}
		if diff := cmp.Diff(want, val, cmp.AllowUnexported(Value[uint64]{}), protocmp.Transform()); diff != "" {
			t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
		}
		_, err = w.Await()
		if d := errdiff.Check(err, "Await already called and Watcher is closed"); d != "" {
			t.Fatalf("Await() returned unexpected diff: %s", d)
		}
	})

	rootPath := testutil.GNMIPath(t, "super-container/leaf-container-struct")
	intPath := testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf")
	enumPath := testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf")
	startTime = time.Now()
	nonLeafQuery := &NonLeafSingletonQuery[*testutil.LeafContainerStruct]{
		nonLeafBaseQuery: nonLeafBaseQuery[*testutil.LeafContainerStruct]{
			dir:     "leaf-container-struct",
			ps:      ygot.NewNodePath([]string{"super-container", "leaf-container-struct"}, nil, ygot.NewDeviceRootBase("")),
			yschema: testutil.GetSchemaStruct()(),
		},
	}

	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantLastVal          *Value[*testutil.LeafContainerStruct]
		wantVals             []*Value[*testutil.LeafContainerStruct]
		wantErr              string
	}{{
		desc: "single notif and pred false",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: intPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 9}},
				}},
			}).Sync()
		},
		wantVals: []*Value[*testutil.LeafContainerStruct]{{
			present:   true,
			val:       &testutil.LeafContainerStruct{Uint64Leaf: ygot.Uint64(9)},
			Timestamp: startTime,
			Path:      rootPath,
		}},
		wantErr:              "EOF",
		wantSubscriptionPath: rootPath,
		wantLastVal: &Value[*testutil.LeafContainerStruct]{
			present:   true,
			val:       &testutil.LeafContainerStruct{Uint64Leaf: ygot.Uint64(9)},
			Timestamp: startTime,
			Path:      rootPath,
		},
	}, {
		desc: "multiple notif and pred true",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: intPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "E_VALUE_FORTY_THREE"}},
				}},
			})
		},
		wantVals: []*Value[*testutil.LeafContainerStruct]{{
			present:   true,
			val:       &testutil.LeafContainerStruct{Uint64Leaf: ygot.Uint64(11)},
			Timestamp: startTime,
			Path:      rootPath,
		}, {
			present: true,
			val: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(11),
				EnumLeaf:   43,
			},
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}},
		wantSubscriptionPath: rootPath,
		wantLastVal: &Value[*testutil.LeafContainerStruct]{
			present: true,
			val: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(11),
				EnumLeaf:   43,
			},
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		},
	}, {
		desc: "multiple notif before sync",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: intPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "E_VALUE_FORTY_THREE"}},
				}},
			}).Sync()
		},
		wantVals: []*Value[*testutil.LeafContainerStruct]{{
			present: true,
			val: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(11),
				EnumLeaf:   43,
			},
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}},
		wantSubscriptionPath: rootPath,
		wantLastVal: &Value[*testutil.LeafContainerStruct]{
			present: true,
			val: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(11),
				EnumLeaf:   43,
			},
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		},
	}, {
		desc: "delete leaf in container",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: intPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 8}},
				}, {
					Path: enumPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "E_VALUE_FORTY_THREE"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Delete:    []*gpb.Path{intPath},
			})
		},
		wantVals: []*Value[*testutil.LeafContainerStruct]{{
			present: true,
			val: &testutil.LeafContainerStruct{
				Uint64Leaf: ygot.Uint64(8),
				EnumLeaf:   43,
			},
			Timestamp: startTime,
			Path:      rootPath,
		}, {
			present: true,
			val: &testutil.LeafContainerStruct{
				EnumLeaf: 43,
			},
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		}},
		wantSubscriptionPath: rootPath,
		wantErr:              "EOF",
		wantLastVal: &Value[*testutil.LeafContainerStruct]{
			present: true,
			val: &testutil.LeafContainerStruct{
				EnumLeaf: 43,
			},
			Timestamp: startTime.Add(time.Millisecond),
			Path:      rootPath,
		},
	}}

	for _, tt := range nonLeafTests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			w := Watch[*testutil.LeafContainerStruct](context.Background(), client, nonLeafQuery, func(v *Value[*testutil.LeafContainerStruct]) bool {
				if len(tt.wantVals) == 0 {
					t.Fatalf("Predicate expected no more values but got: %+v", v)
				}
				if diff := cmp.Diff(tt.wantVals[0], v, cmpopts.IgnoreFields(Value[*testutil.LeafContainerStruct]{}, "RecvTimestamp"), cmp.AllowUnexported(Value[*testutil.LeafContainerStruct]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", diff, v.ComplianceErrors)
				}
				tt.wantVals = tt.wantVals[1:]
				val, present := v.Val()
				return present && val.Uint64Leaf != nil && *val.Uint64Leaf > 10 && val.EnumLeaf == 43
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
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(Value[*testutil.LeafContainerStruct]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestLookupAll(t *testing.T) {
	fakeGNMI, c := getClient(t)
	leafPath := testutil.GNMIPath(t, "super-container/model/a/single-key[key=*]/state/value")
	leafPS := ygot.NewNodePath(
		[]string{"state", "value"},
		nil,
		ygot.NewNodePath([]string{"super-container", "model", "a", "single-key"}, map[string]interface{}{"key": "*"}, ygot.NewDeviceRootBase("")),
	)
	lq := &LeafWildcardQuery[int64]{
		leafBaseQuery: leafBaseQuery[int64]{
			parentDir:  "Model_SingleKey",
			state:      true,
			ps:         leafPS,
			extractFn:  func(vgs ygot.ValidatedGoStruct) int64 { return *((vgs.(*testutil.Model_SingleKey)).Value) },
			goStructFn: func() ygot.ValidatedGoStruct { return new(testutil.Model_SingleKey) },
			yschema:    testutil.GetSchemaStruct()(),
		},
	}
	leaftests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []*Value[int64]
		wantErr              string
	}{{
		desc: "success one value",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantVals: []*Value[int64]{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
			val:       10,
			present:   true,
			Timestamp: time.Unix(0, 100),
		}},
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
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}, {
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Sync()
		},
		wantVals: []*Value[int64]{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
			val:       10,
			present:   true,
			Timestamp: time.Unix(0, 100),
		}, {
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/value"),
			val:       11,
			present:   true,
			Timestamp: time.Unix(0, 100),
		}},
		wantSubscriptionPath: leafPath,
	}, {
		desc: "success multiples value in different notifications",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Sync()
		},
		wantVals: []*Value[int64]{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
			val:       10,
			present:   true,
			Timestamp: time.Unix(0, 100),
		}, {
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/value"),
			val:       11,
			present:   true,
			Timestamp: time.Unix(0, 101),
		}},
		wantSubscriptionPath: leafPath,
	}, {
		desc: "success ignore mismatched paths",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/config/value"),
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
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/value"),
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
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/value"),
					Val:  nil,
				}},
			}).Sync()
		},
		wantErr:              "failed to receive to data",
		wantSubscriptionPath: leafPath,
	}}
	for _, tt := range leaftests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			got, err := LookupAll[int64](context.Background(), c, lq)
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
			if diff := cmp.Diff(tt.wantVals, got, cmp.AllowUnexported(Value[int64]{}), cmpopts.IgnoreFields(Value[int64]{}, "RecvTimestamp"), protocmp.Transform()); diff != "" {
				t.Errorf("LookupAll() returned unexpected diff (-want,+got):\n%s", diff)
			}
		})
	}

	nonLeafPath := testutil.GNMIPath(t, "super-container/model/a/single-key[key=*]")
	nonLeafPS := ygot.NewNodePath([]string{"super-container", "model", "a", "single-key"}, map[string]interface{}{"key": "*"}, ygot.NewDeviceRootBase(""))
	nonLeafQ := &NonLeafWildcardQuery[*testutil.Model_SingleKey]{
		nonLeafBaseQuery: nonLeafBaseQuery[*testutil.Model_SingleKey]{
			dir:     "Model_SingleKey",
			state:   true,
			ps:      nonLeafPS,
			yschema: testutil.GetSchemaStruct()(),
		},
	}
	nonLeaftests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []*Value[*testutil.Model_SingleKey]
		wantErr              string
	}{{
		desc: "non-leaf one value",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantVals: []*Value[*testutil.Model_SingleKey]{{
			Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]"),
			val: &testutil.Model_SingleKey{
				Value: ygot.Int64(10),
			},
			present:   true,
			Timestamp: time.Unix(0, 100),
		}},
		wantSubscriptionPath: nonLeafPath,
	}, {
		desc: "non-leaf multiple values",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 100}},
				}, {
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 101}},
				}, {
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/key"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/key"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 11}},
				}},
			}).Sync()
		},
		wantVals: []*Value[*testutil.Model_SingleKey]{{
			Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]"),
			val: &testutil.Model_SingleKey{
				Value: ygot.Int64(100),
				Key:   ygot.Int32(10),
			},
			present:   true,
			Timestamp: time.Unix(0, 100),
		}, {
			Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]"),
			val: &testutil.Model_SingleKey{
				Value: ygot.Int64(101),
				Key:   ygot.Int32(11),
			},
			present:   true,
			Timestamp: time.Unix(0, 101),
		}},
		wantSubscriptionPath: nonLeafPath,
	}, {
		desc: "no values",
		stub: func(s *testutil.Stubber) {
			s.Sync()
		},
		wantVals:             nil,
		wantSubscriptionPath: nonLeafPath,
	}}
	for _, tt := range nonLeaftests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			got, err := LookupAll[*testutil.Model_SingleKey](context.Background(), c, nonLeafQ)
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
			if diff := cmp.Diff(tt.wantVals, got, cmp.AllowUnexported(Value[*testutil.Model_SingleKey]{}), cmpopts.IgnoreFields(Value[*testutil.Model_SingleKey]{}, "RecvTimestamp"), protocmp.Transform()); diff != "" {
				t.Errorf("LookupAll() returned unexpected diff (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestWatchAll(t *testing.T) {
	fakeGNMI, client := getClient(t)
	leafQueryPath := testutil.GNMIPath(t, "super-container/model/a/single-key[key=*]/state/value")
	key10Path := testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value")
	key11Path := testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/value")

	startTime := time.Now()
	leafPS := ygot.NewNodePath(
		[]string{"state", "value"},
		nil,
		ygot.NewNodePath([]string{"super-container", "model", "a", "single-key"}, map[string]interface{}{"key": "*"}, ygot.NewDeviceRootBase("")),
	)
	lq := &LeafWildcardQuery[int64]{
		leafBaseQuery: leafBaseQuery[int64]{
			parentDir:  "Model_SingleKey",
			state:      true,
			ps:         leafPS,
			extractFn:  func(vgs ygot.ValidatedGoStruct) int64 { return *((vgs.(*testutil.Model_SingleKey)).Value) },
			goStructFn: func() ygot.ValidatedGoStruct { return new(testutil.Model_SingleKey) },
			yschema:    testutil.GetSchemaStruct()(),
		},
	}
	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *Value[int64]
		wantVals             []*Value[int64]
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
		wantSubscriptionPath: leafQueryPath,
		wantVals: []*Value[int64]{{
			Timestamp: startTime,
			Path:      key10Path,
			val:       100,
			present:   true,
		}},
		wantLastVal: &Value[int64]{
			Timestamp: startTime,
			Path:      key10Path,
			val:       100,
			present:   true,
		},
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
		wantVals: []*Value[int64]{{
			Timestamp: startTime,
			Path:      key10Path,
			val:       100,
			present:   true,
		}, {
			Timestamp: startTime.Add(time.Millisecond),
			Path:      key11Path,
			val:       101,
			present:   true,
		}},
		wantLastVal: &Value[int64]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      key11Path,
			val:       101,
			present:   true,
		},
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
		wantVals: []*Value[int64]{{
			Timestamp: startTime,
			Path:      key10Path,
			val:       100,
			present:   true,
		}, {
			Timestamp: startTime,
			Path:      key11Path,
			val:       101,
			present:   true,
		}},
		wantLastVal: &Value[int64]{
			Timestamp: startTime,
			Path:      key11Path,
			val:       101,
			present:   true,
		},
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

			w := WatchAll[int64](ctx, client, lq, func(v *Value[int64]) bool {
				if i > len(tt.wantVals) {
					t.Fatalf("Predicate(%d) expected no more values but got: %+v", i, v)
				}
				if diff := cmp.Diff(tt.wantVals[i], v, cmpopts.IgnoreFields(Value[int64]{}, "RecvTimestamp"), cmp.AllowUnexported(Value[int64]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate(%d) got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", i, diff, v.ComplianceErrors)
				}
				val, present := v.Val()
				key10Cond = key10Cond || (present && proto.Equal(v.Path, key10Path) && val == 100)
				key11Cond = key11Cond || (present && proto.Equal(v.Path, key11Path) && val == 101)
				i++
				return key10Cond && key11Cond
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
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(Value[int64]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}

	nonLeafPath := testutil.GNMIPath(t, "super-container/model/a/single-key[key=*]")
	nonLeafKey10Path := testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]")
	nonLeafKey11Path := testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]")

	nonLeafPS := ygot.NewNodePath([]string{"super-container", "model", "a", "single-key"}, map[string]interface{}{"key": "*"}, ygot.NewDeviceRootBase(""))
	nonLeafQ := &NonLeafWildcardQuery[*testutil.Model_SingleKey]{
		nonLeafBaseQuery: nonLeafBaseQuery[*testutil.Model_SingleKey]{
			dir:     "Model_SingleKey",
			state:   true,
			ps:      nonLeafPS,
			yschema: testutil.GetSchemaStruct()(),
		},
	}
	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *Value[*testutil.Model_SingleKey]
		wantVals             []*Value[*testutil.Model_SingleKey]
		wantErr              string
	}{{
		desc: "non-leaf predicate not true",
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
		wantVals: []*Value[*testutil.Model_SingleKey]{{
			Timestamp: startTime,
			Path:      nonLeafKey10Path,
			val:       &testutil.Model_SingleKey{Value: ygot.Int64(100)},
			present:   true,
		}},
		wantLastVal: &Value[*testutil.Model_SingleKey]{
			Timestamp: startTime,
			Path:      nonLeafKey10Path,
			val:       &testutil.Model_SingleKey{Value: ygot.Int64(100)},
			present:   true,
		},
		wantErr: "EOF",
	}, {
		desc: "non-leaf predicate becomes true",
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
		wantVals: []*Value[*testutil.Model_SingleKey]{{
			Timestamp: startTime,
			Path:      nonLeafKey10Path,
			val:       &testutil.Model_SingleKey{Value: ygot.Int64(100)},
			present:   true,
		}, {
			Timestamp: startTime.Add(time.Millisecond),
			Path:      nonLeafKey11Path,
			val:       &testutil.Model_SingleKey{Value: ygot.Int64(101)},
			present:   true,
		}},
		wantLastVal: &Value[*testutil.Model_SingleKey]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      nonLeafKey11Path,
			val:       &testutil.Model_SingleKey{Value: ygot.Int64(101)},
			present:   true,
		},
	}}
	for _, tt := range nonLeafTests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			i := 0
			ctx, cancel := context.WithTimeout(context.Background(), tt.dur)
			defer cancel()
			var key10Cond, key11Cond bool

			w := WatchAll[*testutil.Model_SingleKey](ctx, client, nonLeafQ, func(v *Value[*testutil.Model_SingleKey]) bool {
				if i > len(tt.wantVals) {
					t.Fatalf("Predicate(%d) expected no more values but got: %+v", i, v)
				}
				if diff := cmp.Diff(tt.wantVals[i], v, cmpopts.IgnoreFields(Value[*testutil.Model_SingleKey]{}, "RecvTimestamp"), cmp.AllowUnexported(Value[*testutil.Model_SingleKey]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate(%d) got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", i, diff, v.ComplianceErrors)
				}
				val, present := v.Val()
				key10Cond = key10Cond || (present && proto.Equal(v.Path, nonLeafKey10Path) && *val.Value == 100)
				key11Cond = key11Cond || (present && proto.Equal(v.Path, nonLeafKey11Path) && *val.Value == 101)
				i++
				return key10Cond && key11Cond
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
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(Value[*testutil.Model_SingleKey]{}), protocmp.Transform()); diff != "" {
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

type fakeGNMISetClient struct {
	gpb.GNMIClient
	// Responses are the gNMI responses to return from calls to Set.
	Responses []*gpb.SetResponse
	// Requests received by the client are stored in the slice.
	Requests []*gpb.SetRequest
	// ResponseErrs are the errors to return from calls to Set.
	ResponseErrs []error
	// i is index current index of the response and error to return.
	i int
}

func (f *fakeGNMISetClient) Reset() {
	f.Requests = nil
	f.Responses = nil
	f.ResponseErrs = nil
	f.i = 0
}

func (f *fakeGNMISetClient) AddResponse(resp *gpb.SetResponse, err error) *fakeGNMISetClient {
	f.Responses = append(f.Responses, resp)
	f.ResponseErrs = append(f.ResponseErrs, err)
	return f
}

func (f *fakeGNMISetClient) Set(_ context.Context, req *gpb.SetRequest, opts ...grpc.CallOption) (*gpb.SetResponse, error) {
	defer func() { f.i++ }()
	f.Requests = append(f.Requests, req)
	return f.Responses[f.i], f.ResponseErrs[f.i]
}

func TestUpdate(t *testing.T) {
	setClient := &fakeGNMISetClient{}
	client := &Client{
		gnmiC:  setClient,
		target: "dut",
	}
	tests := []struct {
		desc         string
		op           func(*Client) (*gpb.SetResponse, error)
		wantErr      string
		wantRequest  *gpb.SetRequest
		stubResponse *gpb.SetResponse
		stubErr      error
	}{{
		desc: "scalar leaf",
		op: func(c *Client) (*gpb.SetResponse, error) {
			q := &LeafConfigQuery[uint64]{
				leafBaseQuery: leafBaseQuery[uint64]{
					state:  false,
					scalar: true,
					ps:     ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "uint64-leaf"}, nil, ygot.NewDeviceRootBase("")),
				},
			}
			return Update[uint64](context.Background(), c, q, 10)
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
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
		op: func(c *Client) (*gpb.SetResponse, error) {
			q := &LeafConfigQuery[testutil.EnumType]{
				leafBaseQuery: leafBaseQuery[testutil.EnumType]{
					state:  false,
					scalar: false,
					ps:     ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "enum-leaf"}, nil, ygot.NewDeviceRootBase("")),
				},
			}
			return Update[testutil.EnumType](context.Background(), c, q, testutil.EnumType(43))
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"E_VALUE_FORTY_THREE\"")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "non leaf",
		op: func(c *Client) (*gpb.SetResponse, error) {
			q := &NonLeafConfigQuery[*testutil.LeafContainerStruct]{
				nonLeafBaseQuery: nonLeafBaseQuery[*testutil.LeafContainerStruct]{
					state: false,
					ps:    ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "enum-leaf"}, nil, ygot.NewDeviceRootBase("")),
				},
			}
			return Update[*testutil.LeafContainerStruct](context.Background(), c, q, &testutil.LeafContainerStruct{Uint64Leaf: ygot.Uint64(10)})
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("{\n  \"state\": {\n    \"uint64-leaf\": \"10\"\n  }\n}")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "server error",
		op: func(c *Client) (*gpb.SetResponse, error) {
			q := &LeafConfigQuery[uint64]{
				leafBaseQuery: leafBaseQuery[uint64]{
					state:  false,
					scalar: true,
					ps:     ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "uint64-leaf"}, nil, ygot.NewDeviceRootBase("")),
				},
			}
			return Update[uint64](context.Background(), c, q, 10)
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Update: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
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
			if diff := cmp.Diff(tt.wantRequest, setClient.Requests[0], protocmp.Transform()); diff != "" {
				t.Errorf("Update() sent unexpected request (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.stubResponse, got, protocmp.Transform()); diff != "" {
				t.Errorf("Update() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	setClient := &fakeGNMISetClient{}
	client := &Client{
		gnmiC:  setClient,
		target: "dut",
	}
	tests := []struct {
		desc         string
		op           func(*Client) (*gpb.SetResponse, error)
		wantErr      string
		wantRequest  *gpb.SetRequest
		stubResponse *gpb.SetResponse
		stubErr      error
	}{{
		desc: "scalar leaf",
		op: func(c *Client) (*gpb.SetResponse, error) {
			q := &LeafConfigQuery[uint64]{
				leafBaseQuery: leafBaseQuery[uint64]{
					state:  false,
					scalar: true,
					ps:     ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "uint64-leaf"}, nil, ygot.NewDeviceRootBase("")),
				},
			}
			return Replace[uint64](context.Background(), c, q, 10)
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
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
		op: func(c *Client) (*gpb.SetResponse, error) {
			q := &LeafConfigQuery[testutil.EnumType]{
				leafBaseQuery: leafBaseQuery[testutil.EnumType]{
					state:  false,
					scalar: false,
					ps:     ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "enum-leaf"}, nil, ygot.NewDeviceRootBase("")),
				},
			}
			return Replace[testutil.EnumType](context.Background(), c, q, testutil.EnumType(43))
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("\"E_VALUE_FORTY_THREE\"")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "non leaf",
		op: func(c *Client) (*gpb.SetResponse, error) {
			q := &NonLeafConfigQuery[*testutil.LeafContainerStruct]{
				nonLeafBaseQuery: nonLeafBaseQuery[*testutil.LeafContainerStruct]{
					state: false,
					ps:    ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "enum-leaf"}, nil, ygot.NewDeviceRootBase("")),
				},
			}
			return Replace[*testutil.LeafContainerStruct](context.Background(), c, q, &testutil.LeafContainerStruct{Uint64Leaf: ygot.Uint64(10)})
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf"),
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("{\n  \"state\": {\n    \"uint64-leaf\": \"10\"\n  }\n}")}},
			}},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "server error",
		op: func(c *Client) (*gpb.SetResponse, error) {
			q := &LeafConfigQuery[uint64]{
				leafBaseQuery: leafBaseQuery[uint64]{
					state:  false,
					scalar: true,
					ps:     ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "uint64-leaf"}, nil, ygot.NewDeviceRootBase("")),
				},
			}
			return Replace[uint64](context.Background(), c, q, 10)
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Replace: []*gpb.Update{{
				Path: testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
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
			if diff := cmp.Diff(tt.wantRequest, setClient.Requests[0], protocmp.Transform()); diff != "" {
				t.Errorf("Replace() sent unexpected request (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.stubResponse, got, protocmp.Transform()); diff != "" {
				t.Errorf("Replace() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	setClient := &fakeGNMISetClient{}
	client := &Client{
		gnmiC:  setClient,
		target: "dut",
	}
	tests := []struct {
		desc         string
		op           func(*Client) (*gpb.SetResponse, error)
		wantErr      string
		wantRequest  *gpb.SetRequest
		stubResponse *gpb.SetResponse
		stubErr      error
	}{{
		desc: "success",
		op: func(c *Client) (*gpb.SetResponse, error) {
			q := &LeafConfigQuery[uint64]{
				leafBaseQuery: leafBaseQuery[uint64]{
					state:  false,
					scalar: true,
					ps:     ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "uint64-leaf"}, nil, ygot.NewDeviceRootBase("")),
				},
			}
			return Delete[uint64](context.Background(), c, q)
		},
		wantRequest: &gpb.SetRequest{
			Prefix: &gpb.Path{
				Target: "dut",
			},
			Delete: []*gpb.Path{
				testutil.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf"),
			},
		},
		stubResponse: &gpb.SetResponse{
			Prefix: &gpb.Path{
				Target: "dut",
			},
		},
	}, {
		desc: "server error",
		op: func(c *Client) (*gpb.SetResponse, error) {
			q := &LeafConfigQuery[uint64]{
				leafBaseQuery: leafBaseQuery[uint64]{
					state:  false,
					scalar: true,
					ps:     ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "uint64-leaf"}, nil, ygot.NewDeviceRootBase("")),
				},
			}
			return Delete[uint64](context.Background(), c, q)
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
			if diff := cmp.Diff(tt.wantRequest, setClient.Requests[0], protocmp.Transform()); diff != "" {
				t.Errorf("Delete() sent unexpected request (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.stubResponse, got, protocmp.Transform()); diff != "" {
				t.Errorf("Delete() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}
}

package ygnmi

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygnmi/testing/fakegnmi"
	"github.com/openconfig/ygnmi/testing/schema"
	"github.com/openconfig/ygot/testutil"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/testing/protocmp"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

func getClient(t testing.TB) (*fakegnmi.FakeGNMI, *Client) {
	fakeGNMI, err := fakegnmi.Start(0)
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
	leafPath := schema.GNMIPath(t, "super-container/leaf-container-struct/uint64-leaf")
	lq := &LeafSingletonQuery[uint64]{
		parentDir:  "leaf-container-struct",
		state:      true,
		ps:         ygot.NewNodePath([]string{"super-container", "leaf-container-struct", "uint64-leaf"}, nil, ygot.NewDeviceRootBase("")),
		extractFn:  func(vgs ygot.ValidatedGoStruct) uint64 { return *(vgs.(*schema.LeafContainerStruct)).Uint64Leaf },
		goStructFn: func() ygot.ValidatedGoStruct { return new(schema.LeafContainerStruct) },
		yschema:    schema.GetSchemaStruct()(),
	}

	leaftests := []struct {
		desc                 string
		stub                 func(s *fakegnmi.Stubber)
		inQuery              SingletonQuery[uint64]
		wantSubscriptionPath *gpb.Path
		wantVal              *Value[uint64]
		wantErr              string
	}{{
		desc:    "success update and sync",
		inQuery: lq,
		stub: func(s *fakegnmi.Stubber) {
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
		stub: func(s *fakegnmi.Stubber) {
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
		stub: func(s *fakegnmi.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 100,
				Prefix:    schema.GNMIPath(t, "super-container"),
				Update: []*gpb.Update{{
					Path: schema.GNMIPath(t, "leaf-container-struct/uint64-leaf"),
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
		stub: func(s *fakegnmi.Stubber) {
			s.Notification(&gpb.Notification{
				Update: []*gpb.Update{},
			}).Notification(&gpb.Notification{
				Timestamp: 100,
				Prefix:    schema.GNMIPath(t, "super-container"),
				Update: []*gpb.Update{{
					Path: schema.GNMIPath(t, "leaf-container-struct/uint64-leaf"),
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
		stub: func(s *fakegnmi.Stubber) {
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
		stub: func(s *fakegnmi.Stubber) {
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
		stub: func(s *fakegnmi.Stubber) {
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
		stub: func(s *fakegnmi.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: schema.GNMIPath(t, "super-container/leaf-container-struct/enum-leaf"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "E_VALUE_FORTY_THREE"}},
				}},
			}).Sync()
		},
		wantErr: "noncompliant data encountered while unmarshalling leaf",
	}, {
		desc:    "error non existant path",
		inQuery: lq,
		stub: func(s *fakegnmi.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: 101,
				Update: []*gpb.Update{{
					Path: schema.GNMIPath(t, "super-container/leaf-container-struct/does-not-exist"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		wantErr: "does-not-exist",
	}, {
		desc:    "error nil update",
		inQuery: lq,
		stub: func(s *fakegnmi.Stubber) {
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
		stub: func(s *fakegnmi.Stubber) {
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

// checks that the received time is just before now
func checkJustReceived(t *testing.T, recvTime time.Time) {
	if diffSecs := time.Now().Sub(recvTime).Seconds(); diffSecs <= 0 && diffSecs > 1 {
		t.Errorf("received time is too far (%v seconds) away from now", diffSecs)
	}
}

// verifySubscriptionPathsSent verifies the paths of the sent subscription requests is the same as wantPaths.
func verifySubscriptionPathsSent(t *testing.T, fakeGNMI *fakegnmi.FakeGNMI, wantPaths ...*gpb.Path) {
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
	if diff := cmp.Diff(wantPaths, gotPaths, protocmp.Transform(), cmpopts.SortSlices(testutil.PathLess)); diff != "" {
		t.Errorf("subscription paths (-want, +got):\n%s", diff)
	}
}

package ygnmi_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygnmi/internal/exampleoc"
	"github.com/openconfig/ygnmi/internal/exampleoc/device"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	ygottestutil "github.com/openconfig/ygot/testutil"
)

func TestLookup(t *testing.T) {
	fakeGNMI, c := getClient(t)
	leafPath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := device.DeviceRoot("").RemoteContainer().ALeaf().State()

	leaftests := []struct {
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
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantErr: "failed to unmarshal",
	}}
	for _, tt := range leaftests {
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

	configQuery := device.DeviceRoot("").Parent().Child().Config()
	stateQuery := device.DeviceRoot("").Parent().Child().State()

	nonLeafTests := []struct {
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

	for _, tt := range nonLeafTests {
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
}

func TestWatch(t *testing.T) {
	fakeGNMI, client := getClient(t)
	leafPath := testutil.GNMIPath(t, "/remote-container/state/a-leaf")
	lq := device.DeviceRoot("").RemoteContainer().ALeaf().State()

	startTime := time.Now()
	leafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[string]
		wantVals             []*ygnmi.Value[string]
		wantErr              string
	}{{
		desc: "single notif and pred true",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			}).Sync()
		},
		dur: time.Second,
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      leafPath,
			}).SetVal("foo")},
		wantSubscriptionPath: leafPath,
		wantLastVal: (&ygnmi.Value[string]{
			Timestamp: startTime,
			Path:      leafPath,
		}).SetVal("foo"),
	}, {
		desc: "single notif and pred false error EOF",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync()
		},
		dur: time.Second,
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      leafPath,
			}).SetVal("bar"),
		},
		wantSubscriptionPath: leafPath,
		wantLastVal: (&ygnmi.Value[string]{
			Timestamp: startTime,
			Path:      leafPath,
		}).SetVal("bar"),
		wantErr: "error receiving gNMI response: EOF",
	}, {
		desc: "multiple notif and pred true",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			})
		},
		dur: time.Second,
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      leafPath,
			}).SetVal("bar"),
			(&ygnmi.Value[string]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      leafPath,
			}).SetVal("foo"),
		},
		wantSubscriptionPath: leafPath,
		wantLastVal: (&ygnmi.Value[string]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      leafPath,
		}).SetVal("foo"),
	}, {
		desc: "multiple notif with deletes",
		stub: func(s *testutil.Stubber) {
			s.Notification(&gpb.Notification{
				Timestamp: startTime.UnixNano(),
				Update: []*gpb.Update{{
					Path: leafPath,
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "bar"}},
				}},
			}).Sync().Notification(&gpb.Notification{
				Timestamp: startTime.Add(time.Millisecond).UnixNano(),
				Delete:    []*gpb.Path{leafPath},
			})
		},
		dur: time.Second,
		wantVals: []*ygnmi.Value[string]{
			(&ygnmi.Value[string]{
				Timestamp: startTime,
				Path:      leafPath,
			}).SetVal("bar"),
			(&ygnmi.Value[string]{
				Timestamp: startTime.Add(time.Millisecond),
				Path:      leafPath,
			}),
		},
		wantSubscriptionPath: leafPath,
		wantLastVal: (&ygnmi.Value[string]{
			Timestamp: startTime.Add(time.Millisecond),
			Path:      leafPath,
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
	for _, tt := range leafTests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.stub(fakeGNMI.Stub())
			i := 0
			ctx, cancel := context.WithTimeout(context.Background(), tt.dur)
			defer cancel()
			w := ygnmi.Watch(ctx, client, lq, func(v *ygnmi.Value[string]) bool {
				if i > len(tt.wantVals) {
					t.Fatalf("Predicate(%d) expected no more values but got: %+v", i, v)
				}
				if diff := cmp.Diff(tt.wantVals[i], v, cmpopts.IgnoreFields(ygnmi.Value[string]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[string]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate(%d) got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", i, diff, v.ComplianceErrors)
				}
				val, present := v.Val()
				i++
				return present && val == "foo"
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
			if diff := cmp.Diff(tt.wantLastVal, val, cmp.AllowUnexported(ygnmi.Value[string]{}), protocmp.Transform()); diff != "" {
				t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
			}
		})
	}

	t.Run("multiple awaits", func(t *testing.T) {
		fakeGNMI.Stub().Sync()
		w := ygnmi.Watch(context.Background(), client, device.DeviceRoot("").RemoteContainer().ALeaf().State(), func(v *ygnmi.Value[string]) bool { return true })
		want := &ygnmi.Value[string]{
			Path: leafPath,
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
	nonLeafQuery := device.DeviceRoot("").Parent().Child().State()

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
			w := ygnmi.Watch(context.Background(), client, nonLeafQuery, func(v *ygnmi.Value[*exampleoc.Parent_Child]) bool {
				if len(tt.wantVals) == 0 {
					t.Fatalf("Predicate expected no more values but got: %+v", v)
				}
				if diff := cmp.Diff(tt.wantVals[0], v, cmpopts.IgnoreFields(ygnmi.Value[*exampleoc.Parent_Child]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[*exampleoc.Parent_Child]{}), protocmp.Transform()); diff != "" {
					t.Errorf("Predicate got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", diff, v.ComplianceErrors)
				}
				tt.wantVals = tt.wantVals[1:]
				val, present := v.Val()
				return present && val.One != nil && *val.One == "foo" && val.Three == exampleoc.Child_Three_ONE
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

func getClient(t testing.TB) (*testutil.FakeGNMI, *ygnmi.Client) {
	fakeGNMI, err := testutil.Start(0)
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

func TestLookupAll(t *testing.T) {
	fakeGNMI, c := getClient(t)
	leafPath := testutil.GNMIPath(t, "super-container/model/a/single-key[key=*]/state/value")
	lq := device.DeviceRoot("").Model().SingleKeyAny().Value().State()

	leaftests := []struct {
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
					Path: testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
					Val:  &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 10}},
				}},
			}).Sync()
		},
		wantVals: []*ygnmi.Value[int64]{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
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
		wantVals: []*ygnmi.Value[int64]{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
			Timestamp: time.Unix(0, 100),
		}, {
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/value"),
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
		wantVals: []*ygnmi.Value[int64]{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]/state/value"),
			Timestamp: time.Unix(0, 100),
		}, {
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]/state/value"),
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

	nonLeafPath := testutil.GNMIPath(t, "super-container/model/a/single-key[key=*]")
	nonLeafQ := device.DeviceRoot("").Model().SingleKeyAny().State()
	nonLeaftests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		wantSubscriptionPath *gpb.Path
		wantVals             []*ygnmi.Value[*exampleoc.Model_SingleKey]
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
		wantVals: []*ygnmi.Value[*exampleoc.Model_SingleKey]{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]"),
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
		wantVals: []*ygnmi.Value[*exampleoc.Model_SingleKey]{{
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=10]"),
			Timestamp: time.Unix(0, 100),
		}, {
			Path:      testutil.GNMIPath(t, "super-container/model/a/single-key[key=11]"),
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
}

func TestWatchAll(t *testing.T) {
	fakeGNMI, client := getClient(t)
	leafQueryPath := testutil.GNMIPath(t, "model/a/single-key[key=*]/state/value")
	key10Path := testutil.GNMIPath(t, "model/a/single-key[key=10]/state/value")
	key11Path := testutil.GNMIPath(t, "model/a/single-key[key=11]/state/value")

	startTime := time.Now()
	lq := device.DeviceRoot("").Model().SingleKeyAny().Value().State()
	tests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[int64]
		wantVals             []*ygnmi.Value[int64]
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
		wantVals: []*ygnmi.Value[int64]{{
			Timestamp: startTime,
			Path:      key10Path,
			val:       100,
			present:   true,
		}},
		wantLastVal: &ygnmi.Value[int64]{
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
		wantVals: []*ygnmi.Value[int64]{{
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
		wantLastVal: &ygnmi.Value[int64]{
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
		wantVals: []*ygnmi.Value[int64]{{
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
		wantLastVal: &ygnmi.Value[int64]{
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

			w := WatchAll[int64](ctx, client, lq, func(v *ygnmi.Value[int64]) bool {
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
			state:   false,
			ps:      nonLeafPS,
			yschema: testutil.GetSchemaStruct()(),
		},
	}
	nonLeafTests := []struct {
		desc                 string
		stub                 func(s *testutil.Stubber)
		dur                  time.Duration
		wantSubscriptionPath *gpb.Path
		wantLastVal          *ygnmi.Value[*exampleoc.Model_SingleKey]
		wantVals             []*ygnmi.Value[*exampleoc.Model_SingleKey]
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
		wantVals: []*ygnmi.Value[*exampleoc.Model_SingleKey]{{
			Timestamp: startTime,
			Path:      nonLeafKey10Path,
			val:       &testutil.Model_SingleKey{Value: ygot.Int64(100)},
			present:   true,
		}},
		wantLastVal: &ygnmi.Value[*exampleoc.Model_SingleKey]{
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
		wantVals: []*ygnmi.Value[*exampleoc.Model_SingleKey]{{
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
		wantLastVal: &ygnmi.Value[*exampleoc.Model_SingleKey]{
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

			w := WatchAll[*testutil.Model_SingleKey](ctx, client, nonLeafQ, func(v *ygnmi.Value[*exampleoc.Model_SingleKey]) bool {
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
	if diffSecs := time.Since(recvTime).Seconds(); diffSecs <= 0 && diffSecs > 1 {
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

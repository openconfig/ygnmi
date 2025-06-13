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
	"strings"
	"testing"
	"time"

	"google3/third_party/golang/cmp/cmp"
	"google3/third_party/golang/cmp/cmpopts/cmpopts"
	"google3/third_party/golang/protobuf/v2/proto/proto"
	"google3/third_party/golang/protobuf/v2/testing/protocmp/protocmp"
	"google3/third_party/golang/ygot/util/util"
	"google3/third_party/openconfig/gnmi/errdiff/errdiff"
	"google3/third_party/openconfig/ygnmi/exampleoc/exampleoc"
	"google3/third_party/openconfig/ygnmi/internal/exampleocconfig/exampleocconfig"
	"google3/third_party/openconfig/ygnmi/internal/gnmitestutil/gnmitestutil"
	"google3/third_party/openconfig/ygnmi/schemaless/schemaless"
	"google3/third_party/openconfig/ygnmi/ygnmi/ygnmi"

	ygottestutil "google3/third_party/golang/ygot/testutil/testutil"

	anypb "google3/google/protobuf/any_go_proto"
	gpb "google3/third_party/openconfig/gnmi/proto/gnmi/gnmi_go_proto"
)

func lookupCheckFn[T any](t *testing.T, fakeGNMI *gnmitestutil.FakeGNMI, c *ygnmi.Client, inQuery ygnmi.SingletonQuery[T], wantErrSubstring string, wantRequestValues *ygnmi.RequestValues, wantSubscriptionPath *gpb.Path, wantVal *ygnmi.Value[T], ygnmiOpts ...ygnmi.Option) {
	t.Helper()
	got, err := ygnmi.Lookup(context.Background(), c, inQuery, ygnmiOpts...)
	if diff := errdiff.Substring(err, wantErrSubstring); diff != "" {
		t.Fatalf("Lookup(ctx, c, %v) returned unexpected diff: %s", inQuery, diff)
	}
	if err != nil {
		return
	}
	verifySubscriptionPathsSent(t, fakeGNMI, wantSubscriptionPath)
	checkJustReceived(t, got.RecvTimestamp)
	wantVal.RecvTimestamp = got.RecvTimestamp

	if diff := cmp.Diff(wantVal, got, cmp.AllowUnexported(ygnmi.Value[T]{}, exampleoc.Model_SingleKey_OrderedList_OrderedMap{}, exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{}), protocmp.Transform()); diff != "" {
		t.Errorf("Lookup(ctx, c, %v) returned unexpected diff (-want,+got):\n %s\nComplianceErrors:\n%v", inQuery, diff, got.ComplianceErrors)
	}

	if wantRequestValues != nil {
		if diff := cmp.Diff(wantRequestValues, fakeGNMI.LastRequestContextValues()); diff != "" {
			t.Errorf("RequestValues (-want, +got):\n%s", diff)
		}
	}
}

func lookupWithGetCheckFn[T any](t *testing.T, fakeGNMI *gnmitestutil.FakeGNMI, c *ygnmi.Client, inQuery ygnmi.SingletonQuery[T], wantErrSubstring string, wantRequestValues *ygnmi.RequestValues, wantRequest *gpb.GetRequest, wantVal *ygnmi.Value[T], nonLeaf bool) {
	t.Helper()
	got, err := ygnmi.Lookup(context.Background(), c, inQuery, ygnmi.WithUseGet())
	if diff := errdiff.Substring(err, wantErrSubstring); diff != "" {
		t.Fatalf("Lookup(ctx, c, %v) returned unexpected diff (-want, +got):\n%s", inQuery, diff)
	}
	if err != nil {
		return
	}
	copts := []cmp.Option{cmp.AllowUnexported(ygnmi.Value[T]{}, exampleoc.Model_SingleKey_OrderedList_OrderedMap{}, exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{}), cmpopts.IgnoreFields(ygnmi.Value[T]{}, "RecvTimestamp"), protocmp.Transform()}
	if nonLeaf {
		copts = append(copts, cmpopts.IgnoreFields(ygnmi.TelemetryError{}, "Err"))
	}
	if diff := cmp.Diff(wantVal, got, copts...); diff != "" {
		t.Errorf("Lookup() returned unexpected diff (-want, +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantRequest, fakeGNMI.GetRequests()[0], protocmp.Transform()); diff != "" {
		t.Errorf("Lookup() GetRequest different from expected (-want, +got):\n%s", diff)
	}

	if wantRequestValues != nil {
		if diff := cmp.Diff(wantRequestValues, fakeGNMI.LastRequestContextValues()); diff != "" {
			t.Errorf("RequestValues (-want, +got):\n%s", diff)
		}
	}
}

func getCheckFn[T any](t *testing.T, fakeGNMI *gnmitestutil.FakeGNMI, c *ygnmi.Client, inQuery ygnmi.SingletonQuery[T], wantErrSubstring string, wantRequestValues *ygnmi.RequestValues, wantSubscriptionPath *gpb.Path, wantVal T, ygnmiOpts ...ygnmi.Option) {
	t.Helper()
	got, err := ygnmi.Get(context.Background(), c, inQuery, ygnmiOpts...)
	if diff := errdiff.Substring(err, wantErrSubstring); diff != "" {
		t.Fatalf("Get(ctx, c, %v) returned unexpected diff: %s", inQuery, diff)
	}
	if err != nil {
		return
	}
	verifySubscriptionPathsSent(t, fakeGNMI, wantSubscriptionPath)

	if diff := cmp.Diff(wantVal, got, cmp.AllowUnexported(exampleoc.Model_SingleKey_OrderedList_OrderedMap{}, exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{})); diff != "" {
		t.Errorf("Get(ctx, c, %v) returned unexpected diff (-want,+got):\n %s", inQuery, diff)
	}

	if wantRequestValues != nil {
		if diff := cmp.Diff(wantRequestValues, fakeGNMI.LastRequestContextValues()); diff != "" {
			t.Errorf("RequestValues (-want, +got):\n%s", diff)
		}
	}
}

func watchCheckFn[T any](t *testing.T, fakeGNMI *gnmitestutil.FakeGNMI, duration time.Duration, c *ygnmi.Client, inQuery ygnmi.SingletonQuery[T],
	valPred func(T) bool, wantErrSubstring string, wantRequestValues *ygnmi.RequestValues, wantSubscriptionPaths []*gpb.Path, wantModes []gpb.SubscriptionMode, wantIntervals []uint64, wantVals []*ygnmi.Value[T], wantLastVal *ygnmi.Value[T], ygnmiOpts ...ygnmi.Option) {
	t.Helper()
	i := 0
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	extractErrorMessage := cmp.Transformer("ExtractErrorString", func(err error) string {
		if err == nil {
			return ""
		}
		return err.Error()
	})

	w := ygnmi.Watch(ctx, c, inQuery, func(v *ygnmi.Value[T]) error {
		if i > len(wantVals) {
			t.Fatalf("Predicate(%d) expected no more values but got: %+v", i, v)
		}
		if diff := cmp.Diff(wantVals[i], v, extractErrorMessage, cmpopts.IgnoreFields(ygnmi.Value[T]{}, "RecvTimestamp", "Timestamp"), cmp.AllowUnexported(ygnmi.Value[T]{}, exampleoc.Model_SingleKey_OrderedList_OrderedMap{}, exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{}), protocmp.Transform()); diff != "" {
			t.Errorf("Predicate(%d) got unexpected input (-want,+got):\n %s\nComplianceErrors:\n%v", i, diff, v.ComplianceErrors)
		}
		val, present := v.Val()
		i++
		if present && valPred(val) {
			return nil
		}
		return ygnmi.Continue
	}, ygnmiOpts...)
	val, err := w.Await()
	if i < len(wantVals) {
		t.Errorf("Predicate received too few values: got %d, want %d", i, len(wantVals))
	}
	if diff := errdiff.Substring(err, wantErrSubstring); diff != "" {
		t.Fatalf("Await() returned unexpected diff: %s", diff)
	}
	if err != nil {
		return
	}
	verifySubscriptionPathsSent(t, fakeGNMI, wantSubscriptionPaths...)
	verifySubscriptionModesSent(t, fakeGNMI, wantModes...)
	verifySubscriptionSampleIntervalsSent(t, fakeGNMI, wantIntervals...)
	if val != nil {
		checkJustReceived(t, val.RecvTimestamp)
		wantLastVal.RecvTimestamp = val.RecvTimestamp
	}
	if diff := cmp.Diff(wantLastVal, val, extractErrorMessage, cmp.AllowUnexported(ygnmi.Value[T]{}, exampleoc.Model_SingleKey_OrderedList_OrderedMap{}, exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{}), protocmp.Transform()); diff != "" {
		t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
	}

	if wantRequestValues != nil {
		if diff := cmp.Diff(wantRequestValues, fakeGNMI.LastRequestContextValues()); diff != "" {
			t.Errorf("RequestValues (-want, +got):\n%s", diff)
		}
	}
}

func collectCheckFn[T any](t *testing.T, fakeGNMI *gnmitestutil.FakeGNMI, c *ygnmi.Client, inQuery ygnmi.SingletonQuery[T], wantErrSubstring string, wantRequestValues *ygnmi.RequestValues, wantSubscriptionPath *gpb.Path, wantVals []*ygnmi.Value[T]) {
	t.Helper()
	vals, err := ygnmi.Collect(context.Background(), c, inQuery).Await()
	if diff := errdiff.Substring(err, wantErrSubstring); diff != "" {
		t.Errorf("Await() returned unexpected diff: %s", diff)
	}
	verifySubscriptionPathsSent(t, fakeGNMI, wantSubscriptionPath)
	for _, val := range vals {
		checkJustReceived(t, val.RecvTimestamp)
	}
	if diff := cmp.Diff(wantVals, vals, cmpopts.IgnoreFields(ygnmi.Value[T]{}, "RecvTimestamp"), cmp.AllowUnexported(ygnmi.Value[T]{}), protocmp.Transform()); diff != "" {
		t.Errorf("Await() returned unexpected value (-want,+got):\n%s", diff)
	}

	if wantRequestValues != nil {
		if diff := cmp.Diff(wantRequestValues, fakeGNMI.LastRequestContextValues()); diff != "" {
			t.Errorf("RequestValues (-want, +got):\n%s", diff)
		}
	}
}

func lookupAllCheckFn[T any](t *testing.T, fakeGNMI *gnmitestutil.FakeGNMI, c *ygnmi.Client, inQuery ygnmi.WildcardQuery[T], wantErrSubstring string, wantRequestValues *ygnmi.RequestValues, wantSubscriptionPath *gpb.Path, wantVals []*ygnmi.Value[T], nonLeaf bool) {
	t.Helper()
	got, err := ygnmi.LookupAll(context.Background(), c, inQuery)
	if diff := errdiff.Substring(err, wantErrSubstring); diff != "" {
		t.Fatalf("LookupAll(ctx, c, %v) returned unexpected diff: %s", inQuery, diff)
	}
	if err != nil {
		return
	}
	verifySubscriptionPathsSent(t, fakeGNMI, wantSubscriptionPath)
	for _, val := range got {
		checkJustReceived(t, val.RecvTimestamp)
	}
	copts := []cmp.Option{
		cmp.AllowUnexported(ygnmi.Value[T]{}, exampleoc.Model_SingleKey_OrderedList_OrderedMap{}, exampleocconfig.Model_SingleKey_OrderedList_OrderedMap{}),
		cmpopts.IgnoreFields(ygnmi.Value[T]{}, "RecvTimestamp"),
		protocmp.Transform(),
	}
	if nonLeaf {
		copts = append(copts, cmpopts.IgnoreFields(ygnmi.TelemetryError{}, "Err"))
	}
	if diff := cmp.Diff(wantVals, got, copts...); diff != "" {
		t.Errorf("LookupAll() returned unexpected diff (-want,+got):\n%s", diff)
	}

	if wantRequestValues != nil {
		if diff := cmp.Diff(wantRequestValues, fakeGNMI.LastRequestContextValues()); diff != "" {
			t.Errorf("RequestValues (-want, +got):\n%s", diff)
		}
	}
}

func configCheckFn(t *testing.T, setClient *gnmitestutil.SetClient, c *ygnmi.Client, op func(*ygnmi.Client) (*ygnmi.Result, error), wantRequest *gpb.SetRequest, stubResponse *gpb.SetResponse, wantErrSubstring string, stubErr error) {
	t.Helper()
	setClient.Reset()
	setClient.AddResponse(stubResponse, stubErr)

	got, err := op(c)
	if diff := errdiff.Substring(err, wantErrSubstring); diff != "" {
		t.Fatalf("config operation returned unexpected diff: %s", diff)
	}
	if err != nil {
		return
	}
	if diff := cmp.Diff(wantRequest, setClient.Requests[0], protocmp.Transform()); diff != "" {
		t.Errorf("config operation sent unexpected request (-want,+got):\n%s", diff)
	}
	want := &ygnmi.Result{
		RawResponse: stubResponse,
		Timestamp:   time.Unix(0, stubResponse.GetTimestamp()),
	}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("config operation returned unexpected value (-want,+got):\n%s", diff)
	}
}

type testStruct struct {
	Val string
}

func mustSchemaless[T any](t testing.TB, path, origin string) ygnmi.ConfigQuery[T] {
	q, err := schemaless.NewConfig[T](path, origin)
	if err != nil {
		t.Fatal(err)
	}
	return q
}

func mustAnyNew(t testing.TB, m proto.Message) *anypb.Any {
	any, err := anypb.New(m)
	if err != nil {
		t.Fatal(err)
	}
	return any
}

func removeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), "")
}

func newClient(t testing.TB) (*gnmitestutil.FakeGNMI, *ygnmi.Client) {
	fakeGNMI, err := gnmitestutil.StartGNMI(0)
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

// checkJustReceived checks that the received time is just before now.
func checkJustReceived(t *testing.T, recvTime time.Time) {
	if diffSecs := time.Since(recvTime).Seconds(); diffSecs <= 0 && diffSecs > 1 {
		t.Errorf("Received time is too far (%v seconds) away from now", diffSecs)
	}
}

// verifySubscriptionPathsSent verifies the paths of the sent subscription requests is the same as wantPaths.
func verifySubscriptionPathsSent(t *testing.T, fakeGNMI *gnmitestutil.FakeGNMI, wantPaths ...*gpb.Path) {
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
func verifySubscriptionModesSent(t *testing.T, fakeGNMI *gnmitestutil.FakeGNMI, wantModes ...gpb.SubscriptionMode) {
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

// verifySubscriptionSampleIntervalsSent verifies the modes of the sent subscription requests is the same as wantModes.
func verifySubscriptionSampleIntervalsSent(t *testing.T, fakeGNMI *gnmitestutil.FakeGNMI, wantIntervals ...uint64) {
	t.Helper()
	requests := fakeGNMI.Requests()
	if len(requests) != 1 {
		t.Errorf("Number of subscription requests sent is not 1: %v", requests)
		return
	}

	var gotIntervals []uint64
	req := requests[0].GetSubscribe()
	for _, sub := range req.GetSubscription() {
		gotIntervals = append(gotIntervals, sub.SampleInterval)
	}
	if diff := cmp.Diff(wantIntervals, gotIntervals); diff != "" {
		t.Errorf("Subscription sample intervals (-want, +got):\n%s", diff)
	}
}

type gnmiS struct {
	gpb.UnimplementedGNMIServer
	errCh chan error
}

func (g *gnmiS) Subscribe(srv gpb.GNMI_SubscribeServer) error {
	if _, err := srv.Recv(); err != nil {
		return err
	}
	if err := srv.Send(&gpb.SubscribeResponse{Response: &gpb.SubscribeResponse_SyncResponse{}}); err != nil {
		return err
	}
	// This send must fail because the client will have cancelled the subscription context.
	time.Sleep(time.Second)
	err := srv.Send(&gpb.SubscribeResponse{Response: &gpb.SubscribeResponse_SyncResponse{}})
	g.errCh <- err
	return nil
}


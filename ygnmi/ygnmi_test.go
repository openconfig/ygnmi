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

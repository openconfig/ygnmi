// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package gnmitestutil implements a fake GNMI server with the ability to stub responses and fake schema.
package gnmitestutil

import (
	"context"
	"io"

	"github.com/openconfig/gnmi/testing/fake/gnmi"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	fpb "github.com/openconfig/gnmi/testing/fake/proto"
)

// FakeGNMI is a running fake GNMI server.
type FakeGNMI struct {
	agent         *gnmi.Agent
	stub          *Stubber
	clientWrapper *clientWrap
}

// StartGNMI launches a new fake GNMI server on the given port
func StartGNMI(port int) (*FakeGNMI, error) {
	gen := &fpb.FixedGenerator{}
	config := &fpb.Config{
		Port:        int32(port),
		Generator:   &fpb.Config_Fixed{Fixed: gen},
		DisableSync: true,  // Require every sync to be stubbed explicitly.
		EnableDelay: true,  // Respect timestamps if present.
		DisableEof:  false, // Close subscribe RPCs when the data runs out.
	}
	agent, err := gnmi.New(config, nil)
	if err != nil {
		return nil, err
	}
	stub := &Stubber{gen: gen}
	return &FakeGNMI{
		agent:         agent,
		stub:          stub,
		clientWrapper: &clientWrap{stub: stub},
	}, nil
}

// Dial dials the fake gNMI client and returns a gNMI client stub.
func (g *FakeGNMI) Dial(ctx context.Context, opts ...grpc.DialOption) (gpb.GNMIClient, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.DialContext(ctx, g.agent.Address(), opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "DialContext(%s, %v)", g.agent.Address(), opts)
	}
	g.clientWrapper.GNMIClient = gpb.NewGNMIClient(conn)
	return g.clientWrapper, nil
}

// Stub reset the stubbed responses to empty and returns a handle to add new ones.
func (g *FakeGNMI) Stub() *Stubber {
	g.stub.gen.Reset()
	return g.stub
}

// Requests returns the set of SubscribeRequests sent to the gNMI server.
func (g *FakeGNMI) Requests() []*gpb.SubscribeRequest {
	return g.agent.Requests()
}

// LastRequestContextValues returns the request-scoped values from the last
// ygnmi request sent to the gNMI server.
func (g *FakeGNMI) LastRequestContextValues() *ygnmi.RequestValues {
	return g.clientWrapper.requestValues
}

// GetRequests returns the set of GetRequests sent to the gNMI server.
//
// They're ordered in reverse time of request.
func (g *FakeGNMI) GetRequests() []*gpb.GetRequest {
	return g.clientWrapper.getRequests
}

// clientWrap adds gNMI Get functionality to a GNMI client.
type clientWrap struct {
	gpb.GNMIClient
	stub          *Stubber
	getRequests   []*gpb.GetRequest
	requestValues *ygnmi.RequestValues
}

// Get is a fake implementation of gnmi.Get, it returns the GetResponse contained in the stub.
func (g *clientWrap) Get(ctx context.Context, req *gpb.GetRequest, _ ...grpc.CallOption) (*gpb.GetResponse, error) {
	g.requestValues = ygnmi.FromContext(ctx)
	g.getRequests = append([]*gpb.GetRequest{req}, g.getRequests...)
	if len(g.stub.getResponses) == 0 {
		return nil, io.EOF
	}
	resp := g.stub.getResponses[0]
	g.stub.getResponses = g.stub.getResponses[1:]
	err := g.stub.getErrs[0]
	g.stub.getErrs = g.stub.getErrs[1:]
	return resp, err
}

// Subscribe is a wrapper of the fake implementation of gnmi.Subscribe.
func (g *clientWrap) Subscribe(ctx context.Context, opts ...grpc.CallOption) (gpb.GNMI_SubscribeClient, error) {
	g.requestValues = ygnmi.FromContext(ctx)
	return g.GNMIClient.Subscribe(ctx, opts...)
}

// Stubber is a handle to add stubbed responses.
type Stubber struct {
	gen          *fpb.FixedGenerator
	getResponses []*gpb.GetResponse
	getErrs      []error
}

// Notification appends the given notification as a stub response.
func (s *Stubber) Notification(n *gpb.Notification) *Stubber {
	s.gen.Responses = append(s.gen.Responses, &gpb.SubscribeResponse{
		Response: &gpb.SubscribeResponse_Update{Update: n},
	})
	return s
}

// GetResponse appends the given GetResponse as a stub response.
func (s *Stubber) GetResponse(gr *gpb.GetResponse, err error) *Stubber {
	s.getResponses = append(s.getResponses, gr)
	s.getErrs = append(s.getErrs, err)
	return s
}

// Sync appends a sync stub response.
func (s *Stubber) Sync() *Stubber {
	s.gen.Responses = append(s.gen.Responses, &gpb.SubscribeResponse{
		Response: &gpb.SubscribeResponse_SyncResponse{},
	})
	return s
}

// SetClient is a test stub for ygnmi APIs exercising gNMI.Set.
type SetClient struct {
	gpb.GNMIClient
	// Responses are the gNMI Responses to return from calls to Set.
	Responses []*gpb.SetResponse
	// Requests received by the client are stored in the slice.
	Requests []*gpb.SetRequest
	// ResponseErrs are the errors to return from calls to Set.
	ResponseErrs []error
	// i is index current index of the response and error to return.
	i int
}

// Reset clears SetClient's internal state.
func (f *SetClient) Reset() {
	f.Requests = nil
	f.Responses = nil
	f.ResponseErrs = nil
	f.i = 0
}

// AddResponse appends a gNMI SetResponse to be returned when gNMI.Set is
// called.
func (f *SetClient) AddResponse(resp *gpb.SetResponse, err error) *SetClient {
	f.Responses = append(f.Responses, resp)
	f.ResponseErrs = append(f.ResponseErrs, err)
	return f
}

// Set is a stub that records the provided request, and returns the next stub
// response.
func (f *SetClient) Set(_ context.Context, req *gpb.SetRequest, opts ...grpc.CallOption) (*gpb.SetResponse, error) {
	defer func() { f.i++ }()
	f.Requests = append(f.Requests, req)
	return f.Responses[f.i], f.ResponseErrs[f.i]
}

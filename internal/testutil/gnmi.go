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

// Package testutil implements a fake GNMI server with the ability to stub responses and fake schema.
package testutil

import (
	"context"
	"io"

	"github.com/openconfig/gnmi/testing/fake/gnmi"
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
	clientWrapper *clientWithGetter
}

// StartGNMI launches a new fake GNMI server on the given port
func StartGNMI(port int) (*FakeGNMI, error) {
	gen := &fpb.FixedGenerator{}
	config := &fpb.Config{
		Port:        int32(port),
		Generator:   &fpb.Config_Fixed{Fixed: gen},
		DisableSync: true, // Require every sync to be stubbed explicitly.
		EnableDelay: true, // Respect timestamps if present.
		DisableEof:  true,
	}
	agent, err := gnmi.New(config, nil)
	if err != nil {
		return nil, err
	}
	stub := &Stubber{gen: gen}
	return &FakeGNMI{
		agent:         agent,
		stub:          stub,
		clientWrapper: &clientWithGetter{stub: stub},
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

// GetRequests returns the set of GetRequests sent to the gNMI server.
func (g *FakeGNMI) GetRequests() []*gpb.GetRequest {
	return g.clientWrapper.getRequests
}

// clientWithGetter adds gNMI Get functionality to a GNMI client.
type clientWithGetter struct {
	gpb.GNMIClient
	stub        *Stubber
	getRequests []*gpb.GetRequest
}

// Get is a fake implementation of gnmi.Get, it returns the GetResponse contained in the stub.
func (g *clientWithGetter) Get(ctx context.Context, req *gpb.GetRequest, _ ...grpc.CallOption) (*gpb.GetResponse, error) {
	g.getRequests = append(g.getRequests, req)
	if len(g.stub.getResponses) == 0 {
		return nil, io.EOF
	}
	resp := g.stub.getResponses[0]
	g.stub.getResponses = g.stub.getResponses[1:]
	return resp, nil
}

// Stubber is a handle to add stubbed responses.
type Stubber struct {
	gen          *fpb.FixedGenerator
	getResponses []*gpb.GetResponse
}

// Notification appends the given notification as a stub response.
func (s *Stubber) Notification(n *gpb.Notification) *Stubber {
	s.gen.Responses = append(s.gen.Responses, &gpb.SubscribeResponse{
		Response: &gpb.SubscribeResponse_Update{Update: n},
	})
	return s
}

// GetResponse appends the given GetResponse as a stub response.
func (s *Stubber) GetResponse(gr *gpb.GetResponse) *Stubber {
	s.getResponses = append(s.getResponses, gr)
	return s
}

// Sync appends a sync stub response.
func (s *Stubber) Sync() *Stubber {
	s.gen.Responses = append(s.gen.Responses, &gpb.SubscribeResponse{
		Response: &gpb.SubscribeResponse_SyncResponse{},
	})
	return s
}

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

func (f *SetClient) Reset() {
	f.Requests = nil
	f.Responses = nil
	f.ResponseErrs = nil
	f.i = 0
}

func (f *SetClient) AddResponse(resp *gpb.SetResponse, err error) *SetClient {
	f.Responses = append(f.Responses, resp)
	f.ResponseErrs = append(f.ResponseErrs, err)
	return f
}

func (f *SetClient) Set(_ context.Context, req *gpb.SetRequest, opts ...grpc.CallOption) (*gpb.SetResponse, error) {
	defer func() { f.i++ }()
	f.Requests = append(f.Requests, req)
	return f.Responses[f.i], f.ResponseErrs[f.i]
}

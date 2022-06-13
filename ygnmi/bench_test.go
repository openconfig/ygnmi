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

//nolint:errcheck
package ygnmi_test

import (
	"context"
	"io"
	"testing"

	"github.com/openconfig/ygnmi/internal/exampleoc"
	"github.com/openconfig/ygnmi/internal/exampleoc/root"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// subClient will send numResponses notifications, then return EOF.
// If numResponses > len(resp), the last notification is sent repeatedly.
type subClient struct {
	gpb.GNMI_SubscribeClient
	resp []*gpb.SubscribeResponse
	// numResponses is number of non-error responses to return (can be larger than len(resp)).
	numResponses int
	recvs        int
}

func (c *subClient) CloseSend() error { return nil }
func (c *subClient) Send(*gpb.SubscribeRequest) error {
	c.recvs = 0
	return nil
}

func (c *subClient) Recv() (*gpb.SubscribeResponse, error) {
	if c.recvs >= c.numResponses {
		return nil, io.EOF
	}
	idx := c.recvs
	if idx >= len(c.resp) {
		idx = len(c.resp) - 1
	}
	c.recvs += 1
	return c.resp[idx], nil
}

type benchmarkClient struct {
	gpb.GNMIClient
	sc *subClient
}

func (c *benchmarkClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (gpb.GNMI_SubscribeClient, error) {
	return c.sc, nil
}

func BenchmarkGet(b *testing.B) {
	bc := &benchmarkClient{sc: &subClient{numResponses: 1}}
	c, err := ygnmi.NewClient(bc)
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	setUpdate := func(path *gpb.Path, val *gpb.TypedValue) {
		bc.sc.resp = []*gpb.SubscribeResponse{{
			Response: &gpb.SubscribeResponse_Update{
				Update: &gpb.Notification{
					Update: []*gpb.Update{{
						Path: path,
						Val:  val,
					}},
				},
			},
		}}
	}

	b.Run("deeply nested leaf", func(b *testing.B) {
		ctx := context.Background()
		q := root.New().A().B().C().D().E().F().G().H().I().J().K().L().M().Foo().State()
		setUpdate(testutil.GNMIPath(b, "/a/b/c/d/e/f/g/h/i/j/k/l/m/state/foo"), &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "sample"}})
		var got string
		for i := 0; i < b.N; i++ {
			got, _ = ygnmi.Get(ctx, c, q)
		}
		b.StopTimer()
		if got != "sample" {
			b.Fatalf("Get unexpected result: got %v want %v", got, "sample")
		}
	})
	b.Run("leaf update into root", func(b *testing.B) {
		ctx := context.Background()
		q := root.New().State()
		setUpdate(testutil.GNMIPath(b, "/a/b/c/d/e/f/g/h/i/j/k/l/m/state/foo"), &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "sample"}})
		var got *exampleoc.Root
		for i := 0; i < b.N; i++ {
			got, _ = ygnmi.Get(ctx, c, q)
		}
		b.StopTimer()
		if got.GetA().GetB().GetC().GetD().GetE().GetF().GetG().GetH().GetI().GetJ().GetK().GetL().GetM().GetFoo() != "sample" {
			b.Fatalf("Get unexpected result: got %v want %v", got, "sample")
		}
	})
	b.Run("list", func(b *testing.B) {
		ctx := context.Background()
		q := root.New().Model().SingleKeyAny().State()
		setUpdate(testutil.GNMIPath(b, "/model/a/single-key[key=\"foo\"]/state/value"), &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 1}})
		var got []*exampleoc.Model_SingleKey
		for i := 0; i < b.N; i++ {
			got, _ = ygnmi.GetAll(ctx, c, q)
		}
		b.StopTimer()
		if *got[0].Value != 1 {
			b.Fatalf("GetAll unexpected result: got %v want %v", *got[0].Value, 1)
		}
	})
}

func BenchmarkWatch(b *testing.B) {
	bc := &benchmarkClient{sc: &subClient{numResponses: 100}}
	c, err := ygnmi.NewClient(bc)
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	setUpdate := func(path *gpb.Path, val *gpb.TypedValue) {
		bc.sc.resp = []*gpb.SubscribeResponse{{
			Response: &gpb.SubscribeResponse_Update{
				Update: &gpb.Notification{
					Update: []*gpb.Update{{
						Path: path,
						Val:  val,
					}},
				},
			},
		}, {
			Response: &gpb.SubscribeResponse_SyncResponse{},
		}, {
			Response: &gpb.SubscribeResponse_Update{
				Update: &gpb.Notification{
					Update: []*gpb.Update{{
						Path: path,
						Val:  val,
					}},
				},
			},
		}}
	}

	b.Run("deeply nested leaf", func(b *testing.B) {
		ctx := context.Background()
		q := root.New().A().B().C().D().E().F().G().H().I().J().K().L().M().Foo().State()
		setUpdate(testutil.GNMIPath(b, "/a/b/c/d/e/f/g/h/i/j/k/l/m/state/foo"), &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "sample"}})
		var got *ygnmi.Value[string]
		for i := 0; i < b.N; i++ {
			got, _ = ygnmi.Watch(ctx, c, q, func(v *ygnmi.Value[string]) error {
				return ygnmi.Continue
			}).Await()
		}
		b.StopTimer()
		if v, _ := got.Val(); v != "sample" {
			b.Fatalf("Watch unexpected result: got %v want %v", v, "sample")
		}
	})
	b.Run("leaf update into root", func(b *testing.B) {
		ctx := context.Background()
		q := root.New().State()
		setUpdate(testutil.GNMIPath(b, "/a/b/c/d/e/f/g/h/i/j/k/l/m/state/foo"), &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "sample"}})
		var got *ygnmi.Value[*exampleoc.Root]
		for i := 0; i < b.N; i++ {
			got, _ = ygnmi.Watch(ctx, c, q, func(v *ygnmi.Value[*exampleoc.Root]) error {
				return ygnmi.Continue
			}).Await()
		}
		b.StopTimer()
		if v, _ := got.Val(); v.GetA().GetB().GetC().GetD().GetE().GetF().GetG().GetH().GetI().GetJ().GetK().GetL().GetM().GetFoo() != "sample" {
			b.Fatalf("Watch unexpected result: got %v want %v", v, "sample")
		}
	})
	b.Run("list", func(b *testing.B) {
		ctx := context.Background()
		q := root.New().Model().SingleKeyAny().State()
		setUpdate(testutil.GNMIPath(b, "/model/a/single-key[key=\"foo\"]/state/value"), &gpb.TypedValue{Value: &gpb.TypedValue_IntVal{IntVal: 1}})
		var got *ygnmi.Value[*exampleoc.Model_SingleKey]
		for i := 0; i < b.N; i++ {
			got, _ = ygnmi.WatchAll(ctx, c, q, func(v *ygnmi.Value[*exampleoc.Model_SingleKey]) error {
				return ygnmi.Continue
			}).Await()
		}
		b.StopTimer()
		if v, _ := got.Val(); *v.Value != 1 {
			b.Fatalf("WatchAll unexpected result: got %v want %v", v.Value, 1)
		}
	})
}

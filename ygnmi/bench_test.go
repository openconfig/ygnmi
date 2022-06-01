package ygnmi_test

import (
	"context"
	"io"
	"testing"

	"github.com/openconfig/ygnmi/internal/exampleoc/root"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

type subClient struct {
	gpb.GNMI_SubscribeClient
	resp *gpb.SubscribeResponse
	sent bool
}

func (c *subClient) CloseSend() error { return nil }
func (c *subClient) Send(*gpb.SubscribeRequest) error {
	c.sent = false
	return nil
}

func (c *subClient) Recv() (*gpb.SubscribeResponse, error) {
	if c.sent {
		return nil, io.EOF
	}
	c.sent = true
	return c.resp, nil
}

type benchmarkClient struct {
	gpb.GNMIClient
	sc *subClient
}

func (c *benchmarkClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (gpb.GNMI_SubscribeClient, error) {
	return c.sc, nil
}

func BenchmarkGet(b *testing.B) {
	bc := &benchmarkClient{sc: &subClient{}}
	c, err := ygnmi.NewClient(bc)
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	setUpdate := func(path *gpb.Path, val *gpb.TypedValue) {
		bc.sc.resp = &gpb.SubscribeResponse{
			Response: &gpb.SubscribeResponse_Update{
				Update: &gpb.Notification{
					Update: []*gpb.Update{{
						Path: path,
						Val:  val,
					}},
				},
			},
		}
	}

	b.Run("deeply nested leaf", func(b *testing.B) {
		ctx := context.Background()
		q := root.New().A().B().C().D().E().F().G().H().I().J().K().L().M().Foo().State()
		setUpdate(testutil.GNMIPath(b, "/a/b/c/d/e/f/g/h/i/j/k/l/m/state/foo"), &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "sample"}})
		for i := 0; i < b.N; i++ {
			ygnmi.Get(ctx, c, q)
		}
	})
	b.Run("leaf update into root", func(b *testing.B) {
		ctx := context.Background()
		q := root.New().State()
		setUpdate(testutil.GNMIPath(b, "/a/b/c/d/e/f/g/h/i/j/k/l/m/state/foo"), &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "sample"}})
		for i := 0; i < b.N; i++ {
			ygnmi.Get(ctx, c, q)
		}
	})
}

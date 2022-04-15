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

package ygnmi

import (
	"context"
	"io"
	"time"

	"github.com/openconfig/gnmi/errlist"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"

	log "github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	closer "github.com/openconfig/gocloser"
)

// subscribe create a gNMI SubscribeClient for the given query.
func subscribe[T any](ctx context.Context, c *Client, q AnyQuery[T], mode gpb.SubscriptionList_Mode) (_ gpb.GNMI_SubscribeClient, rerr error) {
	path, _, errs := ygot.ResolvePath(q.pathStruct())
	if len(errs) > 0 {
		l := errlist.List{}
		l.Add(errs...)
		return nil, l.Err()
	}

	sub, err := c.gnmiC.Subscribe(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "gNMI failed to Subscribe")
	}
	defer closer.Close(&rerr, sub.CloseSend, "error closing gNMI send stream")

	subs := []*gpb.Subscription{{
		Path: &gpb.Path{
			Elem:   path.GetElem(),
			Origin: path.GetOrigin(),
		},
		Mode: gpb.SubscriptionMode_TARGET_DEFINED,
	}}

	sr := &gpb.SubscribeRequest{
		Request: &gpb.SubscribeRequest_Subscribe{
			Subscribe: &gpb.SubscriptionList{
				Prefix: &gpb.Path{
					Target: c.target,
				},
				Subscription: subs,
				Mode:         mode,
				Encoding:     gpb.Encoding_PROTO,
			},
		},
	}

	log.V(1).Info(prototext.Format(sr))
	if err := sub.Send(sr); err != nil {
		return nil, errors.Wrapf(err, "gNMI failed to Send(%+v)", sr)
	}

	return sub, nil
}

// receive processes a single response from the subscription stream. If an "update" response is
// received, those points are appended to the given data and the result of that concatenation is
// the first return value, and the second return value is false. If a "sync" response is received,
// the data is returned as-is and the second return value is true. If Delete paths are present in
// the update, they are appended to the given data before the Update values. If deletesExpected
// is false, however, any deletes received will cause an error.
//nolint:deadcode // TODO(DanG100) remove this once this func is used
func receive(sub gpb.GNMI_SubscribeClient, data []*DataPoint, deletesExpected bool) ([]*DataPoint, bool, error) {
	res, err := sub.Recv()
	if err != nil {
		return data, false, err
	}
	recvTS := time.Now()

	switch v := res.Response.(type) {
	case *gpb.SubscribeResponse_Update:
		n := v.Update
		if !deletesExpected && len(n.Delete) != 0 {
			return data, false, errors.Errorf("unexpected delete updates: %v", n.Delete)
		}
		ts := time.Unix(0, n.GetTimestamp())
		newDataPoint := func(p *gpb.Path, val *gpb.TypedValue) (*DataPoint, error) {
			j, err := util.JoinPaths(n.GetPrefix(), p)
			if err != nil {
				return nil, err
			}
			// Record the deprecated Element field for clearer compliance error messages.
			//lint:ignore SA1019 ignore deprecated check
			if elements := append(append([]string{}, n.GetPrefix().GetElement()...), p.GetElement()...); len(elements) > 0 {
				//lint:ignore SA1019 ignore deprecated check
				j.Element = elements
			}
			// Use the target only for the subscription but exclude from the datapoint construction.
			j.Target = ""
			return &DataPoint{Path: j, Value: val, Timestamp: ts, RecvTimestamp: recvTS}, nil
		}

		// Append delete data before the update values -- per gNMI spec, they
		// should always be processed first if both update types exist in the
		// same notification.
		for _, p := range n.Delete {
			log.V(2).Infof("Received gNMI Delete at path: %s", prototext.Format(p))
			dp, err := newDataPoint(p, nil)
			if err != nil {
				return data, false, err
			}
			log.V(2).Infof("Constructed datapoint for delete: %v", dp)
			data = append(data, dp)
		}
		for _, u := range n.GetUpdate() {
			if u.Path == nil {
				return data, false, errors.Errorf("invalid nil path in update: %v", u)
			}
			if u.Val == nil {
				return data, false, errors.Errorf("invalid nil Val in update: %v", u)
			}
			log.V(2).Infof("Received gNMI Update value %s at path: %s", prototext.Format(u.Val), prototext.Format(u.Path))
			dp, err := newDataPoint(u.Path, u.Val)
			if err != nil {
				return data, false, err
			}
			log.V(2).Infof("Constructed datapoint for update: %v", dp)
			data = append(data, dp)
		}
		return data, false, nil
	case *gpb.SubscribeResponse_SyncResponse:
		log.V(2).Infof("Received gNMI SyncResponse.")
		data = append(data, &DataPoint{
			RecvTimestamp: recvTS,
			Sync:          true,
		})
		return data, true, nil
	default:
		return data, false, errors.Errorf("unexpected response: %v (%T)", v, v)
	}
}

// receiveAll receives data until the context deadline is reached, or when in
func receiveAll(sub gpb.GNMI_SubscribeClient, deletesExpected bool, mode gpb.SubscriptionList_Mode) (data []*DataPoint, err error) {
	for {
		var sync bool
		data, sync, err = receive(sub, data, deletesExpected)
		if err != nil {
			if mode == gpb.SubscriptionList_ONCE && err == io.EOF {
				// TODO(wenbli): It is unclear whether "subscribe ONCE stream closed without sync_response"
				// should be an error, so tolerate both scenarios.
				// See https://github.com/openconfig/reference/pull/156
				log.V(1).Infof("subscribe ONCE stream closed without sync_response.")
				break
			}
			// DeadlineExceeded is expected when collections are complete.
			if st, ok := status.FromError(err); ok && st.Code() == codes.DeadlineExceeded {
				break
			}
			return nil, errors.Wrapf(err, "error receiving gNMI response")
		}
		if mode == gpb.SubscriptionList_ONCE && sync {
			break
		}
	}
	return data, nil
}

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
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/openconfig/gnmi/errlist"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
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
		return nil, fmt.Errorf("gNMI failed to Subscribe: %w", err)
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
		return nil, fmt.Errorf("gNMI failed to Send(%+v): %w", sr, err)
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
			return data, false, fmt.Errorf("unexpected delete updates: %v", n.Delete)
		}
		ts := time.Unix(0, n.GetTimestamp())
		newDataPoint := func(p *gpb.Path, val *gpb.TypedValue) (*DataPoint, error) {
			j, err := util.JoinPaths(n.GetPrefix(), p)
			if err != nil {
				return nil, err
			}
			// Record the deprecated Element field for clearer compliance error messages.
			//nolint:staticcheck // ignore deprecated check
			if elements := append(append([]string{}, n.GetPrefix().GetElement()...), p.GetElement()...); len(elements) > 0 {
				//nolint:staticcheck // ignore deprecated check
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
				return data, false, fmt.Errorf("invalid nil path in update: %v", u)
			}
			if u.Val == nil {
				return data, false, fmt.Errorf("invalid nil Val in update: %v", u)
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
		return data, false, fmt.Errorf("unexpected response: %v (%T)", v, v)
	}
}

// receiveAll receives data until the context deadline is reached, or when in
// ONCE mode, a sync response is received.
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
			return nil, fmt.Errorf("error receiving gNMI response: %w", err)
		}
		if mode == gpb.SubscriptionList_ONCE && sync {
			break
		}
	}
	return data, nil
}

// receiveStream receives an async stream of gNMI notifications and sends them to a chan.
// Note: this does not imply that mode is gpb.SubscriptionList_STREAM (though it usually is).
// If the query is a leaf, each datapoint will be sent the chan individually.
// If the query is a non-leaf, all the datapoints from a SubscriptionResponse are bundled.
func receiveStream[T any](sub gpb.GNMI_SubscribeClient, query AnyQuery[T]) (<-chan []*DataPoint, <-chan error) {
	dataCh := make(chan []*DataPoint)
	errCh := make(chan error)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		var recvData []*DataPoint
		var hasSynced bool
		var sync bool
		var err error
		for {
			recvData, sync, err = receive(sub, recvData, true)
			if err != nil {
				errCh <- fmt.Errorf("error receiving gNMI response: %w", err)
				return
			}
			firstSync := !hasSynced && (sync || query.isLeaf())
			hasSynced = hasSynced || sync || query.isLeaf()
			// Skip conversion and predicate until first sync for non-leaves.
			if !hasSynced {
				continue
			}
			var datas [][]*DataPoint
			if query.isLeaf() {
				for _, datum := range recvData {
					// Add all datapoints except sync datapoints after the first sync.
					if (len(recvData) == 1 && firstSync) || !datum.Sync {
						datas = append(datas, []*DataPoint{datum})
					}
				}
			} else {
				datas = [][]*DataPoint{recvData}
			}
			for _, data := range datas {
				dataCh <- data
			}
			recvData = nil
		}
	}()
	return dataCh, errCh
}

// set configures the target at the query path.
func set[T any](ctx context.Context, c *Client, q ConfigQuery[T], val T, op setOperation) (*gpb.SetResponse, *gpb.Path, error) {
	path, _, errs := ygot.ResolvePath(q.pathStruct())
	if err := errsToErr(errs); err != nil {
		return nil, nil, err
	}

	req := &gpb.SetRequest{}
	var setVal interface{} = val
	if q.isLeaf() && q.isScalar() {
		setVal = &val
	}
	if err := populateSetRequest(req, path, setVal, op); err != nil {
		return nil, nil, err
	}

	req.Prefix = &gpb.Path{
		Target: c.target,
	}
	log.V(1).Info(prettySetRequest(req))
	resp, err := c.gnmiC.Set(ctx, req)
	log.V(1).Infof("SetResponse:\n%s", prototext.Format(resp))

	return resp, path, err
}

// setOperation is an enum representing the different kinds of SetRequest
// operations available.
type setOperation int

const (
	deletePath setOperation = iota
	replacePath
	updatePath
)

// populateSetRequest fills a SetResponse for a val and operation type.
func populateSetRequest(req *gpb.SetRequest, path *gpb.Path, val interface{}, op setOperation) error {
	if req == nil {
		return fmt.Errorf("cannot populate a nil SetRequest")
	}

	switch op {
	case deletePath:
		req.Delete = append(req.Delete, path)
	case replacePath, updatePath:
		// Since the GoStructs are generated using preferOperationalState, we
		// need to turn on preferShadowPath to prefer marshalling config paths.
		js, err := ygot.Marshal7951(val, ygot.JSONIndent("  "), &ygot.RFC7951JSONConfig{AppendModuleName: true, PreferShadowPath: true})
		if err != nil {
			return fmt.Errorf("could not encode value into JSON format: %w", err)
		}
		update := &gpb.Update{
			Path: path,
			Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: js}},
		}
		if op == replacePath {
			req.Replace = append(req.Replace, update)
		} else {
			req.Update = append(req.Update, update)
		}
	default:
		return fmt.Errorf("unknown set operation: %v", op)
	}

	return nil
}

// prettySetRequest returns a string version of a gNMI SetRequest for human
// consumption and ignores errors. Note that the output is subject to change.
// See documentation for prototext.Format.
func prettySetRequest(setRequest *gpb.SetRequest) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "SetRequest:\n%s\n", prototext.Format(setRequest))

	writePath := func(path *gpb.Path) {
		pathStr, err := ygot.PathToString(path)
		if err != nil {
			pathStr = prototext.Format(path)
		}
		fmt.Fprintf(&buf, "%s\n", pathStr)
	}

	writeVal := func(val *gpb.TypedValue) {
		switch v := val.Value.(type) {
		case *gpb.TypedValue_JsonIetfVal:
			fmt.Fprintf(&buf, "%s\n", v.JsonIetfVal)
		default:
			fmt.Fprintf(&buf, "%s\n", prototext.Format(val))
		}
	}

	for i, path := range setRequest.Delete {
		fmt.Fprintf(&buf, "-------delete path #%d------\n", i)
		writePath(path)
	}
	for i, update := range setRequest.Replace {
		fmt.Fprintf(&buf, "-------replace path/value pair #%d------\n", i)
		writePath(update.Path)
		writeVal(update.Val)
	}
	for i, update := range setRequest.Update {
		fmt.Fprintf(&buf, "-------update path/value pair #%d------\n", i)
		writePath(update.Path)
		writeVal(update.Val)
	}
	return buf.String()
}

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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/openconfig/ygnmi/internal/logutil"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	closer "github.com/openconfig/gocloser"
)

// subscribe create a gNMI SubscribeClient for the given query.
func subscribe[T any](ctx context.Context, c *Client, q AnyQuery[T], mode gpb.SubscriptionList_Mode, o *opt) (_ gpb.GNMI_SubscribeClient, rerr error) {
	var queryPaths []*gpb.Path
	var subs []*gpb.Subscription
	for _, path := range q.subPaths() {
		path, err := resolvePath(path)
		if err != nil {
			return nil, err
		}
		queryPaths = append(queryPaths, path)
	}
	if len(queryPaths) > 0 && o.ft != nil {
		if !q.isLeaf() {
			return nil, fmt.Errorf("functional translators only support leaf queries, given %+v", queryPaths)
		}
		if len(queryPaths) != 1 {
			// This should never happen because leaf queries should only have one path. (Batch queries with multiple leaf paths still don't have isLeaf == true)
			return nil, fmt.Errorf("functional translators only support one query path, given %d paths", len(queryPaths))
		}
		// Convert the query path to a schema path by stripping keys because output paths provided to FT
		// OutputToInput() for subscription translation must be schema paths. This is expected to cause
		// the subscription to return more data than just the requested key, but we'll filter the output
		// in receive() to only include paths matching the keys we actually queried.
		schemaPath := proto.Clone(queryPaths[0]).(*gpb.Path)
		for _, elem := range schemaPath.Elem {
			elem.Key = nil
		}
		match, inputs, err := o.ft.OutputToInput(schemaPath)
		if err != nil {
			log.ErrorContextf(ctx, "Received error from FunctionalTranslator.OutputToInput(): %v", err)
			return nil, err
		}
		if !match {
			return nil, fmt.Errorf("FunctionalTranslator.OutputToInput() did not match on path: %s", prototext.Format(schemaPath))
		}

		log.V(2).InfoContextf(ctx, "FunctionalTranslator.OutputToInput() mapped original query path %s to actual subscription paths: %+v", prototext.Format(queryPaths[0]), inputs)
		queryPaths = inputs
	}

	for _, path := range queryPaths {
		subs = append(subs, &gpb.Subscription{
			Path: &gpb.Path{
				Elem:   path.GetElem(),
				Origin: path.GetOrigin(),
			},
			Mode:           o.mode,
			SampleInterval: o.sampleInterval,
		})
	}

	sr := &gpb.SubscribeRequest{
		Request: &gpb.SubscribeRequest_Subscribe{
			Subscribe: &gpb.SubscriptionList{
				Prefix: &gpb.Path{
					Target: c.target,
				},
				Subscription: subs,
				Mode:         mode,
				Encoding:     o.encoding,
			},
		},
	}
	if o.useGet && mode != gpb.SubscriptionList_ONCE {
		return nil, fmt.Errorf("using gnmi.Get is only valid for ONCE subscriptions")
	}

	ctx = NewContext(ctx, q)

	var sub gpb.GNMI_SubscribeClient
	var err error
	if o.useGet {
		dt := gpb.GetRequest_CONFIG
		if q.IsState() {
			dt = gpb.GetRequest_STATE
		}
		sub = &getSubscriber{
			client:   c,
			ctx:      ctx,
			dataType: dt,
		}
	} else {
		sub, err = c.gnmiC.Subscribe(ctx)
		if err != nil {
			return nil, fmt.Errorf("gNMI failed to Subscribe: %w", err)
		}
	}
	defer closer.Close(&rerr, sub.CloseSend, "error closing gNMI send stream")
	if !o.useGet {
		log.V(c.requestLogLevel).InfoContext(ctx, prototext.Format(sr))
	}
	if err := sub.Send(sr); err != nil {
		// If the server closes the RPC with an error, the real error may only be visible on Recv.
		// https://pkg.go.dev/google.golang.org/grpc?utm_source=godoc#ClientStream
		if errors.Is(err, io.EOF) {
			if _, recvErr := sub.Recv(); recvErr != nil {
				err = recvErr
			}
		}
		return nil, fmt.Errorf("gNMI failed to Send(%+v): %w", sr, err)
	}

	return sub, nil
}

// getSubscriber is an implementation of gpb.GNMI_SubscribeClient that uses gpb.Get.
// Send() does the Get call and Recv returns the Get response.
type getSubscriber struct {
	gpb.GNMI_SubscribeClient
	client   *Client
	ctx      context.Context
	notifs   []*gpb.Notification
	dataType gpb.GetRequest_DataType
}

// Send call gnmi.Get with a request equivalent to the SubscribeRequest.
func (gs *getSubscriber) Send(req *gpb.SubscribeRequest) error {
	getReq := &gpb.GetRequest{
		Prefix:   req.GetSubscribe().GetPrefix(),
		Encoding: gpb.Encoding_JSON_IETF,
		Type:     gs.dataType,
	}
	for _, sub := range req.GetSubscribe().GetSubscription() {
		getReq.Path = append(getReq.Path, sub.GetPath())
	}
	log.V(gs.client.requestLogLevel).Info(prototext.Format(getReq))
	resp, err := gs.client.gnmiC.Get(gs.ctx, getReq)
	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound { // Make this behave like Subscribe, where non-existent paths don't return values.
		return nil
	} else if err != nil {
		return fmt.Errorf("gNMI.Get request failed: %w", err)
	}
	gs.notifs = resp.GetNotification()
	return nil
}

// Recv returns the result of the Get request, returning io.EOF after all responses are read.
func (gs *getSubscriber) Recv() (*gpb.SubscribeResponse, error) {
	if len(gs.notifs) == 0 {
		return nil, io.EOF
	}
	resp := &gpb.SubscribeResponse{
		Response: &gpb.SubscribeResponse_Update{
			Update: gs.notifs[0],
		},
	}
	gs.notifs = gs.notifs[1:]
	return resp, nil
}

// CloseSend is noop implementation gRPC subscribe interface.
func (gs *getSubscriber) CloseSend() error {
	return nil
}

// receive processes a single response from the subscription stream. If an "update" response is
// received, those points are appended to the given data and the result of that concatenation is
// the first return value, and the second return value is false. If a "sync" response is received,
// the data is returned as-is and the second return value is true. If Delete paths are present in
// the update, they are appended to the given data before the Update values. If deletesExpected
// is false, however, any deletes received will cause an error.
func receive(sub gpb.GNMI_SubscribeClient, data []*DataPoint, deletesExpected bool, queryPath *gpb.Path, o *opt) ([]*DataPoint, bool, error) {
	res, err := sub.Recv()
	if err != nil {
		return data, false, err
	}
	recvTS := time.Now()

	switch v := res.Response.(type) {
	case *gpb.SubscribeResponse_Update:
		if o.ft != nil {
			out, err := o.ft.Translate(res)
			if err != nil {
				log.Errorf("FunctionalTranslator.Translate() failed to translate notification: %v", err)
				return data, false, nil
			}
			if out == nil {
				log.V(2).Infof("Received nil response from functional translatator with input: %s", prototext.Format(res))
				return data, false, nil
			}
			log.V(2).Infof("FT successfully translated a notification. input: %s, output: %s", prototext.Format(res), prototext.Format(out))
			res = out
		}
		n := res.GetUpdate()
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
			// Filter out paths that don't match the query here as a workaround for the edge case where we
			// query a path including a specific key and the FT subscribes to more data than just that.
			// This extra data would otherwise be filtered out downstream but with a compliance error.
			// This uses the same logic that unmarshal(...) in unmarshal.go uses to check compliance.
			if o.ft != nil && !util.PathMatchesQuery(dp.Path, queryPath) {
				log.V(2).Infof("Skipping delete datapoint that doesn't match the query. path: %s, query: %s", prototext.Format(dp.Path), prototext.Format(queryPath))
				continue
			}
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
			if o.ft != nil && !util.PathMatchesQuery(dp.Path, queryPath) {
				log.V(2).Infof("Skipping update datapoint that doesn't match the query. path: %s, query: %s", prototext.Format(dp.Path), prototext.Format(queryPath))
				continue
			}
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

// receiveAll receives data until the context deadline is reached, or when a sync response is received.
// This func is only used when receiving data from a ONCE subscription.
func receiveAll[T any](sub gpb.GNMI_SubscribeClient, deletesExpected bool, query AnyQuery[T], o *opt) (data []*DataPoint, err error) {
	queryPath, err := resolvePath(query.PathStruct())
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}
	for {
		var sync bool
		data, sync, err = receive(sub, data, deletesExpected, queryPath, o)
		if err != nil {
			if err == io.EOF {
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
		if sync {
			break
		}
	}
	return data, nil
}

// receiveStream receives an async stream of gNMI notifications and sends them to a chan.
// Note: this does not imply that mode is gpb.SubscriptionList_STREAM (though it usually is).
// If the query is a leaf, each datapoint will be sent the chan individually.
// If the query is a non-leaf, all the datapoints from a SubscriptionResponse are bundled.
func receiveStream[T any](ctx context.Context, sub gpb.GNMI_SubscribeClient, query AnyQuery[T], o *opt) (<-chan []*DataPoint, <-chan error) {
	dataCh := make(chan []*DataPoint)
	errCh := make(chan error)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		var recvData []*DataPoint
		var hasSynced bool
		var sync bool
		var err error

		queryPath, err := resolvePath(query.PathStruct())
		if err != nil {
			errCh <- fmt.Errorf("failed to resolve path: %w", err)
			return
		}
		for {
			recvData, sync, err = receive(sub, recvData, true, queryPath, o)
			if err != nil {
				// In the case that the context is cancelled, the reader of errCh
				// may have gone away. In order to avoid this goroutine blocking
				// indefinitely on the channel write, allow ctx being cancelled or
				// done to ensure that it returns.
				select {
				case errCh <- fmt.Errorf("error receiving gNMI response: %w", err):
				case <-ctx.Done():
				}
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

// removeRedundantModPrefix7951 removes any modName prefixes in the first layer
// of the input JSON v.
//
// https://www.rfc-editor.org/rfc/rfc7951#section-4:
//
//	A namespace-qualified member name MUST be used for all members of a
//	top-level JSON object and then also whenever the namespaces of the
//	data node and its parent node are different.  In all other cases, the
//	simple form of the member name MUST be used.
func removeRedundantModPrefix7951(v any, modName string) any {
	switch v := v.(type) {
	case map[string]any:
		newMap := map[string]any{}
		for name, v := range v {
			ps := strings.Split(name, ":")
			if len(ps) != 2 {
				return fmt.Errorf("ygnmi error: when wrapping RFC7951 JSON, got invalid JSON with first-level not qualified: %q", name)
			}
			if ps[0] == modName {
				newMap[ps[1]] = v
			} else {
				newMap[name] = v
			}
		}
		return newMap
	case []any:
		var newSlice []any
		if v != nil {
			newSlice = []any{}
		}
		for _, v := range v {
			newSlice = append(newSlice, removeRedundantModPrefix7951(v, modName))
		}
		return newSlice
	default:
		return v
	}
}

// wrapJSONIETFSinglePath wraps a single layer of JSON-IETF for the given
// JSON input.
func wrapJSONIETFSinglePath(v any, modName, pathName string) any {
	if v == nil {
		return map[string]any{}
	}
	return map[string]any{fmt.Sprintf("%s:%s", modName, pathName): removeRedundantModPrefix7951(v, modName)}
}

// wrapJSONIETF returns the input RFC7951 JSON wrapped in another layer.
//
// If the input TypedValue is not JSON_IETF, then an error is returned.
func wrapJSONIETF(tv *gpb.TypedValue, qualifiedRelPath []string) error {
	bs := tv.GetJsonIetfVal()
	if len(bs) == 0 {
		return fmt.Errorf("TypedValue is not JSON IETF:\n%s", prototext.Format(tv))
	}
	var jv any
	if err := json.Unmarshal(bs, &jv); err != nil {
		return fmt.Errorf("cannot unmarshal IETF JSON: %v", err)
	}
	for i := len(qualifiedRelPath) - 1; i >= 0; i-- {
		qualPathEle := qualifiedRelPath[i]
		ps := strings.Split(qualPathEle, ":")
		if len(ps) != 2 {
			return fmt.Errorf("ygnmi error: got unexpected qualified YANG name: %q", qualPathEle)
		}
		jv = wrapJSONIETFSinglePath(jv, ps[0], ps[1])
	}
	bs, err := json.Marshal(jv)
	if err != nil {
		return fmt.Errorf("ygnmi internal error: cannot marshal wrapped JSON: %v", err)
	}
	tv.Value = &gpb.TypedValue_JsonIetfVal{JsonIetfVal: bs}
	return nil
}

// set configures the target at the query path.
func set[T any](ctx context.Context, c *Client, q ConfigQuery[T], val T, op setOperation, opts ...Option) (*gpb.SetResponse, *gpb.Path, error) {
	path, err := resolvePath(q.PathStruct())
	if err != nil {
		return nil, nil, err
	}

	req := &gpb.SetRequest{}
	var setVal interface{} = val
	if q.isLeaf() && q.isScalar() {
		setVal = &val
	}
	if err := populateSetRequest(req, path, setVal, op, q.isShadowPath(), q.isLeaf(), q.compressInfo(), opts...); err != nil {
		return nil, nil, err
	}

	req.Prefix = &gpb.Path{
		Target: c.target,
	}
	logutil.LogByLine(c.requestLogLevel, prettySetRequest(req))
	resp, err := c.gnmiC.Set(ctx, req)
	log.V(c.requestLogLevel).Infof("SetResponse:\n%s", prototext.Format(resp))

	return resp, path, err
}

// setOperation is an enum representing the different kinds of SetRequest
// operations available.
type setOperation int

const (
	deletePath setOperation = iota
	replacePath
	updatePath
	unionreplacePath
)

// populateSetRequest fills a SetResponse for a val and operation type.
func populateSetRequest(req *gpb.SetRequest, path *gpb.Path, val interface{}, op setOperation, preferShadowPath, isLeaf bool, compressInfo *CompressionInfo, opts ...Option) error {
	if req == nil {
		return fmt.Errorf("cannot populate a nil SetRequest")
	}
	opt := resolveOpts(opts)

	switch op {
	case deletePath:
		req.Delete = append(req.Delete, path)
	case replacePath, updatePath, unionreplacePath:
		var typedVal *gpb.TypedValue
		var err error
		if s, ok := val.(*string); ok && path.Origin == "cli" {
			typedVal = &gpb.TypedValue{Value: &gpb.TypedValue_AsciiVal{AsciiVal: *s}}
		} else if s, ok := val.(string); ok && strings.HasSuffix(path.Origin, "_cli") {
			typedVal = &gpb.TypedValue{Value: &gpb.TypedValue_AsciiVal{AsciiVal: s}}
		} else if opt.preferProto {
			typedVal, err = ygot.EncodeTypedValue(val, gpb.Encoding_JSON_IETF, &ygot.RFC7951JSONConfig{AppendModuleName: opt.appendModuleName, PreferShadowPath: preferShadowPath})
		} else {
			// Since the GoStructs are generated using preferOperationalState, we
			// need to turn on preferShadowPath to prefer marshalling config paths.
			var b []byte
			b, err = ygot.Marshal7951(val, ygot.JSONIndent("  "), &ygot.RFC7951JSONConfig{AppendModuleName: opt.appendModuleName, PreferShadowPath: preferShadowPath})

			// Respect the encoding option.
			switch opt.encoding {
			case gpb.Encoding_JSON:
				typedVal = &gpb.TypedValue{Value: &gpb.TypedValue_JsonVal{JsonVal: b}}
			case gpb.Encoding_JSON_IETF:
				fallthrough
			default:
				typedVal = &gpb.TypedValue{Value: &gpb.TypedValue_JsonIetfVal{JsonIetfVal: b}}
			}
		}

		if err != nil && opt.setFallback && path.Origin != "openconfig" {
			if m, ok := val.(proto.Message); ok {
				any, err := anypb.New(m)
				if err != nil {
					return fmt.Errorf("failed to marshal proto: %v", err)
				}
				typedVal = &gpb.TypedValue{Value: &gpb.TypedValue_AnyVal{AnyVal: any}}
			} else {
				b, err := json.Marshal(val)
				if err != nil {
					return fmt.Errorf("failed to marshal json: %v", err)
				}
				typedVal = &gpb.TypedValue{Value: &gpb.TypedValue_JsonVal{JsonVal: b}}
			}
		} else if err != nil {
			return fmt.Errorf("failed to encode set request: %v", err)
		}

		var modifyTypedValueFn func(*gpb.TypedValue) error
		if !isLeaf && compressInfo != nil && len(compressInfo.PostRelPath) > 0 && len(typedVal.GetJsonIetfVal()) > 0 {
			// When the path struct points to a node that's compressed out,
			// then we know that the type is a node lower than it should be
			// as far as the JSON is concerned.
			modifyTypedValueFn = func(tv *gpb.TypedValue) error { return wrapJSONIETF(tv, compressInfo.PostRelPath) }
		}
		if modifyTypedValueFn != nil {
			if err := modifyTypedValueFn(typedVal); err != nil {
				return fmt.Errorf("failed to modify TypedValue: %v", err)
			}
		}
		update := &gpb.Update{
			Path: path,
			Val:  typedVal,
		}

		switch op {
		case replacePath:
			req.Replace = append(req.Replace, update)
		case updatePath:
			req.Update = append(req.Update, update)
		case unionreplacePath:
			req.UnionReplace = append(req.UnionReplace, update)
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

const (
	// OriginOverride is the key to custom opt that sets the path origin.
	OriginOverride = "origin-override"
)

func resolvePath(q PathStruct) (*gpb.Path, error) {
	path, opts, err := ResolvePath(q)
	if err != nil {
		return nil, err
	}
	if originSetter, ok := q.(interface{ PathOriginName() string }); ok {
		// When the path struct has the method PathOriginName(),
		// then the output of the method is set to the path.Origin
		path.Origin = originSetter.PathOriginName()
	}
	if origin, ok := opts[OriginOverride]; ok {
		path.Origin = origin.(string)
	}

	// TODO: remove when fixed https://github.com/openconfig/ygot/issues/615
	if path.Origin == "" && (len(path.Elem) == 0 || path.Elem[0].Name != "meta") {
		path.Origin = "openconfig"
	}

	return path, nil
}

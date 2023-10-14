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
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/openconfig/gnmi/errlist"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// DataPoint is a value of a gNMI path at a particular time.
type DataPoint struct {
	// Path of the received value.
	Path *gpb.Path
	// Value of the data; nil means delete.
	Value *gpb.TypedValue
	// Timestamp is the time at which the value was updated on the device.
	Timestamp time.Time
	// RecvTimestamp is the time the update was received.
	RecvTimestamp time.Time
	// Sync indicates whether the received datapoint was gNMI sync response.
	Sync bool
}

func (d *DataPoint) String() string {
	if d == nil {
		return ""
	}
	if d.Sync {
		return "gNMI SyncResponse"
	}
	path, err := ygot.PathToString(d.Path)
	if err != nil {
		path = prototext.Format(d.Path)
	}
	valStr := fmt.Sprint(d.Value)
	if jsonietf := d.Value.GetJsonIetfVal(); len(jsonietf) > 0 {
		valStr = string(jsonietf)
	}
	return fmt.Sprintf("%s (timestamp: %v, recvTimestamp: %v, isSync: %v): %s", path, d.Timestamp, d.RecvTimestamp, d.Sync, valStr)
}

// TelemetryError stores the path, value, and error string from unsuccessfully
// unmarshalling a datapoint into a YANG schema.
type TelemetryError struct {
	Path  *gpb.Path
	Value *gpb.TypedValue
	Err   error
}

func (t *TelemetryError) String() string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("Unmarshal %v into %v: %s", t.Value, t.Path, t.Err.Error())
}

// ComplianceErrors contains the compliance errors encountered from an Unmarshal operation.
type ComplianceErrors struct {
	// PathErrors are compliance errors encountered due to an invalid schema path.
	PathErrors []*TelemetryError
	// TypeErrors are compliance errors encountered due to an invalid type.
	TypeErrors []*TelemetryError
	// ValidateErrors are compliance errors encountered while doing schema
	// validation on the unmarshalled data.
	ValidateErrors []error
}

func (c *ComplianceErrors) String() string {
	if c == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString("Noncompliance Errors by category (see https://github.com/openconfig/ygnmi#noncompliance-errors):")
	b.WriteString("\nPath Noncompliance Errors:")
	if len(c.PathErrors) != 0 {
		for _, e := range c.PathErrors {
			b.WriteString("\n\t")
			b.WriteString(e.String())
		}
	} else {
		b.WriteString(" None")
	}
	b.WriteString("\nType Noncompliance Errors:")
	if len(c.TypeErrors) != 0 {
		for _, e := range c.TypeErrors {
			b.WriteString("\n\t")
			b.WriteString(e.String())
		}
	} else {
		b.WriteString(" None")
	}
	b.WriteString("\nValue Restriction Noncompliance Errors:")
	if len(c.ValidateErrors) != 0 {
		for _, e := range c.ValidateErrors {
			b.WriteString("\n\t")
			b.WriteString(e.Error())
		}
	} else {
		b.WriteString(" None")
	}
	b.WriteString("\n")
	return b.String()
}

// unmarshalAndExtract is a wrapper to ytypes.SetNode() and ygot.Validate() that
// unmarshals a given []*DataPoint to its field given an input query and
// verifies that all data conform to the schema. Any errors due to
// the unmarshal operations above are returned in a *ComplianceErrors for the
// caller to choose whether to tolerate, while other errors are returned directly.
// NOTE: The datapoints are applied in order as they are in the input slice,
// *NOT* in order of their timestamps. As such, in order to correctly support
// Collect calls, the input data must be sorted in order of timestamps.
func unmarshalAndExtract[T any](data []*DataPoint, q AnyQuery[T], goStruct ygot.ValidatedGoStruct, opts *opt) (*Value[T], error) {
	queryPath, err := resolvePath(q.PathStruct())
	if err != nil {
		return nil, err
	}
	ret := &Value[T]{
		Path: queryPath,
	}
	if len(data) == 0 {
		return ret, nil
	}
	schema := q.schema()

	if schema == nil { // Handle schema-less queries.
		val, _ := q.extract(nil)
		var setVal any = val
		if q.isScalar() {
			setVal = &val
		}

		delete, err := unmarshalSchemaless(data, setVal)
		if err != nil {
			return ret, err
		}
		ret.Timestamp = data[0].Timestamp
		ret.RecvTimestamp = data[0].RecvTimestamp
		ret.Path = proto.Clone(data[0].Path).(*gpb.Path)
		if !delete {
			ret.SetVal(val)
		}
		return ret, nil
	}

	unmarshalledData, complianceErrs, err := unmarshal(data, schema.SchemaTree[q.dirName()], goStruct, queryPath, schema, q.isLeaf(), !q.IsState(), q.compressInfo(), opts)
	ret.ComplianceErrors = complianceErrs
	if err != nil {
		return ret, err
	}
	if len(unmarshalledData) == 0 {
		return ret, nil
	}

	path := unmarshalledData[0].Path
	if !q.isLeaf() {
		// Apply a mask to the returned path according to which
		// elements were specified in the query.
		path = proto.Clone(unmarshalledData[0].Path).(*gpb.Path)
		path.Elem = path.Elem[:len(queryPath.Elem)]
		path.Origin = queryPath.Origin
		ret.Timestamp = LatestTimestamp(unmarshalledData)
		ret.RecvTimestamp = LatestRecvTimestamp(unmarshalledData)
		if q.isListContainer() {
			path.Elem[len(path.Elem)-1].Key = nil
		}
	} else {
		ret.Timestamp = unmarshalledData[0].Timestamp
		ret.RecvTimestamp = unmarshalledData[0].RecvTimestamp
	}
	ret.Path = path

	for _, data := range unmarshalledData {
		ret.ChangedPaths = append(ret.ChangedPaths, data.Path)
	}

	// For non-leaf config queries, prune all state-only leaves.
	// Note that config/state separation only applies to compressed structs.
	if q.isCompressedSchema() && !q.isLeaf() && !q.IsState() {
		err := ygot.PruneConfigFalse(q.schema().SchemaTree[q.dirName()], goStruct)
		if err != nil {
			return ret, err
		}
	}
	if val, ok := q.extract(goStruct); ok {
		ret.SetVal(val)
	}
	return ret, nil
}

// unmarshalSchemaless unmarshals the datapoint into the value, returning whether the datapoint was a delete.
func unmarshalSchemaless(data []*DataPoint, val any) (bool, error) {
	switch {
	case len(data) > 2:
		return false, fmt.Errorf("got multiple datapoints for schemaless node")
	case len(data) == 2 && !data[1].Sync:
		return false, fmt.Errorf("got multiple datapoints for schemaless node")
	}
	if data[0].Value == nil {
		return true, nil
	}

	rVal := reflect.ValueOf(val).Elem()
	valType := reflect.TypeOf(val).Elem()
	kind := valType.Kind()
	if !rVal.CanSet() {
		return false, fmt.Errorf("value not settable")
	}

	switch dataVal := data[0].Value.Value.(type) {
	case *gpb.TypedValue_StringVal:
		if kind != reflect.String {
			return false, fmt.Errorf("unmarshal err notification type %T, generic type %T", dataVal, val)
		}
		rVal.SetString(dataVal.StringVal)
	case *gpb.TypedValue_AsciiVal:
		if kind != reflect.String {
			return false, fmt.Errorf("unmarshal err notification type %T, generic type %T", dataVal, val)
		}
		rVal.SetString(dataVal.AsciiVal)
	case *gpb.TypedValue_IntVal:
		if kind != reflect.Int && kind != reflect.Int64 {
			return false, fmt.Errorf("unmarshal err notification type %T, generic type %T", dataVal, val)
		}
		rVal.SetInt(dataVal.IntVal)
	case *gpb.TypedValue_UintVal:
		if kind != reflect.Uint && kind != reflect.Uint64 {
			return false, fmt.Errorf("unmarshal err notification type %T, generic type %T", dataVal, val)
		}
		rVal.SetUint(dataVal.UintVal)
	case *gpb.TypedValue_BoolVal:
		if kind != reflect.Bool {
			return false, fmt.Errorf("unmarshal err notification type %T, generic type %T", dataVal, val)
		}
		rVal.SetBool(dataVal.BoolVal)
	case *gpb.TypedValue_DoubleVal:
		if kind != reflect.Float64 {
			return false, fmt.Errorf("unmarshal err notification type %T, generic type %T", dataVal, val)
		}
		rVal.SetFloat(dataVal.DoubleVal)
	case *gpb.TypedValue_LeaflistVal:
		return false, fmt.Errorf("leaf lists not supported")
	case *gpb.TypedValue_AnyVal:
		msg, ok := val.(proto.Message)
		if !ok {
			return false, fmt.Errorf("unmarshal err notification %T parameter %s", val, kind.String())
		}
		return false, dataVal.AnyVal.UnmarshalTo(msg)
	case *gpb.TypedValue_ProtoBytes:
		msg, ok := val.(proto.Message)
		if !ok {
			return false, fmt.Errorf("unmarshal err notification %T parameter %s", val, kind.String())
		}
		return false, proto.Unmarshal(dataVal.ProtoBytes, msg)
	case *gpb.TypedValue_JsonVal:
		return false, json.Unmarshal(dataVal.JsonVal, val)
	case *gpb.TypedValue_JsonIetfVal:
		return false, json.Unmarshal(dataVal.JsonIetfVal, val)
	default:
		return false, fmt.Errorf("unsupported type: %T", dataVal)
	}
	return false, nil
}

// unmarshal unmarshals a given slice of datapoints to its field given a
// containing GoStruct and its schema and verifies that all data conform to the
// schema. The subset of datapoints that successfully unmarshalled into the given GoStruct is returned.
// NOTE: The subset of datapoints includes datapoints that are value restriction noncompliant.
// The second error slice are internal errors, while the returned
// *ComplianceError stores the compliance errors.
func unmarshal(data []*DataPoint, structSchema *yang.Entry, structPtr ygot.ValidatedGoStruct, queryPath *gpb.Path, schema *ytypes.Schema, isLeaf, isConfig bool, compressInfo *CompressionInfo, opts *opt) ([]*DataPoint, *ComplianceErrors, error) {
	queryPathStr := pathToString(queryPath)
	if isLeaf {
		switch {
		case len(data) > 2:
			return nil, &ComplianceErrors{PathErrors: []*TelemetryError{{
				Err: fmt.Errorf("got multiple (%d) data points for leaf node at path %s: %v", len(data), queryPathStr, data),
			}}}, nil
		case len(data) == 2 && !data[1].Sync:
			return nil, &ComplianceErrors{PathErrors: []*TelemetryError{{
				Err: fmt.Errorf("got multiple (%d) data points for leaf node at path %s: %v", len(data), queryPathStr, data),
			}}}, nil
		}
	}

	var unmarshalledDatapoints []*DataPoint
	var pathUnmarshalErrs []*TelemetryError
	var typeUnmarshalErrs []*TelemetryError

	errs := &errlist.List{}
	if !schema.IsValid() {
		errs.Add(fmt.Errorf("input schema for generated code is invalid"))
		return nil, nil, errs.Err()
	}
	// TODO(wenbli): Add fatal check for duplicate paths, as they're not allowed by GET semantics.
	for _, dp := range data {
		var gcopts []ytypes.GetOrCreateNodeOpt
		if isConfig {
			gcopts = append(gcopts, &ytypes.PreferShadowPath{})
		}
		// Sync datapoints don't contain path nor values.
		if dp.Sync {
			continue
		}

		dpPathStr := pathToString(dp.Path)
		// 1a. Check for path compliance by doing a prefix-match, since
		// the given datapoint must be a descendant of the query.
		if !util.PathMatchesQuery(dp.Path, queryPath) {
			var pathErr error
			switch {
			//nolint:staticcheck // ignore deprecated check
			case len(dp.Path.Elem) == 0 && len(dp.Path.Element) > 0:
				pathErr = fmt.Errorf("datapoint path uses deprecated and unsupported Element field: %s", prototext.Format(dp.Path))
			default:
				pathErr = fmt.Errorf("datapoint path %q (value %v) does not match the query path %q", dpPathStr, dp.Value, queryPathStr)
			}
			pathUnmarshalErrs = append(pathUnmarshalErrs, &TelemetryError{
				Path:  dp.Path,
				Value: dp.Value,
				Err:   pathErr,
			})
			continue
		}

		unmarshalPath := dp.Path
		if compressInfo != nil && len(compressInfo.PreRelPath) > 0 && len(dp.Value.GetJsonIetfVal()) > 0 {
			// When the path points to a node that's compressed out, then we must make the
			// unmarshal target the parent since JSON unmarshal currently doesn't support
			// unmarshalling at a compressed-out node.
			// Skip wrapping for other encoding formats because scalars, the other format,
			// doesn't need wrapping.
			if err := wrapJSONIETF(dp.Value, compressInfo.PreRelPath); err != nil {
				errs.Add(fmt.Errorf("failed to wrap received JSON: %v", err))
			} else {
				unmarshalPath = proto.Clone(dp.Path).(*gpb.Path)
				unmarshalPath.Elem = unmarshalPath.Elem[:len(unmarshalPath.Elem)-len(compressInfo.PreRelPath)]
			}
		}

		// 1b. Check for path compliance: by unmarshalling from the
		// root, we check that the path, including the list key,
		// corresponds to an actual schema element.
		if _, _, err := ytypes.GetOrCreateNode(schema.RootSchema(), schema.Root, unmarshalPath, gcopts...); err != nil {
			pathUnmarshalErrs = append(pathUnmarshalErrs, &TelemetryError{Path: dp.Path, Value: dp.Value, Err: fmt.Errorf("path %q is invalid and cannot be matched to a generated GoStruct field: %v", dpPathStr, err)})
			continue
		}
		// The structSchema passed in here is assumed to be the unzipped
		// version from the generated structs file. That schema has a single
		// root entry that all top-level schemas are connected to via their
		// parent pointers. Therefore, we must remove that first element to
		// obtain the sanitized path.

		relPath := util.TrimGNMIPathPrefix(unmarshalPath, util.PathStringToElements(structSchema.Path())[1:])
		if dp.Value == nil {
			var dopts []ytypes.DelNodeOpt
			if isConfig {
				dopts = append(dopts, &ytypes.PreferShadowPath{})
			}
			if err := ytypes.DeleteNode(structSchema, structPtr, relPath, dopts...); err == nil {
				unmarshalledDatapoints = append(unmarshalledDatapoints, dp)
			} else {
				errs.Add(fmt.Errorf("path %q cannot be deleted: %v", dpPathStr, err))
			}
		} else {
			sopts := []ytypes.SetNodeOpt{&ytypes.InitMissingElements{}, &ytypes.TolerateJSONInconsistencies{}}
			if isConfig {
				sopts = append(sopts, &ytypes.PreferShadowPath{})
			}
			// TODO: Fully support partial unmarshaling of JSON blobs.
			if opts.useGet {
				sopts = append(sopts, &ytypes.IgnoreExtraFields{})
			}
			// 2. Check for type compliance (since path should already be compliant).

			if err := ytypes.SetNode(structSchema, structPtr, relPath, dp.Value, sopts...); err == nil {
				unmarshalledDatapoints = append(unmarshalledDatapoints, dp)
			} else {
				typeUnmarshalErrs = append(typeUnmarshalErrs, &TelemetryError{Path: dp.Path, Value: dp.Value, Err: fmt.Errorf("datapoint path %q (value %v) cannot be unmarshalled: %v", dpPathStr, dp.Value, err)})
			}
		}
	}
	// 3. Check for value (restriction) compliance.
	validateErrs := ytypes.Validate(structSchema, structPtr)
	if pathUnmarshalErrs != nil || typeUnmarshalErrs != nil || validateErrs != nil {
		return unmarshalledDatapoints, &ComplianceErrors{PathErrors: pathUnmarshalErrs, TypeErrors: typeUnmarshalErrs, ValidateErrors: validateErrs}, errs.Err()
	}
	return unmarshalledDatapoints, nil, errs.Err()
}

// LatestTimestamp returns the latest timestamp of the input datapoints.
// If datapoints is empty, then the zero time is returned.
func LatestTimestamp(data []*DataPoint) time.Time {
	var latest time.Time
	for _, dp := range data {
		if ts := dp.Timestamp; ts.After(latest) {
			latest = ts
		}
	}
	return latest
}

// LatestRecvTimestamp returns the latest recv timestamp of the input datapoints.
// If datapoints is empty, then the zero time is returned.
func LatestRecvTimestamp(data []*DataPoint) time.Time {
	var latest time.Time
	for _, dp := range data {
		if ts := dp.RecvTimestamp; ts.After(latest) {
			latest = ts
		}
	}
	return latest
}

// pathToString returns a string version of the input path for display during
// debugging.
func pathToString(path *gpb.Path) string {
	pathStr, err := ygot.PathToString(path)
	if err != nil {
		// Use Sprint instead of prototext.Format to avoid newlines.
		pathStr = fmt.Sprint(path)
	}
	return pathStr
}

func bundleDatapoints(datapoints []*DataPoint, prefixLen int) (map[string][]*DataPoint, []string, error) {
	groups := map[string][]*DataPoint{}

	for _, dp := range datapoints {
		if dp.Sync { // Sync datapoints don't have a path, so ignore them.
			continue
		}
		elems := dp.Path.GetElem()
		if len(elems) < prefixLen {
			groups["/"] = append(groups["/"], dp)
			continue
		}
		prefixPath, err := ygot.PathToString(&gpb.Path{Elem: elems[:prefixLen]})
		if err != nil {
			return nil, nil, err
		}
		groups[prefixPath] = append(groups[prefixPath], dp)
	}

	var prefixes []string
	for prefix := range groups {
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(prefixes)

	return groups, prefixes, nil
}

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
	"fmt"
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
	b.WriteString("Noncompliance Errors by category:")
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

// unmarshalToValue is a wrapper to ytypes.SetNode() and ygot.Validate() that
// unmarshals a given []*DataPoint to its field given an input query and
// verifies that all data conform to the schema. Any errors due to
// the unmarshal operations above are returned in a *ComplianceErrors for the
// caller to choose whether to tolerate, while other errors are returned directly.
// NOTE: The datapoints are applied in order as they are in the input slice,
// *NOT* in order of their timestamps. As such, in order to correctly support
// Collect calls, the input data must be sorted in order of timestamps.
//lint:ignore U1000 TODO(DanG100) remove this once this func is used
//nolint:deadcode // see above
func unmarshalToValue[T any](data []*DataPoint, q AnyQuery[T], goStruct ygot.ValidatedGoStruct) (*Value[T], error) {
	queryPath, _, errs := ygot.ResolvePath(q.pathStruct())
	if len(errs) > 0 {
		l := &errlist.List{}
		l.Add(errs...)
		return nil, l.Err()
	}

	ret := &Value[T]{
		Path: queryPath,
	}
	if len(data) == 0 {
		return ret, nil
	}

	unmarshalledData, complianceErrs, err := unmarshal(data, q.schema().SchemaTree[q.fieldName()], goStruct, queryPath, q.schema(), q.isLeaf(), q.reverseShadowPaths())
	ret.ComplianceErrors = complianceErrs
	if err != nil {
		return ret, err
	}
	if len(unmarshalledData) == 0 {
		return ret, nil
	}

	path := unmarshalledData[0].Path
	if !q.isLeaf() {
		path = proto.Clone(unmarshalledData[0].Path).(*gpb.Path)
		path.Elem = path.Elem[:len(queryPath.Elem)]
		ret.Timestamp = LatestTimestamp(unmarshalledData)
		ret.RecvTimestamp = LatestRecvTimestamp(unmarshalledData)
	} else {
		ret.Timestamp = unmarshalledData[0].Timestamp
		ret.RecvTimestamp = unmarshalledData[0].RecvTimestamp
	}
	ret.Path = path

	ret.SetVal(q.extract(goStruct))
	return ret, nil
}

// unmarshal unmarshals a given slice of datapoints to its field given a
// containing GoStruct and its schema and verifies that all data conform to the
// schema. The subset of datapoints that successfully unmarshalled into the given GoStruct is returned.
// NOTE: The subset of datapoints includes datapoints that are value restriction noncompliant.
// The second error slice are internal errors, while the returned
// *ComplianceError stores the compliance errors.
func unmarshal(data []*DataPoint, structSchema *yang.Entry, structPtr ygot.ValidatedGoStruct, queryPath *gpb.Path, schema *ytypes.Schema, isLeaf, reverseShadowPaths bool) ([]*DataPoint, *ComplianceErrors, error) {
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
		if reverseShadowPaths {
			gcopts = append(gcopts, &ytypes.PreferShadowPath{})
		}
		// Sync datapoints don't contain path nor values.
		if dp.Sync {
			continue
		}

		// 1a. Check for path compliance by doing a prefix-match, since
		// the given datapoint must be a descendant of the query.
		if !util.PathMatchesQuery(dp.Path, queryPath) {
			var pathErr error
			dpPathStr := pathToString(dp.Path)
			switch {
			//lint:ignore SA1019 ignore deprecated check
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
		// 1b. Check for path compliance: by unmarshalling from the
		// root, we check that the path, including the list key,
		// corresponds to an actual schema element.
		if _, _, err := ytypes.GetOrCreateNode(schema.RootSchema(), schema.Root, dp.Path, gcopts...); err != nil {
			pathUnmarshalErrs = append(pathUnmarshalErrs, &TelemetryError{Path: dp.Path, Value: dp.Value, Err: err})
			continue
		}
		// The structSchema passed in here is assumed to be the unzipped
		// version from the generated structs file. That schema has a single
		// root entry that all top-level schemas are connected to via their
		// parent pointers. Therefore, we must remove that first element to
		// obtain the sanitized path.
		relPath := util.TrimGNMIPathPrefix(dp.Path, util.PathStringToElements(structSchema.Path())[1:])
		if dp.Value == nil {
			var dopts []ytypes.DelNodeOpt
			if reverseShadowPaths {
				dopts = append(dopts, &ytypes.PreferShadowPath{})
			}
			if err := ytypes.DeleteNode(structSchema, structPtr, relPath, dopts...); err == nil {
				unmarshalledDatapoints = append(unmarshalledDatapoints, dp)
			} else {
				errs.Add(err)
			}
		} else {
			sopts := []ytypes.SetNodeOpt{&ytypes.InitMissingElements{}, &ytypes.TolerateJSONInconsistencies{}}
			if reverseShadowPaths {
				sopts = append(sopts, &ytypes.PreferShadowPath{})
			}
			// 2. Check for type compliance (since path should already be compliant).
			if err := ytypes.SetNode(structSchema, structPtr, relPath, dp.Value, sopts...); err == nil {
				unmarshalledDatapoints = append(unmarshalledDatapoints, dp)
			} else {
				typeUnmarshalErrs = append(typeUnmarshalErrs, &TelemetryError{Path: dp.Path, Value: dp.Value, Err: err})
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

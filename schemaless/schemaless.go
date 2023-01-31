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

// Package schemaless allows the creation of schema-less queries.
// Schema-less queries are not associated with a YANG schema,
// the user is free to choose any serializable type for any particular path,
// although the selection must be consistent to avoid runtime type mismatch.
// These queries have limited functionality compared to standard queries:
// unmarshaling only works if the gNMI server returns the value (or a list entry) in a single gpb.Update,
// this is the standard behavior for leaves. For non-leaves, this can be resolved by requesting JSON encoding
// (if supported by the server).
package schemaless

import (
	"fmt"
	"reflect"

	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// NewConfig creates a config query for the given path and type. The path must be gNMI path.
// See package comment for limitations of this query type.
func NewConfig[T any](path, origin string) (ygnmi.ConfigQuery[T], error) {
	ps, createFn, scalar, err := newQueryField[T](path, origin)
	if err != nil {
		return nil, err
	}

	return ygnmi.NewLeafConfigQuery("",
			false,
			scalar,
			ps,
			createFn,
			func() ygot.ValidatedGoStruct {
				return nil
			},
			func() *ytypes.Schema { return nil }),
		nil
}

// NewWildcard creates a wildcard query for the given path and type. The path must be gNMI path.
// See package comment for limitations of this query type.
func NewWildcard[T any](path, origin string) (ygnmi.WildcardQuery[T], error) {
	ps, createFn, scalar, err := newQueryField[T](path, origin)
	if err != nil {
		return nil, err
	}

	return ygnmi.NewLeafWildcardQuery("",
			false,
			scalar,
			ps,
			createFn,
			func() ygot.ValidatedGoStruct {
				return nil
			},
			func() *ytypes.Schema { return nil }),
		nil
}

func newQueryField[T any](path, origin string) (ygnmi.PathStruct, func(vgs ygot.ValidatedGoStruct) (T, bool), bool, error) {
	root := ygnmi.NewDeviceRootBase()
	protoPath, err := ygot.StringToStructuredPath(path)
	if err != nil {
		return nil, nil, false, err
	}

	var ps ygnmi.PathStruct = root
	for _, elem := range protoPath.Elem {
		keys := map[string]interface{}{}
		for key, val := range elem.Key {
			keys[key] = val
		}
		ps = ygnmi.NewNodePath([]string{elem.Name}, keys, ps)
	}
	root.PutCustomData(ygnmi.OriginOverride, origin)

	createFn := func(vgs ygot.ValidatedGoStruct) (T, bool) {
		return *new(T), true
	}

	var zero T
	scalar := true
	switch paramType := reflect.TypeOf(zero); paramType.Kind() {
	case reflect.Struct:
		return nil, nil, false, fmt.Errorf("struct type not supported, must use struct pointer")
	case reflect.Pointer:
		scalar = false
		createFn = func(vgs ygot.ValidatedGoStruct) (T, bool) {
			return reflect.New(paramType.Elem()).Interface().(T), true
		}
	case reflect.Chan, reflect.Func:
		return nil, nil, false, fmt.Errorf("unsupported parameterize type: %s", paramType.Kind().String())
	}
	return ps, createFn, scalar, nil
}

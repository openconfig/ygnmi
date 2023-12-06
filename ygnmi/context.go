// Copyright 2023 Google Inc.
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
)

// RequestValues contains request-scoped values for ygnmi queries.
type RequestValues struct {
	// StateFiltered is a key type that means that the query is
	// uninterested in /state paths and will filter them out.
	StateFiltered bool
	// ConfigFiltered is a key type that means that the query is
	// uninterested in /config paths and will filter them out.
	ConfigFiltered bool
}

// FromContext extracts certain ygnmi request-scoped values, if present.
func FromContext(ctx context.Context) *RequestValues {
	requestValues, _ := ctx.Value(requestValuesKey{}).(*RequestValues)
	return requestValues
}

// NewContext returns a new Context carrying ygnmi request-scoped values.
func NewContext(ctx context.Context, q UntypedQuery) context.Context {
	return context.WithValue(ctx, requestValuesKey{}, &RequestValues{
		StateFiltered:  q.isCompressedSchema() && !q.IsState(),
		ConfigFiltered: q.isCompressedSchema() && q.IsState(),
	})
}

type requestValuesKey struct{}

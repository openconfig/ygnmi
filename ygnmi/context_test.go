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

package ygnmi_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/ygnmi/exampleoc/exampleocpath"
	"github.com/openconfig/ygnmi/internal/uexampleoc/uexampleocpath"
	"github.com/openconfig/ygnmi/ygnmi"
)

func TestFromContext(t *testing.T) {
	tests := []struct {
		desc              string
		inContext         context.Context
		wantRequestValues *ygnmi.RequestValues
	}{{
		desc:      "compress-config",
		inContext: ygnmi.NewContext(context.Background(), exampleocpath.Root().Parent().Config()),
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  true,
			ConfigFiltered: false,
		},
	}, {
		desc:      "compress-config-leaf",
		inContext: ygnmi.NewContext(context.Background(), exampleocpath.Root().Parent().Child().Five().Config()),
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  true,
			ConfigFiltered: false,
		},
	}, {
		desc:      "compress-state",
		inContext: ygnmi.NewContext(context.Background(), exampleocpath.Root().Parent().State()),
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  false,
			ConfigFiltered: true,
		},
	}, {
		desc:      "compress-state-leaf",
		inContext: ygnmi.NewContext(context.Background(), exampleocpath.Root().Parent().Child().Five().State()),
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  false,
			ConfigFiltered: true,
		},
	}, {
		desc:      "uncompressed-container",
		inContext: ygnmi.NewContext(context.Background(), uexampleocpath.Root().Parent()),
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  false,
			ConfigFiltered: false,
		},
	}, {
		desc:      "uncompressed-config-container",
		inContext: ygnmi.NewContext(context.Background(), uexampleocpath.Root().Parent().Child().Config()),
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  false,
			ConfigFiltered: false,
		},
	}, {
		desc:      "uncompressed-leaf",
		inContext: ygnmi.NewContext(context.Background(), uexampleocpath.Root().Parent().Child().Config().Five()),
		wantRequestValues: &ygnmi.RequestValues{
			StateFiltered:  false,
			ConfigFiltered: false,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := ygnmi.FromContext(tt.inContext)
			if diff := cmp.Diff(tt.wantRequestValues, got); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}
		})
	}
}

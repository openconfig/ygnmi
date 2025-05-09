// Copyright 2019 Google Inc.
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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygnmi/internal/testutil"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/proto"
)

type deviceRoot struct {
	*DeviceRootBase
}

func TestResolvePath(t *testing.T) {
	wantCustomData := map[string]interface{}{"foo": "bar"}
	root := deviceRoot{NewDeviceRootBase()}
	root.PutCustomData("foo", "bar")

	tests := []struct {
		name        string
		in          PathStruct
		wantPathStr string
		wantErr     bool
	}{{
		name: "simple",
		in: &NodePath{
			relSchemaPath: []string{"child"},
			keys:          map[string]interface{}{},
			p: &NodePath{
				relSchemaPath: []string{"parent"},
				keys:          map[string]interface{}{},
				p:             root,
			},
		},
		wantPathStr: "/parent/child",
	}, {
		name:        "root",
		in:          root,
		wantPathStr: "/",
	}, {
		name: "list",
		in: &NodePath{
			relSchemaPath: []string{"values", "value"},
			keys:          map[string]interface{}{"ID": 5},
			p: &NodePath{
				relSchemaPath: []string{"parent"},
				keys:          map[string]interface{}{},
				p:             root,
			},
		},
		wantPathStr: "/parent/values/value[ID=5]",
	}, {
		name: "list with unconvertible key value",
		in: &NodePath{
			relSchemaPath: []string{"values", "value"},
			keys:          map[string]interface{}{"ID": complex(1, 2)},
			p: &NodePath{
				relSchemaPath: []string{"parent"},
				keys:          map[string]interface{}{},
				p:             root,
			},
		},
		wantPathStr: "",
		wantErr:     true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath, err := ygot.StringToStructuredPath(tt.wantPathStr)
			if err != nil {
				t.Fatal(err)
			}

			gotPath, gotCustomData, gotErrs := ResolvePath(tt.in)
			if gotErrs != nil && !tt.wantErr {
				t.Fatal(gotErrs)
			} else if gotErrs == nil && tt.wantErr {
				t.Fatal("expected error but did not receive any")
			}
			if gotErrs != nil {
				return
			}

			if diff := cmp.Diff(wantPath, gotPath, cmp.Comparer(proto.Equal)); diff != "" {
				t.Errorf("ResolvePath returned diff (-want, +got):\n%s", diff)
			}

			if gotPath.Elem == nil {
				t.Errorf("gotPath.PathElem is nil, but should not be")
			}

			if diff := cmp.Diff(wantCustomData, gotCustomData); diff != "" {
				t.Errorf("ResolvePath: customData is not same as expected (-want, +got)\n%s", diff)
			}
		})
	}
}

func TestResolveRelPath(t *testing.T) {
	root := &NodePath{}

	tests := []struct {
		name    string
		in      PathStruct
		want    string
		wantErr bool
	}{{
		name: "simple",
		in: &NodePath{
			relSchemaPath: []string{"child"},
			keys:          map[string]interface{}{},
		},
		want: "child",
	}, {
		name: "root",
		in:   root,
		want: "",
	}, {
		name: "list",
		in: &NodePath{
			relSchemaPath: []string{"values", "value"},
			keys:          map[string]interface{}{"ID": 5},
		},
		want: "values/value[ID=5]",
	}, {
		name: "list with unconvertible key value",
		in: &NodePath{
			relSchemaPath: []string{"values", "value"},
			keys:          map[string]interface{}{"ID": complex(1, 2)},
		},
		want:    "",
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantP, err := ygot.StringToStructuredPath(tt.want)
			wantPath := wantP.Elem
			if err != nil {
				t.Fatal(err)
			}

			gotPath, gotErrs := ResolveRelPath(tt.in)
			if gotErrs != nil && !tt.wantErr {
				t.Fatal(gotErrs)
			} else if gotErrs == nil && tt.wantErr {
				t.Fatal("expected error but did not receive any")
			}

			if diff := cmp.Diff(wantPath, gotPath, cmp.Comparer(proto.Equal)); diff != "" {
				t.Errorf("ResolveRelPath returned diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExtractPathKeys(t *testing.T) {
	samplePath := testutil.GNMIPath(t, "/lists/list[key1=test][key2=test2]/sublists/sublist[index=1]/some/value")

	tests := []struct {
		desc      string
		path      *gpb.Path
		keystruct any
		want      any
		wantErr   string
	}{{
		desc:    "nil",
		path:    samplePath,
		wantErr: "expected pointer to struct",
	}, {
		desc:      "not pointer",
		path:      samplePath,
		keystruct: struct{}{},
		wantErr:   "expected pointer to struct",
	}, {
		desc:      "not pointer to struct",
		path:      samplePath,
		keystruct: ygot.String("foo"),
		wantErr:   "expected pointer to struct",
	}, {
		desc: "wrong keys",
		path: samplePath,
		keystruct: &(struct {
			ListKey1 string `pathkey:"/foo:key1"`
		}{}),
		want: &(struct {
			ListKey1 string `pathkey:"/foo:key1"`
		}{}),
	}, {
		desc: "all keys",
		path: samplePath,
		keystruct: &(struct {
			ListKey1 string `pathkey:"/lists/list:key1"`
			ListKey2 string `pathkey:"/lists/list:key2"`
			SubList  string `pathkey:"/lists/list/sublists/sublist:index"`
		}{}),
		want: &(struct {
			ListKey1 string `pathkey:"/lists/list:key1"`
			ListKey2 string `pathkey:"/lists/list:key2"`
			SubList  string `pathkey:"/lists/list/sublists/sublist:index"`
		}{
			ListKey1: "test",
			ListKey2: "test2",
			SubList:  "1",
		}),
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := ExtractPathKeys(tt.path, tt.keystruct)
			if diff := errdiff.Substring(err, tt.wantErr); diff != "" {
				t.Fatalf("ExtractPathKeys() unexpect error: %s", diff)
			}
			if err != nil {
				return
			}
			if d := cmp.Diff(tt.keystruct, tt.want); d != "" {
				t.Errorf("ExtractPathKeys() unexpected diff: %s", d)
			}
		})
	}
}

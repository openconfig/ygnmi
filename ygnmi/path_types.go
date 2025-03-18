// Copyright 2020 Google Inc.
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
	"reflect"
	"strings"

	"github.com/openconfig/gnmi/errlist"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
)

const (
	// PathStructInterfaceName is the name for the interface implemented by all
	// generated path structs.
	PathStructInterfaceName = "PathStruct"
	// PathBaseTypeName is the type name of the common embedded struct
	// containing the path information for a path struct.
	PathBaseTypeName = "NodePath"
	// FakeRootBaseTypeName is the type name of the fake root struct which
	// should be embedded within the fake root path struct.
	FakeRootBaseTypeName = "DeviceRootBase"
)

// PathStruct is an interface that is implemented by any generated path struct
// type; it allows for generic handling of a path struct at any node.
type PathStruct interface {
	parent() PathStruct
	relPath() ([]*gpb.PathElem, []error)
	schemaPath() []string
	getKeys() map[string]interface{}
}

// NewNodePath is the constructor for NodePath.
func NewNodePath(relSchemaPath []string, keys map[string]interface{}, p PathStruct) *NodePath {
	return &NodePath{relSchemaPath: relSchemaPath, keys: keys, p: p}
}

// NodePath is a common embedded type within all path structs. It
// keeps track of the necessary information to create the relative schema path
// as a []*gpb.PathElem during later processing using the Resolve() method,
// thereby delaying any errors being reported until that time.
type NodePath struct {
	relSchemaPath []string
	keys          map[string]interface{}
	p             PathStruct
}

// fakeRootPathStruct is an interface that is implemented by the fake root path
// struct type.
type fakeRootPathStruct interface {
	PathStruct
	CustomData() map[string]interface{}
}

// NewDeviceRootBase creates a DeviceRootBase.
func NewDeviceRootBase() *DeviceRootBase {
	return &DeviceRootBase{
		NodePath:   &NodePath{},
		customData: map[string]interface{}{},
	}
}

// DeviceRootBase represents the fakeroot for all YANG schema elements.
type DeviceRootBase struct {
	*NodePath
	// customData is meant to store root-specific information that may be
	// useful to know when processing the resolved path. It is meant to be
	// accessible through a user-defined accessor.
	customData map[string]interface{}
}

// CustomData returns the customData field of the DeviceRootBase struct.
func (d *DeviceRootBase) CustomData() map[string]interface{} {
	return d.customData
}

// PutCustomData modifies an entry in the customData field of the DeviceRootBase struct.
func (d *DeviceRootBase) PutCustomData(key string, val interface{}) {
	d.customData[key] = val
}

// ResolvePath is a helper which returns the resolved *gpb.Path of a PathStruct
// node as well as the root node's customData.
func ResolvePath(n PathStruct) (*gpb.Path, map[string]interface{}, error) {
	p := []*gpb.PathElem{}
	err := errlist.List{}
	for ; n.parent() != nil; n = n.parent() {
		rel, es := n.relPath()
		if es != nil {
			err.Add(es...)
			continue
		}
		p = append(rel, p...)
	}
	if err.Err() != nil {
		return nil, nil, err.Err()
	}

	root, ok := n.(fakeRootPathStruct)
	if !ok {
		err.Add(fmt.Errorf("ygot.ResolvePath(PathStruct): got unexpected root of (type, value) (%T, %v)", n, n))
		return nil, nil, err.Err()
	}
	return &gpb.Path{Elem: p}, root.CustomData(), nil
}

// ResolveRelPath returns the partial []*gpb.PathElem representing the
// PathStruct's relative path.
func ResolveRelPath(n PathStruct) ([]*gpb.PathElem, []error) {
	return n.relPath()
}

// ModifyKey updates a NodePath's key value.
func ModifyKey(n *NodePath, name string, value interface{}) {
	n.keys[name] = value
}

// relPath converts the information stored in NodePath into the partial
// []*gpb.PathElem representing the node's relative path.
func (n *NodePath) relPath() ([]*gpb.PathElem, []error) {
	var pathElems []*gpb.PathElem
	for _, name := range n.relSchemaPath {
		pathElems = append(pathElems, &gpb.PathElem{Name: name})
	}
	if len(n.keys) == 0 {
		return pathElems, nil
	}

	var errs []error
	keys := make(map[string]string)
	for name, val := range n.keys {
		var err error
		// TODO(wenbli): It is ideal to also implement leaf restriction validation.
		if keys[name], err = ygot.KeyValueAsString(val); err != nil {
			errs = append(errs, err)
		}
	}
	if errs != nil {
		return nil, errs
	}
	pathElems[len(pathElems)-1].Key = keys
	return pathElems, nil
}

func (n *NodePath) parent() PathStruct { return n.p }

func (n *NodePath) schemaPath() []string { return n.relSchemaPath }

func (n *NodePath) getKeys() map[string]interface{} { return n.keys }

// configToStatePS converts a config path struct it's equivalent state path struct.
// Note: This function only works for OpenConfig style path structs, where the last container is either "config" or "state".
func configToStatePS(ps PathStruct) PathStruct {
	statePath := []*NodePath{}
	configPath := []PathStruct{}
	for ; ps != nil; ps = ps.parent() {
		sp := make([]string, len(ps.schemaPath()))
		copy(sp, ps.schemaPath())
		if len(ps.schemaPath()) == 2 && ps.schemaPath()[0] == configPathElem {
			sp[0] = statePathElem
		}
		configPath = append(configPath, ps)
		statePath = append(statePath, NewNodePath(sp, ps.getKeys(), nil))
	}
	for i := 0; i < len(statePath)-2; i++ {
		statePath[i].p = statePath[i+1]
	}
	// Copy the fake root directly
	statePath[len(statePath)-2].p = configPath[len(configPath)-1]
	return statePath[0]
}

func configToState[T any](cfg AnyQuery[T]) AnyQuery[T] {
	return &SingletonQueryStruct[T]{
		baseQuery: baseQuery[T]{
			goStructName:     cfg.dirName(),
			state:            !cfg.IsState(),
			shadowpath:       !cfg.isShadowPath(),
			ps:               configToStatePS(cfg.PathStruct()),
			leaf:             cfg.isLeaf(),
			scalar:           cfg.isScalar(),
			compressedSchema: cfg.isCompressedSchema(),
			listContainer:    cfg.isListContainer(),
			extractFn:        cfg.extract,
			goStructFn:       cfg.goStruct,
			yschemaFn:        cfg.schema,
			compInfo:         cfg.compressInfo(),
		},
	}
}

const pathKeyTag = "pathkey"

// ExtractPathKeys extracts keys from gNMI path into a custom struct. Use the "pathkey" field tag to specify the keys to extract.
// The format of the pathkey is <path>:<key>
// Example:
//
//	Path: /interfaces/interfaces[name=eth0]/subinterfaces/subinterface[index=0]/config/enabled
//	Struct Fields: type ifacekey struct { Name string `pathkey:"/interfaces/interfaces:name"` }
func ExtractPathKeys(path *gpb.Path, keystruct any) error {
	if keystruct == nil {
		return fmt.Errorf("expected pointer to struct")
	}
	t := reflect.TypeOf(keystruct)
	if ok := util.IsTypeStructPtr(t); !ok {
		return fmt.Errorf("unexpected type %v, expected pointer to struct", t.Kind())
	}
	t = t.Elem()
	v := reflect.ValueOf(keystruct).Elem()
	for i := 0; i < t.NumField(); i++ {
		key, ok := t.Field(i).Tag.Lookup(pathKeyTag)
		if !ok {
			continue
		}
		pathKey := strings.Split(key, ":")
		if len(pathKey) != 2 {
			return fmt.Errorf("invalid struct key format %q, expect <path>:<key>", key)
		}
		elems := strings.Split(strings.TrimPrefix(pathKey[0], "/"), "/")
		match := true
		for keyIdx := 0; keyIdx < len(path.Elem) && keyIdx < len(elems); keyIdx++ {
			if path.Elem[keyIdx].Name != elems[keyIdx] {
				match = false
				break
			}
		}
		if match && len(elems) <= len(path.Elem) {
			v.Field(i).Set(reflect.ValueOf(path.Elem[len(elems)-1].Key[pathKey[1]]))
		}
	}
	return nil
}

// Copyright 2022 Google LLC
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

/*
Package nested is a generated package which contains definitions
of structs which generate gNMI paths for a YANG schema. The generated paths are
based on a compressed form of the schema.

This package was generated by ygnmi version: (devel): (ygot: v0.21.0)
using the following YANG input files:
	- ../../pathgen/testdata/yang/openconfig-simple.yang
	- ../../pathgen/testdata/yang/openconfig-withlistval.yang
	- ../../pathgen/testdata/yang/openconfig-nested.yang
Imported modules were sourced from:
*/
package nested

import (
	oc "github.com/openconfig/ygnmi/internal/exampleoc"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// APath represents the /openconfig-nested/a YANG schema element.
type APath struct {
	*ygot.NodePath
}

// APathAny represents the wildcard version of the /openconfig-nested/a YANG schema element.
type APathAny struct {
	*ygot.NodePath
}

// B (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "b"
// Path from root: "/a/b"
func (n *APath) B() *A_BPath {
	return &A_BPath{
		NodePath: ygot.NewNodePath(
			[]string{"b"},
			map[string]interface{}{},
			n,
		),
	}
}

// B (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "b"
// Path from root: "/a/b"
func (n *APathAny) B() *A_BPathAny {
	return &A_BPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"b"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *APath) State() ygnmi.SingletonQuery[*oc.A] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A](
		"A",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *APathAny) State() ygnmi.WildcardQuery[*oc.A] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A](
		"A",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *APath) Config() ygnmi.ConfigQuery[*oc.A] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A](
		"A",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *APathAny) Config() ygnmi.WildcardQuery[*oc.A] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A](
		"A",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_BPath represents the /openconfig-nested/a/b YANG schema element.
type A_BPath struct {
	*ygot.NodePath
}

// A_BPathAny represents the wildcard version of the /openconfig-nested/a/b YANG schema element.
type A_BPathAny struct {
	*ygot.NodePath
}

// C (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "c"
// Path from root: "/a/b/c"
func (n *A_BPath) C() *A_B_CPath {
	return &A_B_CPath{
		NodePath: ygot.NewNodePath(
			[]string{"c"},
			map[string]interface{}{},
			n,
		),
	}
}

// C (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "c"
// Path from root: "/a/b/c"
func (n *A_BPathAny) C() *A_B_CPathAny {
	return &A_B_CPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"c"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_BPath) State() ygnmi.SingletonQuery[*oc.A_B] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B](
		"A_B",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_BPathAny) State() ygnmi.WildcardQuery[*oc.A_B] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B](
		"A_B",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_BPath) Config() ygnmi.ConfigQuery[*oc.A_B] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B](
		"A_B",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_BPathAny) Config() ygnmi.WildcardQuery[*oc.A_B] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B](
		"A_B",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_CPath represents the /openconfig-nested/a/b/c YANG schema element.
type A_B_CPath struct {
	*ygot.NodePath
}

// A_B_CPathAny represents the wildcard version of the /openconfig-nested/a/b/c YANG schema element.
type A_B_CPathAny struct {
	*ygot.NodePath
}

// D (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "d"
// Path from root: "/a/b/c/d"
func (n *A_B_CPath) D() *A_B_C_DPath {
	return &A_B_C_DPath{
		NodePath: ygot.NewNodePath(
			[]string{"d"},
			map[string]interface{}{},
			n,
		),
	}
}

// D (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "d"
// Path from root: "/a/b/c/d"
func (n *A_B_CPathAny) D() *A_B_C_DPathAny {
	return &A_B_C_DPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"d"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_CPath) State() ygnmi.SingletonQuery[*oc.A_B_C] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C](
		"A_B_C",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_CPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C](
		"A_B_C",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_CPath) Config() ygnmi.ConfigQuery[*oc.A_B_C] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C](
		"A_B_C",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_CPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C](
		"A_B_C",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_C_DPath represents the /openconfig-nested/a/b/c/d YANG schema element.
type A_B_C_DPath struct {
	*ygot.NodePath
}

// A_B_C_DPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d YANG schema element.
type A_B_C_DPathAny struct {
	*ygot.NodePath
}

// E (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "e"
// Path from root: "/a/b/c/d/e"
func (n *A_B_C_DPath) E() *A_B_C_D_EPath {
	return &A_B_C_D_EPath{
		NodePath: ygot.NewNodePath(
			[]string{"e"},
			map[string]interface{}{},
			n,
		),
	}
}

// E (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "e"
// Path from root: "/a/b/c/d/e"
func (n *A_B_C_DPathAny) E() *A_B_C_D_EPathAny {
	return &A_B_C_D_EPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"e"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_DPath) State() ygnmi.SingletonQuery[*oc.A_B_C_D] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C_D](
		"A_B_C_D",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_DPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C_D] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D](
		"A_B_C_D",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_DPath) Config() ygnmi.ConfigQuery[*oc.A_B_C_D] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C_D](
		"A_B_C_D",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_DPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C_D] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D](
		"A_B_C_D",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_C_D_EPath represents the /openconfig-nested/a/b/c/d/e YANG schema element.
type A_B_C_D_EPath struct {
	*ygot.NodePath
}

// A_B_C_D_EPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d/e YANG schema element.
type A_B_C_D_EPathAny struct {
	*ygot.NodePath
}

// F (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "f"
// Path from root: "/a/b/c/d/e/f"
func (n *A_B_C_D_EPath) F() *A_B_C_D_E_FPath {
	return &A_B_C_D_E_FPath{
		NodePath: ygot.NewNodePath(
			[]string{"f"},
			map[string]interface{}{},
			n,
		),
	}
}

// F (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "f"
// Path from root: "/a/b/c/d/e/f"
func (n *A_B_C_D_EPathAny) F() *A_B_C_D_E_FPathAny {
	return &A_B_C_D_E_FPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"f"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_EPath) State() ygnmi.SingletonQuery[*oc.A_B_C_D_E] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C_D_E](
		"A_B_C_D_E",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_EPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C_D_E] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E](
		"A_B_C_D_E",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_EPath) Config() ygnmi.ConfigQuery[*oc.A_B_C_D_E] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C_D_E](
		"A_B_C_D_E",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_EPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C_D_E] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E](
		"A_B_C_D_E",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_C_D_E_FPath represents the /openconfig-nested/a/b/c/d/e/f YANG schema element.
type A_B_C_D_E_FPath struct {
	*ygot.NodePath
}

// A_B_C_D_E_FPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d/e/f YANG schema element.
type A_B_C_D_E_FPathAny struct {
	*ygot.NodePath
}

// G (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "g"
// Path from root: "/a/b/c/d/e/f/g"
func (n *A_B_C_D_E_FPath) G() *A_B_C_D_E_F_GPath {
	return &A_B_C_D_E_F_GPath{
		NodePath: ygot.NewNodePath(
			[]string{"g"},
			map[string]interface{}{},
			n,
		),
	}
}

// G (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "g"
// Path from root: "/a/b/c/d/e/f/g"
func (n *A_B_C_D_E_FPathAny) G() *A_B_C_D_E_F_GPathAny {
	return &A_B_C_D_E_F_GPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"g"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_FPath) State() ygnmi.SingletonQuery[*oc.A_B_C_D_E_F] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C_D_E_F](
		"A_B_C_D_E_F",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_FPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F](
		"A_B_C_D_E_F",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_FPath) Config() ygnmi.ConfigQuery[*oc.A_B_C_D_E_F] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C_D_E_F](
		"A_B_C_D_E_F",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_FPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F](
		"A_B_C_D_E_F",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_C_D_E_F_GPath represents the /openconfig-nested/a/b/c/d/e/f/g YANG schema element.
type A_B_C_D_E_F_GPath struct {
	*ygot.NodePath
}

// A_B_C_D_E_F_GPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d/e/f/g YANG schema element.
type A_B_C_D_E_F_GPathAny struct {
	*ygot.NodePath
}

// H (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "h"
// Path from root: "/a/b/c/d/e/f/g/h"
func (n *A_B_C_D_E_F_GPath) H() *A_B_C_D_E_F_G_HPath {
	return &A_B_C_D_E_F_G_HPath{
		NodePath: ygot.NewNodePath(
			[]string{"h"},
			map[string]interface{}{},
			n,
		),
	}
}

// H (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "h"
// Path from root: "/a/b/c/d/e/f/g/h"
func (n *A_B_C_D_E_F_GPathAny) H() *A_B_C_D_E_F_G_HPathAny {
	return &A_B_C_D_E_F_G_HPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"h"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_GPath) State() ygnmi.SingletonQuery[*oc.A_B_C_D_E_F_G] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C_D_E_F_G](
		"A_B_C_D_E_F_G",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_GPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G](
		"A_B_C_D_E_F_G",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_GPath) Config() ygnmi.ConfigQuery[*oc.A_B_C_D_E_F_G] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C_D_E_F_G](
		"A_B_C_D_E_F_G",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_GPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G](
		"A_B_C_D_E_F_G",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_C_D_E_F_G_HPath represents the /openconfig-nested/a/b/c/d/e/f/g/h YANG schema element.
type A_B_C_D_E_F_G_HPath struct {
	*ygot.NodePath
}

// A_B_C_D_E_F_G_HPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d/e/f/g/h YANG schema element.
type A_B_C_D_E_F_G_HPathAny struct {
	*ygot.NodePath
}

// I (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "i"
// Path from root: "/a/b/c/d/e/f/g/h/i"
func (n *A_B_C_D_E_F_G_HPath) I() *A_B_C_D_E_F_G_H_IPath {
	return &A_B_C_D_E_F_G_H_IPath{
		NodePath: ygot.NewNodePath(
			[]string{"i"},
			map[string]interface{}{},
			n,
		),
	}
}

// I (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "i"
// Path from root: "/a/b/c/d/e/f/g/h/i"
func (n *A_B_C_D_E_F_G_HPathAny) I() *A_B_C_D_E_F_G_H_IPathAny {
	return &A_B_C_D_E_F_G_H_IPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"i"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_HPath) State() ygnmi.SingletonQuery[*oc.A_B_C_D_E_F_G_H] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C_D_E_F_G_H](
		"A_B_C_D_E_F_G_H",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_HPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H](
		"A_B_C_D_E_F_G_H",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_HPath) Config() ygnmi.ConfigQuery[*oc.A_B_C_D_E_F_G_H] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C_D_E_F_G_H](
		"A_B_C_D_E_F_G_H",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_HPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H](
		"A_B_C_D_E_F_G_H",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_C_D_E_F_G_H_IPath represents the /openconfig-nested/a/b/c/d/e/f/g/h/i YANG schema element.
type A_B_C_D_E_F_G_H_IPath struct {
	*ygot.NodePath
}

// A_B_C_D_E_F_G_H_IPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d/e/f/g/h/i YANG schema element.
type A_B_C_D_E_F_G_H_IPathAny struct {
	*ygot.NodePath
}

// J (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "j"
// Path from root: "/a/b/c/d/e/f/g/h/i/j"
func (n *A_B_C_D_E_F_G_H_IPath) J() *A_B_C_D_E_F_G_H_I_JPath {
	return &A_B_C_D_E_F_G_H_I_JPath{
		NodePath: ygot.NewNodePath(
			[]string{"j"},
			map[string]interface{}{},
			n,
		),
	}
}

// J (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "j"
// Path from root: "/a/b/c/d/e/f/g/h/i/j"
func (n *A_B_C_D_E_F_G_H_IPathAny) J() *A_B_C_D_E_F_G_H_I_JPathAny {
	return &A_B_C_D_E_F_G_H_I_JPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"j"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_IPath) State() ygnmi.SingletonQuery[*oc.A_B_C_D_E_F_G_H_I] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C_D_E_F_G_H_I](
		"A_B_C_D_E_F_G_H_I",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_IPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H_I] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H_I](
		"A_B_C_D_E_F_G_H_I",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_IPath) Config() ygnmi.ConfigQuery[*oc.A_B_C_D_E_F_G_H_I] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C_D_E_F_G_H_I](
		"A_B_C_D_E_F_G_H_I",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_IPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H_I] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H_I](
		"A_B_C_D_E_F_G_H_I",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_C_D_E_F_G_H_I_JPath represents the /openconfig-nested/a/b/c/d/e/f/g/h/i/j YANG schema element.
type A_B_C_D_E_F_G_H_I_JPath struct {
	*ygot.NodePath
}

// A_B_C_D_E_F_G_H_I_JPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d/e/f/g/h/i/j YANG schema element.
type A_B_C_D_E_F_G_H_I_JPathAny struct {
	*ygot.NodePath
}

// K (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "k"
// Path from root: "/a/b/c/d/e/f/g/h/i/j/k"
func (n *A_B_C_D_E_F_G_H_I_JPath) K() *A_B_C_D_E_F_G_H_I_J_KPath {
	return &A_B_C_D_E_F_G_H_I_J_KPath{
		NodePath: ygot.NewNodePath(
			[]string{"k"},
			map[string]interface{}{},
			n,
		),
	}
}

// K (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "k"
// Path from root: "/a/b/c/d/e/f/g/h/i/j/k"
func (n *A_B_C_D_E_F_G_H_I_JPathAny) K() *A_B_C_D_E_F_G_H_I_J_KPathAny {
	return &A_B_C_D_E_F_G_H_I_J_KPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"k"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_JPath) State() ygnmi.SingletonQuery[*oc.A_B_C_D_E_F_G_H_I_J] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C_D_E_F_G_H_I_J](
		"A_B_C_D_E_F_G_H_I_J",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_JPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J](
		"A_B_C_D_E_F_G_H_I_J",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_JPath) Config() ygnmi.ConfigQuery[*oc.A_B_C_D_E_F_G_H_I_J] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C_D_E_F_G_H_I_J](
		"A_B_C_D_E_F_G_H_I_J",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_JPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J](
		"A_B_C_D_E_F_G_H_I_J",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_C_D_E_F_G_H_I_J_KPath represents the /openconfig-nested/a/b/c/d/e/f/g/h/i/j/k YANG schema element.
type A_B_C_D_E_F_G_H_I_J_KPath struct {
	*ygot.NodePath
}

// A_B_C_D_E_F_G_H_I_J_KPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d/e/f/g/h/i/j/k YANG schema element.
type A_B_C_D_E_F_G_H_I_J_KPathAny struct {
	*ygot.NodePath
}

// L (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "l"
// Path from root: "/a/b/c/d/e/f/g/h/i/j/k/l"
func (n *A_B_C_D_E_F_G_H_I_J_KPath) L() *A_B_C_D_E_F_G_H_I_J_K_LPath {
	return &A_B_C_D_E_F_G_H_I_J_K_LPath{
		NodePath: ygot.NewNodePath(
			[]string{"l"},
			map[string]interface{}{},
			n,
		),
	}
}

// L (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "l"
// Path from root: "/a/b/c/d/e/f/g/h/i/j/k/l"
func (n *A_B_C_D_E_F_G_H_I_J_KPathAny) L() *A_B_C_D_E_F_G_H_I_J_K_LPathAny {
	return &A_B_C_D_E_F_G_H_I_J_K_LPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"l"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_KPath) State() ygnmi.SingletonQuery[*oc.A_B_C_D_E_F_G_H_I_J_K] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C_D_E_F_G_H_I_J_K](
		"A_B_C_D_E_F_G_H_I_J_K",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_KPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K](
		"A_B_C_D_E_F_G_H_I_J_K",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_KPath) Config() ygnmi.ConfigQuery[*oc.A_B_C_D_E_F_G_H_I_J_K] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C_D_E_F_G_H_I_J_K](
		"A_B_C_D_E_F_G_H_I_J_K",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_KPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K](
		"A_B_C_D_E_F_G_H_I_J_K",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_C_D_E_F_G_H_I_J_K_LPath represents the /openconfig-nested/a/b/c/d/e/f/g/h/i/j/k/l YANG schema element.
type A_B_C_D_E_F_G_H_I_J_K_LPath struct {
	*ygot.NodePath
}

// A_B_C_D_E_F_G_H_I_J_K_LPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d/e/f/g/h/i/j/k/l YANG schema element.
type A_B_C_D_E_F_G_H_I_J_K_LPathAny struct {
	*ygot.NodePath
}

// M (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "m"
// Path from root: "/a/b/c/d/e/f/g/h/i/j/k/l/m"
func (n *A_B_C_D_E_F_G_H_I_J_K_LPath) M() *A_B_C_D_E_F_G_H_I_J_K_L_MPath {
	return &A_B_C_D_E_F_G_H_I_J_K_L_MPath{
		NodePath: ygot.NewNodePath(
			[]string{"m"},
			map[string]interface{}{},
			n,
		),
	}
}

// M (container):
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "m"
// Path from root: "/a/b/c/d/e/f/g/h/i/j/k/l/m"
func (n *A_B_C_D_E_F_G_H_I_J_K_LPathAny) M() *A_B_C_D_E_F_G_H_I_J_K_L_MPathAny {
	return &A_B_C_D_E_F_G_H_I_J_K_L_MPathAny{
		NodePath: ygot.NewNodePath(
			[]string{"m"},
			map[string]interface{}{},
			n,
		),
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_K_LPath) State() ygnmi.SingletonQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L](
		"A_B_C_D_E_F_G_H_I_J_K_L",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_K_LPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L](
		"A_B_C_D_E_F_G_H_I_J_K_L",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_K_LPath) Config() ygnmi.ConfigQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L](
		"A_B_C_D_E_F_G_H_I_J_K_L",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_K_LPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L](
		"A_B_C_D_E_F_G_H_I_J_K_L",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// A_B_C_D_E_F_G_H_I_J_K_L_MPath represents the /openconfig-nested/a/b/c/d/e/f/g/h/i/j/k/l/m YANG schema element.
type A_B_C_D_E_F_G_H_I_J_K_L_MPath struct {
	*ygot.NodePath
}

// A_B_C_D_E_F_G_H_I_J_K_L_MPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d/e/f/g/h/i/j/k/l/m YANG schema element.
type A_B_C_D_E_F_G_H_I_J_K_L_MPathAny struct {
	*ygot.NodePath
}

// A_B_C_D_E_F_G_H_I_J_K_L_M_FooPath represents the /openconfig-nested/a/b/c/d/e/f/g/h/i/j/k/l/m/state/foo YANG schema element.
type A_B_C_D_E_F_G_H_I_J_K_L_M_FooPath struct {
	parent ygot.PathStruct
}

// A_B_C_D_E_F_G_H_I_J_K_L_M_FooPathAny represents the wildcard version of the /openconfig-nested/a/b/c/d/e/f/g/h/i/j/k/l/m/state/foo YANG schema element.
type A_B_C_D_E_F_G_H_I_J_K_L_M_FooPathAny struct {
	parent ygot.PathStruct
}

// Foo corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *A_B_C_D_E_F_G_H_I_J_K_L_MPath) Foo() *A_B_C_D_E_F_G_H_I_J_K_L_M_FooPath {
	return &A_B_C_D_E_F_G_H_I_J_K_L_M_FooPath{
		parent: n,
	}
}

// Foo corresponds to an ambiguous path; use .Config() or .State() to get a resolved path for this leaf.
// Note: The returned struct does not implement the PathStruct interface.
func (n *A_B_C_D_E_F_G_H_I_J_K_L_MPathAny) Foo() *A_B_C_D_E_F_G_H_I_J_K_L_M_FooPathAny {
	return &A_B_C_D_E_F_G_H_I_J_K_L_M_FooPathAny{
		parent: n,
	}
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_K_L_MPath) State() ygnmi.SingletonQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L_M] {
	return ygnmi.NewNonLeafSingletonQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L_M](
		"A_B_C_D_E_F_G_H_I_J_K_L_M",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_K_L_MPathAny) State() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L_M] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L_M](
		"A_B_C_D_E_F_G_H_I_J_K_L_M",
		true,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_K_L_MPath) Config() ygnmi.ConfigQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L_M] {
	return ygnmi.NewNonLeafConfigQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L_M](
		"A_B_C_D_E_F_G_H_I_J_K_L_M",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *A_B_C_D_E_F_G_H_I_J_K_L_MPathAny) Config() ygnmi.WildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L_M] {
	return ygnmi.NewNonLeafWildcardQuery[*oc.A_B_C_D_E_F_G_H_I_J_K_L_M](
		"A_B_C_D_E_F_G_H_I_J_K_L_M",
		false,
		n,
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "state/foo"
// Path from root: "/a/b/c/d/e/f/g/h/i/j/k/l/m/state/foo"
func (n *A_B_C_D_E_F_G_H_I_J_K_L_M_FooPath) State() ygnmi.SingletonQuery[string] {
	return ygnmi.NewLeafSingletonQuery[string](
		"A_B_C_D_E_F_G_H_I_J_K_L_M",
		true,
		true,
		ygot.NewNodePath(
			[]string{"state", "foo"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.A_B_C_D_E_F_G_H_I_J_K_L_M).Foo
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.A_B_C_D_E_F_G_H_I_J_K_L_M) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}

// State returns a Query that can be used in gNMI operations.
// ----------------------------------------
// Defining module: "openconfig-nested"
// Instantiating module: "openconfig-nested"
// Path from parent: "state/foo"
// Path from root: "/a/b/c/d/e/f/g/h/i/j/k/l/m/state/foo"
func (n *A_B_C_D_E_F_G_H_I_J_K_L_M_FooPathAny) State() ygnmi.WildcardQuery[string] {
	return ygnmi.NewLeafWildcardQuery[string](
		"A_B_C_D_E_F_G_H_I_J_K_L_M",
		true,
		true,
		ygot.NewNodePath(
			[]string{"state", "foo"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.A_B_C_D_E_F_G_H_I_J_K_L_M).Foo
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.A_B_C_D_E_F_G_H_I_J_K_L_M) },
		&ytypes.Schema{
			Root:       &oc.Root{},
			SchemaTree: oc.SchemaTree,
			Unmarshal:  oc.Unmarshal,
		},
	)
}
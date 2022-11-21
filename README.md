# ygnmi
![Build Status](https://github.com/openconfig/ygnmi/workflows/Go/badge.svg?branch=main)
[![Coverage Status](https://coveralls.io/repos/github/openconfig/ygnmi/badge.svg?branch=main)](https://coveralls.io/github/openconfig/ygnmi?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/openconfig/ygnmi.svg)](https://pkg.go.dev/github.com/openconfig/ygnmi)
## Introduction

ygnmi is a A Go gNMI client library based on [ygot](github.com/openconfig/ygot)-generated code. It includes a generator whose input is a set of YANG modules and output is ygot Go structs and a path library that can be used for making gNMI queries.

The library supports querying telemetry and unmarshaling it into generated structs and setting config. Only gnmi.Subscribe and gnmi.Set RPC are supported by this library.

This project is under active development and does not provide any compatibility or stability guarantees.

Note: This is not an official Google product.

## Generation

The ygnmi generator can be installed by running: `go install github.com/openconfig/ygnmi/app/ygnmi@latest`.

For the most up-to-date documentation of the generator commands and flags: use the built-in `help` command. ygnmi can be configured using (in order of precedence): flags, environment variables, or a config file. 
An example generation script is located at internal/exampleoc/gen.sh.

Not all ygot generator flags are supported by ygnmi. Notably ygnmi makes two important assumptions about the generated code:

1. Path compression is enabled.
2. PreferOperationState is selected.

Note: the supported flags may evolve over time to include these options.

### Output

Calling the generation with `--base_import_path=<somepath>/exampleoc` flag will output:

* exampleoc
    * This package contains the structs, enums, unions, and schema.
    * These correspond to **values** that can be returned or set.
* exampleoc/\<module\>
    * For every YANG module (that defines at least one container) one Go package is generated.
    * Each package contains PathStructs: structs that represent a gNMI **path** that can queried or set.
    * Each PathStruct has a State() method that returns Query for path. It may also have a Config() method. 
* exampleoc/root
    * This package contains a special "fakeroot" struct.
        * It is called the fakeroot because there is no YANG container that corresponds to this struct.
    * The package also contains a batch struct.

## gNMI Client Library

The ygnmi client library uses the generated code to perform schema compliant subscriptions and set gNMI RPCs. 

### gNMI path to ygnmi path

ygnmi paths mimic with gNMI paths with a few transformations applied.

1. Names are CamelCased: `network-instance` -> `NetworkInstance`
2. YANG module names are omitted, and added to a root struct: `/openconfig-network-instance/network-interfaces/` -> `ocpath.Root().NetworkInstance()`
3. Lists are compressed: `network-instances/network-instance[name=DEFAULT]` -> `NetworkInstance("DEFAULT")`
4. List keys can be specified several ways:
    1. Fully by specifying all keys: `protocols/protocol[identifier=BGP][name=test]` -> `Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "test")`
    2. Specifying no keys: `protocols/protocol[identifier=*][name=*]` -> `ProtocolAny()`
    3. Specifying some keys: `protocols/protocol[identifier=BGP][name=*]` -> `ProtocolAny().WithIdentifier(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP)`
5. State or Config are specified at the end: `interface[name=eth0]/state/name` -> `Interface("eth0").Name().State()`

Examples:

|gNMI Path|ygnmi Call|
|---------|----------|
|`/network-instances/network-instance[name=DEFAULT]/protocols/protocol[identifier=BGP][name=test]/bgp/neighbors/neighbor[neighbor-address=localhost]/state/description`|`ocpath.Root().NetworkInstance("DEFAULT").Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "test").Bgp().Neighbor("localhost").Description().State()`
| `/interfaces/interface[name=*]/config/name` | `ocpath.Root().InterfaceAny().Name().Config()`

Note: It is highly recommended to use this library with an IDE or autocomplete configured.

### Queries

The ygnmi library uses generic queries to represent a gNMI path, the value type, and schema. Queries should never be constructed directly.
Instead, they are returned by calling .Config() or .State() on the generated code. There are several query types that allow type safety when running an operation.
The relationship of the query types is:

![Query Diagram](doc/queries.svg)

* Singleton: Lookup, Get, Watch, Await, Collect
* Config: Update, Replace, Delete, BatchUpdate, BatchReplace, BatchDelete
* Wildcard: LookupAll, GetAll, WatchAll, CollectAll

## Additional Reference

* See [ygot](github.com/openconfig/ygot) for more information on how YANG is mapped to Go code.
* See [gNMI](github.com/openconfig/gnmi) and [gNMI Reference](https://github.com/openconfig/reference/tree/master/rpc/gnmi) for more information on the gNMI protocol and spec.
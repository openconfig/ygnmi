# ygnmi
## Introduction

ygnmi is a set of tools used to:

* A gNMI client library compatible with [ygot](github.com/openconfig/ygot) generated Go code.
* A generator that creates Go structs and helpers for a set of YANG modules.

This project is under active development and does not provide any compatibility or stability guarantees.

Note: This is not an official Google product.

## Generation

The ygnmi generator can be installed by running: `go install github.com/openconfig/ygnmi/app/ygnmi@latest`.

For the most up-to-date documentation of the generator commands and flags: use the built-in `help` command.
An example generation script is located at internal/exampleoc/gen.sh.

Not all configuration options are supported by ygnmi. Notably: ygnmi makes to two important assumptions about the generated code:

1. Path compression is enabled.
2. PreferOperationState is selected.

## gNMI Client Library

ygnmi is also a client library for gNMI. The library supports querying telemetry and unmarshaling it into generated structs and setting config. Only gnmi.Subscribe and gnmi.Set RPC are supported by this library.
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
Package exampleoc is a generated package which contains definitions
of structs which represent a YANG schema. The generated schema can be
compressed by a series of transformations (compression was true
in this case).

This package was generated by ygnmi version: (devel): (ygot: v0.23.0)
using the following YANG input files:
  - ../../pathgen/testdata/yang/openconfig-simple.yang
  - ../../pathgen/testdata/yang/openconfig-withlistval.yang
  - ../../pathgen/testdata/yang/openconfig-nested.yang

Imported modules were sourced from:
*/
package exampleoc

var (
	// ySchema is a byte slice contain a gzip compressed representation of the
	// YANG schema from which the Go code was generated. When uncompressed the
	// contents of the byte slice is a JSON document containing an object, keyed
	// on the name of the generated struct, and containing the JSON marshalled
	// contents of a goyang yang.Entry struct, which defines the schema for the
	// fields within the struct.
	ySchema = []byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x5d, 0x5d, 0x73, 0x9b, 0xba,
		0x16, 0x7d, 0xcf, 0xaf, 0xf0, 0xf0, 0xdc, 0x14, 0x9b, 0x38, 0x71, 0x92, 0xb7, 0x34, 0x6d, 0x6f,
		0xef, 0x4d, 0xd3, 0xdb, 0x69, 0xcf, 0x9c, 0x97, 0xd3, 0x33, 0x0c, 0xc6, 0xb2, 0xad, 0x06, 0x43,
		0x06, 0x70, 0xd3, 0xcc, 0x19, 0xfe, 0xfb, 0x19, 0x83, 0x21, 0xf1, 0x17, 0xe8, 0x5b, 0xc2, 0xec,
		0xc7, 0xc4, 0x08, 0x84, 0xf6, 0x5a, 0x6b, 0x4b, 0x5b, 0x5b, 0x9b, 0x7f, 0x4e, 0x7a, 0xbd, 0x5e,
		0xcf, 0xfa, 0xe2, 0x2d, 0x90, 0x75, 0xdd, 0xb3, 0xe2, 0x28, 0x4a, 0xad, 0x37, 0xc5, 0xff, 0xee,
		0x70, 0x38, 0xb1, 0xae, 0x7b, 0x83, 0xf5, 0x9f, 0xb7, 0x51, 0x38, 0xc5, 0x33, 0xeb, 0xba, 0xd7,
		0x5f, 0xff, 0xe3, 0x3d, 0x8e, 0xad, 0xeb, 0x5e, 0x71, 0x83, 0xfc, 0x1f, 0xde, 0xc6, 0x9f, 0x1b,
		0xf7, 0xf5, 0xd6, 0x37, 0xad, 0x7e, 0xd8, 0xbc, 0x79, 0xf5, 0xef, 0xed, 0x87, 0x54, 0x3f, 0x7c,
		0x8d, 0xd1, 0x14, 0xff, 0xde, 0x79, 0xc0, 0xc6, 0x43, 0x22, 0x3f, 0xdc, 0x7a, 0x4c, 0xfe, 0xf3,
		0xf7, 0x68, 0x19, 0xfb, 0x68, 0x6f, 0xd3, 0xa2, 0x2b, 0xe8, 0xf9, 0x29, 0x8a, 0x57, 0xbd, 0xb1,
		0x1e, 0x8b, 0xa7, 0xbc, 0xd9, 0x7f, 0xe1, 0x27, 0x2f, 0xb9, 0x89, 0x67, 0xcb, 0x05, 0x0a, 0x53,
		0xeb, 0xba, 0x97, 0xc6, 0x4b, 0x74, 0xe0, 0xc2, 0x57, 0x57, 0xe5, 0x9d, 0xda, 0xb9, 0x2a, 0xdb,
		0xf8, 0x4f, 0xb6, 0xf5, 0xae, 0xdb, 0x03, 0x5b, 0xfd, 0x30, 0x3e, 0xfc, 0x12, 0xe5, 0x18, 0x8c,
		0x0f, 0x75, 0x7e, 0xff, 0x80, 0x37, 0x0e, 0x3c, 0x89, 0x01, 0x08, 0x0d, 0x41, 0x6a, 0x10, 0x6a,
		0xc3, 0x50, 0x1b, 0x88, 0xdc, 0x50, 0xfb, 0x0d, 0x76, 0xc0, 0x70, 0x8d, 0x06, 0xac, 0x2e, 0xf0,
		0x9b, 0x5f, 0xbe, 0x1c, 0x4b, 0xbf, 0xe9, 0xa5, 0xeb, 0x0d, 0x4b, 0x6c, 0x60, 0x1a, 0x43, 0x53,
		0x1a, 0x9c, 0xd6, 0xf0, 0xcc, 0x00, 0x60, 0x06, 0x02, 0x3d, 0x20, 0xea, 0x81, 0xd1, 0x00, 0x10,
		0x62, 0xa0, 0x54, 0x17, 0x4e, 0xc8, 0x07, 0xad, 0xb4, 0xc9, 0x84, 0x74, 0xb0, 0xc8, 0x00, 0x44,
		0x0d, 0x24, 0x16, 0x40, 0x31, 0x02, 0x8b, 0x15, 0x60, 0xdc, 0x40, 0xe3, 0x06, 0x1c, 0x3b, 0xf0,
		0xc8, 0x00, 0x48, 0x08, 0x44, 0x6a, 0x40, 0x56, 0x0d, 0x18, 0x06, 0xbb, 0xb4, 0x2d, 0xa2, 0x1d,
		0x64, 0x3a, 0xa0, 0x32, 0x03, 0x96, 0x07, 0xb8, 0x9c, 0x00, 0xe6, 0x05, 0xb2, 0x30, 0x40, 0x0b,
		0x03, 0x36, 0x3f, 0xc0, 0xe9, 0x80, 0x4e, 0x09, 0x78, 0x66, 0xe0, 0x57, 0x0d, 0xa7, 0xec, 0x46,
		0x2a, 0x31, 0x32, 0x65, 0x35, 0x0e, 0x1b, 0x21, 0xb8, 0x89, 0x21, 0x82, 0x20, 0x82, 0x88, 0x22,
		0x8a, 0x30, 0xc2, 0x89, 0x23, 0x9c, 0x40, 0xe2, 0x88, 0xc4, 0x46, 0x28, 0x46, 0x62, 0x71, 0x13,
		0xac, 0xba, 0xc1, 0x8c, 0xdf, 0xb8, 0x25, 0xd6, 0x66, 0xbc, 0x46, 0xe5, 0x23, 0x9e, 0x30, 0x02,
		0x8a, 0x24, 0xa2, 0x60, 0x42, 0x8a, 0x26, 0xa6, 0x34, 0x82, 0x4a, 0x23, 0xaa, 0x78, 0xc2, 0xf2,
		0x11, 0x97, 0x93, 0xc0, 0xc2, 0x88, 0x5c, 0xdd, 0x68, 0x2e, 0x0e, 0x14, 0x25, 0x66, 0xe7, 0xa2,
		0xc0, 0x20, 0x86, 0xe0, 0xc2, 0x89, 0x2e, 0x83, 0xf0, 0x92, 0x88, 0x2f, 0x4b, 0x00, 0xa4, 0x0b,
		0x81, 0x74, 0x41, 0x90, 0x27, 0x0c, 0x62, 0x04, 0x42, 0x90, 0x50, 0x08, 0x17, 0x8c, 0xea, 0x86,
		0x58, 0x3c, 0x98, 0x4a, 0xec, 0x63, 0xd1, 0x20, 0x12, 0x2b, 0x24, 0xd2, 0x04, 0x45, 0xa6, 0xb0,
		0x48, 0x16, 0x18, 0xd9, 0x42, 0xa3, 0x4c, 0x70, 0x94, 0x09, 0x8f, 0x7c, 0x01, 0x12, 0x2b, 0x44,
		0x82, 0x05, 0x49, 0x9a, 0x30, 0x55, 0x37, 0xfe, 0x29, 0x0f, 0x84, 0x25, 0x87, 0x7e, 0xca, 0x02,
		0x9f, 0x1c, 0xc1, 0x92, 0x2e, 0x5c, 0x2a, 0x04, 0x4c, 0x91, 0x90, 0xa9, 0x12, 0x34, 0xe5, 0xc2,
		0xa6, 0x5c, 0xe0, 0xd4, 0x09, 0x9d, 0x1c, 0xc1, 0x93, 0x24, 0x7c, 0xd2, 0x05, 0xb0, 0x7a, 0xc0,
		0x83, 0x7c, 0xf0, 0x96, 0x5c, 0x7c, 0x90, 0x0d, 0x5a, 0xb9, 0xc2, 0xa8, 0x4c, 0x20, 0x55, 0x0a,
		0xa5, 0x62, 0xc1, 0x54, 0x2d, 0x9c, 0xda, 0x04, 0x54, 0x9b, 0x90, 0xaa, 0x17, 0x54, 0xb9, 0xc2,
		0x2a, 0x59, 0x60, 0x95, 0x09, 0x6d, 0xf5, 0xa0, 0x40, 0x1d, 0xe8, 0x4b, 0x4e, 0x07, 0xaa, 0xc0,
		0xae, 0x46, 0x80, 0x95, 0x0b, 0xb1, 0x0e, 0x41, 0xd6, 0x24, 0xcc, 0xba, 0x04, 0x5a, 0xbb, 0x50,
		0x6b, 0x17, 0x6c, 0x7d, 0xc2, 0xad, 0x46, 0xc0, 0x15, 0x09, 0xb9, 0x72, 0x41, 0xaf, 0x1e, 0xb8,
		0x50, 0x4f, 0x96, 0x52, 0x1b, 0x16, 0xaa, 0x49, 0xa2, 0x56, 0xe8, 0xb5, 0x09, 0xbe, 0x4e, 0xe1,
		0xd7, 0xec, 0x00, 0x74, 0x3b, 0x02, 0x63, 0x1c, 0x82, 0x31, 0x8e, 0x41, 0xbf, 0x83, 0x50, 0xeb,
		0x28, 0x14, 0x3b, 0x0c, 0x6d, 0x8e, 0xa3, 0x7a, 0x70, 0x92, 0x7a, 0xa9, 0x46, 0xa2, 0x95, 0x3a,
		0x53, 0x74, 0x43, 0x13, 0xb6, 0xf5, 0x38, 0x96, 0x5d, 0x07, 0xe3, 0x68, 0xea, 0x80, 0x46, 0x47,
		0x63, 0x88, 0xc3, 0x31, 0xc5, 0xf1, 0x18, 0xe7, 0x80, 0x8c, 0x73, 0x44, 0xe6, 0x38, 0x24, 0x3d,
		0x8e, 0x49, 0x93, 0x83, 0xd2, 0xee, 0xa8, 0xaa, 0x0e, 0x4c, 0xa3, 0x48, 0x3f, 0x3d, 0xab, 0x14,
		0xed, 0x28, 0xd2, 0x4d, 0xcc, 0xb5, 0xf3, 0xea, 0x6b, 0xee, 0x86, 0xae, 0x55, 0x92, 0x89, 0xce,
		0xcc, 0x30, 0xa7, 0x66, 0x9a, 0x73, 0x33, 0xd6, 0xc9, 0x19, 0xeb, 0xec, 0xcc, 0x73, 0x7a, 0x7a,
		0x9d, 0x9f, 0x66, 0x27, 0x58, 0x99, 0xe3, 0x8f, 0xe7, 0x47, 0x64, 0x96, 0xd2, 0x24, 0x69, 0x8c,
		0xc3, 0x99, 0x09, 0x62, 0x53, 0x2e, 0xaa, 0x2e, 0x4f, 0xba, 0x89, 0xcf, 0x6e, 0x4d, 0x0b, 0x6f,
		0xc2, 0x30, 0x4a, 0xbd, 0x14, 0x47, 0xa1, 0xde, 0xd9, 0x61, 0xe2, 0xcf, 0xd1, 0xc2, 0x7b, 0xf4,
		0xd2, 0xf9, 0x8a, 0x0d, 0x76, 0xf4, 0x88, 0x42, 0x3f, 0x9f, 0x99, 0x9c, 0x86, 0x28, 0x49, 0xd1,
		0xc4, 0xf6, 0xec, 0xb1, 0xed, 0xdb, 0x13, 0x1b, 0xd9, 0x53, 0x7b, 0x66, 0xcf, 0x6d, 0x6c, 0xff,
		0xb4, 0x1f, 0xec, 0xc0, 0x5e, 0xd8, 0x45, 0x08, 0xe2, 0xa4, 0x1b, 0x58, 0x39, 0xee, 0x68, 0x9a,
		0x66, 0x34, 0x72, 0xa0, 0x50, 0x47, 0xb0, 0x3d, 0x49, 0xe3, 0xa5, 0x9f, 0x86, 0x6b, 0x0f, 0x72,
		0xe3, 0xbe, 0x73, 0x6f, 0xdd, 0xf7, 0xee, 0x07, 0xf7, 0xa3, 0xfb, 0x1f, 0xf7, 0x93, 0xfb, 0x5f,
		0xf7, 0x7f, 0xee, 0x9d, 0xfb, 0xd9, 0xbd, 0x57, 0x4b, 0x8d, 0x0c, 0x76, 0x20, 0xdb, 0x04, 0x7d,
		0x66, 0xc8, 0xab, 0x4c, 0x2f, 0x20, 0x02, 0xba, 0x1a, 0x98, 0x67, 0x90, 0x19, 0x65, 0x12, 0x8c,
		0x19, 0xe1, 0xab, 0x22, 0x69, 0x91, 0x00, 0xb4, 0x72, 0x21, 0x9b, 0x41, 0x56, 0xb4, 0x4a, 0x48,
		0x32, 0x41, 0x51, 0xe6, 0x71, 0x83, 0x46, 0x00, 0xca, 0x81, 0x5f, 0xd6, 0xb1, 0x53, 0x47, 0x92,
		0xe1, 0xc5, 0x00, 0x2b, 0x19, 0x87, 0xf1, 0x1a, 0xc0, 0x24, 0x16, 0x4a, 0xd9, 0x91, 0x9e, 0x9c,
		0x95, 0x04, 0x15, 0x6a, 0x88, 0x88, 0x3c, 0x0e, 0x5e, 0x0b, 0x0c, 0x31, 0xb0, 0xc8, 0x5a, 0x5e,
		0x61, 0x41, 0xb0, 0xd9, 0x29, 0xcd, 0x2d, 0xa2, 0xe8, 0x47, 0x8d, 0x91, 0xf9, 0x4c, 0x9c, 0xb5,
		0xa4, 0xda, 0x8d, 0x20, 0x13, 0x52, 0x99, 0x8e, 0xa7, 0x7c, 0xd2, 0x41, 0x83, 0xb1, 0x99, 0x2b,
		0x33, 0xac, 0xba, 0x17, 0xa7, 0x39, 0x28, 0xcc, 0xc0, 0x52, 0xec, 0xed, 0xc0, 0xe0, 0xd3, 0x0d,
		0x7d, 0xa6, 0xa9, 0x62, 0x20, 0xe3, 0xd0, 0x12, 0x0f, 0x29, 0x4d, 0xf9, 0xc7, 0xbd, 0x03, 0x49,
		0x36, 0x8c, 0x99, 0xe4, 0x4a, 0x9f, 0x94, 0xc3, 0x44, 0x38, 0x3c, 0x24, 0x45, 0x57, 0xf7, 0x0c,
		0x4a, 0xfd, 0x90, 0x64, 0x82, 0xaa, 0xe0, 0x12, 0xbe, 0x32, 0xd1, 0xab, 0xd6, 0x95, 0x13, 0xde,
		0x79, 0xc1, 0xfd, 0xaf, 0x97, 0x51, 0x56, 0x61, 0x6e, 0xe8, 0x3e, 0x41, 0xb7, 0xf7, 0x15, 0xa3,
		0xde, 0xea, 0xec, 0x66, 0x57, 0x5f, 0x3a, 0xf4, 0xaa, 0x33, 0xd6, 0x22, 0x9a, 0xa0, 0xe0, 0x70,
		0x69, 0xed, 0xe2, 0x67, 0x25, 0xe5, 0xb5, 0x4d, 0xac, 0xae, 0x2d, 0xac, 0xb8, 0xb6, 0xd7, 0x5c,
		0x5c, 0xdb, 0x33, 0xa3, 0xb8, 0x76, 0x1b, 0x6b, 0x6b, 0x2b, 0x2b, 0xad, 0x9d, 0xe0, 0x70, 0x16,
		0xa0, 0xd3, 0x07, 0xf4, 0x4c, 0x5e, 0x63, 0xfb, 0x55, 0x9b, 0x56, 0x14, 0xdb, 0x3e, 0xc6, 0x5a,
		0xdb, 0xc6, 0x94, 0xda, 0xf6, 0x4b, 0xdb, 0x51, 0xd6, 0xdb, 0x5e, 0xb7, 0x3b, 0x8a, 0xa2, 0xdb,
		0x5d, 0xac, 0xb9, 0xdd, 0x9a, 0x92, 0xdb, 0x24, 0xd2, 0x76, 0xd0, 0xb6, 0xcd, 0x1a, 0x77, 0x08,
		0xaa, 0xfd, 0x76, 0x95, 0xdd, 0x86, 0xaa, 0xdb, 0x3c, 0x10, 0x57, 0xb3, 0x2c, 0x67, 0x4e, 0xab,
		0xe3, 0x4f, 0x83, 0x63, 0x4c, 0x5b, 0xa3, 0x58, 0x65, 0x53, 0x48, 0xe8, 0x2f, 0x2f, 0x58, 0x72,
		0x14, 0xd2, 0x2f, 0x9a, 0x03, 0xab, 0x81, 0xd5, 0x47, 0xc2, 0x6a, 0x1c, 0xa6, 0x17, 0x43, 0x0e,
		0x52, 0x0f, 0x19, 0x9a, 0x7e, 0xf3, 0xc2, 0xd9, 0xea, 0xe1, 0x7f, 0x31, 0x99, 0x84, 0x23, 0xc4,
		0x7b, 0x8f, 0xf9, 0xc3, 0xfc, 0xd6, 0x9f, 0x6b, 0x05, 0xb9, 0x72, 0x9c, 0xb3, 0xb3, 0x91, 0xd3,
		0x3f, 0xbb, 0xb8, 0x3c, 0x1f, 0x8e, 0x46, 0xe7, 0x97, 0xfd, 0x4b, 0xce, 0xcd, 0x88, 0x8f, 0xb1,
		0xe7, 0xa7, 0x38, 0x0a, 0xdf, 0xe3, 0x19, 0x4e, 0x13, 0x01, 0xb5, 0xca, 0xbf, 0xa0, 0x99, 0x97,
		0xe2, 0x5f, 0x68, 0x4d, 0x01, 0xf6, 0x80, 0x3e, 0x47, 0xb8, 0xfb, 0xde, 0xfb, 0x2d, 0x73, 0xcc,
		0x47, 0x06, 0x8f, 0xf9, 0xd4, 0x0b, 0x12, 0xa4, 0x7a, 0x17, 0x85, 0xba, 0xd5, 0xdf, 0x1d, 0x8d,
		0x7d, 0x3f, 0xe1, 0x74, 0x1e, 0xe0, 0x24, 0xfd, 0xe5, 0x05, 0x76, 0x1e, 0xd8, 0xb3, 0x3d, 0xfb,
		0x25, 0x3a, 0x61, 0xaf, 0xd7, 0x98, 0xa2, 0x22, 0xda, 0x04, 0x61, 0x04, 0x9a, 0xc5, 0x06, 0xc3,
		0x22, 0x83, 0x72, 0x1a, 0x02, 0xeb, 0xe0, 0xe3, 0x5d, 0x07, 0x53, 0x4f, 0x1b, 0x5e, 0x2a, 0x90,
		0x21, 0x6f, 0x1a, 0x23, 0x9a, 0xfd, 0xcf, 0x6a, 0xf6, 0x4f, 0xa1, 0xd4, 0xd6, 0xd7, 0x35, 0x5f,
		0xdf, 0xbe, 0x2d, 0x8e, 0x10, 0xd8, 0x2b, 0x98, 0x2b, 0xa4, 0x22, 0x5d, 0x05, 0x07, 0xa6, 0x8a,
		0x0b, 0xcc, 0x61, 0x29, 0x07, 0xe8, 0x08, 0x61, 0x29, 0x08, 0x4b, 0xc1, 0x02, 0x16, 0x16, 0xb0,
		0x5a, 0xc2, 0x52, 0x8a, 0xf3, 0x61, 0x9e, 0x67, 0x51, 0x7a, 0x1a, 0xf9, 0xa7, 0x7e, 0xb4, 0x78,
		0x8c, 0x51, 0x92, 0xa0, 0xc9, 0xe9, 0xca, 0x09, 0xaf, 0x6e, 0x96, 0x41, 0x3c, 0x0d, 0xe4, 0x08,
		0xe4, 0x08, 0xe2, 0x69, 0x10, 0x4f, 0x83, 0x78, 0x1a, 0xc4, 0xd3, 0xda, 0x19, 0x4f, 0xeb, 0xc0,
		0x74, 0xa2, 0xdd, 0x81, 0x40, 0x8a, 0x1a, 0x06, 0xb2, 0x33, 0x5b, 0xef, 0xf2, 0x75, 0x1f, 0xc1,
		0x0a, 0xce, 0xfa, 0x8c, 0x93, 0xf4, 0x26, 0x4d, 0x09, 0xf3, 0x70, 0xee, 0x71, 0xf8, 0x21, 0x40,
		0x2b, 0xb7, 0x4c, 0xc8, 0xd8, 0x95, 0x0a, 0xbd, 0x6a, 0x31, 0xb8, 0x1c, 0x0e, 0x2f, 0x46, 0xc3,
		0x61, 0x7f, 0x74, 0x36, 0xea, 0x5f, 0x9d, 0x9f, 0x0f, 0x2e, 0x06, 0xe7, 0x04, 0x37, 0xf9, 0x7f,
		0x3c, 0x41, 0x31, 0x9a, 0xbc, 0x5b, 0xbd, 0x54, 0xb8, 0x0c, 0x02, 0x43, 0xb2, 0x7e, 0xeb, 0xf1,
		0x40, 0x9f, 0x07, 0x7c, 0xbf, 0xba, 0x87, 0xfb, 0x3d, 0xbf, 0xc3, 0x5d, 0x53, 0x1c, 0xcb, 0x9c,
		0x8c, 0xe0, 0x3d, 0xc3, 0x40, 0x9c, 0xed, 0xbb, 0x27, 0x79, 0x75, 0xdc, 0x9c, 0xf3, 0x39, 0x86,
		0x9c, 0x4f, 0xe3, 0x73, 0x3e, 0x17, 0xcb, 0x20, 0xc5, 0x74, 0x29, 0x9f, 0x2f, 0x4d, 0x20, 0xe3,
		0x13, 0x32, 0x3e, 0x21, 0xe3, 0x13, 0x42, 0xeb, 0xea, 0x66, 0x8b, 0x2c, 0xa1, 0xf5, 0x01, 0x57,
		0x6c, 0x7d, 0x00, 0xd1, 0x2c, 0x88, 0x66, 0x1d, 0x4b, 0x34, 0x6b, 0x89, 0xc3, 0xf4, 0xcc, 0xe1,
		0x08, 0x67, 0x8d, 0x3a, 0x1b, 0xce, 0xea, 0x1f, 0x6f, 0x20, 0xc5, 0x94, 0xe8, 0xd5, 0xd0, 0xb9,
		0x1a, 0x5e, 0x5d, 0x8c, 0x9c, 0xab, 0x73, 0x08, 0x5a, 0xa9, 0x0d, 0x5a, 0x49, 0xd9, 0x4a, 0x7a,
		0x40, 0xcf, 0x0e, 0x97, 0xeb, 0x75, 0xc0, 0xf5, 0x82, 0xeb, 0x3d, 0x26, 0xd7, 0xcb, 0xb5, 0x93,
		0x74, 0x09, 0xae, 0x17, 0x5c, 0xaf, 0x2c, 0xd7, 0xcb, 0x16, 0xf1, 0x05, 0x27, 0xac, 0xc7, 0x09,
		0xb7, 0x6e, 0x03, 0x66, 0x6c, 0x57, 0x41, 0x43, 0x3d, 0x89, 0xd8, 0x03, 0xa6, 0x4c, 0xec, 0x01,
		0xa4, 0x62, 0x43, 0x80, 0x8a, 0x8e, 0x4d, 0xad, 0x4c, 0xc5, 0x1e, 0x28, 0x66, 0xa3, 0xc3, 0xc4,
		0x46, 0x07, 0xd8, 0x08, 0x6c, 0xec, 0x00, 0x1b, 0x1d, 0x38, 0x19, 0x01, 0x27, 0x23, 0x8e, 0x93,
		0x8f, 0xb0, 0x7d, 0x03, 0x31, 0x24, 0x88, 0x21, 0x71, 0xc5, 0x90, 0x60, 0xfb, 0x06, 0x62, 0x48,
		0x66, 0xc6, 0x90, 0x60, 0xfb, 0x46, 0x57, 0xe4, 0x08, 0x8e, 0x30, 0x31, 0x2f, 0x30, 0x79, 0x16,
		0x9a, 0x30, 0x67, 0x80, 0x39, 0x03, 0xec, 0x3b, 0xc1, 0x9c, 0x01, 0xe6, 0x0c, 0xcc, 0x43, 0x0c,
		0xfb, 0x4e, 0x30, 0x7b, 0x80, 0x13, 0x4b, 0xf5, 0x1b, 0x66, 0x46, 0x1e, 0x58, 0x1a, 0xfc, 0x08,
		0x7f, 0x84, 0x3d, 0x82, 0xc9, 0x12, 0x1c, 0x5d, 0x92, 0x05, 0x0c, 0xd6, 0x93, 0x4b, 0xf7, 0xab,
		0x1b, 0xb4, 0xfa, 0xe0, 0xd2, 0x58, 0xf7, 0x67, 0x0a, 0x76, 0xba, 0xd4, 0xfc, 0xbd, 0x82, 0x7c,
		0xec, 0x49, 0xbe, 0x59, 0xf0, 0xe8, 0xc5, 0xc5, 0x74, 0xf4, 0xc0, 0x47, 0x0b, 0xd6, 0xbf, 0x2b,
		0xf9, 0x6a, 0x41, 0x62, 0xe2, 0x67, 0x0b, 0x12, 0x61, 0xdf, 0x2d, 0xf0, 0xe7, 0x38, 0x98, 0x34,
		0x9f, 0x63, 0x2b, 0x2e, 0x33, 0xe2, 0x2c, 0x5b, 0xd2, 0xc6, 0xc3, 0x6c, 0x89, 0xb2, 0xd3, 0x6c,
		0x84, 0xc7, 0x91, 0xe8, 0x8e, 0x21, 0x19, 0x72, 0x8e, 0x2d, 0x39, 0xc6, 0x83, 0x6c, 0x89, 0x31,
		0x27, 0xd9, 0xa6, 0xd1, 0x32, 0xa6, 0xdf, 0x07, 0xcd, 0x5b, 0x1d, 0x45, 0x5a, 0x42, 0xd2, 0xc5,
		0x7d, 0xd0, 0xa4, 0x7d, 0x89, 0x09, 0x63, 0x1c, 0x7a, 0xf1, 0x33, 0x43, 0x5e, 0xc2, 0x95, 0xc2,
		0xac, 0x82, 0x28, 0x64, 0xc8, 0x29, 0x58, 0x35, 0x02, 0x2a, 0x01, 0x95, 0x94, 0x51, 0x89, 0xba,
		0xd4, 0x1c, 0x65, 0x89, 0x39, 0x31, 0x5c, 0x4a, 0xe7, 0x31, 0x62, 0x60, 0x53, 0xd1, 0x0c, 0xf8,
		0x04, 0x7c, 0x52, 0xc6, 0x27, 0x14, 0x2e, 0x17, 0x28, 0x2e, 0x16, 0xd4, 0x0c, 0xa4, 0xa2, 0x28,
		0x95, 0x66, 0x7d, 0x08, 0x97, 0x0b, 0xf2, 0x88, 0x61, 0x7b, 0xbe, 0x50, 0x99, 0xe0, 0xc5, 0x63,
		0x80, 0xec, 0x62, 0x8d, 0x6f, 0xe7, 0x0b, 0x4f, 0xa2, 0x84, 0xf9, 0x9a, 0x68, 0x4d, 0xed, 0xe7,
		0x1f, 0x49, 0x92, 0xff, 0xa8, 0x92, 0xfe, 0xa8, 0x17, 0x4b, 0x0e, 0x2c, 0x96, 0x60, 0xb1, 0x04,
		0x1e, 0x09, 0x3c, 0x52, 0x1b, 0x17, 0x4b, 0x92, 0x77, 0x8c, 0xb8, 0xb7, 0xca, 0x60, 0x35, 0x07,
		0x5c, 0x87, 0xd5, 0x9c, 0xa0, 0xd5, 0x5c, 0x37, 0xc8, 0x0e, 0xcb, 0x4d, 0x20, 0x3c, 0x2c, 0x37,
		0xf9, 0x96, 0x9b, 0x1d, 0x51, 0x8a, 0xa7, 0x88, 0x41, 0x27, 0x9e, 0x22, 0x50, 0x09, 0x50, 0x89,
		0x4e, 0x05, 0x79, 0x0d, 0x8e, 0x2d, 0x11, 0xe4, 0x96, 0x99, 0x93, 0x08, 0xb4, 0xe7, 0x05, 0xac,
		0xda, 0xe8, 0xd6, 0xeb, 0xfc, 0x9b, 0xaf, 0x79, 0x2b, 0xf7, 0x36, 0x6f, 0xa5, 0x39, 0x7d, 0x68,
		0xe3, 0x45, 0x9a, 0x53, 0x87, 0x8a, 0xae, 0x93, 0xe4, 0x0e, 0xc5, 0x68, 0x11, 0xa5, 0xe8, 0xd4,
		0x8f, 0xc2, 0xd4, 0xc3, 0x21, 0x8a, 0x0f, 0x67, 0x11, 0xed, 0x5c, 0xa9, 0x24, 0x9f, 0x28, 0x36,
		0x31, 0x9f, 0x28, 0x16, 0x97, 0x4f, 0x54, 0x9f, 0x7e, 0x42, 0x96, 0x76, 0xa2, 0x38, 0xa3, 0x28,
		0x6e, 0x63, 0x46, 0x51, 0xac, 0x2c, 0xa3, 0xc8, 0xab, 0xa6, 0x56, 0x64, 0x31, 0xf2, 0xf5, 0xf5,
		0x64, 0x41, 0xf2, 0xbe, 0xde, 0x8c, 0xa2, 0xf8, 0x18, 0x83, 0xe4, 0xb1, 0xec, 0x20, 0x39, 0xf1,
		0x54, 0x84, 0x7e, 0x0a, 0x42, 0x38, 0xf5, 0x30, 0xce, 0x23, 0x6f, 0x3b, 0x93, 0xda, 0x2d, 0x2b,
		0xb2, 0xcf, 0x0b, 0xd4, 0x6f, 0x4d, 0x11, 0x6d, 0x49, 0x11, 0x0b, 0xa9, 0x03, 0x42, 0x0a, 0x42,
		0x0a, 0x42, 0xda, 0x41, 0x21, 0x15, 0xbc, 0xf6, 0x62, 0x0e, 0xcc, 0x98, 0xaf, 0xe8, 0x35, 0x0b,
		0x45, 0xe5, 0xeb, 0xa6, 0x86, 0xb5, 0xcb, 0x9e, 0x15, 0xd4, 0xb7, 0xbc, 0xc5, 0x6d, 0xd5, 0xe0,
		0xd0, 0x52, 0xea, 0xe4, 0x55, 0x7f, 0x0f, 0xf5, 0xd3, 0xc2, 0xc9, 0x6d, 0x65, 0xdf, 0xef, 0x79,
		0x5f, 0x77, 0x68, 0x66, 0xe1, 0xe4, 0xa3, 0xf7, 0x80, 0xbe, 0x45, 0xd1, 0x2e, 0x05, 0xb7, 0xdf,
		0xcf, 0x7a, 0xfd, 0xd3, 0x66, 0xaf, 0x57, 0xcd, 0x8b, 0x2e, 0x9d, 0x64, 0xff, 0x02, 0x00, 0x00,
		0xff, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0x46, 0xb3, 0xff, 0x6a, 0x6e, 0xee, 0x00, 0x00,
	}
)

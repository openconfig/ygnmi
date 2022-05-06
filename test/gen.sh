#!/bin/bash
#
# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cd "$(dirname "$0")"

go run github.com/openconfig/ygot/generator \
  -generate_structs \
  -generate_path_structs=false \
  -exclude_state=false \
  -prefer_operational_state \
  -list_builder_key_threshold=4 \
  -output_file=telemetry.go \
  -package_name=test \
  -compress_paths \
  -generate_fakeroot \
  -fakeroot_name=device \
  -ignore_shadow_schema_paths \
  -generate_simple_unions \
  -shorten_enum_leaf_names \
  -typedef_enum_with_defmod \
  -enum_suffix_for_simple_union_enums \
  -trim_enum_openconfig_prefix \
  -include_schema \
  -generate_append \
  -generate_getters \
  -generate_rename \
  -generate_delete \
  -generate_leaf_getters \
  -generate_populate_defaults \
  -path=../pathgen/testdata/yang \
  ../pathgen/testdata/yang/openconfig-simple.yang \
  ../pathgen/testdata/yang/openconfig-withlist.yang

mkdir -p device

go run ../app/ygnmi generator \
  --output_dir=device \
  --base_import_path=github.com/openconfig/ygnmi/test/device \
  --schema_struct_path=github.com/openconfig/ygnmi/test \
  ../pathgen/testdata/yang/openconfig-simple.yang \
  ../pathgen/testdata/yang/openconfig-withlist.yang

goimports -w .
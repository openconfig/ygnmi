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
set -e

cd "$(dirname "$0")"

# Codes with the function PathOriginName() returning "test-origin" are generated.
# e.g. setpathoriginpath/setpathoriginpath.go
# func (n *APath) PathOriginName() string {
#    return "test-origin"
# }
go run ../../app/ygnmi generator \
  --trim_module_prefix=openconfig \
  --prefer_operational_state=false \
  --base_package_path=github.com/openconfig/ygnmi/internal/examplesetorigin/setpathorigin \
  --output_dir=setpathorigin \
  --split_package_paths=/model/a=a,/model/b \
  --split_top_level_packages=false \
  --set_path_origin=test-origin \
  ../../pathgen/testdata/yang/openconfig-simple.yang \
  ../../pathgen/testdata/yang/openconfig-withlistval.yang \
  ../../pathgen/testdata/yang/openconfig-nested.yang

# Codes with the function PathOriginName() returning the YANG module name are generated.
# e.g. usemodulenameaspathoriginoriginpath/usemodulenameaspathoriginpath.go
# func (n *APath) PathOriginName() string {
#    return "openconfig-nested"
# }
go run ../../app/ygnmi generator \
  --trim_module_prefix=openconfig \
  --prefer_operational_state=false \
  --base_package_path=github.com/openconfig/ygnmi/internal/examplesetorigin/usemodulenameaspathorigin \
  --output_dir=usemodulenameaspathorigin  \
  --split_package_paths=/model/a=a,/model/b \
  --split_top_level_packages=false \
  --use_module_name_as_path_origin \
  ../../pathgen/testdata/yang/openconfig-simple.yang \
  ../../pathgen/testdata/yang/openconfig-withlistval.yang \
  ../../pathgen/testdata/yang/openconfig-nested.yang

go install golang.org/x/tools/cmd/goimports@latest
goimports -w .

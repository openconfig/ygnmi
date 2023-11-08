#!/bin/bash
# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SRC_TEST_NAME="ygnmi_test.go"
TARGET_TEST_NAME="ygnmi_preferconfig_test.go"

printf "// Code generated from ygnmi_test.go. DO NOT EDIT.\n\n" > "${TARGET_TEST_NAME}"
cat "${SRC_TEST_NAME}" >> "${TARGET_TEST_NAME}"
sed -i -e 's/exampleoc/exampleocconfig/g' "${TARGET_TEST_NAME}"
sed -i -e 's/github.com\/openconfig\/ygnmi\/exampleocconfig/github.com\/openconfig\/ygnmi\/internal\/exampleocconfig/g' "${TARGET_TEST_NAME}"
sed -i -e 's/^func Test/func TestPreferConfig/' "${TARGET_TEST_NAME}"
sed -i -e 's/getSample/getSamplePreferConfig/' "${TARGET_TEST_NAME}"
gofmt -w -s "${TARGET_TEST_NAME}"

#!/usr/bin/env bash
# Copyright 2026 The Karmada Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script regenerates the command-line flag reference files under
# docs/command-flags/ by running the extract-flags tool against every
# Karmada component.
#
# Usage: hack/update-command-flags.sh

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${REPO_ROOT}"/hack/util.sh

util::verify_go_version

cd "${REPO_ROOT}"

# Remove stale files so that renamed/removed commands don't leave orphans.
rm -rf docs/command-flags
mkdir -p docs/command-flags

echo "Generating command flag reference files..."
go run hack/tools/extract-flags/main.go docs/command-flags
echo "Done."


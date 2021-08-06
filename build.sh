#!/bin/bash

#
# Copyright © 2021 peizhaoyou <peizhaoyou@4paradigm.com>
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
#

set -e

function usage() {
  echo "Usage: ${0} { oapi-codegen billing-oapi-codegen billing-client-codegen webterminal-oapi-codegen }"
  exit 1
}

if [ "$#" -lt "1" ]; then
  usage
fi


pushd "$(dirname "${0}")" > /dev/null || return

command=$1; shift

case $command in
  oapi-codegen)
    # for go1.16 or above, using go install
    # ref: https://github.com/golang/go/issues/40276
    go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.5
    mkdir -p ./src/pineapple/apigen
    mkdir -p ./src/pineapple/apigen/internalapigen
    oapi-codegen -include-tags finished -generate types,server -package apigen ./doc/api/main.yaml > ./src/pineapple/apigen/serverapi.gen.go
    oapi-codegen -include-tags finished -generate types,server -package internalapigen ./doc/api/internal-api.yaml > ./src/pineapple/apigen/internalapigen/internalapi.gen.go
  ;;

  run-local-dev)
    mkdir -p ./build
    ( cd src && go build -o ../build/pineapple ./pineapple )
    args=$(tr "\\n" " " < ./configs/local-dev.conf)
    env -vS "${args[@]}" ./build/pineapple --storage-root "${HOME}/shared" --debug
  ;;

  billing-oapi-codegen)
    go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.5
    mkdir -p ./src/billing/apigen
    oapi-codegen -generate types,server -package apigen ./doc/api/billing.yaml > ./src/billing/apigen/billingapi.gen.go
  ;;

  billing-client-codegen)
    go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.5
    mkdir -p ./src/internal/billingclient/apigen
    oapi-codegen -generate types,client -package apigen ./doc/api/billing.yaml > ./src/internal/billingclient/apigen/billingclient.gen.go
  ;;

  webterminal-oapi-codegen)
    go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.5
    mkdir -p ./src/webterminal/server/apigen
    oapi-codegen -include-tags finished -generate types,server -package apigen ./doc/api/webterminal.yaml > ./src/webterminal/server/apigen/serverapi.gen.go
  ;;

  billing-run-local-dev)
    mkdir -p ./build
    ( cd src && go build -o ../build/billing ./billing )
    args=$(tr "\\n" " " < ./configs/billing-local-dev.conf)
    env -vS "${args[@]}" ./build/billing
  ;;
  all)
  ;;
  *)
    usage
  ;;
esac

popd > /dev/null || return


#!/bin/bash

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
  *)
    usage
  ;;
esac

popd > /dev/null || return


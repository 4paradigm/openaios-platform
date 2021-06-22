#!/bin/bash

# go setup
export GO111MODULE=on
export GOPROXY=http://goproxy.cn,direct
export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"

# gen openapi 
go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.5
mkdir -p ./test/openapi/billing/apigen/billingclient
mkdir -p ./test/openapi/main/apigen/restclient
oapi-codegen -generate types,client -package billingclient ./doc/api/billing.yaml > ./test/openapi/billing/apigen/billingclient/billingclient.gen.go
oapi-codegen -generate types,client -package restclient ./doc/api/main.yaml > ./test/openapi/main/apigen/restclient/restclient.gen.go

# kubectl port-forward
kubectl config view
kubectl port-forward -n pineapple svc/pineapple-pineapple-billing 4321:80 > kubectl_forward.log 2>&1 &

# run ginkgo
go get -u github.com/onsi/ginkgo/ginkgo
(cd test/system && ginkgo -v)
(cd test/integration && ginkgo -v)

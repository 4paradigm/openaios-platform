GO ?= go
GO111MODULE = on
GOPROXY ?= https://goproxy.cn,direct
CGO_ENABLED = 0
OUTPUT_DIR ?= bin
VERSION ?= unknown
BUILDARGS ?= -ldflags '-s -w -X github.com/4paradigm/openaios-platform/src/internal/version.version=$(VERSION)'

all: pineapple billing webhook web-terminal gotty

tidy: deps
	cd src && $(GO) mod tidy

oapi_codegen:
	$(GO) get 'github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.5'

apigen: oapi_codegen
	mkdir -p ./src/pineapple/apigen
	oapi-codegen -include-tags finished -generate types,server -package apigen ./doc/api/main.yaml > ./src/pineapple/apigen/serverapi.gen.go

internalapigen: oapi_codegen
	mkdir -p ./src/pineapple/apigen/internalapigen
	oapi-codegen -include-tags finished -generate types,server -package internalapigen ./doc/api/internal-api.yaml > ./src/pineapple/apigen/internalapigen/internalapi.gen.go

pineapple: apigen internalapigen billing_client_codegen tidy
	cd src && $(GO) build $(BUILDARGS) -o ../$(OUTPUT_DIR)/pineapple ./pineapple

billing_oapi_codegen: oapi_codegen
	mkdir -p ./src/billing/apigen
	oapi-codegen -generate types,server -package apigen ./doc/api/billing.yaml > ./src/billing/apigen/billingapi.gen.go

billing_client_codegen: oapi_codegen
	mkdir -p ./src/internal/billingclient/apigen
	oapi-codegen -generate types,client -package apigen ./doc/api/billing.yaml > ./src/internal/billingclient/apigen/billingclient.gen.go

billing: billing_oapi_codegen billing_client_codegen tidy
	cd src && $(GO) build $(BUILDARGS) -o ../$(OUTPUT_DIR)/billing ./billing

webhook: billing_client_codegen tidy
	cd src && $(GO) build $(BUILDARGS) -o ../$(OUTPUT_DIR)/webhook ./webhook

webterminal_oapi_codegen: oapi_codegen
	mkdir -p ./src/webterminal/server/apigen
	oapi-codegen -include-tags finished -generate types,server -package apigen ./doc/api/webterminal.yaml > ./src/webterminal/server/apigen/serverapi.gen.go

web-terminal: webterminal_oapi_codegen tidy
	cd src && $(GO) build $(BUILDARGS) -o ../$(OUTPUT_DIR)/web-terminal ./webterminal/server

gotty: webterminal_oapi_codegen tidy
	cd src && $(GO) build $(BUILDARGS) -o ../$(OUTPUT_DIR)/gotty ./webterminal/gotty


deps: apigen internalapigen billing_client_codegen billing_oapi_codegen webterminal_oapi_codegen

lint: deps
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	cd src && staticcheck ./...

vet:
	cd src && go vet ./...


clean:
	rm -rf $(OUTPUT_DIR) src/pineapple/apigen src/billing/apigen src/internal/billingclient/apigen src/webterminal/server/apigen

.PHONY: all clean

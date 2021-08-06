ARG BUILDBASE=golang:1.16.6-buster
FROM $BUILDBASE AS build

WORKDIR /build
COPY . .
ARG GOPROXY=https://goproxy.cn,direct
ARG GOMODCACHE=/go/pkg/mod
RUN cd src && go mod download -x
RUN make webhook


FROM alpine AS target
WORKDIR /root
COPY --from=build /build/bin/webhook /root/webhook

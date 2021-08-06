ARG BUILDBASE=golang:1.16.6-buster
FROM $BUILDBASE AS build

WORKDIR /build
COPY . .
ARG GOPROXY=https://goproxy.cn,direct
RUN cd src && go mod download -x
RUN make billing

#FROM broothie/redoc-cli:0.9.8 AS doc
#WORKDIR /root
#COPY ./doc/api/billing.yaml ./
#RUN redoc-cli bundle -o billing.html billing.yaml


FROM alpine AS target
WORKDIR /root
COPY --from=build /build/bin/billing /root/billing

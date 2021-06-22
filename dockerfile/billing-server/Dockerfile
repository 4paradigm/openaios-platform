FROM golang:1.14-alpine AS build
WORKDIR /src/
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk --no-cache add ca-certificates git bash
ENV GO111MODULE on
ENV GOPROXY http://goproxy.cn,direct

# fetch dependencies
RUN mkdir -p ./src
COPY ./src/go.mod ./src/go.mod
RUN cd src && go mod download -x

RUN mkdir -p ./doc
COPY ./doc/api ./doc/api
COPY ./build.sh ./
RUN ./build.sh billing-oapi-codegen

# build
COPY ./src ./src/
RUN cd src && CGO_ENABLED=0 go build -v -o /bin/billing ./billing

FROM broothie/redoc-cli:0.9.8 AS doc
WORKDIR /root
COPY ./doc/api/billing.yaml ./
RUN redoc-cli bundle -o billing.html billing.yaml

FROM alpine AS target
WORKDIR /root
COPY --from=build /bin/billing /root/billing

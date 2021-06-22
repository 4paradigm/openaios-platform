FROM golang:1.14-alpine AS build
WORKDIR /src/
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk --no-cache add ca-certificates git bash
ENV GO111MODULE on
ENV GOPROXY http://goproxy.cn,direct

RUN mkdir -p ./src
COPY ./src/go.mod ./src/go.mod
RUN cd src && go mod download -x

RUN mkdir -p ./doc
COPY ./doc/api ./doc/api
COPY ./build.sh ./
RUN ./build.sh billing-client-codegen

COPY ./src ./src/
RUN cd src && CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/webhook ./webhook

FROM alpine AS target
WORKDIR /root
COPY --from=build /bin/webhook /root/webhook

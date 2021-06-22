FROM golang:1.14-alpine AS build
WORKDIR /src/
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk --no-cache add git bash
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

# fetch dependencies
RUN mkdir -p ./src
COPY ./src/go.mod ./src/go.mod
RUN cd src && go mod download -x


#build api
RUN mkdir -p ./doc
COPY ./doc/api ./doc/api
COPY ./build.sh ./
RUN ./build.sh webterminal-oapi-codegen

# build
COPY ./src ./src
RUN cd src && CGO_ENABLED=0 go build -v -o /bin/web-terminal ./webterminal/server
RUN cd src && CGO_ENABLED=0 go build -v -o /bin/gotty ./webterminal/gotty


FROM ubuntu:18.04
USER root
RUN apt-get update \
    && apt-get install -y apt-transport-https ca-certificates curl \
    && curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg \
    && echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list \
    && apt-get update
    && apt-get install -y kubectl

WORKDIR /root

COPY --from=build /bin/gotty /root/gotty
COPY --from=build /bin/web-terminal /root/web-terminal


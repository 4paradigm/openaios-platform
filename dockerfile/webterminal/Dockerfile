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

FROM bitnami/kubectl:1.20.9

WORKDIR /root

COPY --from=build /bin/gotty /root/gotty
COPY --from=build /bin/web-terminal /root/web-terminal


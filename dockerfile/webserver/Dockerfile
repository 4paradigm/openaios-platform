FROM golang:1.14-alpine AS build
WORKDIR /src/
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk --no-cache add ca-certificates git bash
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

# fetch dependencies
RUN mkdir -p ./src
COPY ./src/go.mod ./src/go.mod
RUN cd src && go mod download -x

RUN mkdir -p ./doc
COPY ./doc/api ./doc/api
COPY ./build.sh ./
RUN ./build.sh oapi-codegen
RUN ./build.sh billing-client-codegen

# build
## COPY the content of src(sub dir of file) into container src directory
COPY ./src ./src/
RUN cd src && CGO_ENABLED=0 go build -v -o /bin/pineapple ./pineapple


FROM broothie/redoc-cli:0.9.8 AS doc
WORKDIR /root
COPY ./doc/api/main.yaml ./
RUN redoc-cli bundle -o index.html main.yaml


FROM alpine AS target
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add attr
WORKDIR /root

COPY --from=build /bin/pineapple /root/pineapple

RUN mkdir -p /root/static
COPY --from=doc /root/index.html /root/static/index.html
ENV PINEAPPLE_STATIC_DIR /root/static

COPY ./charts /root/charts
ENV PINEAPPLE_ENV_CHARTSDIR /root/charts

# ENTRYPOINT ["/root/pineapple"]

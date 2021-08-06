ARG BUILDBASE=golang:1.16.6-buster
FROM $BUILDBASE AS build

WORKDIR /build
COPY . .
ARG GOPROXY=https://goproxy.cn,direct
ARG GOMODCACHE=/go/pkg/mod
RUN cd src && go mod download -x
RUN make pineapple

FROM broothie/redoc-cli:0.9.8 AS doc
WORKDIR /root
COPY ./doc/api/main.yaml ./
RUN redoc-cli bundle -o index.html main.yaml


FROM alpine AS target
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add attr
WORKDIR /root
COPY --from=build /build/bin/pineapple /root/pineapple
COPY --from=doc /root/index.html /root/static/index.html
ENV PINEAPPLE_STATIC_DIR /root/static
COPY ./charts /root/charts
ENV PINEAPPLE_ENV_CHARTSDIR /root/charts

# ENTRYPOINT ["/root/pineapple"]

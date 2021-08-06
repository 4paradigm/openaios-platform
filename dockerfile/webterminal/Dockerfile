ARG BUILDBASE=golang:1.16.6-buster
FROM $BUILDBASE AS build

WORKDIR /build
COPY . .
ARG GOPROXY=https://goproxy.cn,direct
RUN make gotty web-terminal


FROM bitnami/kubectl:1.20.9
USER root
WORKDIR /root

COPY --from=build /build/bin/gotty /root/gotty
COPY --from=build /build/bin/web-terminal /root/web-terminal


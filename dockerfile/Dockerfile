ARG BUILDBASE=golang:1.16.6-buster
FROM $BUILDBASE AS build

WORKDIR /build
COPY . .
ARG GOPROXY=https://goproxy.cn,direct
RUN make

FROM bitnami/kubectl:1.20.9 as kubectl
WORKDIR /usr/bin

FROM alpine AS target
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add attr
WORKDIR /
COPY --from=kubectl /usr/bin/kubectl /usr/bin/kubectl
COPY --from=build /build/bin/pineapple /root/pineapple
COPY --from=build /build/bin/webhook /root/webhook
COPY --from=build /build/bin/gotty /root/gotty
COPY --from=build /build/bin/web-terminal /root/web-terminal
COPY --from=build /build/bin/billing /root/billing
COPY ./charts /root/charts
WORKDIR /root
ENV PINEAPPLE_ENV_CHARTSDIR /root/charts

ARG BUILDBASE=golang:1.16.7-buster
ARG KUBECTLBASE=bitnami/kubectl:1.20.9
ARG TARGETBASE=ubuntu:20.04

FROM $BUILDBASE AS build
WORKDIR /build
COPY . .
ARG VERSION=unknown
ARG GOPROXY=https://goproxy.cn,direct
RUN make

FROM $KUBECTLBASE as kubectl

FROM $TARGETBASE AS target
RUN apt update
RUN apt install -y attr
WORKDIR /
COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/bin/kubectl
COPY --from=build /build/bin/pineapple /root/pineapple
COPY --from=build /build/bin/webhook /root/webhook
COPY --from=build /build/bin/gotty /root/gotty
COPY --from=build /build/bin/web-terminal /root/web-terminal
COPY --from=build /build/bin/billing /root/billing
COPY ./charts /root/charts
COPY ./LICENSE /root/LICENSE
WORKDIR /root
ENV PINEAPPLE_ENV_CHARTSDIR /root/charts

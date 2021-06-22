FROM ubuntu:18.04 AS build
WORKDIR /test
ENV GO111MODULE on
ENV GOPROXY http://goproxy.cn,direct

# install kubectl
RUN curl -LO ftp://172.27.232.67/pub/pineapple/kubectl
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# install ginkgo cmd
RUN go get -u github.com/onsi/ginkgo/ginkgo

# Update the CA store
COPY pineapple.crt /etc/ssl/certs
RUN update-ca-certificates
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


FROM docker.4pd.io/adc/infras/kubectl:latest
USER root
#RUN apt-get --allow-unauthenticated update \
#    && apt-get --allow-unauthenticated install -y wget \
#    && apt-get --allow-unauthenticated update && apt-get --allow-unauthenticated install -y apt-transport-https gnupg2 curl \
#    && curl -k -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
#    && echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list \
#    && apt-get --allow-unauthenticated update \
#    && apt-get --allow-unauthenticated install -y kubectl
#RUN apt-get update && apt-get install -y apt-transport-https gnupg2 curl wget
#RUN wget -q -O apt-key.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg && apt-key add apt-key.gpg
#RUN echo "deb https://apt.kuberqnetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list
#RUN apt-get update
#RUN sudo apt-get install -y kubectl


WORKDIR /root

COPY --from=build /bin/gotty /root/gotty
COPY --from=build /bin/web-terminal /root/web-terminal


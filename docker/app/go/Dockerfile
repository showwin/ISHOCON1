ARG BASE_IMAGE=ishocon1-app-base:latest
FROM ${BASE_IMAGE}

ENV GO_VERSION=1.23.1

# Go のインストール
ARG TARGETARCH
RUN sudo wget -q https://dl.google.com/go/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz && \
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-${TARGETARCH}.tar.gz && \
    sudo rm go${GO_VERSION}.linux-${TARGETARCH}.tar.gz
ENV PATH $PATH:/usr/local/go/bin
ENV GOROOT /usr/local/go
ENV GOPATH /home/ishocon/.local/go
ENV PATH $PATH:$GOROOT/bin

# アプリケーション
COPY --chown=ishocon:ishocon webapp/ /home/ishocon/webapp
RUN cd ~/webapp/go && \
    mkdir -p /tmp/go && \
    GOARCH=${TARGETARCH} go build -x -o /tmp/go/webapp *.go

WORKDIR /home/ishocon

COPY run.sh /home/ishocon/run.sh

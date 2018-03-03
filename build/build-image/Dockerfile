FROM golang:1.10-stretch

ENV USE_DOCKER=false

WORKDIR /tmp

# Install tree
RUN apt-get update
RUN apt-get install -y tree

# Install Helm
RUN curl -sLO https://kubernetes-helm.storage.googleapis.com/helm-v2.8.0-linux-amd64.tar.gz && \
    tar -zxvf helm-v2.8.0-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/

# Install Draft
RUN curl -sLO https://azuredraft.blob.core.windows.net/draft/draft-v0.10.0-linux-amd64.tar.gz && \
    tar -zxvf draft-v0.10.0-linux-amd64.tar.gz && \
    mv linux-amd64/draft /usr/local/bin/

# Install kubectl
RUN curl -sLO https://storage.googleapis.com/kubernetes-release/release/v1.9.2/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl

# Install svcat
RUN curl -sLO https://download.svcat.sh/cli/latest/linux/amd64/svcat && \
    chmod +x ./svcat && \
    mv ./svcat /usr/local/bin/

# Install dep
ENV DEP_VERSION=v0.4.1
RUN curl -sLO https://github.com/golang/dep/releases/download/$DEP_VERSION/dep-linux-amd64 && \
    chmod +x ./dep-linux-amd64 && \
    mv ./dep-linux-amd64 /usr/local/bin/dep

WORKDIR /go/src/github.com/Azure/service-catalog-templates

FROM golang:1.9.3-stretch

ENV USE_DOCKER=false

WORKDIR /tmp

# Install Helm
RUN curl -sLO https://kubernetes-helm.storage.googleapis.com/helm-v2.8.0-linux-amd64.tar.gz
RUN tar -zxvf helm-v2.8.0-linux-amd64.tar.gz
RUN mv linux-amd64/helm /usr/local/bin/

# Install Draft
RUN curl -sLO https://azuredraft.blob.core.windows.net/draft/draft-v0.10.0-linux-amd64.tar.gz
RUN tar -zxvf draft-v0.10.0-linux-amd64.tar.gz
RUN mv linux-amd64/draft /usr/local/bin/

# Install kubectl
RUN curl -sLO https://storage.googleapis.com/kubernetes-release/release/v1.9.2/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl

# Install svcat
RUN go get github.com/kubernetes-incubator/service-catalog/cmd/svcat

# Install dep
ENV DEP_VERSION=v0.4.1
RUN curl -sSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/$DEP_VERSION/dep-linux-amd64 && \
    chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/Azure/service-catalog-templates

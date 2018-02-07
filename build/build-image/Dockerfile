FROM golang:1.10-rc-stretch

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

WORKDIR /go/src/github.com/Azure/service-catalog-templates

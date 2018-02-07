PKG = github.com/Azure/service-catalog-templates
DOCKER_IMG = svcatt-build

USE_DOCKER ?= true

ifeq ($(USE_DOCKER),true)
  DO = docker run --rm -it -v $$HOME/.kube:/root/.kube -v $$HOME/.minikube:$$HOME/.minikube -v $$(pwd):/go/src/$(PKG) $(DOCKER_IMG)
else
  DO =
endif

default: build

.PHONY: buildimage build run create-cluster test

buildimage:
	docker build -t $(DOCKER_IMG) ./hack/buildimage

#build: buildimage
#	$(DO) ./hack/build.sh

#run: buildimage
#	$(DO) ./hack/run.sh
#	$(DO) svcat get brokers

create-cluster:
	./hack/create-cluster.sh

test: buildimage
	$(DO) ./hack/test.sh

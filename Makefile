PKG = github.com/Azure/service-catalog-templates
DOCKER_IMG = svcatt-build

USE_DOCKER ?= true

ifeq ($(USE_DOCKER),true)
  DO = docker run --rm -it -v $$HOME/.kube:/root/.kube -v $$HOME/.minikube:$$HOME/.minikube -v $$(pwd):/go/src/$(PKG) $(DOCKER_IMG)
else
  DO =
endif

default: build

.PHONY: build-image build create-cluster test

build-image:
	docker build -t $(DOCKER_IMG) ./build/build-image

build: build-image
	$(DO) ./build/build.sh

run: build-image
	$(DO) ./hack/run.sh

create-cluster:
	./hack/create-cluster.sh

test: build-image
	$(DO) ./hack/test.sh

PKG = github.com/Azure/service-catalog-templates
BUILD_IMG = svcatt-build
RUNTIME_IMG ?= carolynvs/service-catalog-templates

PKG_PATH = /go/src/$(PKG)
GOCACHE = $(PKG_PATH)/.gocache
BINDIR ?= build/service-catalog-templates

DO = docker run --rm -it -e GOCACHE=$(GOCACHE) -e BINDIR=$(BINDIR) -v $$HOME/.kube:/root/.kube -v $$HOME/.minikube:$$HOME/.minikube -v $$(pwd):$(PKG_PATH) $(BUILD_IMG)

default: build-image build

build-image:
	docker build -t $(BUILD_IMG) ./build/build-image

TYPES_FILES = $(shell find pkg/apis -name types.go)
codegen: pkg/client
pkg/client: $(TYPES_FILES)
	$(DO) ./build/update-codegen.sh

build/service-catalog-templates/service-catalog-templates: build
build: pkg/client
	echo $(BINDIR)
	$(DO) ./build/build.sh

push: build
	docker build -t $(RUNTIME_IMG) ./build/service-catalog-templates
	docker push $(RUNTIME_IMG)

deploy: push
	helm upgrade --install svcatt-crd --namespace svcatt charts/svcatt-crd
	helm upgrade --install svcatt-osba --namespace svcatt charts/svcatt-osba
	helm upgrade --install svcatt --namespace svcatt \
		--recreate-pods --force charts/svcatt \
		--set image.Repository="$(RUNTIME_IMAGE)",image.tag="latest",image.pullPolicy="Always",deploymentStrategy="Recreate"

create-cluster:
	./hack/create-cluster.sh

test:
	$(DO) ./hack/test.sh

.PHONY: default build-image codegen build build-runtime runtime-image push run deploy create-cluster test

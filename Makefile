SHELL := /bin/bash

PKG = github.com/Azure/service-catalog-templates
GOPATH = $(shell go env GOPATH)
PKG_PATH = /go/src/$(PKG)
GOCACHE = $(PKG_PATH)/.gocache

# Tagged commits should publish to latest, otherwise canary
PERMALINK ?= $(shell git name-rev --name-only --tags --no-undefined HEAD &> /dev/null && echo latest || echo canary)
VERSION ?= $(shell git describe --tags --dirty='+dev' 2> /dev/null || echo v0)
LDFLAGS = -w -X $(PKG)/pkg.Version=$(VERSION)
XBUILD = CGO_ENABLED=0 go build -a -tags netgo -ldflags '$(LDFLAGS)'
RELEASE_DIR = bin/cli/$(VERSION)
BINDIR ?= build/service-catalog-templates

BUILD_IMG = svcatt-build
RUNTIME_IMG ?= carolynvs/service-catalog-templates

ifdef NO_DOCKER
DO =
else
DO = docker run --rm -it -e GOCACHE=$(GOCACHE) -e BINDIR=$(BINDIR) -v $$HOME/.kube:/root/.kube -v $$HOME/.minikube:$$HOME/.minikube -v $$(pwd):$(PKG_PATH) $(BUILD_IMG)
endif

default: build-image build svcatt

build-image:
	docker build -t $(BUILD_IMG) ./build/build-image

TYPES_FILES = $(shell find pkg/apis -name types.go)
codegen: pkg/client
pkg/client: $(TYPES_FILES)
	$(DO) ./build/update-codegen.sh

build/service-catalog-templates/service-catalog-templates: build
build: pkg/client
	$(DO) ./build/build.sh

push: build
	docker build -t $(RUNTIME_IMG) ./build/service-catalog-templates
	docker push $(RUNTIME_IMG)

deploy: push
	helm upgrade --install svcatt-crd --namespace svcatt charts/svcatt-crd
	helm upgrade --install svcatt-osba --namespace svcatt charts/svcatt-osba
	helm upgrade --install svcatt --namespace svcatt \
		--recreate-pods --force charts/svcatt \
		--set image.Repository="$(RUNTIME_IMAGE)",image.tag="$(PERMALINK)",image.pullPolicy="Always",deploymentStrategy="Recreate"

create-cluster:
	./hack/create-cluster.sh

test-unit: pkg/client
	$(DO) go test ./...

test-integration:
	$(DO) ./hack/test.sh

check-dep:
	@if [ -z "$$(which dep)" ]; then \
		echo 'Missing `dep` client which is required for development'; \
		exit 2; \
	else \
		dep version; \
	fi

get-dep:
	# Install the latest release of dep
	go get -d -u github.com/golang/dep
	cd $(GOPATH)/src/github.com/golang/dep && \
	DEP_TAG=$$(git describe --abbrev=0 --tags) && \
	git checkout $$DEP_TAG && \
	go install -ldflags="-X main.version=$$DEP_TAG" ./cmd/dep; \
	git checkout master # Make go get happy by switching back to master

verify-vendor: check-dep
	dep ensure --vendor-only
	@if [ -n "$$(git status --porcelain vendor)" ]; then \
		echo 'vendor/ is out-of-date: run `dep ensure --vendor-only`'; \
		git status --porcelain vendor; \
		exit 2; \
	fi

clean:
	-rm -r bin
	-rm build/service-catalog-templates/service-catalog-templates

svcatt: pkg/client
	go build -o bin/svcatt -ldflags '$(LDFLAGS)' ./cmd/svcatt

svcatt-linux:
	GOOS=linux GOARCH=amd64 $(XBUILD) -o $(RELEASE_DIR)/Linux/x86_64/svcatt ./cmd/svcatt
	cd $(RELEASE_DIR)/Linux/x86_64 && shasum -a 256 svcatt > svcatt.sha256

svcatt-darwin:
	GOOS=darwin GOARCH=amd64 $(XBUILD) -o $(RELEASE_DIR)/Darwin/x86_64/svcatt ./cmd/svcatt
	cd $(RELEASE_DIR)/Darwin/x86_64 && shasum -a 256 svcatt > svcat.sha256

svcatt-windows:
	GOOS=windows GOARCH=amd64 $(XBUILD) -o $(RELEASE_DIR)/Windows/x86_64/svcatt.exe ./cmd/svcatt
	cd $(RELEASE_DIR)/Windows/x86_64 && shasum -a 256 svcatt.exe > svcatt.exe.sha256

svcatt-all: svcatt-linux svcatt-darwin svcatt-windows

svcatt-install: svcatt
	cp bin/svcatt /usr/local/bin/

publish-cli: clean svcatt-all
	cp -R $(RELEASE_DIR) bin/cli/$(PERMALINK)/
	# AZURE_STORAGE_CONNECTION_STRING will be used for auth in the following command
	az storage blob upload-batch -d cli -s bin/cli

publish-charts: clean
	./build/publish-charts.sh

.PHONY: svcatt

SHELL := /bin/bash

REPO_ROOT := $(shell git rev-parse --show-toplevel)

BINARY := bucket-proxy

IMAGE_REPO := ghcr.io/vietanhduong/bucket-proxy

EXTRA_DOCKER_BUILD_ARGS ?=

ifneq ($(BUILDPLATFORM),)
	EXTRA_DOCKER_BUILD_ARGS += --platform $(BUILDPLATFORM)
endif

include Makefile.defs

.PHONY: build
build:
	$(GO_BUILD) -o $(BINARY) $(REPO_ROOT)/cmd

.PHONY: install
install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 $(BINARY) $(DESTDIR)$(BINDIR)
	@rm -f $(BINARY)

.PHONY: docker-build
docker-build:
	@docker buildx create --use --name=crossplat --node=crossplat
	docker buildx build $(REPO_ROOT) \
		--build-arg=VERSION=$(VERSION) \
		--build-arg=GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg=BUILD_DATE=$(BUILD_DATE) \
		--build-arg=NOOPT=$(NOOPT) \
		--build-arg=NOSTRIP=$(NOSTRIP) \
		--build-arg=NOSTRIP=$(NOSTRIP) \
		$(EXTRA_DOCKER_BUILD_ARGS) \
		--output "type=docker,push=false" \
		--tag $(IMAGE_REPO):$(VERSION) \
		--file $(REPO_ROOT)/Dockerfile

.PHONY: docker-push
docker-push:
	@docker buildx create --use --name=crossplat --node=crossplat
	docker buildx build $(REPO_ROOT) \
		--build-arg=VERSION=$(VERSION) \
		--build-arg=GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg=BUILD_DATE=$(BUILD_DATE) \
		--build-arg=NOOPT=$(NOOPT) \
		--build-arg=NOSTRIP=$(NOSTRIP) \
		--build-arg=NOSTRIP=$(NOSTRIP) \
		$(EXTRA_DOCKER_BUILD_ARGS) \
		--output "type=image,push=true" \
		--tag $(IMAGE_REPO):$(VERSION) \
		--file $(REPO_ROOT)/Dockerfile

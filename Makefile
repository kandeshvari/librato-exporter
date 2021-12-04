NAME := librato-exporter
VERSION := 0.1.0
RELEASE := 1
REGISTRY := kandeshvari

IMAGE := $(REGISTRY)/$(NAME)

BUILD_DATE := $(shell LANG=c date)
GO_VERSION := $(shell go version | awk '{print $$3" "$$4}')
GIT_REVISION := $(shell git rev-list -1 HEAD)

.PHONY: docker-build docker-push docker-run

docker-build:
	docker build --progress=plain . \
		--network host \
		-f Dockerfile \
		-t $(IMAGE):$(VERSION)-$(RELEASE) \
		--build-arg VERSION="$(VERSION)" \
		--build-arg RELEASE="$(RELEASE)" \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		--build-arg GIT_REVISION="$(GIT_REVISION)"

docker-push:
	docker push $(IMAGE):$(VERSION)-$(RELEASE)

docker-run:
	docker run -it --rm -p 9800:9800 \
		$(IMAGE):$(VERSION)-$(RELEASE)

#VERSION ?= $(shell cat VERSION.txt || git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "v0.0.0-$(CI_COMMIT_SHA)")
#LGFLAGS := -X "github.com.de/xyz/gcp-housekeeper/internal/version.Version=$(VERSION:v%=%)"

#DOCKER := docker
#ifdef DOCKER_CMD
#DOCKER = $(DOCKER_CMD)
#endif

#DOCKER_BUILD_ARGS := --build-arg http_proxy=$(http_proxy) --build-arg https_proxy=$(https_proxy)
#DOCKER_RUN_ARGS := -e https_proxy=$(https_proxy) -e https_proxy=$(https_proxy)

#.PHONY: build
#build: gox
#	CGO_ENABLED=0 gox -ldflags "$(LGFLAGS)" -osarch="linux/amd64" -osarch="windows/amd64" -osarch="darwin/amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"

#.PHONY: gox
#gox:
#	GO111MODULE=off go get -u github.com/mitchellh/gox

#.PHONY: docker-e2e-image
#docker-e2e-image: check-env build
#	$(DOCKER) build $(DOCKER_BUILD_ARGS) --tag gcp-housekeeper-e2e -f e2e/Dockerfile .

#.PHONY: e2e
#e2e: docker-e2e-image
#	$(DOCKER) run -ti --rm -e gcp-housekeeper-e2e /usr/local/bin/...


#.PHONY: check-env
#check-env:
#ifndef GCP_HOUSEKEEPER_USERNAME
#	$(error GCP_HOUSEKEEPER_USERNAME is undefined)
#endif

.PHONY: test
test:
	go fmt ./...
	go vet ./...
	go test ./...

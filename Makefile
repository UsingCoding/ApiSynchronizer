export APP_CMD_NAME = apisynchronizer
export REGISTRY = vadimmakerov
export DOCKER_IMAGE_NAME = $(REGISTRY)/$(APP_CMD_NAME):master

all: build test

.PHONY: build
build: modules
	bin/go-build.sh "cmd" "bin/$(APP_CMD_NAME)" $(APP_CMD_NAME) .env

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run

.PHONY: publish
publish:
	docker build . --tag=$(DOCKER_IMAGE_NAME)
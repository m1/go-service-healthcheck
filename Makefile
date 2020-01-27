APP_NAME:=go_service_healthcheck
VERSION := `cat VERSION`
BUILD_DATE := `date +%FT%T%z`
GIT_COMMIT := `git describe --always --long --dirty`

VM_LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE} -X main.GitCommit=${GIT_COMMIT}"

.PHONY: help
.DEFAULT_GOAL := help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

docker-build: ## Build the cont${LDFLAGS_f1}ainer
	docker build -t $(APP_NAME) .

docker-build-nc: ## Build the container without caching
	docker build --no-cache -t $(APP_NAME) .

docker-run: ## Run container on port configured in `.env`
	docker run -i -t --rm --env-file=.env -p=8282:$(API_PORT) --name="$(APP_NAME)" $(APP_NAME)

docker-up: docker-build docker-run ## Run container on port configured in `config.env` (Alias to run)

stop: ## Stop and remove a running container
	docker stop $(APP_NAME); docker rm $(APP_NAME)

local-build:
	go build -ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitCommit=$(GIT_COMMIT)" -o bin/api

local-run: local-build
	./bin/api run

vm-build:
	GO111MODULE=on go mod vendor
	GO111MODULE=on go build ${LDFLAGS_f1} -o bin/api -mod vendor -race
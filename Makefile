# This is adapted from https://habr.com/ru/post/542410/
PROJECT = harbor-scan-report
BINARY = hsr

COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
TAG=$(shell git describe --tags |cut -d- -f1)

LDFLAGS = -ldflags "-X main.gitTag=${TAG} -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH}"
LDFLAGS_SMALL = -ldflags "-s -w"
.PHONY: help clean dep build binary upx small-binary build-small

.DEFAULT_GOAL := help

help: ## Display this help screen.
      @echo "Makefile available targets:"
      @grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'

dep: ## Download the dependencies.
	go mod download && go get github.com/kyberorg/${PROJECT}/cmd/${PROJECT}

build: dep ## Build executable.
	mkdir -p ./bin
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY} ./cmd/${PROJECT}

binary: build ## Alias to `build`

build-small:
	mkdir -p ./bin
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS_SMALL} -o bin/${BINARY} ./cmd/${PROJECT}

small-binary: build-small upx

upx: ## Compress with upx
	upx --brute bin/${BINARY}

clean: ## Clean build directory.
	rm -f ./bin/${BINARY}
	rmdir ./bin
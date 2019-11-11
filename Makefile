VERSION := $(shell git describe --tags 2> /dev/null || echo no-tag)
BRANCH := $(shell git symbolic-ref -q --short HEAD)
COMMIT := $(shell git rev-parse HEAD)

# Go related variables.
# GO_PKG ?= $(shell go list -e -f "{{ .ImportPath }}") -- для старых версий golang
GO_PKG ?= "github.com/Madredix/clickadu"
export GO_BIN := $(GGDEB_BIN_NAME)

# Default value for params
ifndef GO_BIN
override GO_BIN = "app"
endif

# Use linker flags to provide version/build settings
LDFLAGS := -ldflags "-X $(GO_PKG)/src/cmd.version=$(VERSION) -X $(GO_PKG)/src/cmd.commit=$(COMMIT) -X $(GO_PKG)/src/cmd.branch=$(BRANCH) -X $(GO_PKG)/src/cmd.buildTime=`date '+%Y-%m-%d_%H:%M:%S_%Z'`"

all: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

doc: ## Open swagger api documentation
	@nohup swagger serve "./swagger.json" -p1234 > /dev/null &
	@sleep 3 && killall swagger

swagger: ## Generate swagger server
	@rm -rf models
	@rm -rf restapi
	@swagger generate server -f ./swagger.xml > /dev/null 2>&1
	@rm -rf cmd

build: swagger ## Build app
	@go build $(LDFLAGS) -o $(GO_BIN) ./src/main.go


testall: swagger test-lint test-unit ## Run all tests
test-unit: ## Run unit tests
	@go clean -testcache ./src/... && go test -count 1 -parallel 9 -cover -coverprofile=cover.out ./src/...
	@go tool cover -html=cover.out

test-lint: ## Run golangci linter
	@golangci-lint run ./src/...

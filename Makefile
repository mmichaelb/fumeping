PROJECT_NAME=fumeping

GIT_VERSION=$(shell git describe --always)
GIT_BRANCH=$(shell git branch --show-current)
GIT_DEFAULT_BRANCH=main

LD_FLAGS = -X main.GitVersion=${GIT_VERSION} -X main.GitBranch=${GIT_BRANCH} -X main.GitDefaultBranch=${GIT_DEFAULT_BRANCH}

OUTPUT_SUFFIX=$(go env GOEXE)

# builds and formats the project with the built-in Golang tool
.PHONY: build
build:
	OUTPUT_NAME=bin/${PROJECT_NAME}-${GIT_VERSION}-${GOOS}-${GOARCH}${OUTPUT_SUFFIX}
	@go build -ldflags '${LD_FLAGS}' -o "${OUTPUT_NAME}" ./cmd/fumeping/main.go

# build go application for docker usage
.PHONY: build-docker
build-docker:
	OUTPUT_NAME=bin/${PROJECT_NAME}-${GIT_VERSION}-docker
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '${LD_FLAGS}' -o "${OUTPUT_NAME}" ./cmd/fumeping/main.go

# installs and formats the project with the built-in Golang tool
install:
	@go install -ldflags '${LD_FLAGS}' ./cmd/fumeping/main.go

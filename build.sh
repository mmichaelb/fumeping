# create binary directory
mkdir -p bin

PROJECT_NAME=fumeping

GIT_VERSION=$(git describe --always)
GIT_BRANCH=$(git branch --show-current)
GIT_DEFAULT_BRANCH=main

OUTPUT_SUFFIX=$(go env GOEXE)

OUTPUT_NAME=bin/${PROJECT_NAME}-${GIT_VERSION}-${GOOS}-${GOARCH}${OUTPUT_SUFFIX}

CGO_ENABLED=0 go build -ldflags "-X main.GitVersion=${GIT_VERSION} -X main.GitBranch=${GIT_BRANCH} -X main.GitDefaultBranch=${GIT_DEFAULT_BRANCH}" -o "${OUTPUT_NAME}" ./cmd/fumeping/main.go

EXIT_STATUS=$?

if [ $EXIT_STATUS == 0 ]; then
  echo "Build succeeded"
else
  echo "Build failed"
fi

exit $EXIT_STATUS
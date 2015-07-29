CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
TAG=${1:-latest}
REPO=ehazlett/docker-demo

all: build

test:
	@go test -v ./...

deps:
	@go get -d ./...

build:
	@go build -a -tags 'netgo' -ldflags '-w -linkmode external -extldflags -static' .

image: build
	@docker build -t $(REPO):$(TAG) .

.PHONY: build

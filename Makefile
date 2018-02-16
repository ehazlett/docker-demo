CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
TAG?=latest
REPO=ehazlett/docker-demo

all: build

test:
	@go test -v ./...

build:
	@docker build -t ${REPO}:${TAG} .

.PHONY: build

CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
TAG=${1:-latest}

all: build

deps:
	@go get -d ./...

build: deps
	@go build -o demo -a -tags 'netgo' -ldflags '-w -linkmode external -extldflags -static' .

image: build
	@docker build -t ehazlett/docker-demo-pg:$(TAG) .

.PHONY: build

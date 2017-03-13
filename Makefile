CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
TAG?=latest
REPO=ehazlett/docker-demo
MEDIA_SRCS=$(shell find ui/ -type f \
	-not -path "ui/semantic/dist/*" \
	-not -path "ui/node_modules/*")

all: media-semantic build

test:
	@go test -v ./...

deps:
	@go get -d ./...

build:
	@go build -a -tags 'netgo' -ldflags '-w -linkmode external -extldflags -static' .

dev-setup:
	@echo "This could take a while..."
	@npm install --loglevel verbose -g gulp browserify babelify
	@cd ui && npm install --loglevel verbose
	@cd ui/node_modules/semantic-ui && gulp install

media-semantic: static/dist/.bundle_timestamp
static/dist/.bundle_timestamp: $(MEDIA_SRCS)
	@cp -f ui/semantic.theme.config ui/semantic/src/theme.config
	@mkdir -p ui/semantic/src/themes/app && cp -rf ui/semantic.theme/* ui/semantic/src/themes/app/
	@cd ui/semantic && gulp build
	@mkdir -p static/dist
	@rm -rf static/dist/semantic* static/dist/themes
	@cp -f ui/semantic/dist/semantic.min.css static/dist/semantic.min.css
	@cp -f ui/semantic/dist/semantic.min.js static/dist/semantic.min.js
	@mkdir -p static/dist/themes/default && cp -r ui/semantic/dist/themes/default/assets static/dist/themes/default/
	@touch static/dist/.bundle_timestamp

image: build
	@docker build -t $(REPO):$(TAG) .

clean:
	@rm -rf docker-demo
	@rm -rf static/dist/.bundle_timestamp

.PHONY: build deps clean image

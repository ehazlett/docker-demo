FROM node:latest as ui
RUN npm install -g gulp browserify babelify
COPY ui/package.json /tmp/
COPY ui/semantic.json /tmp/
RUN cd /tmp && npm install && \
    mkdir -p /usr/src/app/ui && \
    cp -rf /tmp/node_modules /usr/src/app/ui/
WORKDIR /usr/src/app
COPY . /usr/src/app
RUN cd ui/node_modules/fomantic-ui && npx gulp install
RUN cp -f ui/semantic.theme.config ui/semantic/src/theme.config && \
    mkdir -p ui/semantic/src/themes/app && \
    cp -rf ui/semantic.theme/* ui/semantic/src/themes/app
RUN cd ui/semantic && npx gulp build

FROM golang:1.14-alpine as app
RUN apk add -U build-base git
COPY . /go/src/app
WORKDIR /go/src/app
ENV GO111MODULE=on
RUN go build -a -v -tags 'netgo' -ldflags '-w -linkmode external -extldflags -static' -o docker-demo .

FROM alpine:latest
RUN apk add -U --no-cache curl
COPY static /static
COPY --from=ui /usr/src/app/ui/semantic/dist/semantic.min.css static/dist/semantic.min.css
COPY --from=ui /usr/src/app/ui/semantic/dist/semantic.min.js static/dist/semantic.min.js
COPY --from=ui /usr/src/app/ui/semantic/dist/themes/default/assets static/dist/themes/default/
COPY --from=app /go/src/app/docker-demo /bin/docker-demo
COPY templates /templates
EXPOSE 8080
ENTRYPOINT ["/bin/docker-demo"]

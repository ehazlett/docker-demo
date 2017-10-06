FROM alpine:latest
MAINTAINER "Evan Hazlett <ejhazlett@gmail.com>"
COPY docker-demo /bin/docker-demo
COPY static /static
COPY templates /templates
EXPOSE 8080
ENTRYPOINT ["/bin/docker-demo"]

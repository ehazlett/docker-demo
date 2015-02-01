# Docker Demo Application
This is a Go demo application used for demonstrating Docker.

# Demo
The `demo` directory contains the demo application.

## Environment Variables

* `TITLE`: sets title in demo app

## Build

- `make`
- `docker build -t docker-demo .`

## Run

`docker run -P --rm docker-demo`

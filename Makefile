
SHELL=/bin/bash

default: build

build:
	go build -o bin/kka-cli

static:
	CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o bin/kka-cli

install:
	go install

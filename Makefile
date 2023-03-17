
SHELL=/bin/bash

default: build

build:
	go build -o bin/kka-cli

install:
	go install

update-bindata:
	go-bindata -pkg bindata -o src/utils/bindata/bindata.go docker/Dockerfile
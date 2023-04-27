
SHELL=/bin/bash

default: build

build:
	go build -o bin/kka-cli

install:
	go install

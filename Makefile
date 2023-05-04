
PROJECT_DIRECTORY := $(PWD)
SHELL=/bin/bash


default: help

help: ## help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## build cli app as binary
	go build -o bin/kka-cli

project_build: ## build a project using the cli
	@go run main.go --project-directory PROJECT_DIRECTORY=$(PROJECT_DIRECTORY) build

project_push: ## push a project to an registry
	@go run main.go push

static: ## build static
	CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o bin/kka-cli

install: ## install dependencies
	go install

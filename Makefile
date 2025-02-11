VERSION := $(shell git describe --tags --always --dirty)
SOURCES := $(shell find . -type f -name '*.go')

define HELP_HEADER
Usage:	make <target>

Targets:
endef

export HELP_HEADER

.PHONY: help
help: ## List all targets.
	@echo "$$HELP_HEADER"
	@grep -E '^[a-zA-Z0-9%_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

build: bin/labctl ## Build the application.

bin/labctl: $(SOURCES) ## Build labctl binary.
	go build -ldflags "-X main.version=$(VERSION) -s -w" -o $@ cmd/$(@F)/main.go

.PHONY: plan
plan: ## Plan the infrastructure changes.
	tofu -chdir=deploy/tofu init
	tofu -chdir=deploy/tofu plan | tee tofu.log
	@sed -i 's/\x1b\[[0-9;]*m//g' tofu.log

.PHONY: apply
apply: ## Apply the infrastructure changes.
	tofu -chdir=deploy/tofu init
	tofu -chdir=deploy/tofu apply -auto-approve

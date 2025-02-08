define HELP_HEADER
Usage:	make <target>

Targets:
endef

export HELP_HEADER

.PHONY: help
help: ## List all targets.
	@echo "$$HELP_HEADER"
	@grep -E '^[a-zA-Z0-9%_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: plan
plan: ## Plan the infrastructure changes.
	terraform -chdir=deploy/tofu plan

.PHONY: apply
apply: ## Apply the infrastructure changes.
	terraform -chdir=deploy/tofu apply -auto-approve

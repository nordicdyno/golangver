.PHONY: help
help:
	@./scripts/make-help.sh $(MAKEFILE_LIST)

.PHONY: test
test: ## run linter
	./scripts/linter.sh

.PHONY: install
install: ## go install golangver
	go install -v ./
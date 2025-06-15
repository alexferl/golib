.PHONY: dev audit cover fmt lint pre-commit test tidy update-deps

.DEFAULT: help
help:
	@echo "make dev"
	@echo "	setup development environment"
	@echo "make audit"
	@echo "	conduct quality checks"
	@echo "make cover"
	@echo "	generate coverage report"
	@echo "make fmt"
	@echo "	fix code format issues"
	@echo "make lint"
	@echo "	run lint checks"
	@echo "make pre-commit"
	@echo "	run pre-commit hooks"
	@echo "make test"
	@echo "	execute all tests"
	@echo "make tidy"
	@echo "	clean and tidy dependencies"
	@echo "make update-deps"
	@echo "	update dependencies"

DIRS = $(shell find . -name 'go.mod' -exec dirname {} \;)
define FOREACH
	for DIR in $(DIRS); do \
  		cd $(CURDIR)/$$DIR && $(1); \
  	done
endef

GOTESTSUM := go run gotest.tools/gotestsum@latest -f testname -- ./... -race -count=1
TESTFLAGS := -shuffle=on
COVERFLAGS := -covermode=atomic
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6

check-pre-commit:
ifeq (, $(shell which pre-commit))
	$(error "pre-commit not in $(PATH), pre-commit (https://pre-commit.com) is required")
endif

dev: check-pre-commit
	pre-commit install

audit:
	$(call FOREACH, go mod verify)
	$(call FOREACH, go run golang.org/x/vuln/cmd/govulncheck@latest ./...)

cover:
	$(call FOREACH,	$(GOTESTSUM) $(TESTFLAGS) $(COVERFLAGS))

fmt:
	$(call FOREACH, $(GOLANGCI_LINT) fmt)

lint:
	$(call FOREACH, $(GOLANGCI_LINT) run)

pre-commit: check-pre-commit
	pre-commit run --all-files

test:
	$(call FOREACH,	$(GOTESTSUM) $(TESTFLAGS))

tidy:
	$(call FOREACH,go mod tidy)

update-deps: tidy
	$(call FOREACH,go get -u)

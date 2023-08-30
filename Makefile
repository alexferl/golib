.PHONY: dev test cover tidy fmt

.DEFAULT: help
help:
	@echo "make dev"
	@echo "	setup development environment"
	@echo "make test"
	@echo "	run go test"
	@echo "make cover"
	@echo "	run go test with -cover"
	@echo "make tidy"
	@echo "	run go mod tidy"
	@echo "make fmt"
	@echo "	run gofumpt"
	@echo "make pre-commit"
	@echo "	run pre-commit hooks"

DIRS := $(wildcard */)
define FOREACH
	for DIR in $(DIRS); do \
  		cd $$DIR && $(1) && cd $(CURDIR); \
  	done
endef

check-gofumpt:
ifeq (, $(shell which gofumpt))
	$(error "No gofumpt in $(PATH), gofumpt (https://pkg.go.dev/mvdan.cc/gofumpt) is required")
endif

check-pre-commit:
ifeq (, $(shell which pre-commit))
	$(error "No pre-commit in $(PATH), pre-commit (https://pre-commit.com) is required")
endif

checks: check-gofumpt check-pre-commit

dev: checks
	pre-commit install

test:
	$(call FOREACH,go test -v ./...)

cover:
	$(call FOREACH,go test -cover -v ./...)

tidy:
	$(call FOREACH,go mod tidy)

fmt: check-gofumpt
	gofumpt -l -w .

pre-commit: check-pre-commit
	pre-commit

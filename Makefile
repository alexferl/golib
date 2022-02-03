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
	@echo "	run gofmt"
	@echo "make pre-commit"
	@echo "	run pre-commit hooks"

SUBDIRS := $(wildcard */)
define FOREACH
	for DIR in $(SUBDIRS); do \
  		cd $$DIR && $(1) && cd $(CURDIR); \
  	done
endef

check-pre-commit:
 ifeq (, $(shell which pre-commit))
 $(error "No pre-commit in $(PATH), pre-commit (https://pre-commit.com) is required")
 endif

dev: check-pre-commit
	pre-commit install

test:
	$(call FOREACH,go test -v)

cover:
	$(call FOREACH,go test -cover -v)

tidy:
	$(call FOREACH,go mod tidy -compat=1.17)

fmt:
	gofmt -s -w .

pre-commit: check-pre-commit
	pre-commit

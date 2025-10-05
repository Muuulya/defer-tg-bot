SHELL := /bin/bash

.PHONY: lint test tidy run

LINT_VERSION := v1.60.3
GOLANGCI_LINT := $(shell command -v golangci-lint 2>/dev/null)

lint:
	@echo "🔍 Проверка кода линтером..."
	@if [ -z "$(GOLANGCI_LINT)" ]; then \
		echo "Устанавливаю golangci-lint $(LINT_VERSION)…"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINT_VERSION); \
	fi
	golangci-lint run

# test:
# 	scripts/test.sh

tidy:
	go mod tidy

run:
	go run .

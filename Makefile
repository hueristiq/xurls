SHELL = /bin/bash

# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
# --- Git Hooks Install ----------------------------------------------------------------------------------------------------------------------------------------------------------------------
# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

.PHONY: lefthook-install
lefthook-install:
	(command -v lefthook || go install github.com/evilmartians/lefthook@latest) && lefthook install

# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
# --- Go(Golang) -----------------------------------------------------------------------------------------------------------------------------------------------------------------------------
# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

.PHONY: go-mod-clean
go-mod-clean:
	go clean -modcache

.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy

.PHONY: go-mod-update
go-mod-update:
	go get -f -t -u ./...
	go get -f -u ./...

.PHONY: go-fmt
go-fmt:
	(command -v golangci-lint || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2) && golangci-lint fmt ./...

.PHONY: go-lint
go-lint: go-fmt
	(command -v golangci-lint || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2) && golangci-lint run ./...

.PHONY: go-test
go-test:
	go test -v -race ./...

.PHONY: go-build
go-build:
	go build -v -ldflags '-s -w' -o bin/xurls cmd/xurls/main.go

.PHONY: go-install
go-install:
	go install -v ./...

# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
# --- Help -----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

.PHONY: help
help:
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo " Setup Commands:"
	@echo ""
	@echo "  lefthook-install ......... Install Git hooks."
	@echo ""
	@echo " Go Commands:"
	@echo ""
	@echo "  go-mod-clean ............. Clean Go module cache."
	@echo "  go-mod-tidy .............. Tidy Go modules."
	@echo "  go-mod-update ............ Update Go modules."
	@echo "  go-fmt ................... Format Go code."
	@echo "  go-lint .................. Lint Go code."
	@echo "  go-test .................. Run Go tests."
	@echo "  go-build ................. Build Go program."
	@echo "  go-install ............... Install Go program."
	@echo ""
	@echo " Help Commands:"
	@echo ""
	@echo "  help ..................... Display this help information."
	@echo ""

.DEFAULT_GOAL = help
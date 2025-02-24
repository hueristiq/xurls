# Specifies the shell to be used for executing commands. In this case, it's set to `/bin/bash`.
# Bash is chosen for its advanced scripting capabilities, including string manipulation and conditional checks.
SHELL = /bin/bash

# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------
# --- Prepare | Setup ------------------------------------------------------------------------------------------------------------------------------------------------------
# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------

.PHONY: git-hooks-install
# Target: git-hooks-install
# Purpose:
#   Installs and configures Git hooks using Lefthook, a Git hooks manager.
# Details:
#   - First, the target checks if the `lefthook` command is available in the system PATH.
#   - If `lefthook` is not installed, it installs the latest version using Go.
#   - Finally, it executes `lefthook install` to set up Git hooks based on the repository's configuration.
git-hooks-install:
	@command -v lefthook || go install github.com/evilmartians/lefthook@latest; lefthook install

# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------
# --- Go (Golang) ----------------------------------------------------------------------------------------------------------------------------------------------------------
# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------

.PHONY: go-mod-clean
# Target: go-mod-clean
# Purpose:
#   Cleans the Go module cache to remove any cached module files.
# Details:
#   - This target runs `go clean -modcache` which is useful if you encounter issues with outdated or corrupt module cache.
go-mod-clean:
	go clean -modcache

.PHONY: go-mod-tidy
# Target: go-mod-tidy
# Purpose:
#   Tidies up the go.mod file by adding missing and removing unused modules.
# Details:
#   - Running `go mod tidy` ensures that the go.mod file accurately reflects the dependencies used in the project.
go-mod-tidy:
	go mod tidy

.PHONY: go-mod-update
# Target: go-mod-update
# Purpose:
#   Updates all Go modules to their latest versions.
# Details:
#   - First, updates test dependencies with the flags: -f (force), -t (include test packages), and -u (update).
#   - Then, updates all other dependencies.
go-mod-update:
	go get -f -t -u ./...
	go get -f -u ./...

.PHONY: go-fmt
# Target: go-fmt
# Purpose:
#   Formats the Go source code.
# Details:
#   - Uses `go fmt ./...` to format all Go source files across the module, ensuring consistent code style.
go-fmt:
	go fmt ./...

.PHONY: go-lint
# Target: go-lint
# Purpose:
#   Lints the Go source code to catch potential issues and enforce code quality.
# Details:
#   - The target first ensures that the code is properly formatted by invoking the `go-fmt` target.
#   - Then, it checks if `golangci-lint` is available; if not, it installs a specific version.
#   - Finally, it runs the linter across all packages.
go-lint: go-fmt
	@(command -v golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5) && golangci-lint run ./...

.PHONY: go-test
# Target: go-test
# Purpose:
#   Executes the test suite for the Go project.
# Details:
#   - Runs tests in verbose mode (`-v`) and with race condition detection (`-race`) to ensure thread safety.
#   - The tests are executed for all packages in the module.
go-test:
	go test -v -race ./...

.PHONY: go-build
go-build:
	go build -v -ldflags '-s -w' -o bin/xurlunpack3r cmd/xurlunpack3r/main.go

.PHONY: go-install
go-install:
	go install -v ./...

# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------
# --- Docker ---------------------------------------------------------------------------------------------------------------------------------------------------------------
# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------

DOCKERCMD = docker
DOCKERBUILD = $(DOCKERCMD) build

DOCKERFILE := ./Dockerfile

IMAGE_NAME = hueristiq/xurlunpack3r
IMAGE_TAG = $(shell cat internal/configuration/configuration.go | grep "VERSION =" | sed 's/.*VERSION = "\([0-9.]*\)".*/\1/')
IMAGE = $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: docker-build
docker-build:
	@$(DOCKERBUILD) -f $(DOCKERFILE) -t $(IMAGE) -t $(IMAGE_NAME):latest .

# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------
# --- Help -----------------------------------------------------------------------------------------------------------------------------------------------------------------
# --------------------------------------------------------------------------------------------------------------------------------------------------------------------------

.PHONY: help
# Target: help
# Purpose:
#   Displays an overview of available targets along with their descriptions.
# Details:
#   - When no target is provided, the default action (set by .DEFAULT_GOAL) is to show this help text.
#   - This target prints categorized sections for environment management, Go commands, Docker commands, and help.
help:
	@echo ""
	@echo "*****************************************************************************"
	@echo ""
	@echo "PROJECT : $(PROJECT)"
	@echo ""
	@echo "*****************************************************************************"
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo " Git Hooks:"
	@echo "  git-hooks-install ........ Install Git hooks."
	@echo ""
	@echo " Go Commands:"
	@echo "  go-mod-clean ............. Clean Go module cache."
	@echo "  go-mod-tidy .............. Tidy Go modules."
	@echo "  go-mod-update ............ Update Go modules."
	@echo "  go-fmt ................... Format Go code."
	@echo "  go-lint .................. Lint Go code."
	@echo "  go-test .................. Run Go tests."
	@echo "  go-build ................. Build Go program."
	@echo "  go-install ............... Install Go program."
	@echo ""
	@echo " Docker Commands:"
	@echo "  docker-build ............. Build Docker image."
	@echo ""
	@echo " Help Commands:"
	@echo "  help ..................... Display this help information."
	@echo ""

# Set the default target to the help command.
# This ensures that running `make` without arguments provides a summary of available targets.
.DEFAULT_GOAL = help
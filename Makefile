TARGET ?= oracle
BUILD_DIR ?= $(CURDIR)/out
FLAGS ?= ""
ARCH ?= $(shell go env GOARCH)
.PHONY: all build clean run

all: lint build

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILD_DIR)/

run:
	@GOARCH=$(ARCH) go run ./cmd/$(TARGET) $(FLAGS)

$(BUILD_TARGETS): go.sum $(BUILD_DIR)/
	@echo "--> $@ "
	@GOARCH=$(ARCH) go $@ -mod=readonly $(BUILD_ARGS) ./...

# make BUILD_DIR=./bin
$(BUILD_DIR)/:
	mkdir -p $(BUILD_DIR)/

go.sum: go.mod
	echo "Ensure dependencies have not been modified" >&2
	GOPRIVATE=github.com/Carina-labs go mod verify
	GOPRIVATE=github.com/Carina-labs go mod tidy

pusher: go.sum $(BUILD_DIR)/
	@echo "--> Generate pusher for test"
	@GOARCH=$(ARCH) go build -mod=readonly -o out/ ./test/pusher

#######################################################
###                     Linting                     ###
#######################################################

.PHONY: lint

LintCon=PrivGolint
lint:
	@echo "--> Running linter"
	@golangci-lint run --out-format=tab

loc:
	@tokei .
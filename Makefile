TARGET ?= oracle
BUILD_DIR ?= $(CURDIR)/build

.PHONY: all build clean run

all: build lint

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILD_DIR)/

run:
	@go run ./cmd/$(TARGET)

$(BUILD_TARGETS): go.sum $(BUILD_DIR)/
	@go $@ -mod=readonly $(BUILD_ARGS) ./...

# make BUILD_DIR=./bin
$(BUILD_DIR)/:
	mkdir -p $(BUILD_DIR)/

go.sum: go.mod
	echo "Ensure dependencies have not been modified" >&2
	go mod verify
	go mod tidy

#######################################################
###                     Linting                     ###
#######################################################

.PHONY: lint

lint:
	@golangci-lint run --out-format=tab


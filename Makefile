TARGET ?= oracle
BUILD_DIR ?= $(CURDIR)/build
CUSTOM_ORGS ?= ""

.PHONY: all build clean run

all: lint build

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILD_DIR)/

run:
	@go run ./cmd/$(TARGET) $(CUSTOM_ORGS)

$(BUILD_TARGETS): go.sum $(BUILD_DIR)/
	@echo ">>>>>>>>>>>> $@ <<<<<<<<<<<<"
	@go $@ -mod=readonly $(BUILD_ARGS) ./...

# make BUILD_DIR=./bin
$(BUILD_DIR)/:
	mkdir -p $(BUILD_DIR)/

go.sum: go.mod
	echo "Ensure dependencies have not been modified" >&2
	GOPRIVATE=github.com/Carina-labs go mod verify
	GOPRIVATE=github.com/Carina-labs go mod tidy

#######################################################
###                     Linting                     ###
#######################################################

.PHONY: lint

lint:
	@echo ">>>> Execute golangci-lint <<<<"
	@golangci-lint run --out-format=tab

loc:
	@tokei .
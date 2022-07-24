TARGET ?= oracle
BUILD_DIR ?= $(CURDIR)/build
FLAGS ?= ""

.PHONY: all build clean run

all: lint build

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILD_DIR)/

run:
	@go run ./cmd/$(TARGET) $(FLAGS)

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

LintCon=PrivGolint
lint:
	@echo "--> Running linter"
	@if docker ps -a --format {{.Names}} | grep $(LintCon) > /dev/null; \
	then docker start -a $(LintCon); \
	else docker run -a stdout -a stderr -v $(CURDIR):/app -w /app --name $(LintCon) golint:v1.46.2 golangci-lint run --out-format=tab --timeout=10m; fi

loc:
	@tokei .
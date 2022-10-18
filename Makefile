TARGET ?= hal
ACTION ?= oracle
BUILD_DIR ?= $(CURDIR)/build
FLAGS ?= ""
ARCH ?= $(shell go env GOARCH)
.PHONY: all build clean run tester

all: lint build

BUILD_TARGETS := build install

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
	build_tags = muslc
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))
BUILD_FLAGS := -tags '$(build_tags)' -ldflags '$(ldflags)'

build: BUILD_ARGS=-o $(BUILD_DIR)/

run:
	GOARCH=$(ARCH) go run ./cmd/$(TARGET) $(ACTION) $(FLAGS)

$(BUILD_TARGETS): go.sum $(BUILD_DIR)/
	@echo " 🛠 $@ "
	GOARCH=$(ARCH) go $@ -mod=readonly $(BUILD_ARGS) $(BUILD_FLAGS) ./cmd/...

# make BUILD_DIR=./bin
$(BUILD_DIR)/:
	mkdir -p $(BUILD_DIR)/

go.sum: go.mod
	@echo "Ensure dependencies have not been modified" >&2
	GOPRIVATE=github.com/Carina-labs go mod verify
	GOPRIVATE=github.com/Carina-labs go mod tidy

tester: go.sum $(BUILD_DIR)/
	@echo "--> Generate tester for test"
	GOARCH=$(ARCH) go build -mod=readonly -o out/ ./tester/cmd/...

#######################################################
###                     Linting                     ###
#######################################################

.PHONY: lint

lint:
	@echo " 👻 Running linter"
	golangci-lint run --out-format=tab

loc:
	@tokei .
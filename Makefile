TARGET ?= oracle
.PHONY: build clean run

build:
	@go build -o ./bin/$(TARGET) ./cmd/$(TARGET)

run:
	@go run ./cmd/$(TARGET)
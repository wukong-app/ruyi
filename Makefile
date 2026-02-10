APP_ROOT := $(CURDIR)
INTERNAL_PACKAGE := "$(APP_ROOT)/internal"
CMD_RUYI := "$(APP_ROOT)/cmd/ruyi"
OUTPUT_DIR := "$(APP_ROOT)/output"

export GO111MODULE=on

.DEFAULT: all
all: build

prepare: wire
	@mkdir -p $(OUTPUT_DIR)

build: prepare
	@echo $(APP_ROOT)
	@echo "do build"
	@mkdir -p $(OUTPUT_DIR)/ruyi
	@go build -o $(OUTPUT_DIR)/ruyi/ruyi $(CMD_RUYI)
	@echo "build done"

test:
	@echo "do test"
	@go test -v -cover ./...
	@echo "test done"

fmt:
	@echo "do fmt"
	@go fmt ./...
	@echo "fmt done"

clean:
	@echo "do clean"
	@rm -rf $(OUTPUT_DIR)
	@echo "clean done"

wire:
	@echo "do wire"
	@cd $(INTERNAL_PACKAGE) && go run -mod=mod github.com/google/wire/cmd/wire
	@echo "wire done"

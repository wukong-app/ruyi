APP_ROOT := $(CURDIR)
INTERNAL_PACKAGE := "$(APP_ROOT)/internal"

export GO111MODULE=on

.DEFAULT: all
all: build

prepare: wire

build: prepare

wire:
	@echo "do wire"
	@cd $(INTERNAL_PACKAGE) && go run -mod=mod github.com/google/wire/cmd/wire
	@echo "wire done"

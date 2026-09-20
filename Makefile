# Makefile for buffaring project

BIN_NAME := buffaring
OUT_DIR  := bin
CMD_PKG  := ./cmd/buffaring
LDFLAGS  := -s -w

# Default: build for host platform
.PHONY: all
all: build

# Build for the host platform
.PHONY: build
build:
	@mkdir -p $(OUT_DIR)
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(OUT_DIR)/$(BIN_NAME) $(CMD_PKG)

# Run (pass args via RUN_ARGS, e.g. make run RUN_ARGS="send hello")
.PHONY: run
run: build
	./$(OUT_DIR)/$(BIN_NAME) $(RUN_ARGS)

# Release: cross-compile for linux/darwin/windows × amd64/arm64
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

.PHONY: release
release: clean
	@mkdir -p $(OUT_DIR)/release
	@$(foreach platform,$(PLATFORMS),\
		$(eval GOOS   := $(word 1,$(subst /, ,$(platform)))) \
		$(eval GOARCH := $(word 2,$(subst /, ,$(platform)))) \
		$(eval EXT    := $(if $(filter windows,$(GOOS)),.exe,)) \
		$(eval BIN    := $(BIN_NAME)-$(GOOS)-$(GOARCH)$(EXT)) \
		echo "  → $(GOOS)/$(GOARCH)" && \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -ldflags "$(LDFLAGS)" \
			-o $(OUT_DIR)/release/$(BIN) $(CMD_PKG) && \
		tar -czf $(OUT_DIR)/release/$(BIN).tar.gz -C $(OUT_DIR)/release $(BIN) && \
	) true
	@echo "Release artifacts in $(OUT_DIR)/release/"

# Clean generated files
.PHONY: clean
clean:
	rm -rf $(OUT_DIR)

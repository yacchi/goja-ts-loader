PACKAGE := goja_ts_loader

ASSETS_DIR := assets
ASSET_FILES := $(shell find $(ASSETS_DIR) -name "*.js")
TYPESCRIPT_LIB := $(ASSETS_DIR)/typescript.js
TYPESCRIPT_LIB_SRC := node_modules/typescript/lib/typescript.js
ASSET_SRC := $(ASSETS_DIR)/filesystem.go

.PHONY: build
build:
	yarn build

$(TYPESCRIPT_LIB): $(TYPESCRIPT_LIB_SRC)
	mkdir -p $(dir $(TYPESCRIPT_LIB))
	yarn terser --compress --mangle -o $(TYPESCRIPT_LIB) $(TYPESCRIPT_LIB_SRC)

update-typescript: $(TYPESCRIPT_LIB)

$(ASSET_SRC): $(ASSET_FILES)
	go generate $@
	
.PHONY: assets
assets: $(ASSET_SRC)
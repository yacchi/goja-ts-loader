PACKAGE := goja_ts_loader

TYPESCRIPT_LIB := transpiler/node_modules/typescript.js
TYPESCRIPT_LIB_SRC := node_modules/typescript/lib/typescript.js

all:

update-typescript: $(TYPESCRIPT_LIB)
	mkdir -p $(dir $(TYPESCRIPT_LIB))
	yarn terser --compress --mangle -o $(TYPESCRIPT_LIB) $(TYPESCRIPT_LIB_SRC)

test:
	go test ./...
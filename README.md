# goja-ts-loader

This is a Typescript transpiler use Official Typescript Compiler API, run under [goja](https://github.com/dop251/goja).

## usage

```go
package main

import (
	"github.com/yacchi/goja-ts-loader/transpiler"
	"io/ioutil"
	"os"
)

func main() {
	ts, err := transpiler.NewTranspiler()
	if err != nil {
		panic(err)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	} else {
		defer f.Close()
	}
	src, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	es5src, err := ts.Transpile("moduleName", string(src))
	if err != nil {
		panic(err)
	}
	println(es5src)
}
```

## goja SourceLoader

This package provides `SourceLoader` for goja_require in [goja_nodejs](https://github.com/dop251/goja_nodejs).

```go
package main

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/yacchi/goja-ts-loader"
)

func main() {
	vm := goja.New()

	reg := require.NewRegistryWithLoader(loader.TSLoader(require.DefaultSourceLoader))
	reg.Enable(vm)

	var out interface{}
	if _, err := vm.RunString(`var addr = require("./examples/goja_loader")`); err != nil {
		panic(err)
	}

	v := vm.Get("addr")
	if err := vm.ExportTo(v, &out); err != nil {
		panic(err)
	}
	fmt.Println(out)

	var addr func(a, b int) int
	if err := vm.ExportTo(v.ToObject(vm).Get("Add"), &addr); err != nil {
		panic(err)
	}
	fmt.Println(addr(1, 2))
}
```
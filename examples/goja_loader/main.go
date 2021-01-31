package main

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/yacchi/go-ts-transpiler/goja_loader"
	"github.com/yacchi/go-ts-transpiler/transpiler"
)

func main() {
	ts, err := transpiler.NewTranspiler()
	if err != nil {
		panic(err)
	}

	vm := goja.New()

	reg := require.NewRegistryWithLoader(goja_loader.TSLoader(ts, require.DefaultSourceLoader))
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

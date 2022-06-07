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

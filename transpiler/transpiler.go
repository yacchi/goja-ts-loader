package transpiler

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"io/ioutil"
	"os"
	"sync"
	"syscall"
)

type Transformer func(string) (string, error)

type cacheEntry struct {
	passThrough bool
	transpiled  string
}

type Transpiler struct {
	vm    *goja.Runtime
	m     sync.Mutex
	toES5 Transformer
	cache map[string]*cacheEntry
}

func (t *Transpiler) Transpile(module, src string) (string, error) {
	t.m.Lock()
	defer t.m.Unlock()
	if cached, ok := t.cache[module]; ok {
		if cached.passThrough {
			return src, nil
		} else {
			return cached.transpiled, nil
		}
	}
	if transpiled, err := t.toES5(src); err != nil {
		return "", err
	} else {
		ce := &cacheEntry{}
		if src == transpiled {
			ce.passThrough = true
		} else {
			ce.transpiled = transpiled
		}
		t.cache[module] = ce
		return transpiled, err
	}
}

func NewTranspiler() (*Transpiler, error) {
	vm := goja.New()

	registry := require.NewRegistryWithLoader(func(path string) ([]byte, error) {
		if f, err := FS.Open(path); err != nil {
			if os.IsNotExist(err) || errors.Is(err, syscall.EISDIR) {
				err = require.ModuleFileDoesNotExistError
			}
			return nil, err
		} else {
			defer f.Close()
			return ioutil.ReadAll(f)
		}
	})

	req := registry.Enable(vm)

	loader, err := req.Require("ts-loader")
	if err != nil {
		return nil, fmt.Errorf("transpiler: can not load loader module: %w", err)
	}

	var t Transformer
	if err := vm.ExportTo(loader.ToObject(vm).Get("TransformES5"), &t); err != nil {
		return nil, fmt.Errorf("transpiler: can not load transform method: %w", err)
	}

	return &Transpiler{
		vm:    vm,
		toES5: t,
		cache: map[string]*cacheEntry{},
	}, nil
}

package goja_loader

import (
	"errors"
	"fmt"
	"github.com/dop251/goja_nodejs/require"
	"github.com/yacchi/go-ts-transpiler/transpiler"
	"path/filepath"
)

var (
	DefaultOverrideExtensions = map[string]string{
		".js": ".ts",
	}
	DefaultTargetExtensions = []string{
		".js", ".ts",
	}
	DefaultSourceSizeLimit = map[string]int64{
		".js": 1 * 1024 * 1024, // 1MiB
	}
)

type loaderOption struct {
	overrides    map[string]string
	extensions   []string
	srcSizeLimit map[string]int64
}

type LoaderOption func(opt *loaderOption)

func OverrideExtensions(overrides map[string]string) LoaderOption {
	return func(opt *loaderOption) {
		opt.overrides = overrides
	}
}

func TargetExtensions(extensions []string) LoaderOption {
	return func(opt *loaderOption) {
		opt.extensions = extensions
	}
}

func SourceSizeLimit(limit map[string]int64) LoaderOption {
	return func(opt *loaderOption) {
		opt.srcSizeLimit = limit
	}
}

func TSLoader(transpiler *transpiler.Transpiler, base require.SourceLoader, opts ...LoaderOption) require.SourceLoader {
	o := &loaderOption{
		overrides:    DefaultOverrideExtensions,
		extensions:   DefaultTargetExtensions,
		srcSizeLimit: DefaultSourceSizeLimit,
	}

	for _, opt := range opts {
		opt(o)
	}

	allowExt := map[string]struct{}{}
	for _, ext := range o.extensions {
		allowExt[ext] = struct{}{}
	}

	var loader require.SourceLoader

	loader = func(path string) ([]byte, error) {
		ext := filepath.Ext(path)
		if newExt, ok := o.overrides[ext]; ok {
			srcName := path[:len(path)-len(ext)] + newExt
			if src, err := loader(srcName); err != nil {
				if !errors.Is(err, require.ModuleFileDoesNotExistError) {
					return nil, err
				}
			} else {
				return src, nil
			}
		}
		src, err := base(path)
		if err != nil {
			return nil, err
		} else if _, ok := allowExt[ext]; !ok {
			return src, nil
		}
		if limit, ok := o.srcSizeLimit[ext]; ok {
			if limit < int64(len(src)) {
				return src, nil
			}
		}
		if transpiled, err := transpiler.Transpile(path, string(src)); err != nil {
			return nil, fmt.Errorf("%s: %w", err, require.InvalidModuleError)
		} else {
			return []byte(transpiled), nil
		}
	}

	return loader
}

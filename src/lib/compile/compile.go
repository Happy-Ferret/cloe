package compile

import (
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/desugar"
	"github.com/tisp-lang/tisp/src/lib/parse"
	"github.com/tisp-lang/tisp/src/lib/util"
)

// MainModule compiles a main module of a path into outputs of thunks.
func MainModule(path string) []Output {
	module, err := parse.MainModule(util.ReadFileOrStdin(path))

	if err != nil {
		util.PanicError(err)
	}

	c := newCompiler()
	return c.compile(desugar.Desugar(module))
}

// SubModule compiles a sub module of a path into a map of names to thunks.
func SubModule(path string) map[string]*core.Thunk {
	return subModule(util.ReadFileOrStdin(path))
}

func subModule(filename, source string) map[string]*core.Thunk {
	module, err := parse.SubModule(filename, source)

	if err != nil {
		util.PanicError(err)
	}

	c := newCompiler()
	c.compile(desugar.Desugar(module))
	return c.env.toMap()
}

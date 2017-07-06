package ir

import "github.com/tisp-lang/tisp/src/lib/core"

// CompileFunction compiles a function in IR into a thunk.
func CompileFunction(s core.Signature, vars []interface{}, expr interface{}) *core.Thunk {
	return core.NewLazyFunction(
		s,
		func(ts ...*core.Thunk) core.Value {
			args := append(make([]*core.Thunk, 0, len(ts)+len(vars)), ts...)

			for _, v := range vars {
				args = append(args, interpretExpression(args, v))
			}

			return interpretExpression(args, expr)
		})
}

func interpretExpression(args []*core.Thunk, expr interface{}) *core.Thunk {
	switch x := expr.(type) {
	case int:
		return args[x]
	case *core.Thunk:
		return x
	case App:
		return x.interpret(args)
	}

	panic("Unreachable")
}

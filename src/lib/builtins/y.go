package builtins

import "github.com/coel-lang/coel/src/lib/core"

// Y is Y combinator which takes a function whose first argument is itself
// applied to the combinator.
//
// THE COMMENT BELOW MAY BE OUTDATED because we moved from a lambda calculus
// based combinator to an implementation based on a recursive function in Go.
//
// Using Y combinator to define built-in functions in Go source is dangerous
// because top-level recursive functions generate infinitely nested closures.
// (i.e. closure{f, x} where x will also be evaluated as closure{f, x}.)
var Y = core.NewLazyFunction(
	core.NewSignature([]string{"function"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		return y(ts[0])
	})

func y(f *core.Thunk) core.Value {
	return core.RawFunctionType(func(args core.Arguments) core.Value {
		return core.App(
			f,
			core.NewPositionalArguments(core.Normal(y(f))).Merge(args))
	})
}

package builtins

import "github.com/coel-lang/coel/src/lib/core"

// Ys is Y* combinator which takes functions whose first arguments are a list
// of themselves applied to the combinator.
var Ys = core.NewLazyFunction(
	core.NewSignature(nil, nil, "functions", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		t := ts[0]

		return core.PApp(xx, core.NewLazyFunction(
			core.NewSignature([]string{"x"}, nil, "", nil, nil, ""),
			func(ts ...*core.Thunk) core.Value {
				s := ts[0]

				applyF := core.NewLazyFunction(
					core.NewSignature([]string{"f"}, nil, "args", nil, nil, "kwargs"),
					func(ts ...*core.Thunk) core.Value {
						return core.App(ts[0], core.NewArguments(
							[]core.PositionalArgument{
								core.NewPositionalArgument(core.PApp(s, s), false),
								core.NewPositionalArgument(ts[1], true),
							},
							nil,
							[]*core.Thunk{ts[2]}))
					})

				return createNewFuncs(t, applyF)
			}))
	})

func createNewFuncs(olds, applyF *core.Thunk) *core.Thunk {
	v := core.PApp(core.Equal, olds, core.EmptyList).Eval()
	b, ok := v.(core.BoolType)

	if !ok {
		return core.NotBoolError(v)
	} else if b {
		return core.EmptyList
	}

	return core.PApp(core.Prepend,
		core.PApp(core.Partial, applyF, core.PApp(core.First, olds)),
		createNewFuncs(core.PApp(core.Rest, olds), applyF))
}

var xx = core.NewLazyFunction(
	core.NewSignature([]string{"x"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		return core.PApp(ts[0], ts[0])
	})

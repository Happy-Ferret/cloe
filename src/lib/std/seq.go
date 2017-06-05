package std

import "github.com/tisp-lang/tisp/src/lib/core"

// Seq runs arguments of outputs sequentially.
var Seq = core.NewLazyFunction(
	core.NewSignature(
		nil, nil, "outputs",
		nil, nil, "",
	),
	func(ts ...*core.Thunk) core.Value {
		t := ts[0]

		for {
			out := core.PApp(core.First, t)
			if err, ok := out.EvalOutput().(core.ErrorType); ok {
				return err
			}

			t = core.PApp(core.Rest, t)

			v := core.PApp(core.Equal, t, core.EmptyList).Eval()
			b, ok := v.(core.BoolType)

			if !ok {
				return core.NotBoolError(v)
			} else if b {
				return out
			}
		}
	})

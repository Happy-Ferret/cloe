package builtins

import (
	"github.com/coel-lang/coel/src/lib/core"
)

func createSeqFunction(f func(t core.Value) core.Value) core.Value {
	return core.NewLazyFunction(
		core.NewSignature(nil, nil, "args", nil, nil, ""),
		func(ts ...core.Value) core.Value {
			l := ts[0]

			for {
				t := core.PApp(core.First, l)
				l = core.PApp(core.Rest, l)

				if v := core.ReturnIfEmptyList(l, t); v != nil {
					return v
				}

				if err, ok := f(t).(core.ErrorType); ok {
					return err
				}
			}
		})
}

// Seq runs arguments of pure values sequentially and returns the last one.
var Seq = createSeqFunction(func(t core.Value) core.Value { return t.Eval() })

// EffectSeq runs arguments of effects sequentially and returns the last one.
var EffectSeq = createSeqFunction(func(t core.Value) core.Value { return t.EvalEffect() })

package builtins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/coel-lang/coel/src/lib/core"
)

func TestReadWithStdin(t *testing.T) {
	assert.Equal(t, core.StringType(""), core.PApp(Read).Eval())
}

func TestWrite(t *testing.T) {
	for _, a := range []core.Arguments{
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			nil,
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("sep", core.NewString(","))},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("end", core.NewString(""))},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("file", core.NewString("/tmp/coel.test"))},
			nil),
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.Nil, false)},
			[]core.KeywordArgument{core.NewKeywordArgument("mode", core.NewNumber(0775))},
			nil),
	} {
		assert.Equal(t, core.Nil.Eval(), core.App(Write, a).EvalEffect())
	}
}

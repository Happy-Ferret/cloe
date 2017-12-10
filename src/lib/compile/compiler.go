package compile

import (
	"fmt"
	"path"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/consts"
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/desugar"
	"github.com/tisp-lang/tisp/src/lib/ir"
	"github.com/tisp-lang/tisp/src/lib/modules"
	"github.com/tisp-lang/tisp/src/lib/parse"
)

type compiler struct {
	env   environment
	cache modulesCache
}

func newCompiler(e environment, c modulesCache) compiler {
	return compiler{e, c}
}

func (c *compiler) compileModule(m []interface{}) ([]Effect, error) {
	var err error
	es := []Effect{}

	for _, s := range m {
		switch x := s.(type) {
		case ast.LetVar:
			c.env.set(x.Name(), c.exprToThunk(x.Expr()))
		case ast.LetFunction:
			sig := x.Signature()
			ls := x.Lets()

			vars := make([]interface{}, 0, len(ls))
			varToIndex := sig.NameToIndex()

			for _, l := range ls {
				v := l.(ast.LetVar)
				vars = append(vars, c.exprToIR(varToIndex, v.Expr()))
				varToIndex[v.Name()] = len(varToIndex)
			}

			c.env.set(
				x.Name(),
				ir.CompileFunction(
					c.compileSignature(sig),
					vars,
					c.exprToIR(varToIndex, x.Body())))
		case ast.Effect:
			es = append(es, NewEffect(c.exprToThunk(x.Expr()), x.Expanded()))
		case ast.Import:
			m, ok := modules.Modules[x.Path()]

			if !ok && c.cache != nil {
				if cm, cached, err := c.cache.Get(x.Path()); err != nil {
					return nil, err
				} else if cached {
					m = cm
				} else {
					m, err = c.compileSubModule(x.Path() + consts.FileExtension)

					if err != nil {
						return nil, err
					}

					if err := c.cache.Set(x.Path(), m); err != nil {
						return nil, err
					}
				}
			} else if !ok {
				m, err = c.compileSubModule(x.Path() + consts.FileExtension)

				if err != nil {
					return nil, err
				}
			}

			for k, v := range m {
				c.env.set(path.Base(x.Path())+"."+k, v)
			}
		default:
			panic(fmt.Errorf("Invalid type: %#v", x))
		}
	}

	return es, nil
}

func (c *compiler) compileSubModule(path string) (module, error) {
	p, s, err := readFileOrStdin(path)

	if err != nil {
		return nil, err
	}

	m, err := parse.SubModule(p, s)

	if err != nil {
		return nil, err
	}

	cc := newCompiler(builtinsEnvironment(), c.cache)
	c = &cc
	_, err = c.compileModule(desugar.Desugar(m))

	if err != nil {
		return nil, err
	}

	return c.env.toMap(), nil
}

func (c *compiler) exprToThunk(expr interface{}) *core.Thunk {
	return core.PApp(ir.CompileFunction(
		core.NewSignature(nil, nil, "", nil, nil, ""),
		nil,
		c.exprToIR(nil, expr)))
}

func (c *compiler) compileSignature(sig ast.Signature) core.Signature {
	return core.NewSignature(
		sig.PosReqs(), c.compileOptionalArguments(sig.PosOpts()), sig.PosRest(),
		sig.KeyReqs(), c.compileOptionalArguments(sig.KeyOpts()), sig.KeyRest(),
	)
}

func (c *compiler) compileOptionalArguments(os []ast.OptionalArgument) []core.OptionalArgument {
	ps := make([]core.OptionalArgument, 0, len(os))

	for _, o := range os {
		ps = append(ps, core.NewOptionalArgument(o.Name(), c.exprToThunk(o.DefaultValue())))
	}

	return ps
}

func (c *compiler) exprToIR(varToIndex map[string]int, expr interface{}) interface{} {
	switch x := expr.(type) {
	case string:
		if i, ok := varToIndex[x]; ok {
			return i
		}

		return c.env.get(x)
	case ast.App:
		args := x.Arguments()

		ps := make([]ir.PositionalArgument, 0, len(args.Positionals()))
		for _, p := range args.Positionals() {
			ps = append(ps, ir.NewPositionalArgument(c.exprToIR(varToIndex, p.Value()), p.Expanded()))
		}

		ks := make([]ir.KeywordArgument, 0, len(args.Keywords()))
		for _, k := range args.Keywords() {
			ks = append(ks, ir.NewKeywordArgument(k.Name(), c.exprToIR(varToIndex, k.Value())))
		}

		ds := make([]interface{}, 0, len(args.ExpandedDicts()))
		for _, d := range args.ExpandedDicts() {
			ds = append(ds, c.exprToIR(varToIndex, d))
		}

		return ir.NewApp(
			c.exprToIR(varToIndex, x.Function()),
			ir.NewArguments(ps, ks, ds),
			x.DebugInfo())
	case ast.Switch:
		cs := make([]ir.Case, 0, len(x.Cases()))

		for _, k := range x.Cases() {
			cs = append(cs, ir.NewCase(
				c.env.get(k.Pattern()),
				c.exprToIR(varToIndex, k.Value())))
		}

		d := interface{}(nil)

		if x.DefaultCase() != nil {
			d = c.exprToIR(varToIndex, x.DefaultCase())
		}

		return ir.NewSwitch(c.exprToIR(varToIndex, x.Value()), cs, d)
	}

	panic(fmt.Errorf("Invalid type: %#v", expr))
}

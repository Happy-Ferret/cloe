package desugar

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/gensym"
)

func flattenStatement(s interface{}) []interface{} {
	switch s := s.(type) {
	case ast.LetFunction:
		return flattenLetFunction(s)
	default:
		return []interface{}{s}
	}
}

func flattenStatements(old []interface{}) []interface{} {
	new := make([]interface{}, 0)

	for _, s := range old {
		new = append(new, flattenStatement(s)...)
	}

	return new
}

func flattenLetFunction(f ast.LetFunction) []interface{} {
	f = flattenInnerStatements(f)

	ss := make([]interface{}, 0)
	ls := make([]interface{}, 0)
	ns := signatureToNames(f.Signature())

	for _, l := range f.Lets() {
		switch l := l.(type) {
		case ast.LetVar:
			ls = append(ls, l)
			ns.add(l.Name())
		case ast.LetFunction:
			args := ns.findInFunction(l).slice()
			n := gensym.GenSym(f.Name(), l.Name())

			ss = append(ss, letFlattenedFunction(l, n, args))
			ls = append(ls, letClosure(l, n, args))

			ns.add(l.Name())
		default:
			panic(fmt.Errorf("Invalid value: %#v", l))
		}
	}

	return append(ss, ast.NewLetFunction(f.Name(), f.Signature(), ls, f.Body(), f.DebugInfo()))
}

func letFlattenedFunction(f ast.LetFunction, n string, args []string) ast.LetFunction {
	return ast.NewLetFunction(
		n,
		prependPosReqsToSig(f.Signature(), args),
		f.Lets(),
		f.Body(),
		f.DebugInfo())
}

func letClosure(f ast.LetFunction, n string, args []string) ast.LetVar {
	return ast.NewLetVar(
		f.Name(),
		ast.NewApp(
			"$partial",
			ast.NewArguments(namesToPosArgs(append([]string{n}, args...)), nil, nil),
			f.DebugInfo()))
}

func namesToPosArgs(ns []string) []ast.PositionalArgument {
	ps := make([]ast.PositionalArgument, len(ns))

	for i, n := range ns {
		ps[i] = ast.NewPositionalArgument(n, false)
	}

	return ps
}

func flattenInnerStatements(f ast.LetFunction) ast.LetFunction {
	return ast.NewLetFunction(
		f.Name(),
		f.Signature(),
		flattenStatements(f.Lets()),
		f.Body(),
		f.DebugInfo())
}

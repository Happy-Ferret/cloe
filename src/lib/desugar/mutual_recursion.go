package desugar

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
	"github.com/tisp-lang/tisp/src/lib/gensym"
)

func desugarMutualRecursionStatement(s interface{}) []interface{} {
	switch s := s.(type) {
	case ast.MutualRecursion:
		return desugarMutualRecursion(s)
	default:
		return []interface{}{s}
	}
}

func desugarMutualRecursion(mr ast.MutualRecursion) []interface{} {
	fs := mr.LetFunctions()
	unrecs := make([]interface{}, 0, len(fs))

	for _, f := range fs {
		unrecs = append(unrecs, createUnrecursiveFunction(indexLetFunctions(fs...), f))
	}

	recsList := gensym.GenSym("ys", "mr", "functions")
	recs := make([]interface{}, 0, len(fs))

	for i, f := range fs {
		recs = append(
			recs,
			ast.NewLetVar(
				f.Name(),
				ast.NewPApp(recsList, []interface{}{fmt.Sprint(i)}, f.DebugInfo())))
	}

	return append(
		unrecs,
		append(
			[]interface{}{ast.NewLetVar(
				recsList,
				ast.NewPApp("$ys", stringsToAnys(letStatementsToNames(unrecs)), mr.DebugInfo()))},
			recs...)...)
}

func createUnrecursiveFunction(n2i map[string]int, f ast.LetFunction) ast.LetFunction {
	arg := gensym.GenSym("mr", "functions", "argument")
	n2i = copyNameToIndex(n2i)

	ls := make([]interface{}, 0, len(f.Lets()))

	for _, l := range f.Lets() {
		switch l := l.(type) {
		case ast.LetFunction:
			delete(n2i, l.Name())
		case ast.LetVar:
			delete(n2i, l.Name())
		default:
			panic("Unreachable")
		}

		ls = append(ls, replaceNames(arg, n2i, l, f.DebugInfo()))
	}

	return ast.NewLetFunction(
		gensym.GenSym("mr", "unrec", f.Name()),
		prependPosReqsToSig([]string{arg}, f.Signature()),
		ls,
		replaceNames(arg, n2i, f.Body(), f.DebugInfo()),
		f.DebugInfo())
}

func indexLetFunctions(fs ...ast.LetFunction) map[string]int {
	n2i := make(map[string]int)

	for i, f := range fs {
		n2i[f.Name()] = i
	}

	if len(n2i) != len(fs) {
		panic(fmt.Errorf("Duplicate names were found among mutually-recursive functions"))
	}

	return n2i
}

func replaceNames(funcList string, n2i map[string]int, x interface{}, di debug.Info) interface{} {
	replace := func(x interface{}) interface{} {
		return replaceNames(funcList, n2i, x, di)
	}

	switch x := x.(type) {
	case ast.LetFunction:
		n2i := copyNameToIndex(n2i)

		for n := range signatureToNames(x.Signature()) {
			delete(n2i, n)
		}

		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			switch l := l.(type) {
			case ast.LetFunction:
				delete(n2i, l.Name())
			case ast.LetVar:
				delete(n2i, l.Name())
			default:
				panic("Unreachable")
			}

			ls = append(ls, replaceNames(funcList, n2i, l, di))
		}

		return ast.NewLetFunction(
			x.Name(),
			x.Signature(),
			ls,
			replaceNames(funcList, n2i, x.Body(), di),
			x.DebugInfo())
	case ast.LetVar:
		return ast.NewLetVar(x.Name(), replace(x.Expr()))
	case ast.App:
		return ast.NewApp(replace(x.Function()), replace(x.Arguments()).(ast.Arguments), x.DebugInfo())
	case ast.Arguments:
		ps := make([]ast.PositionalArgument, 0, len(x.Positionals()))

		for _, p := range x.Positionals() {
			ps = append(ps, ast.NewPositionalArgument(replace(p.Value()), p.Expanded()))
		}

		ks := make([]ast.KeywordArgument, 0, len(x.Keywords()))

		for _, k := range x.Keywords() {
			ks = append(ks, ast.NewKeywordArgument(k.Name(), replace(k.Value())))
		}

		ds := make([]interface{}, 0, len(x.ExpandedDicts()))

		for _, d := range x.ExpandedDicts() {
			ds = append(ds, replace(d))
		}

		return ast.NewArguments(ps, ks, ds)
	case string:
		if i, ok := n2i[x]; ok {
			return ast.NewPApp(funcList, []interface{}{fmt.Sprint(i)}, di)
		}

		return x
	}

	panic(fmt.Errorf("Invalid value: %#v", x))
}

func copyNameToIndex(n2i map[string]int) map[string]int {
	new := make(map[string]int)

	for n, i := range n2i {
		new[n] = i
	}

	return new
}

func deleteNamesDefinedByLets(n2i map[string]int, ls []interface{}) map[string]int {
	n2i = copyNameToIndex(n2i)

	for _, n := range letStatementsToNames(ls) {
		delete(n2i, n)
	}

	return n2i
}

func letStatementsToNames(ls []interface{}) []string {
	ns := make([]string, 0, len(ls))

	for _, l := range ls {
		switch l := l.(type) {
		case ast.LetFunction:
			ns = append(ns, l.Name())
		case ast.LetVar:
			ns = append(ns, l.Name())
		default:
			panic("Unreachable")
		}
	}

	return ns
}

func stringsToAnys(ss []string) []interface{} {
	xs := make([]interface{}, 0, len(ss))

	for _, s := range ss {
		xs = append(xs, s)
	}

	return xs
}

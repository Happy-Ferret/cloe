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
				ast.NewPApp("$ys", letStatementsToNames(unrecs), mr.DebugInfo()))},
			recs...)...)
}

func createUnrecursiveFunction(n2i map[string]int, f ast.LetFunction) ast.LetFunction {
	arg := gensym.GenSym("mr", "functions", "argument")

	return replaceNames(
		arg,
		n2i,
		ast.NewLetFunction(
			gensym.GenSym("mr", "unrec", f.Name()),
			prependPosReqsToSig([]string{arg}, f.Signature()),
			f.Lets(),
			f.Body(),
			f.DebugInfo()),
		f.DebugInfo()).(ast.LetFunction)
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
	return ast.Convert(func(x interface{}) interface{} {
		switch x := x.(type) {
		case string:
			if i, ok := n2i[x]; ok {
				return ast.NewPApp(funcList, []interface{}{fmt.Sprint(i)}, di)
			}

			return x
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

				ls = append(ls, replaceNames(funcList, n2i, l, x.DebugInfo()))
			}

			return ast.NewLetFunction(
				x.Name(),
				x.Signature(),
				ls,
				replaceNames(funcList, n2i, x.Body(), x.DebugInfo()),
				x.DebugInfo())
		}

		return nil
	}, x)
}

func copyNameToIndex(n2i map[string]int) map[string]int {
	new := make(map[string]int)

	for n, i := range n2i {
		new[n] = i
	}

	return new
}

func letStatementsToNames(ls []interface{}) []interface{} {
	ns := make([]interface{}, 0, len(ls))

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

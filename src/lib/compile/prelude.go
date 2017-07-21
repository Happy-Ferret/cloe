package compile

import (
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/scalar"
	"github.com/tisp-lang/tisp/src/lib/std"
)

var prelude = func() environment {
	e := newEnvironment(scalar.Convert)

	for _, nv := range []struct {
		name  string
		value *core.Thunk
	}{
		{"if", core.If},
		{"$if", core.If},

		{"partial", core.Partial},

		{"first", core.First},
		{"$first", core.First},
		{"rest", core.Rest},
		{"$rest", core.Rest},
		{"prepend", core.Prepend},

		{"typeOf", core.TypeOf},
		{"$typeOf", core.TypeOf},
		{"isOrdered", core.IsOrdered},

		{"+", core.Add},
		{"-", core.Sub},
		{"*", core.Mul},
		{"/", core.Div},
		{"//", core.FloorDiv},
		{"mod", core.Mod},
		{"**", core.Pow},

		{"=", std.Equal},
		{"$=", std.Equal},
		{"<", std.Less},
		{"<=", std.LessEq},
		{">", std.Greater},
		{">=", std.GreaterEq},

		{"toStr", core.ToString},
		{"dump", core.Dump},

		{"delete", core.Delete},
		{"$delete", core.Delete},
		{"include", core.Include},
		{"$include", core.Include},
		{"insert", core.Insert},
		{"merge", core.Merge},
		{"size", core.Size},
		{"toList", core.ToList},

		{"dict", std.Dictionary},
		{"$dict", std.Dictionary},
		{"error", core.Error},

		{"y", std.Y},
		{"$y", std.Y},
		{"ys", std.Ys},
		{"$ys", std.Ys},

		{"par", std.Par},
		{"seq", std.Seq},
		{"rally", std.Rally},

		{"read", std.Read},
		{"write", std.Write},

		{"$matchError", core.NewError("MatchError", "A value didn't match with any pattern.")},
	} {
		e.set(nv.name, nv.value)
	}

	return e
}()

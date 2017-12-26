package desugar

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/debug"
	"github.com/stretchr/testify/assert"
)

func TestNamesFindInLetVar(t *testing.T) {
	n := "x"
	assert.True(t, newNames(n).findInLetVar(ast.NewLetVar(n, n)).include(n))
}

func TestNamesFindInLetFunction(t *testing.T) {
	n := "x"

	for _, test := range []struct {
		letFunc ast.LetFunction
		answer  bool
	}{
		{
			ast.NewLetFunction(
				n,
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				nil,
				n,
				debug.NewGoInfo(0)),
			true,
		},
		{
			ast.NewLetFunction(
				n,
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				[]interface{}{ast.NewLetVar(n, "y")},
				n,
				debug.NewGoInfo(0)),
			false,
		},
	} {
		assert.Equal(t, test.answer, newNames(n).findInLetFunction(test.letFunc).include(n))
	}
}

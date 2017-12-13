package match

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/stretchr/testify/assert"
)

func TestPatternRenamerRenameFail(t *testing.T) {
	defer func() {
		assert.NotEqual(t, nil, recover())
	}()

	newPatternRenamer().rename(
		ast.NewSwitch("nil", []ast.SwitchCase{ast.NewSwitchCase("nil", "true")}, false))
}

package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	root, err := ioutil.TempDir("", "")

	assert.Nil(t, err)

	d := filepath.Join(root, "foo")

	_, ok := core.PApp(createDirectory, core.NewString(d)).EvalEffect().(core.NilType)
	assert.True(t, ok)

	e, ok := core.PApp(createDirectory, core.NewString(d)).EvalEffect().(core.ErrorType)
	assert.True(t, ok)
	assert.Equal(t, "FileSystemError", e.Name())

	_, ok = core.App(
		createDirectory,
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.NewString(d), false)},
			[]core.KeywordArgument{core.NewKeywordArgument("existOk", core.True)},
			nil),
	).EvalEffect().(core.NilType)
	assert.True(t, ok)

	os.Remove(root)
}

func TestCreateDirectoryWithInvalidArguments(t *testing.T) {
	_, ok := core.PApp(createDirectory, core.Nil).EvalEffect().(core.ErrorType)
	assert.True(t, ok)

	_, ok = core.App(
		createDirectory,
		core.NewArguments(
			[]core.PositionalArgument{core.NewPositionalArgument(core.NewString("foo"), false)},
			[]core.KeywordArgument{core.NewKeywordArgument("existOk", core.Nil)},
			nil),
	).EvalEffect().(core.ErrorType)
	assert.True(t, ok)

}
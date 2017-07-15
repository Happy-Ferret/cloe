package ast

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/debug"
)

// Import represents an import of a sub module.
type Import struct {
	path string
	info debug.Info
}

// NewImport creates an Import.
func NewImport(path string, info debug.Info) Import {
	return Import{path, info}
}

// Path returns a path to an imported sub module.
func (i Import) Path() string {
	return i.path
}

func (i Import) String() string {
	return fmt.Sprintf("(import %v)", i.path)
}

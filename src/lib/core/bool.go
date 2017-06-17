package core

// BoolType represents a boolean values in the language.
type BoolType bool

// True is a true value.
var True = NewBool(true)

// False is a false value.
var False = NewBool(false)

// NewBool converts a Go boolean value into BoolType.
func NewBool(b bool) *Thunk {
	return Normal(rawBool(b))
}

func rawBool(b bool) BoolType {
	return BoolType(b)
}

// If returns the second argument when the first one is true or the third one
// otherwise.
var If = NewLazyFunction(
	NewSignature(
		[]string{"condition", "then", "else"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		b, ok := v.(BoolType)

		if !ok {
			return NotBoolError(v)
		}

		if b {
			return ts[1]
		}

		return ts[2]
	})

func (b BoolType) compare(c comparable) int {
	if b == c.(BoolType) {
		return 0
	} else if b {
		return 1
	}

	return -1
}

func (b BoolType) string() Value {
	return sprint(b)
}

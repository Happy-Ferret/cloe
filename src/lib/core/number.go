package core

import "math"

// NumberType represents a number in the language.
// It will perhaps be represented by DEC64 in the future release.
type NumberType float64

// NewNumber creates a thunk containing a number value.
func NewNumber(n float64) *Thunk {
	return Normal(NumberType(n))
}

// Add sums up numbers of arguments.
var Add = NewLazyFunction(
	NewSignature(
		nil, nil, "nums",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		ts, err := l.ToValues()

		if err != nil {
			return err
		}

		sum := NumberType(0)

		for _, t := range ts {
			v := t.Eval()
			n, ok := v.(NumberType)

			if !ok {
				return NotNumberError(v)
			}

			sum += n
		}

		return sum
	})

// Sub subtracts arguments of the second to the last from the first one as numbers.
var Sub = NewLazyFunction(
	NewSignature(
		[]string{"minuend"}, nil, "subtrahends",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		n0, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		v = ts[1].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		ts, err := l.ToValues()

		if err != nil {
			return err
		}

		if len(ts) == 0 {
			return NumArgsError("sub", ">= 1")
		}

		for _, t := range ts {
			v := t.Eval()
			n, ok := v.(NumberType)

			if !ok {
				return NotNumberError(v)
			}

			n0 -= n
		}

		return n0
	})

// Mul multiplies numbers of arguments.
var Mul = NewLazyFunction(
	NewSignature(
		nil, nil, "nums",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		ts, err := l.ToValues()

		if err != nil {
			return err
		}

		prod := NumberType(1)

		for _, t := range ts {
			v := t.Eval()
			n, ok := v.(NumberType)

			if !ok {
				return NotNumberError(v)
			}

			prod *= n
		}

		return prod
	})

// Div divides the first argument by arguments of the second to the last one by one.
var Div = NewLazyFunction(
	NewSignature(
		[]string{"dividend"}, nil, "divisors",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		n0, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		v = ts[1].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		ts, err := l.ToValues()

		if err != nil {
			return err
		}

		if len(ts) == 0 {
			return NumArgsError("div", ">= 1")
		}

		for _, t := range ts {
			v := t.Eval()
			n, ok := v.(NumberType)

			if !ok {
				return NotNumberError(v)
			}

			n0 /= n
		}

		return n0
	})

// FloorDiv divides the first argument by arguments of the second to the last one by one.
var FloorDiv = NewLazyFunction(
	NewSignature(
		[]string{"dividend"}, nil, "divisors",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		n0, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		v = ts[1].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		ts, err := l.ToValues()

		if err != nil {
			return err
		}

		if len(ts) == 0 {
			return NumArgsError("floorDiv", ">= 1")
		}

		for _, t := range ts {
			v := t.Eval()
			n, ok := v.(NumberType)

			if !ok {
				return NotNumberError(v)
			}

			n0 = NumberType(math.Floor(float64(n0 / n)))
		}

		return n0
	})

// Mod calculate a remainder of a division of the first argument by the second one.
var Mod = NewStrictFunction(
	NewSignature(
		[]string{"dividend", "divisor"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		n1, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		v = ts[1].Eval()
		n2, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		return NewNumber(math.Mod(float64(n1), float64(n2)))
	})

// Pow calculates an exponentiation from a base of the first argument and an
// exponent of the second argument.
var Pow = NewStrictFunction(
	NewSignature(
		[]string{"base", "exponent"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		n1, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		v = ts[1].Eval()
		n2, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		return NewNumber(math.Pow(float64(n1), float64(n2)))
	})

var isInt = NewLazyFunction(
	NewSignature(
		[]string{"number"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		n, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		return NewBool(math.Mod(float64(n), 1) == 0)
	})

func (n NumberType) compare(c comparable) int {
	if n < c.(NumberType) {
		return -1
	} else if n > c.(NumberType) {
		return 1
	}

	return 0
}

func (NumberType) ordered() {}

func (n NumberType) string() Value {
	return sprint(n)
}

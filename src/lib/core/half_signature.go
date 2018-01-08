package core

type halfSignature struct {
	requireds []string
	optionals []OptionalArgument
	rest      string
}

func (hs halfSignature) arity() int {
	n := len(hs.requireds) + len(hs.optionals)

	if hs.rest != "" {
		n++
	}

	return n
}

func (hs halfSignature) bindPositionals(args *Arguments) ([]*Thunk, *Thunk) {
	ts := make([]*Thunk, 0, hs.arity())
	restAreKeywords := false

	for _, name := range hs.requireds {
		t := (*Thunk)(nil)

		if !restAreKeywords {
			t = args.nextPositional()
		}

		if t == nil {
			t = args.searchKeyword(name)
			restAreKeywords = true
		}

		if t == nil {
			return nil, argumentError("Could not bind a required positional argument.")
		}

		ts = append(ts, t)
	}

	for _, o := range hs.optionals {
		t := (*Thunk)(nil)

		if !restAreKeywords {
			t = args.nextPositional()
		}

		if t == nil {
			t = args.searchKeyword(o.name)
			restAreKeywords = true
		}

		if t == nil {
			t = o.defaultValue
		}

		ts = append(ts, t)
	}

	if hs.rest != "" {
		ts = append(ts, args.restPositionals())
	}

	return ts, nil
}

func (hs halfSignature) bindKeywords(args *Arguments) ([]*Thunk, *Thunk) {
	ts := make([]*Thunk, 0, hs.arity())

	for _, name := range hs.requireds {
		t := args.searchKeyword(name)

		if t == nil {
			return nil, argumentError("Could not bind a required keyword argument.")
		}

		ts = append(ts, t)
	}

	for _, opt := range hs.optionals {
		t := args.searchKeyword(opt.name)

		if t == nil {
			t = opt.defaultValue
		}

		ts = append(ts, t)
	}

	if hs.rest != "" {
		ts = append(ts, args.restKeywords())
	}

	return ts, nil
}

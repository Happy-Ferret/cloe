package core

// Arguments represents a structured set of arguments passed to a predicate.
// It allows destructive operations to internal properties because it is
// guaranteed by Thunks that arguments objects are never reused as a function
// call creates a Thunk.
type Arguments struct {
	positionals   []*Thunk
	expandedList  *Thunk
	keywords      []KeywordArgument
	expandedDicts []*Thunk
}

// NewArguments creates a new Arguments.
func NewArguments(
	ps []PositionalArgument,
	ks []KeywordArgument,
	ds []*Thunk) Arguments {
	ts := make([]*Thunk, 0, len(ps))
	l := (*Thunk)(nil)

	for i, p := range ps {
		if p.expanded {
			l = mergeRestPositionalArgs(ps[i].value, ps[i+1:]...)
			break
		}

		ts = append(ts, p.value)
	}

	return Arguments{ts, l, ks, ds}
}

func mergeRestPositionalArgs(t *Thunk, ps ...PositionalArgument) *Thunk {
	for _, p := range ps {
		if p.expanded {
			t = PApp(Merge, t, p.value)
		} else {
			t = PApp(
				NewLazyFunction(appendFuncSignature, appendFunc), // Avoid initialization loop
				t, p.value)
		}
	}

	return t
}

func (args *Arguments) nextPositional() *Thunk {
	if len(args.positionals) != 0 {
		t := args.positionals[0]
		args.positionals = args.positionals[1:]
		return t
	}

	if args.expandedList == nil {
		return nil
	}

	l := args.expandedList
	args.expandedList = PApp(Rest, l)
	return PApp(First, l)
}

func (args *Arguments) restPositionals() *Thunk {
	ts := args.positionals
	l := args.expandedList
	args.positionals = nil
	args.expandedList = nil

	if l == nil {
		return NewList(ts...)
	}

	return PApp(Merge, NewList(ts...), l)
}

func (args *Arguments) searchKeyword(s string) *Thunk {
	for i, k := range args.keywords {
		if s == k.name {
			args.keywords = append(args.keywords[:i], args.keywords[i+1:]...)
			return k.value
		}
	}

	for i, t := range args.expandedDicts {
		v := t.Eval()
		d, ok := v.(DictionaryType)

		if !ok {
			return NotDictionaryError(v)
		}

		k := StringType(s)

		if v, ok := d.Search(k); ok {
			ds := make([]*Thunk, len(args.expandedDicts))
			copy(ds, args.expandedDicts)
			ds[i] = Normal(d.Remove(k))
			args.expandedDicts = ds
			return v
		}
	}

	return nil
}

func (args *Arguments) restKeywords() *Thunk {
	ks := args.keywords
	ds := args.expandedDicts
	args.keywords = nil
	args.expandedDicts = nil

	t := EmptyDictionary

	for _, k := range ks {
		t = PApp(Insert, t, NewString(k.name), k.value)
	}

	return PApp(Merge, append([]*Thunk{t}, ds...)...)
}

// Merge merges 2 sets of arguments into one.
func (args Arguments) Merge(old Arguments) Arguments {
	var ps []*Thunk
	var l *Thunk

	if args.expandedList == nil {
		ps = append(args.positionals, old.positionals...)
		l = old.expandedList
	} else {
		ps = args.positionals
		l = PApp(Append, append([]*Thunk{args.expandedList}, old.positionals...)...)

		if old.expandedList != nil {
			l = PApp(Merge, l, old.expandedList)
		}
	}

	return Arguments{
		ps,
		l,
		append(args.keywords, old.keywords...),
		append(args.expandedDicts, old.expandedDicts...),
	}
}

func (args Arguments) empty() *Thunk {
	if len(args.positionals) > 0 {
		return argumentError("%d positional arguments are left", len(args.positionals))
	}

	// Testing args.expandedList is impossible because we cannot know its length
	// without evaluating it.

	n := 0

	for _, t := range args.expandedDicts {
		v := t.Eval()
		d, ok := v.(DictionaryType)

		if !ok {
			return NotDictionaryError(v)
		}

		n += d.Size()
	}

	if n != 0 || args.keywords != nil && len(args.keywords) > 0 {
		return argumentError("%d keyword arguments are left", len(args.keywords)+n)
	}

	return nil
}

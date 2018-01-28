package re

import "github.com/coel-lang/coel/src/lib/core"

func evaluateStringArguments(ts []core.Value) ([]string, core.Value) {
	ss := make([]string, 0, len(ts))

	for _, t := range ts {
		s, err := t.EvalString()

		if err != nil {
			return nil, err
		}

		ss = append(ss, string(s))
	}

	return ss, nil
}

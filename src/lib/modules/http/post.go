package http

import (
	"net/http"
	"strings"

	"github.com/coel-lang/coel/src/lib/core"
)

var post = core.NewLazyFunction(
	core.NewSignature(
		[]string{"url", "body"}, nil, "",
		nil,
		[]core.OptionalArgument{
			core.NewOptionalArgument("contentType", core.NewString("text/plain")),
			core.NewOptionalArgument("error", core.True),
		},
		"",
	),
	func(ts ...core.Value) core.Value {
		ss := make([]string, 0, 3)

		for i := 0; i < cap(ss); i++ {
			s, err := ts[i].EvalString()

			if err != nil {
				return err
			}

			ss = append(ss, string(s))
		}

		r, err := http.Post(ss[0], ss[2], strings.NewReader(ss[1]))

		return handleMethodResult(r, err, ts[3])
	})

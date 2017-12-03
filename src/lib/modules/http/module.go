package http

import "github.com/tisp-lang/tisp/src/lib/core"

// Module is a module in the language.
var Module = map[string]*core.Thunk{
	"Get":         get,
	"GetRequests": getRequests,
	"Post":        post,
}

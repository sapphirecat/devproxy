package devproxy

// Copyright (c) 2013, Sapphire Cat <https://github.com/sapphirecat>.  All
// rights reserved.  See the accompanying LICENSE file for license terms.

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

// Bitfield of verbosity.
const (
	VerboseNone      = 0
	VerboseRuleMatch = 1
	VerboseGoProxy   = 2
	VerboseAll       = 3
)

func NewServer(r Ruleset, verbosity int) http.Handler {
	proxy := goproxy.NewProxyHttpServer()
	if (verbosity & VerboseGoProxy) != 0 {
		proxy.Verbose = true
	}

	rules_verbose := false
	if (verbosity & VerboseRuleMatch) != 0 {
		rules_verbose = true
	}
	// Set up rules *in goproxy* that apply the ruleset to each request.
	SetDefaultRules(proxy, r, rules_verbose)
	return proxy
}

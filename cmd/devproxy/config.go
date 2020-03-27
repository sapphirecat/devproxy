package main

import (
	"regexp"

	"github.com/sapphirecat/devproxy"
)

// Interception configuration
func ConfigureRules(r RouteArgs) devproxy.Ruleset {
	// Full ruleset constructed here; capacity=1 is an optimization hint.
	// Try to match it to the number of rules, but it's not critical.  The
	// Add() function will allocate internally if needed.
	rules := devproxy.NewRuleset(1)

	// Add a rule to the ruleset.
	rules.Add(devproxy.Rule{
		// Matcher: may be matched against a hostname only (http + default port)
		// or may include a ":port" section (http + other port, https + any port)
		regexp.MustCompile("^(?:.+\\.)?example\\.(com|net|org)(?::\\d+)?$"),

		// Action: a destination to forward to.  Represented as a function that
		// returns either "" (meaning declined, try the next matcher) or a
		// "hostname[:port]" to connect.  :port MUST be added for TLS!
		//
		// There are some pre-defined functions to build Actions: SendHttpTo,
		// SendTlsTo, and SendAllTo all take a host, and send traffic to ports 80
		// and/or 443 as appropriate.  Each of those has a Send*ToPort variant
		// that takes a host and port, and sends traffic to the specified port.
		// SendAllToPort sends both HTTP and TLS traffic to the _same_ port; for
		// different ports, use SendAllToDualPorts(host, httpPort, tlsPort).
		devproxy.SendAllTo(r.Target)})

	// More rules can be added here, if needed.

	// Return the completed ruleset.
	return rules
}

// vim:ft=go

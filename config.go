package main

import "regexp"

// Interception configuration
func ConfigureRules(r RouteArgs) Ruleset {
	// Full ruleset constructed here; capacity=1 because there will be 1 rule.
	rules := NewRuleset(1)
	rules.Add(Rule{
		// Matcher: may be matched against a hostname only (http + default port)
		// or may include a ":port" section (http + other port, https + any port)
		regexp.MustCompile("^(?:.+\\.)?example\\.(com|net|org)(?::\\d+)?$"),

		// Action: a func(host,mode)string that returns either "" (meaning
		// declined, try the next matcher) or a "hostname[:port]" to connect.
		// :port MUST be added for TLS!
		//
		// We also offer a set of functions that build Actions: SendHttpTo,
		// SendTlsTo, and SendAllTo all take a host, and send traffic to ports 80
		// and/or 443 as appropriate.  Each of those has a Send*ToPort variant
		// that takes a host and port, and sends traffic to the specified port.
		//
		// Yes, SendAllToPort sends both HTTP and TLS traffic to the _same_ port.
		SendAllTo(r.Target)})

	return rules
}

// vim:ft=go

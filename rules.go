package devproxy

// Copyright (c) 2013, Sapphire Cat <https://github.com/sapphirecat>.  All
// rights reserved.  See the accompanying LICENSE file for license terms.

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/sapphirecat/devproxy/_third_party/github.com/elazarl/goproxy"
)

const (
	RuleForHttp = iota
	RuleForTls
)

type Rule struct {
	Matcher *regexp.Regexp
	Actor   Action // function called when match hits -> string [hostname[:port]]
}

type Mode int
type Action func(string, Mode) string
type Ruleset struct {
	items []Rule
}

func (r *Ruleset) Add(a Rule) {
	r.items = append(r.items, a)
}

func (r *Ruleset) Length() int {
	return len(r.items)
}

func NewRuleset(capacity int) Ruleset {
	return Ruleset{
		make([]Rule, 0, capacity),
	}
}

func redirectPort(really bool, host string, port int) string {
	if really == true {
		return fmt.Sprintf("%s:%d", host, port)
	} else {
		return ""
	}
}

func SendHttpTo(host string) Action {
	return func(matched_host string, mode Mode) string {
		if mode == RuleForHttp {
			return host
		} else {
			return ""
		}
	}
}

func SendHttpToPort(host string, port int) Action {
	return func(matched_host string, mode Mode) string {
		return redirectPort(mode == RuleForHttp, host, port)
	}
}

func SendAllTo(host string) Action {
	return func(matched_host string, mode Mode) string {
		if mode == RuleForHttp {
			return host
		} else {
			return redirectPort(true, host, 443)
		}
	}
}

func SendAllToPort(host string, port int) Action {
	return func(matched_host string, mode Mode) string {
		return redirectPort(true, host, port)
	}
}

func SendTlsTo(host string) Action {
	return func(matched_host string, mode Mode) string {
		return redirectPort(mode == RuleForTls, host, 443)
	}
}

func SendTlsToPort(host string, port int) Action {
	return func(matched_host string, mode Mode) string {
		return redirectPort(mode == RuleForTls, host, port)
	}
}

func getTarget(rules Ruleset, hostname string, mode Mode) string {
	for _, rule := range rules.items {
		if rule.Matcher.MatchString(hostname) == true {
			if result := rule.Actor(hostname, mode); result != "" {
				return result
			}
		}
	}

	return ""
}

func NewDefaultHttpsRule(ruleset Ruleset, verbose bool) func(string, *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		target := getTarget(ruleset, host, RuleForTls)
		if target == "" {
			target = host
			if verbose {
				log.Println("!match HTTPS", host)
			}
		} else if verbose {
			log.Println("+HTTPS", host, ctx.Req.URL.Path)
		}

		return goproxy.OkConnect, target
	}
}

func NewDefaultHttpRule(ruleset Ruleset, verbose bool) func(*http.Request, *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		host := r.URL.Host
		target := getTarget(ruleset, host, RuleForHttp)
		if target != "" {
			r.URL.Host = target
			if verbose {
				log.Println("+plain", host, r.URL.Path)
			}
		} else if verbose {
			log.Println("!match plain", r.URL.Host)
		}

		return r, nil
	}
}

func SetDefaultRules(proxy *goproxy.ProxyHttpServer, rules Ruleset, verbose bool) {
	proxy.OnRequest().HandleConnectFunc(NewDefaultHttpsRule(rules, verbose))
	proxy.OnRequest().DoFunc(NewDefaultHttpRule(rules, verbose))
}

package main

// Copyright (c) 2013, Sapphire Cat <https://github.com/sapphirecat>.  All
// rights reserved.  See the accompanying LICENSE file for license terms.

import (
	"log"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
)

const (
	RuleForHttp = iota
	RuleForTls
)

type RuleAction struct {
	Matcher      *regexp.Regexp
	HttpUpstream string // HTTP target "host:port"
	TlsUpstream  string // HTTPS target "host:port"
}

var rules = []*RuleAction{}

func RuleAdd(act *RuleAction) {
	rules = append(rules, act)
}

func RuleRemove(act *RuleAction) int {
	// someday stdlib should have a filter function for slices....
	var cur, end, hits int

	for end = len(rules); cur < end; cur++ {
		// this could do multiple copies if you added the same rule many times.
		// so don't do that.
		for rules[cur] == act && cur < end {
			hits++
			end--
			rules[cur] = rules[end]
		}
	}

	if hits > 0 {
		rules = rules[0:end]
	}

	return hits
}

func getTarget(hostname string, mode int) string {
	for _, actor := range rules {
		if actor.Matcher.MatchString(hostname) == true {
			if mode == RuleForHttp && actor.HttpUpstream != "" {
				return actor.HttpUpstream
			} else if mode == RuleForTls && actor.TlsUpstream != "" {
				return actor.TlsUpstream
			}
		}
	}

	return ""
}

func NewDefaultHttpsRule(verbose bool) func(string, *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		target := getTarget(host, RuleForTls)
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

func NewDefaultHttpRule(verbose bool) func(*http.Request, *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		host := r.URL.Host
		target := getTarget(host, RuleForHttp)
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

func SetDefaultRules(proxy *goproxy.ProxyHttpServer, verbose bool) {
	proxy.OnRequest().HandleConnectFunc(NewDefaultHttpsRule(verbose))
	proxy.OnRequest().DoFunc(NewDefaultHttpRule(verbose))
}

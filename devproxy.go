package main

// Copyright (c) 2013, Sapphire Cat <https://github.com/sapphirecat>.  All
// rights reserved.  See the accompanying LICENSE file for license terms.

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

type RouteArgs struct {
	Bind   string
	Port   int
	Target string
}

type Args struct {
	RouteArgs
	verbose_self bool
	verbose_guts bool
}

func ParseFlags() Args {
	var listen = flag.String("listen", "127.0.0.1", "IP address on which to listen for connections")
	var port = flag.Int("port", 8111, "Port on which to listen for connections")
	var dest = flag.String("target", "127.0.0.1", "IP address to direct interceptions to")
	var verbose = flag.Bool("verbose", false, "Enables logging of devproxy interceptions")
	var debug = flag.Bool("debug", false, "Enables excessive logging in goproxy")

	flag.Parse()
	return Args{RouteArgs{*listen, *port, *dest}, *verbose, *debug}
}

func DevProxy(a Args) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = a.verbose_guts

	ruleset := ConfigureRules(a.RouteArgs)
	SetDefaultRules(proxy, ruleset, a.verbose_self)

	listen_addr := fmt.Sprintf("%s:%d", a.Bind, a.Port)
	log.Println("listening on", listen_addr, "with", ruleset.Length(), "interception rules active")
	log.Fatal(http.ListenAndServe(listen_addr, proxy))
}

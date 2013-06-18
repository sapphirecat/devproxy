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

var port = flag.Int("port", 8111, "Port on which to listen for connections")
var verbose = flag.Bool("verbose", false, "Enables logging of devproxy interceptions")
var debug = flag.Bool("debug", false, "Enables excessive logging in goproxy")

func DevProxy(listen_addr string, verbose, debug bool) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = debug

	SetDefaultRules(proxy, verbose)

	log.Println("listening on", listen_addr, "with", len(rules), "interception rules active")
	log.Fatal(http.ListenAndServe(listen_addr, proxy))
}

func main() {
	flag.Parse()
	listen_on := fmt.Sprintf("%s:%d", LISTEN_INTERFACE, *port)
	DevProxy(listen_on, *verbose, *debug)
}

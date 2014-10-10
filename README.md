# DevProxy

A proxy that connects requests made to certain hosts to an arbitrary backend
server, called the doppelgänger.  Example the first (note that the “real”
address and port where `example.com` would be served from is irrelevant):

1. devproxy is configured with 127.0.0.1:8080 as the HTTP doppelgänger for
   example.com
2. The client requests http://example.com/gnorc
3. devproxy connects to 127.0.0.1:8080
4. devproxy requests `/gnorc`
5. devproxy copies the request headers to the doppelgänger, including
   `Host: example.com`
6. devproxy copies the doppelgänger’s full response to the client

Example the second:

1. devproxy is configured with 127.0.0.1:8443 as the TLS doppelgänger for
   example.com
2. The client requests a connection to `example.com:443`
3. devproxy connects to 127.0.0.1:8443 and tells the client “OK”
4. The client and the doppelgänger do a TLS handshake, which succeeds
   without warnings if the doppelgänger has the `example.com` certificate and
   key available to it
5. The client makes an request with `Host: example.com` inside the tunnel
6. The doppelgänger processes the request and returns a response via the
   tunnel

Swapping the backend server is the extent of devproxy’s meddling, so it has no
need to MITM the secure connection.  It doesn’t need to modify that traffic,
it just shuttles bytes.


# Using it

1. Copy config.go.sample to config.go
2. Edit as desired
3. `go install github.com/elazarl/goproxy` if you don’t have it
4. `go build`
5. Run devproxy (Linux/OS X) or devproxy.exe (Windows)
6. Set your web proxy to 127.0.0.1:8111

I use [FoxyProxy Basic](http://getfoxyproxy.org) with Firefox so that I can
easily switch between using devproxy or not, and see at a glance whether I
_am_ using it.

I also use a virtual machine as the doppelgänger with the Web server on the
default ports, so that the application can hard-code production URLs and have
them actually be self-referential.


# Why?

A staging environment should be as close to production as physically possible.
Every single difference may be a source of bugs in production that *cannot* be
detected in staging.

Modifying `/etc/hosts` all the time to point the production DNS at the staging
server (or comment said redirection back out) is tedious, invisible, and
requires frequent privilege escalations.

With FoxyProxy, switching the browser between truth and lie is improved in all
respects: a fast user-level action with a status indicator.  All that you need
is a proxy to transparently connect the staging backend when a production URL
is requested… and that proxy is devproxy.


# License

See LICENSE, but the tl;dr is “3-clause BSD.”
package auth

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strings"
)

type AuthOptions struct {
	Realm                string
	AuthenticationMethod func(string, string) bool
	UnauthorizedHandler  http.HandlerFunc
}

type basicAuth struct {
	opts    *AuthOptions
	handler http.Handler
}

func BasicAuth(options *AuthOptions) *basicAuth {
	if options.UnauthorizedHandler == nil {
		options.UnauthorizedHandler = defaultUnauthorizedHandler
	}

	return &basicAuth{options, nil}
}

func (self *basicAuth) Wrap(handler http.HandlerFunc) http.HandlerFunc {
	self.handler = handler

	return func(w http.ResponseWriter, req *http.Request) {
		self.ServeHTTP(w, req)
	}
}

func (self *basicAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !self.authenticate(req) {
		self.requestAuth(w, req)
		return
	}

	// and call the next handler
	if self.handler != nil {
		self.handler.ServeHTTP(w, req)
	}
}

func (self *basicAuth) requestAuth(w http.ResponseWriter, req *http.Request) {
	// request auth
	w.Header().Set("WWW-Authenticate", "Basic realm="+self.opts.Realm)

	// and send the response
	self.opts.UnauthorizedHandler(w, req)
}

func (self *basicAuth) authenticate(req *http.Request) bool {
	const basicScheme string = "Basic "

	auth := req.Header.Get("Authorization")
	if !strings.HasPrefix(auth, basicScheme) {
		return false
	}

	// Get the plain-text username and password from the request
	// The first six characters are skipped - e.g. "Basic ".
	str, err := base64.StdEncoding.DecodeString(auth[len(basicScheme):])
	if err != nil {
		return false
	}

	creds := bytes.SplitN(str, []byte(":"), 2)

	if len(creds) != 2 {
		return false
	}

	return self.opts.AuthenticationMethod(string(creds[0]), string(creds[1]))
}

func defaultUnauthorizedHandler(w http.ResponseWriter, req *http.Request) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

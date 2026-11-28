package helpers

import (
	"net/http"
	"net/url"
)

func ValidateOrigin(origin string, request *http.Request) bool {
	if request.Method == http.MethodGet || request.Method == http.MethodHead {
		return true
	}
	// An empty Origin is not treated as proof of a same-origin browser request.
	if origin != "" && request.Header.Get("Origin") == origin {
		return true
	}
	// The local Vite proxy can be opened as either localhost or 127.0.0.1.
	// Treat only these two loopback spellings as equivalent, on the same port
	// and scheme, when the configured allowed origin itself is loopback.
	allowed, allowedErr := url.Parse(origin)
	incoming, incomingErr := url.Parse(request.Header.Get("Origin"))
	if allowedErr != nil ||
		incomingErr != nil ||
		allowed.Scheme != "http" ||
		incoming.Scheme != "http" ||
		allowed.Port() == "" ||
		allowed.Port() != incoming.Port() {
		return false
	}

	isLoopback := func(host string) bool {
		return host == "localhost" || host == "127.0.0.1"
	}
	return isLoopback(allowed.Hostname()) && isLoopback(incoming.Hostname())
}

package contracts

import "net/http"

// Module describes the lifecycle of an independently bootstrapped service.
type Module interface {
	RegisterHTTP(mux *http.ServeMux)
	Close()
}

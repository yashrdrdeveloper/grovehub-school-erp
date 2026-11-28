package routes

import (
	"net/http"

	"grovehub/launcher/internal/src/constants/errorcode"
	corehttp "grovehub/libs/core/http"
)

func RegisterFallback(mux *http.ServeMux) {
	mux.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_ = corehttp.WriteErrorContext(
			request.Context(),
			writer,
			http.StatusNotFound,
			errorcode.RouteNotFound,
		)
	}))
}

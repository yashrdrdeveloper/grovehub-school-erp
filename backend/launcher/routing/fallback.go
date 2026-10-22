package routing

import (
	"net/http"

	httpresponse "grovehub/libs/core/http/response"
	"grovehub/launcher/constants/errorcode"
)

func RegisterFallback(mux *http.ServeMux) {
	mux.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
		_ = httpresponse.WriteError(
			writer,
			http.StatusNotFound,
			errorcode.RouteNotFound,
		)
	})
}

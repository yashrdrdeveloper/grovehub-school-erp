package routing

import (
	"net/http"

	httpresponse "grovehub/libs/core/http/response"
	"grovehub/launcher/constants/errorcode"
)

func RegisterFallback(mux *http.ServeMux) {
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_ = httpresponse.WriteErrorContext(
			request.Context(),
			writer,
			http.StatusNotFound,
			errorcode.RouteNotFound,
		)
	})
}

package routing

import (
	"net/http"

	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
	"grovehub/launcher/constants/errorcode"
)

func RegisterFallback(mux *http.ServeMux) {
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_ = corehttp.WriteErrorContext(
			request.Context(),
			writer,
			http.StatusNotFound,
			errorcode.RouteNotFound,
		)
	})
}

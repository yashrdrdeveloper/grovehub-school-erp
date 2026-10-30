package routes

import (
	"net/http"

	"grovehub/services/infrastructure/controller"
)

func RegisterHTTP(mux *http.ServeMux, dkimController controller.DKIMController) {
	mux.HandleFunc("POST /internal/infrastructure/mail/dkim/ensure", dkimController.Ensure)
}

package protected

import "net/http"

// @Routes /audit
type AuditController interface {
	// @POST /search
	Search(writer http.ResponseWriter, request *http.Request)
}

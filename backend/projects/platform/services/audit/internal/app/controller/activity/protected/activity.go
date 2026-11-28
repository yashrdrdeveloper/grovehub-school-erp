package protected

import "net/http"

// @Routes /
type ActivityController interface {
	// @POST /search
	Search(writer http.ResponseWriter, request *http.Request)
}

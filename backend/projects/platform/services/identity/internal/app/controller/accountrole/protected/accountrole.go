package protected

import "net/http"

// @Routes /account-roles
type AccountRoleController interface {
	// @GET /{accountId}
	Get(writer http.ResponseWriter, request *http.Request)

	// @POST /{accountId}/preview
	Preview(writer http.ResponseWriter, request *http.Request)
}

package protected

import "net/http"

// @Routes /accounts
type AccountController interface {
	// @POST /
	Create(writer http.ResponseWriter, request *http.Request)

	// @POST /search
	Search(writer http.ResponseWriter, request *http.Request)

	// @GET /roles
	Roles(writer http.ResponseWriter, request *http.Request)

	// @GET /{accountId}
	Get(writer http.ResponseWriter, request *http.Request)

	// @PATCH /{accountId}
	Update(writer http.ResponseWriter, request *http.Request)

	// @PATCH /{accountId}/suspend
	Suspend(writer http.ResponseWriter, request *http.Request)

	// @PATCH /{accountId}/activate
	Activate(writer http.ResponseWriter, request *http.Request)

	// @PATCH /{accountId}/deactivate
	Deactivate(writer http.ResponseWriter, request *http.Request)

	// @POST /{accountId}/resend-invitation
	Resend(writer http.ResponseWriter, request *http.Request)
}

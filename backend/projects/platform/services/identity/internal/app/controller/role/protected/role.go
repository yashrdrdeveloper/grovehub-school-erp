package protected

import "net/http"

// @Routes /roles
type RoleController interface {
	// @POST /
	CreateDraft(writer http.ResponseWriter, request *http.Request)

	// @PATCH /{roleId}
	UpdateMetadata(writer http.ResponseWriter, request *http.Request)

	// @POST /search
	Search(writer http.ResponseWriter, request *http.Request)

	// @GET /{roleId}
	Get(writer http.ResponseWriter, request *http.Request)

	// @GET /{roleId}/permissions
	Permissions(writer http.ResponseWriter, request *http.Request)

	// @POST /{roleId}/permissions/preview
	PreviewPermissions(writer http.ResponseWriter, request *http.Request)
}

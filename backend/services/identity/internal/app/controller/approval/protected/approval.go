package protected

import "net/http"

// @Routes /approvals
type ApprovalController interface {
	// @POST /history/search
	History(writer http.ResponseWriter, request *http.Request)

	// @POST /inbox/search
	Inbox(writer http.ResponseWriter, request *http.Request)

	// @GET /inbox/count
	PendingCount(writer http.ResponseWriter, request *http.Request)

	// @GET /{requestId}
	Detail(writer http.ResponseWriter, request *http.Request)

	// @GET /{requestId}/audit
	Audit(writer http.ResponseWriter, request *http.Request)
}

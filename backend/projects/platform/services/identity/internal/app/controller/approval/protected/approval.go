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

	// @POST /lifecycle
	SubmitLifecycle(writer http.ResponseWriter, request *http.Request)

	// @POST /permissions
	SubmitPermissions(writer http.ResponseWriter, request *http.Request)

	// @POST /assignments
	SubmitAssignment(writer http.ResponseWriter, request *http.Request)

	// @POST /{requestId}/approve
	Approve(writer http.ResponseWriter, request *http.Request)

	// @POST /{requestId}/reject
	Reject(writer http.ResponseWriter, request *http.Request)

	// @POST /{requestId}/cancel
	Cancel(writer http.ResponseWriter, request *http.Request)

	// @POST /{requestId}/execute
	Execute(writer http.ResponseWriter, request *http.Request)
}

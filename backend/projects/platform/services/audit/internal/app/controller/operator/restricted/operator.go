package restricted

import "net/http"

// @Routes /
type OperatorController interface {
	// @POST /operator-activities
	Write(writer http.ResponseWriter, request *http.Request)
	// @POST /operator-activities/search
	Search(writer http.ResponseWriter, request *http.Request)
	// @GET /schools/{schoolId}/activity-check
	SchoolCheck(writer http.ResponseWriter, request *http.Request)
	// @POST /school-rename-receipts
	Receipt(writer http.ResponseWriter, request *http.Request)
	// @POST /school-status-receipts
	StatusReceipt(writer http.ResponseWriter, request *http.Request)
	// @POST /school-mutation-receipts/search
	ReceiptSearch(writer http.ResponseWriter, request *http.Request)
}

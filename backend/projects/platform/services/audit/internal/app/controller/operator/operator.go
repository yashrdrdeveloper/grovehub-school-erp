package operator

import "net/http"

type Controller interface {
	Write(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)
	SchoolCheck(http.ResponseWriter, *http.Request)
	Receipt(http.ResponseWriter, *http.Request)
	StatusReceipt(http.ResponseWriter, *http.Request)
	ReceiptSearch(http.ResponseWriter, *http.Request)
}

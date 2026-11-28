package operator

import (
	"grovehub/libs/core/request"
	"grovehub/services/audit/internal/app/dto/producer"
)

type WriteRequest struct {
	request.Request
	RequestBody producer.OperatorActivity
}

type RenameReceiptRequest struct {
	request.Request
	RequestBody producer.SchoolRenameReceipt
}

type StatusReceiptRequest struct {
	request.Request
	RequestBody producer.SchoolStatusReceipt
}

type SearchRequest struct {
	request.Request
	RequestBody SearchInput
}

type ReceiptSearchRequest struct {
	request.Request
	RequestBody ReceiptSearchInput
}

type SchoolCheckRequest struct {
	request.Request
	Params SchoolCheckParams
}

func (r RenameReceiptRequest) SchoolID() string {
	return r.RequestBody.SchoolID
}

func (r StatusReceiptRequest) SchoolID() string {
	return r.RequestBody.SchoolID
}

package operator

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
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

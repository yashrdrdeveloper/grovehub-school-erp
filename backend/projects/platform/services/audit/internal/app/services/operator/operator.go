package operator

import (
	"context"

	dto "grovehub/services/audit/internal/app/dto/operator"
	"grovehub/services/audit/internal/app/dto/producer"
)

type Service interface {
	Store(context.Context, producer.OperatorActivity) error
	StoreSchoolReceipt(context.Context, producer.SchoolRenameReceipt) error
	StoreStatusReceipt(context.Context, producer.SchoolStatusReceipt) error
	SchoolCheck(context.Context, string) error
	Search(context.Context, dto.SearchInput) (dto.SearchPage, error)
	SearchFilteredReceipts(context.Context, dto.ReceiptSearchInput) (dto.ReceiptPage, error)
}

package operator

import (
	"time"

	"github.com/google/uuid"

	"grovehub/services/audit/internal/app/constants/operator"
	"grovehub/services/audit/internal/app/dto/producer"
	dto "grovehub/services/audit/internal/app/dto/operator"
	"grovehub/services/audit/internal/app/validators/operator/helpers"
)

func ValidID(value string) bool {
	id, err := uuid.Parse(value)
	return err == nil && id != uuid.Nil && id.String() == value
}

func ValidSearch(input dto.SearchInput) bool {
	return input.PageNumber >= operator.MinPageNumber && input.PageNumber <= operator.MaxPageNumber &&
		input.RowSize >= operator.MinRowSize && input.RowSize <= operator.MaxRowSize
}

func ValidReceiptSearch(input dto.ReceiptSearchInput) bool {
	if !ValidSearch(input.SearchInput) || (input.SchoolID != "" && !ValidID(input.SchoolID)) {
		return false
	}

	switch input.Operation {
	case "", operator.Rename, operator.Activate, operator.Disable:
		return true
	}

	return false
}

func validOccurrenceTime(value time.Time) bool {
	return !value.IsZero() && !value.After(time.Now().Add(5*time.Minute))
}

func validOptionalID(value *string) bool {
	return value == nil || ValidID(*value)
}

func ValidActivity(a producer.OperatorActivity) bool {
	return ValidID(a.EventID) && ValidID(a.ActorID) && ValidID(a.ResourceID) &&
		helpers.ValidOperation(a.Operation) && validOptionalID(a.CorrelationID) && validOccurrenceTime(a.OccurredAt)
}

func ValidRename(a producer.SchoolRenameReceipt) bool {
	return ValidID(a.ExecutionID) && ValidID(a.SchoolID) && a.PreviousVersion >= 1 &&
		validOptionalID(a.CorrelationID) && validOccurrenceTime(a.OccurredAt)
}

func ValidStatus(a producer.SchoolStatusReceipt) bool {
	return ValidID(a.ExecutionID) && ValidID(a.SchoolID) && a.PreviousVersion >= 1 &&
		validOptionalID(a.CorrelationID) && validOccurrenceTime(a.OccurredAt) &&
		(a.Operation == operator.Activate || a.Operation == operator.Disable)
}

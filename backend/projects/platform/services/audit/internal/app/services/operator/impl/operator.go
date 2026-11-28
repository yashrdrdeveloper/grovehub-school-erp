package impl

import (
	"context"
	"strconv"
	"time"

	coreerrors "grovehub/libs/core/errors"
	definitions "grovehub/services/audit/internal/app/constants/error"
	operations "grovehub/services/audit/internal/app/constants/operator"
	dto "grovehub/services/audit/internal/app/dto/operator"
	"grovehub/services/audit/internal/app/dto/producer"
	"grovehub/services/audit/internal/app/entity"
	repository "grovehub/services/audit/internal/app/repository/operator"
	"grovehub/services/audit/internal/app/services/operator"
	validator "grovehub/services/audit/internal/app/validators/operator"
)

type service struct {
	repository repository.Repository
}

func New(repo repository.Repository) operator.Service {
	return &service{repository: repo}
}

func (s *service) Store(ctx context.Context, a producer.OperatorActivity) error {
	if !validator.ValidActivity(a) {
		return coreerrors.NewError(definitions.OperatorInvalid)
	}
	return s.repository.InsertActivity(ctx, entity.OperatorActivity{
		EventID: a.EventID, ActorID: a.ActorID, ResourceID: a.ResourceID,
		Operation: a.Operation, OccurredAt: a.OccurredAt.UTC().Truncate(time.Microsecond), CorrelationID: a.CorrelationID,
	})
}

func (s *service) StoreSchoolReceipt(ctx context.Context, a producer.SchoolRenameReceipt) error {
	if !validator.ValidRename(a) {
		return coreerrors.NewError(definitions.OperatorInvalid)
	}
	return s.storeReceipt(ctx, a, operations.Rename)
}

func (s *service) StoreStatusReceipt(ctx context.Context, a producer.SchoolStatusReceipt) error {
	if !validator.ValidStatus(a) {
		return coreerrors.NewError(definitions.OperatorInvalid)
	}
	return s.storeReceipt(ctx, producer.SchoolRenameReceipt{
		ExecutionID: a.ExecutionID, SchoolID: a.SchoolID, PreviousVersion: a.PreviousVersion,
		CorrelationID: a.CorrelationID, OccurredAt: a.OccurredAt,
	}, a.Operation)
}

func (s *service) storeReceipt(ctx context.Context, a producer.SchoolRenameReceipt, operation string) error {
	return s.repository.InsertReceipt(ctx, entity.SchoolMutationReceipt{
		ExecutionID: a.ExecutionID, SchoolID: a.SchoolID, PreviousVersion: a.PreviousVersion,
		CorrelationID: a.CorrelationID, OccurredAt: a.OccurredAt.UTC().Truncate(time.Microsecond), Operation: operation,
	})
}

func (s *service) SchoolCheck(ctx context.Context, school string) (bool, error) {
	return s.repository.SchoolExists(ctx, school)
}

func (s *service) Search(ctx context.Context, in dto.SearchInput) (dto.SearchPage, error) {
	if !validator.ValidSearch(in) {
		return dto.SearchPage{}, coreerrors.NewError(definitions.OperatorInvalid)
	}

	records, err := s.repository.SearchActivities(ctx, in.RowSize+1, (in.PageNumber-1)*in.RowSize)
	if err != nil {
		return dto.SearchPage{}, err
	}

	page := dto.SearchPage{Records: make([]producer.OperatorActivity, 0), PageNumber: in.PageNumber, RowSize: in.RowSize}
	if len(records) > in.RowSize {
		page.HasMore = true
		records = records[:in.RowSize]
	}

	for _, a := range records {
		page.Records = append(page.Records, producer.OperatorActivity{
			EventID: a.EventID, ActorID: a.ActorID, ResourceID: a.ResourceID,
			Operation: a.Operation, OccurredAt: a.OccurredAt, CorrelationID: a.CorrelationID,
		})
	}
	return page, nil
}

func (s *service) SearchFilteredReceipts(ctx context.Context, in dto.ReceiptSearchInput) (dto.ReceiptPage, error) {
	if !validator.ValidReceiptSearch(in) {
		return dto.ReceiptPage{}, coreerrors.NewError(definitions.OperatorInvalid)
	}

	records, err := s.repository.SearchReceipts(ctx, in.SchoolID, in.Operation, in.RowSize+1, (in.PageNumber-1)*in.RowSize)
	if err != nil {
		return dto.ReceiptPage{}, err
	}

	page := dto.ReceiptPage{Records: make([]dto.ReceiptView, 0), PageNumber: in.PageNumber, RowSize: in.RowSize}
	if len(records) > in.RowSize {
		page.HasMore = true
		records = records[:in.RowSize]
	}

	for _, a := range records {
		page.Records = append(page.Records, dto.ReceiptView{
			ExecutionID: a.ExecutionID, SchoolID: a.SchoolID, Operation: a.Operation,
			PreviousVersion: strconv.FormatInt(a.PreviousVersion, 10), CorrelationID: a.CorrelationID, OccurredAt: a.OccurredAt,
		})
	}
	return page, nil
}

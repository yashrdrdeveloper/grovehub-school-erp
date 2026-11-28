package impl

import (
	"context"
	"net/http"

	"grovehub/libs/core/pipeline"
	definitions "grovehub/services/audit/internal/app/constants/error"
	"grovehub/services/audit/internal/app/controller/operator/constants"
	"grovehub/services/audit/internal/app/controller/operator/helpers"
	restricted "grovehub/services/audit/internal/app/controller/operator/restricted"
	dto "grovehub/services/audit/internal/app/dto/operator"
	service "grovehub/services/audit/internal/app/services/operator"
	validator "grovehub/services/audit/internal/app/validators/operator"
)

type Controller struct {
	service service.Service
}

func New(service service.Service) restricted.OperatorController {
	return &Controller{service: service}
}

func (controller *Controller) Write(w http.ResponseWriter, r *http.Request) {
	helpers.Run(w, r, helpers.Policy{Subject: constants.WriterSubject, Scope: constants.Scope}, pipeline.StructConfig[dto.WriteRequest, pipeline.Empty]{
		Validate: validator.ValidateWrite,
		Service: func(ctx context.Context, input dto.WriteRequest) (pipeline.Result[pipeline.Empty], error) {
			return pipeline.NoContent(), controller.service.Store(ctx, input.RequestBody)
		},
	})
}

func (controller *Controller) Receipt(w http.ResponseWriter, r *http.Request) {
	helpers.Run(w, r, helpers.Policy{Scope: constants.ReceiptScope}, pipeline.StructConfig[dto.RenameReceiptRequest, pipeline.Empty]{
		Validate: helpers.ValidateReceipt(r, validator.ValidateRename, dto.RenameReceiptRequest.SchoolID, constants.ReceiptIssuer, constants.ReceiptScope),
		Service: func(ctx context.Context, input dto.RenameReceiptRequest) (pipeline.Result[pipeline.Empty], error) {
			return pipeline.NoContent(), controller.service.StoreSchoolReceipt(ctx, input.RequestBody)
		},
	})
}

func (controller *Controller) StatusReceipt(w http.ResponseWriter, r *http.Request) {
	helpers.Run(w, r, helpers.Policy{Scope: constants.StatusReceiptScope}, pipeline.StructConfig[dto.StatusReceiptRequest, pipeline.Empty]{
		Validate: helpers.ValidateReceipt(r, validator.ValidateStatus, dto.StatusReceiptRequest.SchoolID, constants.StatusReceiptIssuer, constants.StatusReceiptScope),
		Service: func(ctx context.Context, input dto.StatusReceiptRequest) (pipeline.Result[pipeline.Empty], error) {
			return pipeline.NoContent(), controller.service.StoreStatusReceipt(ctx, input.RequestBody)
		},
	})
}

func (controller *Controller) Search(w http.ResponseWriter, r *http.Request) {
	helpers.Run(w, r, helpers.Policy{Subject: constants.ReaderSubject, Scope: constants.ReadScope}, pipeline.StructConfig[dto.SearchRequest, dto.SearchPage]{
		Validate: validator.ValidateSearch,
		Service: func(ctx context.Context, input dto.SearchRequest) (pipeline.Result[dto.SearchPage], error) {
			page, err := controller.service.Search(ctx, input.RequestBody)
			return pipeline.OK(page), err
		},
	})
}

func (controller *Controller) ReceiptSearch(w http.ResponseWriter, r *http.Request) {
	helpers.Run(w, r, helpers.Policy{Subject: constants.ReceiptReader, Scope: constants.ReceiptReadScope}, pipeline.StructConfig[dto.ReceiptSearchRequest, dto.ReceiptPage]{
		Validate: validator.ValidateReceiptSearch,
		Service: func(ctx context.Context, input dto.ReceiptSearchRequest) (pipeline.Result[dto.ReceiptPage], error) {
			page, err := controller.service.SearchFilteredReceipts(ctx, input.RequestBody)
			return pipeline.OK(page), err
		},
	})
}

func (controller *Controller) SchoolCheck(w http.ResponseWriter, r *http.Request) {
	policy := helpers.Policy{
		Subject: constants.SchoolSubject + r.PathValue(constants.SchoolIDParam),
		Scope:   constants.SchoolCheckScope,
		Timeout: constants.SchoolCheckTimeout,
	}
	helpers.Run(w, r, policy, pipeline.StructConfig[dto.SchoolCheckRequest, pipeline.Empty]{
		InvalidParams: definitions.OperatorForbidden,
		Validate:      validator.ValidateSchoolCheck,
		Service: func(ctx context.Context, input dto.SchoolCheckRequest) (pipeline.Result[pipeline.Empty], error) {
			return pipeline.NoContent(), controller.service.SchoolCheck(ctx, input.Params.SchoolID)
		},
	})
}

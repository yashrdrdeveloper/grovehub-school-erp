package impl

import (
	"context"
	"net/http"

	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
	corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	definitions "grovehub/services/identity/internal/app/constants/error"
	protected "grovehub/services/identity/internal/app/controller/approval/protected"
	"grovehub/services/identity/internal/app/controller/helpers"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/services/approval"
)

type Controller struct {
	service approval.Reader
	mutations approval.Service
	executor approval.Executor
}

var _ protected.ApprovalController = (*Controller)(nil)

func New(service approval.Reader) protected.ApprovalController {
	return NewWithMutations(service, nil, nil)
}

// A nil executor intentionally keeps effect execution disabled during cutover.
func NewWithMutations(service approval.Reader, mutations approval.Service, executor approval.Executor) protected.ApprovalController {
	return &Controller{service: service, mutations: mutations, executor: executor}
}

func (controller *Controller) History(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.SearchRequest, coresearch.Result[dto.Summary]]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.SearchRequest) (corepipeline.Result[coresearch.Result[dto.Summary]], error) {
			result, err := controller.service.History(ctx, input.RequestBody)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Inbox(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.SearchRequest, coresearch.Result[dto.Summary]]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.SearchRequest) (corepipeline.Result[coresearch.Result[dto.Summary]], error) {
			result, err := controller.service.Inbox(ctx, input.RequestBody)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) PendingCount(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.CountRequest, dto.Count]{
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.CountRequest) (corepipeline.Result[dto.Count], error) {
			result, err := controller.service.PendingCount(ctx)
			return corepipeline.OK(dto.Count{Pending: result}), helpers.Map(err)
		},
	})
}

func (controller *Controller) Detail(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.GetRequest, dto.Detail]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (corepipeline.Result[dto.Detail], error) {
			result, err := controller.service.Detail(ctx, input.Params.RequestID)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

package impl

import (
	"context"
	"net/http"

	httprequest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/request"
	rest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/rest"
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
	rest.RunStruct(writer, request, rest.StructConfig[dto.SearchRequest, coresearch.Result[dto.Summary]]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.SearchRequest) (rest.Result[coresearch.Result[dto.Summary]], error) {
			result, err := controller.service.History(ctx, input.RequestBody)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Inbox(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.SearchRequest, coresearch.Result[dto.Summary]]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.SearchRequest) (rest.Result[coresearch.Result[dto.Summary]], error) {
			result, err := controller.service.Inbox(ctx, input.RequestBody)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) PendingCount(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.CountRequest, dto.Count]{
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.CountRequest) (rest.Result[dto.Count], error) {
			result, err := controller.service.PendingCount(ctx)
			return rest.OK(dto.Count{Pending: result}), helpers.Map(err)
		},
	})
}

func (controller *Controller) Detail(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.GetRequest, dto.Detail]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (rest.Result[dto.Detail], error) {
			result, err := controller.service.Detail(ctx, input.Params.RequestID)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Audit(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.GetRequest, []dto.AuditEvent]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (rest.Result[[]dto.AuditEvent], error) {
			result, err := controller.service.Audit(ctx, input.Params.RequestID)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

package impl

import (
	"context"
	"net/http"

	headercontext "grovehub/libs/core/headercontext"
	httprequest "grovehub/libs/core/http/request"
	rest "grovehub/libs/core/pipeline/rest"
	coresearch "grovehub/libs/core/search"
	definitions "grovehub/services/identity/internal/app/constants/error"
	protected "grovehub/services/identity/internal/app/controller/account/protected"
	"grovehub/services/identity/internal/app/controller/helpers"
	dto "grovehub/services/identity/internal/app/dto/account"
	account "grovehub/services/identity/internal/app/services/account"
)

type Controller struct {
	service account.Service
}

var _ protected.AccountController = (*Controller)(nil)

func New(service account.Service) protected.AccountController {
	return &Controller{service: service}
}

func (controller *Controller) Create(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.CreateRequest, dto.CreateResult]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.CreateRequest) (rest.Result[dto.CreateResult], error) {
			result, err := controller.service.Create(ctx, input.RequestBody, headercontext.IdempotencyKey(ctx))
			return rest.Created(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Search(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.SearchRequest, coresearch.Result[dto.Account]]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.SearchRequest) (rest.Result[coresearch.Result[dto.Account]], error) {
			result, err := controller.service.Search(ctx, input.RequestBody)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Get(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.GetRequest, dto.Account]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (rest.Result[dto.Account], error) {
			result, err := controller.service.Get(ctx, input.Params.AccountID)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Update(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.UpdateRequest, dto.Account]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidParams: definitions.InvalidRequest,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.UpdateRequest) (rest.Result[dto.Account], error) {
			result, err := controller.service.Update(ctx, input.Params.AccountID, input.RequestBody)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) transition(writer http.ResponseWriter, request *http.Request, action string) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.GetRequest, dto.Account]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (rest.Result[dto.Account], error) {
			result, err := controller.service.Transition(ctx, input.Params.AccountID, action)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Suspend(writer http.ResponseWriter, request *http.Request) {
	controller.transition(writer, request, "suspend")
}

func (controller *Controller) Activate(writer http.ResponseWriter, request *http.Request) {
	controller.transition(writer, request, "activate")
}

func (controller *Controller) Deactivate(writer http.ResponseWriter, request *http.Request) {
	controller.transition(writer, request, "deactivate")
}

func (controller *Controller) Resend(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.GetRequest, map[string]any]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (rest.Result[map[string]any], error) {
			err := controller.service.Resend(ctx, input.Params.AccountID)
			return rest.Result[map[string]any]{
				StatusCode: http.StatusAccepted,
				Data: map[string]any{
					"accountId": input.Params.AccountID,
					"invitationStatus": "QUEUED",
				},
			}, helpers.Map(err)
		},
	})
}

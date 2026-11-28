package impl

import (
	"context"
	"net/http"

	corehttp "grovehub/libs/core/http"
	corepipeline "grovehub/libs/core/pipeline"
	definitions "grovehub/services/identity/internal/app/constants/error"
	public "grovehub/services/identity/internal/app/controller/activation/public"
	"grovehub/services/identity/internal/app/controller/helpers"
	dto "grovehub/services/identity/internal/app/dto/account"
	authdto "grovehub/services/identity/internal/app/dto/authentication"
	account "grovehub/services/identity/internal/app/services/account"
	"grovehub/services/identity/internal/app/transport/cookies"
)

type Controller struct {
	cookies cookies.Manager
	service account.Service
}

var _ public.ActivationController = (*Controller)(nil)

func New(service account.Service, manager cookies.Manager) public.ActivationController {
	return &Controller{service: service, cookies: manager}
}

func (controller *Controller) SetupStatus(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(cookies.HeaderCacheControl, cookies.DirectiveNoStore)
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[authdto.EmptyRequest, dto.SetupStatusResult]{
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, _ authdto.EmptyRequest) (corepipeline.Result[dto.SetupStatusResult], error) {
			result, err := controller.service.SetupStatus(ctx, controller.cookies.Read(request, cookies.SetupCookie))
			if err != nil {
				return corepipeline.Result[dto.SetupStatusResult]{}, helpers.Map(err)
			}
			return corepipeline.OK(result), nil
		},
	})
}

func (controller *Controller) Setup(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.SetupRequest, dto.SetupResult]{
		Decode:             corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback:           definitions.Internal,
		Service: func(ctx context.Context, input dto.SetupRequest) (corepipeline.Result[dto.SetupResult], error) {
			result, err := controller.service.Setup(ctx, controller.cookies.Read(request, cookies.SetupCookie), input.RequestBody)
			if err != nil {
				return corepipeline.Result[dto.SetupResult]{}, helpers.Map(err)
			}
			controller.cookies.ClearGrant(writer, cookies.SetupCookie, cookies.SetupGrantPath)
			return corepipeline.OK(result), nil
		},
	})
}

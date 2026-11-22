package impl

import (
	"context"
	"net/http"

	httprequest "grovehub/libs/core/http/request"
	rest "grovehub/libs/core/pipeline/rest"
	definitions "grovehub/services/identity/internal/app/constants/error"
	"grovehub/services/identity/internal/app/controller/helpers"
	public "grovehub/services/identity/internal/app/controller/activation/public"
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
  rest.RunStruct(writer, request, rest.StructConfig[authdto.EmptyRequest, dto.SetupStatusResult]{
    Fallback: definitions.Internal,
    Service: func(ctx context.Context, _ authdto.EmptyRequest) (rest.Result[dto.SetupStatusResult], error) {
      result, err := controller.service.SetupStatus(ctx, controller.cookies.Read(request, cookies.SetupCookie))
      if err != nil { return rest.Result[dto.SetupStatusResult]{}, helpers.Map(err) }
      return rest.OK(result), nil
    },
  })
}

func (controller *Controller) Setup(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.SetupRequest, dto.SetupResult]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.SetupRequest) (rest.Result[dto.SetupResult], error) {
			result, err := controller.service.Setup(ctx, controller.cookies.Read(request, cookies.SetupCookie), input.RequestBody)
			if err != nil {
				return rest.Result[dto.SetupResult]{}, helpers.Map(err)
			}
			controller.cookies.ClearGrant(writer, cookies.SetupCookie, cookies.SetupGrantPath)
			return rest.OK(result), nil
		},
	})
}

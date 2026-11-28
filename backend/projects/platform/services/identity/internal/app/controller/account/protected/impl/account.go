package impl

import (
	"context"
	"net/http"

	headercontext "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext"
	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
	corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	definitions "grovehub/services/identity/internal/app/constants/error"
	protected "grovehub/services/identity/internal/app/controller/account/protected"
	"grovehub/services/identity/internal/app/controller/helpers"
	"grovehub/services/identity/internal/app/services/authorization"
	dto "grovehub/services/identity/internal/app/dto/account"
	account "grovehub/services/identity/internal/app/services/account"
	roles "grovehub/services/identity/internal/app/services/account/operations/roles"
	profile "grovehub/services/identity/internal/app/services/account/operations/personprofile"
)

type Controller struct {
	service account.Service
	catalog *roles.Operation
	profile *profile.Operation
}

var _ protected.AccountController = (*Controller)(nil)

func New(service account.Service, catalogs ...*roles.Operation) protected.AccountController {
	controller := &Controller{service: service}
	if len(catalogs) == 1 {
		controller.catalog = catalogs[0]
	}
	return controller
}

func NewWithProfile(service account.Service, catalog *roles.Operation, operation *profile.Operation) protected.AccountController {
 return &Controller{service:service,catalog:catalog,profile:operation}
}

func (controller *Controller) PersonProfile(writer http.ResponseWriter, request *http.Request) {
 writer.Header().Set("Cache-Control","no-store")
 corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.GetRequest,dto.PersonProfile]{
  InvalidParams: definitions.InvalidRequest,
  Fallback: definitions.Internal,
  Service: func(ctx context.Context,input dto.GetRequest)(corepipeline.Result[dto.PersonProfile],error){
   if controller.profile==nil{return corepipeline.Result[dto.PersonProfile]{},helpers.Map(authorization.ErrUnavailable)}
   result,err:=controller.profile.Execute(ctx,input.Params.AccountID)
   return corepipeline.OK(result),helpers.Map(err)
  },
 })
}

func (controller *Controller) Create(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.CreateRequest, dto.CreateResult]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.CreateRequest) (corepipeline.Result[dto.CreateResult], error) {
			result, err := controller.service.Create(ctx, input.RequestBody, headercontext.IdempotencyKey(ctx))
			return corepipeline.Created(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Search(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.SearchRequest, coresearch.Result[dto.Account]]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.SearchRequest) (corepipeline.Result[coresearch.Result[dto.Account]], error) {
			result, err := controller.service.Search(ctx, input.RequestBody)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Roles(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.RolesRequest, []dto.Role]{
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.RolesRequest) (corepipeline.Result[[]dto.Role], error) {
			if controller.catalog == nil {
				return corepipeline.Result[[]dto.Role]{}, helpers.Map(authorization.ErrUnavailable)
			}
			result, err := controller.catalog.Execute(ctx)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Get(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.GetRequest, dto.Account]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (corepipeline.Result[dto.Account], error) {
			result, err := controller.service.Get(ctx, input.Params.AccountID)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Update(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.UpdateRequest, dto.Account]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidParams: definitions.InvalidRequest,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.UpdateRequest) (corepipeline.Result[dto.Account], error) {
			result, err := controller.service.Update(ctx, input.Params.AccountID, input.RequestBody)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) transition(writer http.ResponseWriter, request *http.Request, action string) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.GetRequest, dto.Account]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (corepipeline.Result[dto.Account], error) {
			result, err := controller.service.Transition(ctx, input.Params.AccountID, action)
			return corepipeline.OK(result), helpers.Map(err)
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
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.GetRequest, map[string]any]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (corepipeline.Result[map[string]any], error) {
			err := controller.service.Resend(ctx, input.Params.AccountID)
			return corepipeline.Result[map[string]any]{
				StatusCode: http.StatusAccepted,
				Data: map[string]any{
					"accountId": input.Params.AccountID,
					"invitationStatus": "QUEUED",
				},
			}, helpers.Map(err)
		},
	})
}

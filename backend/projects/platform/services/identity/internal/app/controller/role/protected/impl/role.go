package impl

import (
	"context"
	"net/http"

	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
	rest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/rest"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	definitions "grovehub/services/identity/internal/app/constants/error"
	protected "grovehub/services/identity/internal/app/controller/role/protected"
	"grovehub/services/identity/internal/app/controller/helpers"
	dto "grovehub/services/identity/internal/app/dto/role"
	"grovehub/services/identity/internal/app/services/role"
)

type Controller struct {
	service role.Service
}

var _ protected.RoleController = (*Controller)(nil)

func New(service role.Service) protected.RoleController {
	return &Controller{service: service}
}

func (controller *Controller) Search(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.SearchRequest, coresearch.Result[dto.Detail]]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.SearchRequest) (rest.Result[coresearch.Result[dto.Detail]], error) {
			result, err := controller.service.Search(ctx, input.RequestBody)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Get(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.GetRequest, dto.Detail]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (rest.Result[dto.Detail], error) {
			result, err := controller.service.Get(ctx, input.Params.RoleID)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Permissions(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.GetRequest, dto.PermissionSet]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (rest.Result[dto.PermissionSet], error) {
			result, err := controller.service.Permissions(ctx, input.Params.RoleID)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) PreviewPermissions(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.PreviewRequest, dto.PermissionPreview]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.PreviewRequest) (rest.Result[dto.PermissionPreview], error) {
			result, err := controller.service.PreviewPermissions(ctx, input.Params.RoleID, input.RequestBody)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) CreateDraft(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.CreateRequest, dto.Detail]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.CreateRequest) (rest.Result[dto.Detail], error) {
			result, err := controller.service.CreateDraft(ctx, input.RequestBody)
			return rest.Created(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) UpdateMetadata(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.UpdateRequest, dto.Detail]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.UpdateRequest) (rest.Result[dto.Detail], error) {
			result, err := controller.service.UpdateMetadata(ctx, input.Params.RoleID, input.RequestBody)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

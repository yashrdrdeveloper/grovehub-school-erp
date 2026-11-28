package impl

import (
	"context"
	"net/http"

	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
	corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
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
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.SearchRequest, coresearch.Result[dto.Detail]]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.SearchRequest) (corepipeline.Result[coresearch.Result[dto.Detail]], error) {
			result, err := controller.service.Search(ctx, input.RequestBody)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Get(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.GetRequest, dto.Detail]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (corepipeline.Result[dto.Detail], error) {
			result, err := controller.service.Get(ctx, input.Params.RoleID)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Permissions(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.GetRequest, dto.PermissionSet]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (corepipeline.Result[dto.PermissionSet], error) {
			result, err := controller.service.Permissions(ctx, input.Params.RoleID)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) PreviewPermissions(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.PreviewRequest, dto.PermissionPreview]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.PreviewRequest) (corepipeline.Result[dto.PermissionPreview], error) {
			result, err := controller.service.PreviewPermissions(ctx, input.Params.RoleID, input.RequestBody)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) CreateDraft(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.CreateRequest, dto.Detail]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.CreateRequest) (corepipeline.Result[dto.Detail], error) {
			result, err := controller.service.CreateDraft(ctx, input.RequestBody)
			return corepipeline.Created(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) UpdateMetadata(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.UpdateRequest, dto.Detail]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.UpdateRequest) (corepipeline.Result[dto.Detail], error) {
			result, err := controller.service.UpdateMetadata(ctx, input.Params.RoleID, input.RequestBody)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

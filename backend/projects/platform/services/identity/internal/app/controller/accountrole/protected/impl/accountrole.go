package impl

import (
	"context"
	"net/http"

	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
	corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
	definitions "grovehub/services/identity/internal/app/constants/error"
	protected "grovehub/services/identity/internal/app/controller/accountrole/protected"
	"grovehub/services/identity/internal/app/controller/helpers"
	dto "grovehub/services/identity/internal/app/dto/accountrole"
	"grovehub/services/identity/internal/app/services/accountrole"
)

type Controller struct {
	service accountrole.Service
}

var _ protected.AccountRoleController = (*Controller)(nil)

func New(service accountrole.Service) protected.AccountRoleController {
	return &Controller{service: service}
}

func (controller *Controller) Get(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.GetRequest, dto.Assignment]{
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.GetRequest) (corepipeline.Result[dto.Assignment], error) {
			result, err := controller.service.Get(ctx, input.Params.AccountID)
			return corepipeline.OK(dto.Assignment{AccountID: result.AccountID, Status: result.Status, Version: result.Version, PrimaryRoleID: result.PrimaryRoleID, SecondaryRoleIDs: result.SecondaryRoleIDs}), helpers.Map(err)
		},
	})
}

func (controller *Controller) Preview(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.PreviewRequest, dto.Preview]{
		Decode: corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		InvalidParams: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.PreviewRequest) (corepipeline.Result[dto.Preview], error) {
			result, err := controller.service.Preview(ctx, input.Params.AccountID, input.RequestBody)
			return corepipeline.OK(result), helpers.Map(err)
		},
	})
}

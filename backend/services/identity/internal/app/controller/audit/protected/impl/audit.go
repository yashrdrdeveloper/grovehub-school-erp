package impl

import (
	"context"
	"net/http"

	httprequest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/request"
	rest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/rest"
	definitions "grovehub/services/identity/internal/app/constants/error"
	protected "grovehub/services/identity/internal/app/controller/audit/protected"
	"grovehub/services/identity/internal/app/controller/helpers"
	dto "grovehub/services/identity/internal/app/dto/audit"
	"grovehub/services/identity/internal/app/services/audit"
)

type Controller struct { service audit.Service }

func New(service audit.Service) protected.AuditController {
	return &Controller{service: service}
}

func (controller *Controller) Search(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Cache-Control", "no-store")
	rest.RunStruct(writer, request, rest.StructConfig[dto.SearchRequest, dto.Page]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.SearchRequest) (rest.Result[dto.Page], error) {
			result, err := controller.service.Search(ctx, input.RequestBody)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

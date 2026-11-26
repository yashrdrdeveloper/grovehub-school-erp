package impl

import (
	"context"
	"errors"
	"net/http"

	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	httprequest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/request"
	rest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/rest"
	definitions "grovehub/services/audit/internal/app/constants/error"
	protected "grovehub/services/audit/internal/app/controller/activity/protected"
	dto "grovehub/services/audit/internal/app/dto/activity"
	"grovehub/services/audit/internal/app/services/activity"
)

type Controller struct { service activity.Service }

func New(service activity.Service) protected.ActivityController {
	return &Controller{service:service}
}

func (controller *Controller) Search(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Cache-Control","no-store")
	rest.RunStruct(writer,request,rest.StructConfig[dto.SearchRequest,dto.Page]{
		Decode:httprequest.DecodeStrictJSON,
		InvalidRequestBody:definitions.InvalidRequest,
		Fallback:definitions.Internal,
		Service:func(ctx context.Context,input dto.SearchRequest)(rest.Result[dto.Page],error){
			result,err:=controller.service.Search(ctx,input.RequestBody)
			switch {
			case errors.Is(err,activity.ErrForbidden): err=coreerrors.NewError(definitions.Forbidden)
			case errors.Is(err,activity.ErrInvalidSearch): err=coreerrors.NewError(definitions.InvalidRequest)
			case err!=nil: err=coreerrors.NewError(definitions.Internal)
			}
			return rest.OK(result),err
		},
	})
}

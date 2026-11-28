package impl

import (
	"context"
	"errors"
	"net/http"

	coreerrors "grovehub/libs/core/errors"
	corehttp "grovehub/libs/core/http"
	corepipeline "grovehub/libs/core/pipeline"
	definitions "grovehub/services/audit/internal/app/constants/error"
	protected "grovehub/services/audit/internal/app/controller/activity/protected"
	dto "grovehub/services/audit/internal/app/dto/activity"
	"grovehub/services/audit/internal/app/services/activity"
)

type Controller struct{ service activity.Service }

func New(service activity.Service) protected.ActivityController {
	return &Controller{service: service}
}

func (controller *Controller) Search(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Cache-Control", "no-store")
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.SearchRequest, dto.Page]{
		Decode:             corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback:           definitions.Internal,
		Service: func(ctx context.Context, input dto.SearchRequest) (corepipeline.Result[dto.Page], error) {
			result, err := controller.service.Search(ctx, input.RequestBody)
			switch {
			case errors.Is(err, activity.ErrForbidden):
				err = coreerrors.NewError(definitions.Forbidden)
			case errors.Is(err, activity.ErrInvalidSearch):
				err = coreerrors.NewError(definitions.InvalidRequest)
			case err != nil:
				err = coreerrors.NewError(definitions.Internal)
			}
			return corepipeline.OK(result), err
		},
	})
}

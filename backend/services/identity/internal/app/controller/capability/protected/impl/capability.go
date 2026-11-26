package impl

import (
	"context"
	"net/http"

	rest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/rest"
	definitions "grovehub/services/identity/internal/app/constants/error"
	protected "grovehub/services/identity/internal/app/controller/capability/protected"
	"grovehub/services/identity/internal/app/controller/helpers"
	authdto "grovehub/services/identity/internal/app/dto/authentication"
	dto "grovehub/services/identity/internal/app/dto/capability"
	"grovehub/services/identity/internal/app/services/capability"
)

type Controller struct {
	service capability.Service
}

var _ protected.CapabilityController = (*Controller)(nil)

func New(service capability.Service) protected.CapabilityController {
	return &Controller{service: service}
}

func (controller *Controller) Discover(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Cache-Control", "no-store")
	rest.RunStruct(writer, request, rest.StructConfig[authdto.EmptyRequest, dto.Discovery]{
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, _ authdto.EmptyRequest) (rest.Result[dto.Discovery], error) {
			result, err := controller.service.Discover(ctx)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) Catalog(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Cache-Control", "no-store")
	rest.RunStruct(writer, request, rest.StructConfig[authdto.EmptyRequest, dto.Catalog]{
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, _ authdto.EmptyRequest) (rest.Result[dto.Catalog], error) {
			result, err := controller.service.Catalog(ctx)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

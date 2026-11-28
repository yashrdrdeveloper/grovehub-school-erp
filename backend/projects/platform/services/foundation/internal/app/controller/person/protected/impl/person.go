package impl

import (
	"context"
	"net/http"

	corepipeline "grovehub/libs/core/pipeline"
	coresearch "grovehub/libs/core/search"
	personerror "grovehub/services/foundation/internal/app/constants/error"
	protected "grovehub/services/foundation/internal/app/controller/person/protected"
	persondto "grovehub/services/foundation/internal/app/dto/person"
	foundationservices "grovehub/services/foundation/internal/app/services/person"
	foundationvalidators "grovehub/services/foundation/internal/app/validators/person"
)

type PersonController struct {
	service   foundationservices.PersonService
	validator foundationvalidators.PersonValidator
}

func NewPersonController(service foundationservices.PersonService, validator foundationvalidators.PersonValidator) protected.PersonController {
	return &PersonController{service: service, validator: validator}
}

func (controller *PersonController) Create(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[persondto.CreatePersonRequest, persondto.Person]{
		Rules:              controller.validator.CreateRules(),
		InvalidRequestBody: personerror.InvalidPersonRequest,
		Fallback:           personerror.PersonCreateFailed,
		Service: func(ctx context.Context, input persondto.CreatePersonRequest) (corepipeline.Result[persondto.Person], error) {
			person, created, err := controller.service.Create(ctx, input)
			return corepipeline.CreatedOrOK(person, created), err
		},
	})
}

func (controller *PersonController) GetByGID(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[persondto.GetPersonByGIDRequest, persondto.Person]{
		InvalidParams: personerror.InvalidPersonGID,
		Fallback:      personerror.PersonReadFailed,
		Service: func(ctx context.Context, input persondto.GetPersonByGIDRequest) (corepipeline.Result[persondto.Person], error) {
			person, err := controller.service.GetByGID(ctx, input)
			return corepipeline.OK(person), err
		},
	})
}

func (controller *PersonController) Update(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[persondto.UpdatePersonRequest, persondto.Person]{
		Rules:              controller.validator.UpdateRules(),
		InvalidParams:      personerror.InvalidPersonGID,
		InvalidRequestBody: personerror.InvalidPersonRequest,
		Fallback:           personerror.PersonUpdateFailed,
		Service: func(ctx context.Context, input persondto.UpdatePersonRequest) (corepipeline.Result[persondto.Person], error) {
			person, err := controller.service.Update(ctx, input)
			return corepipeline.OK(person), err
		},
	})
}

func (controller *PersonController) Search(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[persondto.SearchPersonsRequest, coresearch.Result[persondto.Person]]{
		Rules:              controller.validator.SearchRules(),
		InvalidRequestBody: personerror.InvalidPersonSearch,
		Fallback:           personerror.PersonSearchFailed,
		Service: func(ctx context.Context, input persondto.SearchPersonsRequest) (corepipeline.Result[coresearch.Result[persondto.Person]], error) {
			result, err := controller.service.Search(ctx, input)
			return corepipeline.OK(result), err
		},
	})
}

func (controller *PersonController) Deactivate(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[persondto.DeactivatePersonRequest, persondto.Person]{
		InvalidParams: personerror.InvalidPersonGID,
		Fallback:      personerror.PersonUpdateFailed,
		Service: func(ctx context.Context, input persondto.DeactivatePersonRequest) (corepipeline.Result[persondto.Person], error) {
			person, err := controller.service.Deactivate(ctx, input)
			return corepipeline.OK(person), err
		},
	})
}

package impl

import (
	"context"
	"net/http"

	rest "grovehub/libs/core/pipeline/rest"
	coresearch "grovehub/libs/core/search"
	personerror "grovehub/services/foundation/constants/error"
	protected "grovehub/services/foundation/controller/protected"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	foundationservices "grovehub/services/foundation/services"
	foundationvalidators "grovehub/services/foundation/validators"
)

type PersonController struct {
	service foundationservices.PersonService
	validator foundationvalidators.PersonValidator
}

func NewPersonController(service foundationservices.PersonService, validator foundationvalidators.PersonValidator) protected.PersonController {
	return &PersonController{service: service, validator: validator}
}

func (controller *PersonController) Create(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[foundationrequest.CreatePersonRequest, foundationresponse.Person]{
		Rules: controller.validator.CreateRules(),
		InvalidRequestBody: personerror.InvalidPersonRequest,
		Fallback: personerror.PersonCreateFailed,
		Service: func(ctx context.Context, input foundationrequest.CreatePersonRequest) (rest.Result[foundationresponse.Person], error) {
			person, created, err := controller.service.Create(ctx, input)
			return rest.CreatedOrOK(person, created), err
		},
	})
}

func (controller *PersonController) GetByGID(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[foundationrequest.GetPersonByGIDRequest, foundationresponse.Person]{
		InvalidParams: personerror.InvalidPersonGID,
		Fallback: personerror.PersonReadFailed,
		Service: func(ctx context.Context, input foundationrequest.GetPersonByGIDRequest) (rest.Result[foundationresponse.Person], error) {
			person, err := controller.service.GetByGID(ctx, input)
			return rest.OK(person), err
		},
	})
}

func (controller *PersonController) Update(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[foundationrequest.UpdatePersonRequest, foundationresponse.Person]{
		Rules: controller.validator.UpdateRules(),
		InvalidParams: personerror.InvalidPersonGID,
		InvalidRequestBody: personerror.InvalidPersonRequest,
		Fallback: personerror.PersonUpdateFailed,
		Service: func(ctx context.Context, input foundationrequest.UpdatePersonRequest) (rest.Result[foundationresponse.Person], error) {
			person, err := controller.service.Update(ctx, input)
			return rest.OK(person), err
		},
	})
}

func (controller *PersonController) Search(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[foundationrequest.SearchPersonsRequest, coresearch.Result[foundationresponse.Person]]{
		Rules: controller.validator.SearchRules(),
		InvalidRequestBody: personerror.InvalidPersonSearch,
		Fallback: personerror.PersonSearchFailed,
		Service: func(ctx context.Context, input foundationrequest.SearchPersonsRequest) (rest.Result[coresearch.Result[foundationresponse.Person]], error) {
			result, err := controller.service.Search(ctx, input)
			return rest.OK(result), err
		},
	})
}

func (controller *PersonController) Deactivate(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[foundationrequest.DeactivatePersonRequest, foundationresponse.Person]{
		InvalidParams: personerror.InvalidPersonGID,
		Fallback: personerror.PersonUpdateFailed,
		Service: func(ctx context.Context, input foundationrequest.DeactivatePersonRequest) (rest.Result[foundationresponse.Person], error) {
			person, err := controller.service.Deactivate(ctx, input)
			return rest.OK(person), err
		},
	})
}

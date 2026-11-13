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
	rest.Run(writer, request, rest.Config[rest.Empty, rest.Empty, foundationrequest.CreatePerson, foundationresponse.Person]{
		Body: controller.validator.CreateRules(),
		InvalidBody: personerror.InvalidPersonRequest,
		Fallback: personerror.PersonCreateFailed,
		Service: func(ctx context.Context, input rest.Input[rest.Empty, rest.Empty, foundationrequest.CreatePerson]) (rest.Result[foundationresponse.Person], error) {
			person, created, err := controller.service.Create(ctx, input.Body)
			return rest.CreatedOrOK(person, created), err
		},
	})
}

func (controller *PersonController) GetByGID(writer http.ResponseWriter, request *http.Request) {
	rest.Run(writer, request, rest.Config[foundationrequest.PersonByGID, rest.Empty, rest.Empty, foundationresponse.Person]{
		InvalidPath: personerror.InvalidPersonGID,
		Fallback: personerror.PersonReadFailed,
		Service: func(ctx context.Context, input rest.Input[foundationrequest.PersonByGID, rest.Empty, rest.Empty]) (rest.Result[foundationresponse.Person], error) {
			person, err := controller.service.GetByGID(ctx, input.Path)
			return rest.OK(person), err
		},
	})
}

func (controller *PersonController) Update(writer http.ResponseWriter, request *http.Request) {
	rest.Run(writer, request, rest.Config[foundationrequest.PersonByGID, rest.Empty, foundationrequest.UpdatePerson, foundationresponse.Person]{
		Body: controller.validator.UpdateRules(),
		InvalidPath: personerror.InvalidPersonGID,
		InvalidBody: personerror.InvalidPersonRequest,
		Fallback: personerror.PersonUpdateFailed,
		Service: func(ctx context.Context, input rest.Input[foundationrequest.PersonByGID, rest.Empty, foundationrequest.UpdatePerson]) (rest.Result[foundationresponse.Person], error) {
			person, err := controller.service.Update(ctx, input.Path, input.Body)
			return rest.OK(person), err
		},
	})
}

func (controller *PersonController) Search(writer http.ResponseWriter, request *http.Request) {
	rest.Run(writer, request, rest.Config[rest.Empty, rest.Empty, coresearch.Request, coresearch.Result[foundationresponse.Person]]{
		Body: controller.validator.SearchRules(),
		InvalidBody: personerror.InvalidPersonSearch,
		Fallback: personerror.PersonSearchFailed,
		Service: func(ctx context.Context, input rest.Input[rest.Empty, rest.Empty, coresearch.Request]) (rest.Result[coresearch.Result[foundationresponse.Person]], error) {
			result, err := controller.service.Search(ctx, input.Body)
			return rest.OK(result), err
		},
	})
}

func (controller *PersonController) Deactivate(writer http.ResponseWriter, request *http.Request) {
	rest.Run(writer, request, rest.Config[foundationrequest.PersonByGID, rest.Empty, rest.Empty, foundationresponse.Person]{
		InvalidPath: personerror.InvalidPersonGID,
		Fallback: personerror.PersonUpdateFailed,
		Service: func(ctx context.Context, input rest.Input[foundationrequest.PersonByGID, rest.Empty, rest.Empty]) (rest.Result[foundationresponse.Person], error) {
			person, err := controller.service.Deactivate(ctx, input.Path)
			return rest.OK(person), err
		},
	})
}

package impl

import (
	"context"
	"net/http"
	"time"

	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
	definitions "grovehub/services/audit/internal/app/constants/error"
	"grovehub/services/audit/internal/app/controller/operator"
	"grovehub/services/audit/internal/app/controller/operator/constants"
	"grovehub/services/audit/internal/app/controller/operator/helpers"
	dto "grovehub/services/audit/internal/app/dto/operator"
	service "grovehub/services/audit/internal/app/services/operator"
	validator "grovehub/services/audit/internal/app/validators/operator"
)

type controller struct {
	service service.Service
}

func New(s service.Service) operator.Controller {
	return &controller{service: s}
}

func authorize(w http.ResponseWriter, r *http.Request, subject, scope string) bool {
	w.Header().Set("Cache-Control", "no-store")
	if !helpers.Authorized(r, subject, scope) {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}

func noContent(err error) (pipeline.Result[pipeline.Empty], error) {
	return pipeline.Result[pipeline.Empty]{StatusCode: http.StatusNoContent}, helpers.Failure(err)
}

func (c *controller) Write(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r, constants.WriterSubject, constants.Scope) {
		return
	}

	pipeline.RunStruct(w, r, pipeline.StructConfig[dto.WriteRequest, pipeline.Empty]{
		Decode: helpers.Decode(w),
		InvalidRequestBody: definitions.OperatorInvalid,
		Fallback: definitions.OperatorUnavailable,
		Writer: helpers.Writer[pipeline.Empty](),
		Validate: func(input dto.WriteRequest) error {
			if !validator.ValidActivity(input.RequestBody) {
				return coreerrors.NewError(definitions.OperatorInvalid)
			}
			return nil
		},
		Service: func(ctx context.Context, input dto.WriteRequest) (pipeline.Result[pipeline.Empty], error) {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			return noContent(c.service.Store(ctx, input.RequestBody))
		},
	})
}

func (c *controller) Receipt(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r, "", constants.ReceiptScope) {
		return
	}

	pipeline.RunStruct(w, r, pipeline.StructConfig[dto.RenameReceiptRequest, pipeline.Empty]{
		Decode: helpers.Decode(w),
		InvalidRequestBody: definitions.OperatorInvalid,
		Fallback: definitions.OperatorUnavailable,
		Writer: helpers.Writer[pipeline.Empty](),
		Validate: func(input dto.RenameReceiptRequest) error {
			if !validator.ValidRename(input.RequestBody) {
				return coreerrors.NewError(definitions.OperatorInvalid)
			}
			return nil
		},
		Service: func(ctx context.Context, input dto.RenameReceiptRequest) (pipeline.Result[pipeline.Empty], error) {
			if !helpers.Authorized(r, constants.ReceiptIssuer+":"+input.RequestBody.SchoolID, constants.ReceiptScope) {
				return noContent(coreerrors.NewError(definitions.OperatorForbidden))
			}
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			return noContent(c.service.StoreSchoolReceipt(ctx, input.RequestBody))
		},
	})
}

func (c *controller) StatusReceipt(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r, "", constants.StatusReceiptScope) {
		return
	}

	pipeline.RunStruct(w, r, pipeline.StructConfig[dto.StatusReceiptRequest, pipeline.Empty]{
		Decode: helpers.Decode(w),
		InvalidRequestBody: definitions.OperatorInvalid,
		Fallback: definitions.OperatorUnavailable,
		Writer: helpers.Writer[pipeline.Empty](),
		Validate: func(input dto.StatusReceiptRequest) error {
			if !validator.ValidStatus(input.RequestBody) {
				return coreerrors.NewError(definitions.OperatorInvalid)
			}
			return nil
		},
		Service: func(ctx context.Context, input dto.StatusReceiptRequest) (pipeline.Result[pipeline.Empty], error) {
			if !helpers.Authorized(r, constants.StatusReceiptIssuer+":"+input.RequestBody.SchoolID, constants.StatusReceiptScope) {
				return noContent(coreerrors.NewError(definitions.OperatorForbidden))
			}
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			return noContent(c.service.StoreStatusReceipt(ctx, input.RequestBody))
		},
	})
}

func (c *controller) Search(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r, constants.ReaderSubject, constants.ReadScope) {
		return
	}

	pipeline.RunStruct(w, r, pipeline.StructConfig[dto.SearchRequest, dto.SearchPage]{
		Decode: helpers.Decode(w),
		InvalidRequestBody: definitions.OperatorInvalid,
		Fallback: definitions.OperatorUnavailable,
		Writer: helpers.Writer[dto.SearchPage](),
		Validate: func(input dto.SearchRequest) error {
			if !validator.ValidSearch(input.RequestBody) {
				return coreerrors.NewError(definitions.OperatorInvalid)
			}
			return nil
		},
		Service: func(ctx context.Context, input dto.SearchRequest) (pipeline.Result[dto.SearchPage], error) {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			page, err := c.service.Search(ctx, input.RequestBody)
			return pipeline.OK(page), helpers.Failure(err)
		},
	})
}

func (c *controller) ReceiptSearch(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r, constants.ReceiptReader, constants.ReceiptReadScope) {
		return
	}

	pipeline.RunStruct(w, r, pipeline.StructConfig[dto.ReceiptSearchRequest, dto.ReceiptPage]{
		Decode: helpers.Decode(w),
		InvalidRequestBody: definitions.OperatorInvalid,
		Fallback: definitions.OperatorUnavailable,
		Writer: helpers.Writer[dto.ReceiptPage](),
		Validate: func(input dto.ReceiptSearchRequest) error {
			if !validator.ValidReceiptSearch(input.RequestBody) {
				return coreerrors.NewError(definitions.OperatorInvalid)
			}
			return nil
		},
		Service: func(ctx context.Context, input dto.ReceiptSearchRequest) (pipeline.Result[dto.ReceiptPage], error) {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			page, err := c.service.SearchFilteredReceipts(ctx, input.RequestBody)
			return pipeline.OK(page), helpers.Failure(err)
		},
	})
}

func (c *controller) SchoolCheck(w http.ResponseWriter, r *http.Request) {
	school := r.PathValue(constants.SchoolIDParam)
	if !authorize(w, r, constants.SchoolSubject+school, constants.SchoolCheckScope) {
		return
	}
	if !validator.ValidID(school) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	exists, err := c.service.SchoolCheck(ctx, school)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

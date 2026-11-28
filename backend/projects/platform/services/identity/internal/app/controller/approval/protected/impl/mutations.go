package impl

import (
	"context"
	"net/http"

	headercontext "grovehub/libs/core/headercontext"
	corehttp "grovehub/libs/core/http"
	corepipeline "grovehub/libs/core/pipeline"
	kind "grovehub/services/identity/internal/app/constants/approval"
	definitions "grovehub/services/identity/internal/app/constants/error"
	"grovehub/services/identity/internal/app/controller/helpers"
	dto "grovehub/services/identity/internal/app/dto/approval"
	approvalhelpers "grovehub/services/identity/internal/app/services/approval/helpers"
	"grovehub/services/identity/internal/app/services/authorization"
)

func (controller *Controller) SubmitLifecycle(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.LifecycleRequest, dto.Summary]{
		Decode:             corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback:           definitions.Internal,
		Service: func(ctx context.Context, input dto.LifecycleRequest) (corepipeline.Result[dto.Summary], error) {
			if controller.mutations == nil {
				return corepipeline.Result[dto.Summary]{}, helpers.Map(authorization.ErrUnavailable)
			}
			input.RequestBody.IdempotencyKey = headercontext.IdempotencyKey(ctx)
			result, err := controller.mutations.CreateLifecycle(ctx, input.RequestBody)
			return corepipeline.Created(approvalhelpers.ToSummary(result)), helpers.Map(err)
		},
	})
}

func (controller *Controller) SubmitPermissions(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.PermissionRequest, dto.Summary]{
		Decode:             corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback:           definitions.Internal,
		Service: func(ctx context.Context, input dto.PermissionRequest) (corepipeline.Result[dto.Summary], error) {
			if controller.mutations == nil {
				return corepipeline.Result[dto.Summary]{}, helpers.Map(authorization.ErrUnavailable)
			}
			input.RequestBody.IdempotencyKey = headercontext.IdempotencyKey(ctx)
			result, err := controller.mutations.CreatePermissions(ctx, input.RequestBody)
			return corepipeline.Created(approvalhelpers.ToSummary(result)), helpers.Map(err)
		},
	})
}

func (controller *Controller) SubmitAssignment(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.AssignmentRequest, dto.Summary]{
		Decode:             corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback:           definitions.Internal,
		Service: func(ctx context.Context, input dto.AssignmentRequest) (corepipeline.Result[dto.Summary], error) {
			if controller.mutations == nil {
				return corepipeline.Result[dto.Summary]{}, helpers.Map(authorization.ErrUnavailable)
			}
			input.RequestBody.IdempotencyKey = headercontext.IdempotencyKey(ctx)
			result, err := controller.mutations.CreateAssignment(ctx, input.RequestBody)
			return corepipeline.Created(approvalhelpers.ToSummary(result)), helpers.Map(err)
		},
	})
}

func (controller *Controller) Approve(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.DecisionRequest, dto.Summary]{
		Decode:             corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		InvalidParams:      definitions.InvalidRequest,
		Fallback:           definitions.Internal,
		Service: func(ctx context.Context, input dto.DecisionRequest) (corepipeline.Result[dto.Summary], error) {
			if controller.mutations == nil {
				return corepipeline.Result[dto.Summary]{}, helpers.Map(authorization.ErrUnavailable)
			}
			result, err := controller.mutations.Decide(ctx, dto.Review{RequestID: input.Params.RequestID, Version: input.RequestBody.Version, Decision: kind.Approve, Comment: input.RequestBody.Comment})
			return corepipeline.OK(approvalhelpers.ToSummary(result)), helpers.Map(err)
		},
	})
}

func (controller *Controller) Reject(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.DecisionRequest, dto.Summary]{
		Decode:             corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		InvalidParams:      definitions.InvalidRequest,
		Fallback:           definitions.Internal,
		Service: func(ctx context.Context, input dto.DecisionRequest) (corepipeline.Result[dto.Summary], error) {
			if controller.mutations == nil {
				return corepipeline.Result[dto.Summary]{}, helpers.Map(authorization.ErrUnavailable)
			}
			result, err := controller.mutations.Decide(ctx, dto.Review{RequestID: input.Params.RequestID, Version: input.RequestBody.Version, Decision: kind.Reject, Comment: input.RequestBody.Comment})
			return corepipeline.OK(approvalhelpers.ToSummary(result)), helpers.Map(err)
		},
	})
}

func (controller *Controller) Cancel(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.VersionRequest, dto.Summary]{
		Decode:             corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		InvalidParams:      definitions.InvalidRequest,
		Fallback:           definitions.Internal,
		Service: func(ctx context.Context, input dto.VersionRequest) (corepipeline.Result[dto.Summary], error) {
			if controller.mutations == nil {
				return corepipeline.Result[dto.Summary]{}, helpers.Map(authorization.ErrUnavailable)
			}
			result, err := controller.mutations.Cancel(ctx, input.Params.RequestID, input.RequestBody.Version)
			return corepipeline.OK(approvalhelpers.ToSummary(result)), helpers.Map(err)
		},
	})
}

func (controller *Controller) Execute(writer http.ResponseWriter, request *http.Request) {
	corepipeline.RunStruct(writer, request, corepipeline.StructConfig[dto.VersionRequest, dto.Summary]{
		Decode:             corehttp.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		InvalidParams:      definitions.InvalidRequest,
		Fallback:           definitions.Internal,
		Service: func(ctx context.Context, input dto.VersionRequest) (corepipeline.Result[dto.Summary], error) {
			if controller.executor == nil {
				return corepipeline.Result[dto.Summary]{}, helpers.Map(authorization.ErrUnavailable)
			}
			result, err := controller.executor.Execute(ctx, input.Params.RequestID, input.RequestBody.Version)
			return corepipeline.OK(approvalhelpers.ToSummary(result)), helpers.Map(err)
		},
	})
}

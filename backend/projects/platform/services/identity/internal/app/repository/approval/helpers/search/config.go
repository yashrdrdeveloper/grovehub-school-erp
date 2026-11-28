package search

import (
	coresearch "grovehub/libs/core/search"
	dto "grovehub/services/identity/internal/app/dto/approval"
)

var Config = coresearch.Config{
	SearchExpression: "concat_ws(' ', r.request_id::text, r.request_type, r.status)",
	CursorExpression: "r.request_id",
	ScalarFields: map[string]coresearch.ScalarField{
		dto.FieldRequestID: {Expression: "r.request_id::text"},
		dto.FieldRequestType: {Expression: "r.request_type"},
		dto.FieldStatus: {Expression: "r.status"},
		dto.FieldCreatedAt: {Expression: "r.created_at", Cast: "::timestamptz"},
		dto.FieldUpdatedAt: {Expression: "r.updated_at", Cast: "::timestamptz"},
		dto.FieldExpiresAt: {Expression: "r.expires_at", Cast: "::timestamptz"},
	},
	SortFields: map[string]coresearch.SortField{
		dto.FieldRequestID: {Expression: "r.request_id", Cast: "::uuid"},
		dto.FieldRequestType: {Expression: "r.request_type", Cast: "::text"},
		dto.FieldStatus: {Expression: "r.status", Cast: "::text"},
		dto.FieldCreatedAt: {Expression: "r.created_at", Cast: "::timestamptz"},
		dto.FieldUpdatedAt: {Expression: "r.updated_at", Cast: "::timestamptz"},
		dto.FieldExpiresAt: {Expression: "r.expires_at", Cast: "::timestamptz"},
	},
}

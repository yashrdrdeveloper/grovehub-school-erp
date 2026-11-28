package search

import (
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	roledto "grovehub/services/identity/internal/app/dto/role"
)

var RoleConfig = coresearch.Config{
	SearchExpression: "r.search_details",
	CursorExpression: "r.role_id",
	ScalarFields: map[string]coresearch.ScalarField{
		roledto.FieldRoleID:    {Expression: "r.role_id::text"},
		roledto.FieldName:      {Expression: "r.normalized_name"},
		roledto.FieldStatus:    {Expression: "r.status"},
		roledto.FieldCreatedAt: {Expression: "r.created_at", Cast: "::timestamptz"},
		roledto.FieldUpdatedAt: {Expression: "r.updated_at", Cast: "::timestamptz"},
	},
	RelationFields: map[string]coresearch.RelationField{
		roledto.FieldPersonas: {
			Table: "identity.role_personas", Alias: "rp",
			ParentCondition: "rp.tenant_id = r.tenant_id AND rp.role_id = r.role_id",
			ValueExpression: "rp.persona_code",
		},
	},
	SortFields: map[string]coresearch.SortField{
		roledto.FieldRoleID:    {Expression: "r.role_id", Cast: "::uuid"},
		roledto.FieldName:      {Expression: "r.normalized_name", Cast: "::text"},
		roledto.FieldStatus:    {Expression: "r.status", Cast: "::text"},
		roledto.FieldCreatedAt: {Expression: "r.created_at", Cast: "::timestamptz"},
		roledto.FieldUpdatedAt: {Expression: "r.updated_at", Cast: "::timestamptz"},
	},
}

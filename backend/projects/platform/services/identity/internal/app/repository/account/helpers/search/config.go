package search

import (
	coresearch "grovehub/libs/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
)

// AccountConfig is the resource-owned allowlist of trusted SQL expressions.
// Untrusted JSON field names must never be interpolated as raw SQL.
var AccountConfig = coresearch.Config{
	SearchExpression: "a.search_details",
	CursorExpression: "a.account_id",
	ScalarFields: map[string]coresearch.ScalarField{
		accountdto.FieldAccountID: {Expression: "a.account_id::text"},
		accountdto.FieldPersonGID: {Expression: "a.person_gid::text"},
		accountdto.FieldUsername:  {Expression: "LOWER(COALESCE(a.username,''))"},
		accountdto.FieldEmail:     {Expression: "LOWER(a.email)"},
		accountdto.FieldStatus:    {Expression: "a.status"},
		accountdto.FieldPrimaryRoleID: {Expression: "COALESCE(a.primary_role_id,'')"},
		accountdto.FieldCreatedAt: {Expression: "a.created_at", Cast: "::timestamptz"},
		accountdto.FieldUpdatedAt: {Expression: "a.updated_at", Cast: "::timestamptz"},
	},
	SortFields: map[string]coresearch.SortField{
		accountdto.FieldAccountID: {Expression: "a.account_id", Cast: "::uuid"},
		accountdto.FieldUsername:  {Expression: "LOWER(COALESCE(a.username,''))", Cast: "::text"},
		accountdto.FieldStatus:    {Expression: "a.status", Cast: "::text"},
		accountdto.FieldCreatedAt: {Expression: "a.created_at", Cast: "::timestamptz"},
		accountdto.FieldUpdatedAt: {Expression: "a.updated_at", Cast: "::timestamptz"},
	},
}

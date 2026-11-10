package search

import (
	coresearch "grovehub/libs/core/search"
	"grovehub/services/foundation/entity"
	repositoryhelpers "grovehub/services/foundation/repository/helpers"
)

/**
 * PersonConfig defines the trusted search columns, relation filters and cursor
 * expressions for Foundation Person. All user-supplied fields are allowlisted
 * through these maps before SQL execution.
 */
var PersonConfig = coresearch.Config{
	SearchExpression: "p.search_details",
	CursorExpression: "p.gid",
	ScalarFields: map[string]coresearch.ScalarField{
		"gid":        {Expression: "p.gid::text"},
		"firstName":  {Expression: "LOWER(p.first_name)"},
		"middleName": {Expression: "LOWER(COALESCE(p.middle_name, ''))"},
		"lastName":   {Expression: "LOWER(p.last_name)"},
		"status":     {Expression: "p.status"},
		"createdAt":  {Expression: "p.created_at", Cast: "::timestamptz"},
		"updatedAt":  {Expression: "p.updated_at", Cast: "::timestamptz"},
	},
	RelationFields: map[string]coresearch.RelationField{
		"contacts.email":         repositoryhelpers.ContactSearchField(entity.ContactTypeEmail),
		"contacts.phone":         repositoryhelpers.ContactSearchField(entity.ContactTypePhone),
		"contacts.primaryEmail":  repositoryhelpers.ContactSearchField(entity.ContactTypeEmail, coresearch.FixedCondition{Expression: "contact.is_primary = ", Value: true}),
		"contacts.primaryPhone":  repositoryhelpers.ContactSearchField(entity.ContactTypePhone, coresearch.FixedCondition{Expression: "contact.is_primary = ", Value: true}),
		"contacts.verifiedEmail": repositoryhelpers.ContactSearchField(entity.ContactTypeEmail, coresearch.FixedCondition{Expression: "contact.is_verified = ", Value: true}),
		"contacts.verifiedPhone": repositoryhelpers.ContactSearchField(entity.ContactTypePhone, coresearch.FixedCondition{Expression: "contact.is_verified = ", Value: true}),
	},
	SortFields: map[string]coresearch.SortField{
		"gid":       {Expression: "p.gid", Cast: "::uuid"},
		"firstName": {Expression: "LOWER(p.first_name)", Cast: "::text"},
		"lastName":  {Expression: "LOWER(p.last_name)", Cast: "::text"},
		"createdAt": {Expression: "p.created_at", Cast: "::timestamptz"},
		"updatedAt": {Expression: "p.updated_at", Cast: "::timestamptz"},
	},
}

package search
import coresearch "grovehub/libs/core/search"
var AccountConfig=coresearch.Config{
 SearchExpression:"a.search_details",CursorExpression:"a.account_id",
 ScalarFields:map[string]coresearch.ScalarField{
 "accountId":{Expression:"a.account_id::text"},"personGid":{Expression:"a.person_gid::text"},
 "username":{Expression:"LOWER(COALESCE(a.username,''))"},
 "email":{Expression:"LOWER(a.email)"},
 "status":{Expression:"a.status"},
 "createdAt":{Expression:"a.created_at",Cast:"::timestamptz"},
 "updatedAt":{Expression:"a.updated_at",Cast:"::timestamptz"},
 },
 SortFields:map[string]coresearch.SortField{
 "accountId":{Expression:"a.account_id",Cast:"::uuid"},
 "username":{Expression:"LOWER(COALESCE(a.username,''))",Cast:"::text"},
 "status":{Expression:"a.status",Cast:"::text"},
 "createdAt":{Expression:"a.created_at",Cast:"::timestamptz"},
 "updatedAt":{Expression:"a.updated_at",Cast:"::timestamptz"},
 },
}

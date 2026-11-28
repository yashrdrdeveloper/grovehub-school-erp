package helpers

import coresearch "grovehub/libs/core/search"

// Resource-owned sort field allowlist. Generic cursor validators are in core.
var personSortFields = map[string]coresearch.CursorValueValidator{
    "gid": coresearch.ValidateUUIDv7Cursor,
    "firstName": coresearch.ValidateNonEmptyCursor,
    "lastName": coresearch.ValidateNonEmptyCursor,
    "createdAt": coresearch.ValidateTimeCursor,
    "updatedAt": coresearch.ValidateTimeCursor,
}

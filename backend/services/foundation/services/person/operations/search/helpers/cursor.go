// Package helpers owns cursor creation for the Person Search operation.
package helpers

import (
    coresearch "grovehub/libs/core/search"
    "grovehub/services/foundation/entity"
)

// BuildNextCursor encodes only the last item on a page with a successor.
func BuildNextCursor(result coresearch.QueryResult[entity.Person], query coresearch.Query[string], cursor coresearch.CursorPayload) (string, error) {
    if !result.HasNext || len(result.Items) == 0 { return "", nil }
    last := result.Items[len(result.Items)-1]
    cursor.Value = personCursorValue(last, query.Sort.Field)
    cursor.CursorId = last.GID
    return coresearch.EncodeCursor(cursor)
}

package flows

import (
    coresearch "grovehub/libs/core/search"
    persondto "grovehub/services/foundation/client/person/dto"
    "grovehub/services/foundation/internal/app/entity"
)

type Search struct {
    Request persondto.SearchPersonsRequest
    Query coresearch.Query[string]
    Cursor coresearch.CursorPayload
    Result coresearch.QueryResult[entity.Person]
}

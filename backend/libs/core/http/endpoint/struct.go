package endpoint

import (
	"context"
	"net/http"
	"net/url"
	corerequest "grovehub/libs/core/request"
)

// CallStruct invokes the registered controller route from an endpoint-specific
// DTO. Only present exported Params, Query and RequestBody fields are used;
// no empty placeholder types or transport wrapper are required.
func CallStruct[Response any, Request corerequest.Contract](ctx context.Context, handler http.Handler, method, route string, input Request) (Response, int, error) {
	var zero Response
	sections, err := corerequest.Inspect(input)
	if err != nil { return zero, 0, err }
	var params any
	var queryValues any
	var body any
	if sections.HasParams { params = sections.Params.Interface() }
	if sections.HasQuery { queryValues = sections.Query.Interface() }
	if sections.HasRequestBody { body = sections.RequestBody.Interface() }
	encodedQuery := make(url.Values)
	if queryValues != nil {
		var err error
		encodedQuery, err = encodeQuery(queryValues)
		if err != nil { return zero, 0, err }
	}
	return call[Response](ctx, handler, method, route, params, encodedQuery, body)
}

package person

// Person endpoints currently use the shared coresearch.Request as a JSON
// request body for searching, rather than defining a query string contract.
// When an endpoint needs URL query parameters, add a typed struct here and
// declare Query on its operation request. Omitted Query fields are ignored.

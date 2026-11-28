package validators

import "net/http"

// HTTPStatus validates a declaration and preserves its fail-fast semantics.
func HTTPStatus(statuses []int) int {
	switch len(statuses) {
	case 0:
		return http.StatusInternalServerError
	case 1:
		status := statuses[0]
		if status < http.StatusBadRequest || status > 599 {
			panic("error definition requires an HTTP 4xx or 5xx status")
		}

		return status
	default:
		panic("error definition accepts at most one HTTP status")
	}
}

package errorcode_test

import (
	"testing"

	"grovehub/launcher/constants/errorcode"
)

func TestRouteNotFoundDefinition(t *testing.T) {
	if errorcode.RouteNotFound.Code() != "PLATFORM_ROUTE_NOT_FOUND" {
		t.Fatalf("unexpected code: %s", errorcode.RouteNotFound.Code())
	}
	if errorcode.RouteNotFound.Message() != "API route not found" {
		t.Fatalf("unexpected message: %s", errorcode.RouteNotFound.Message())
	}
}

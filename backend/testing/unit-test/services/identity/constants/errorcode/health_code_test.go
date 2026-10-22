package errorcode_test

import (
	"testing"

	"grovehub/services/identity/constants/errorcode"
)

func TestServiceDownDefinition(t *testing.T) {
	t.Parallel()

	if errorcode.ServiceDown.Code() != "IDENTITY_SERVICE_DOWN" {
		t.Fatalf("unexpected code: %q", errorcode.ServiceDown.Code())
	}
	if errorcode.ServiceDown.Message() != "Identity service is unavailable" {
		t.Fatalf("unexpected message: %q", errorcode.ServiceDown.Message())
	}
}

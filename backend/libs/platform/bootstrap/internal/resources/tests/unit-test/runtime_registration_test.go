package bootstrap_test

import (
	"testing"

	"grovehub/libs/platform/bootstrap"
)

func TestRuntimeRegistrationBuild(t *testing.T) {
	registration := bootstrap.RuntimeRegistration[int, int, int]{
		Service:      bootstrap.Registration[int]{ServiceName: "test"},
		BuildRuntime: func(input int) (int, error) { return input + 1, nil },
	}
	var service bootstrap.RuntimeService[int, int, int] = registration
	if service.Name() != "test" {
		t.Fatal("operator registration not retained")
	}
	result, err := service.Build(41)
	if err != nil || result != 42 {
		t.Fatalf("unexpected runtime: %v, %v", result, err)
	}
}

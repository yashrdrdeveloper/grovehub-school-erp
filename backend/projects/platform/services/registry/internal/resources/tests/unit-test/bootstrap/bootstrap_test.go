package bootstrap_test

import (
	"context"
	"testing"

	"grovehub/libs/platform/bootstrap"
	"grovehub/services/registry"
)

func TestBootstrap(t *testing.T) {
	service := registry.Bootstrap()
	migration := service.Migration()
	if service.Name() != "registry" || migration.Target != bootstrap.BusinessDatabase || migration.Run == nil {
		t.Fatal("invalid registry bootstrap")
	}
	if err := migration.Run(context.Background(), nil); err == nil {
		t.Fatal("nil migration database accepted")
	}
}

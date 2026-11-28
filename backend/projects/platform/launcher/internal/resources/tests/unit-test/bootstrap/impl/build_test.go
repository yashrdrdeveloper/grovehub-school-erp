package impl_test

import (
 "testing"
 "grovehub/launcher/internal/src/bootstrap/impl"
)

func TestBuildRejectsMissingTenantConfiguration(t *testing.T){
 t.Setenv("PLATFORM_TENANT_ACCESS_REQUIRED","true")
 t.Setenv("PLATFORM_TENANT_ACCESS_ORIGIN","")
 t.Setenv("PLATFORM_TENANT_ACCESS_KEY_ID","")
 t.Setenv("PLATFORM_TENANT_ACCESS_PRIVATE_KEY","")
 runtime,err:=impl.New(nil).Build()
 if runtime!=nil||err==nil{t.Fatal("missing required tenant verifier accepted")}
}

func TestBuildRejectsMissingDatabase(t *testing.T){
 t.Setenv("PLATFORM_TENANT_ACCESS_REQUIRED","false")
 t.Setenv("PLATFORM_TENANT_ACCESS_ORIGIN","")
 t.Setenv("PLATFORM_TENANT_ACCESS_KEY_ID","")
 t.Setenv("PLATFORM_TENANT_ACCESS_PRIVATE_KEY","")
 t.Setenv("DATABASE_URL","")
 runtime,err:=impl.New(nil).Build()
 if runtime!=nil||err==nil{t.Fatal("missing database URL accepted")}
}

package config_test
import (
 "testing"
 "grovehub/launcher/internal/src/config"
)
func TestFromEnv(t *testing.T){
 t.Setenv("GROOVEHUB_HTTP_ADDR","")
 t.Setenv("PLATFORM_TENANT_ACCESS_REQUIRED","true")
 if c:=config.FromEnv();c.Address!="127.0.0.1:8080"||!c.TenantRequired{t.Fatal(c.Address,c.TenantRequired)}
 t.Setenv("GROOVEHUB_HTTP_ADDR","127.0.0.1:9090")
 if c:=config.FromEnv();c.Address!="127.0.0.1:9090"{t.Fatal(c.Address)}
}

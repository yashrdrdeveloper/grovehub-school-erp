package config
import "os"
type Config struct {Address,AllowedOrigin,Workload,TenantOrigin,TenantKeyID,TenantPrivateKey string;TenantRequired bool}
func FromEnv() Config {
 address:=os.Getenv("GROOVEHUB_HTTP_ADDR");if address==""{address="127.0.0.1:8080"}
 return Config{Address:address,AllowedOrigin:os.Getenv("IDENTITY_ALLOWED_ORIGIN"),Workload:os.Getenv("PLATFORM_WORKLOAD_CONFIG"),TenantOrigin:os.Getenv("PLATFORM_TENANT_ACCESS_ORIGIN"),TenantKeyID:os.Getenv("PLATFORM_TENANT_ACCESS_KEY_ID"),TenantPrivateKey:os.Getenv("PLATFORM_TENANT_ACCESS_PRIVATE_KEY"),TenantRequired:os.Getenv("PLATFORM_TENANT_ACCESS_REQUIRED")=="true"}
}

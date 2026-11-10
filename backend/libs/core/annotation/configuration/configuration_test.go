package configuration

import (
    "strings"
    "testing"
    "time"
)

type sample struct {
    Client string `env:"CLIENT" required:"true"`
    Secret string `env:"SECRET" required:"true" secret:"true"`
    Tenant string `env:"TENANT" default:"common"`
    Enabled bool `env:"ENABLED" default:"true"`
    Timeout time.Duration `env:"TIMEOUT" default:"5s"`
}

func TestDecodeRequiredAndDefaults(t *testing.T) {
    env:=map[string]string{"CLIENT":"  demo  ","SECRET":"  safe-key "}
    got,err:=Decode[sample](func(key string)(string,bool){v,ok:=env[key];return v,ok})
    if err!=nil{t.Fatal(err)}
    if got.Client!="demo"||got.Secret!="safe-key"||got.Tenant!="common"||!got.Enabled||got.Timeout!=5*time.Second{
        t.Fatalf("unexpected typed config values: client=%q tenant=%q enabled=%v duration=%v",got.Client,got.Tenant,got.Enabled,got.Timeout)
    }
}

func TestMissingRequiredAndSecretsAreNotExposed(t *testing.T) {
    _,err:=Decode[sample](func(key string)(string,bool) {
        if key=="CLIENT"{return "visible",true}
        return "",false
    })
    if err==nil{t.Fatal("expected missing secret error")}
    if strings.Contains(err.Error(),"safe-key"){t.Fatal("secret appears in error")}
}

func TestInvalidTypedValues(t *testing.T) {
    env:=map[string]string{"CLIENT":"x","SECRET":"secret","ENABLED":"maybe"}
    _,err:=Decode[sample](func(key string)(string,bool){v,ok:=env[key];return v,ok})
    if err==nil||!strings.Contains(err.Error(),"invalid boolean"){t.Fatalf("unexpected error: %v",err)}
}

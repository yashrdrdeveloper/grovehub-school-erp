// Package configuration provides an explicit, startup-only typed environment
// loader for @Configuration contract structs. It never logs secret values.
package configuration

import (
    "fmt"
    "os"
    "reflect"
    "strconv"
    "strings"
    "time"
)

// Lookup is injectable so configuration parsing can be tested without
// changing process environment. Returned bool indicates whether a key exists.
type Lookup func(string)(string,bool)

// FromEnv binds the env/default/required tags on an exported struct.
func FromEnv[T any]() (T,error) {return Decode[T](os.LookupEnv)}

// Decode validates a typed config contract and reads its declared environment
// sources. Reflection occurs once at startup; request paths never use it.
func Decode[T any](lookup Lookup) (T,error) {
    var value T
    if lookup==nil{return value,fmt.Errorf("configuration lookup is nil")}
    v:=reflect.ValueOf(&value).Elem()
    if v.Kind()!=reflect.Struct{return value,fmt.Errorf("configuration target must be a struct")}
    t:=v.Type()
    for i:=0;i<t.NumField();i++ {
        field:=t.Field(i)
        key,has:=field.Tag.Lookup("env")
        if !has{continue}
        if field.PkgPath!=""{return value,fmt.Errorf("configuration field %s must be exported",field.Name)}
        if strings.TrimSpace(key)==""{return value,fmt.Errorf("configuration field %s has an empty environment key",field.Name)}
        text,ok:=lookup(key)
        text=strings.TrimSpace(text)
        if (!ok||text=="") {
            if def,exists:=field.Tag.Lookup("default");exists{text=def}
        }
        if text==""&&field.Tag.Get("required")=="true"{
            return value,fmt.Errorf("required configuration %s is missing",key)
        }
        if text==""{continue}
        dest:=v.Field(i)
        switch {
        case dest.Kind()==reflect.String:
            dest.SetString(text)
        case dest.Type()==reflect.TypeOf(time.Duration(0)):
            duration,err:=time.ParseDuration(text);if err!=nil{return value,fmt.Errorf("%s: invalid duration",key)}
            dest.SetInt(int64(duration))
        case dest.Kind()==reflect.Bool:
            parsed,err:=strconv.ParseBool(text);if err!=nil{return value,fmt.Errorf("%s: invalid boolean",key)}
            dest.SetBool(parsed)
        case dest.Kind()>=reflect.Int&&dest.Kind()<=reflect.Int64:
            parsed,err:=strconv.ParseInt(text,10,dest.Type().Bits());if err!=nil{return value,fmt.Errorf("%s: invalid integer",key)}
            dest.SetInt(parsed)
        default:
            return value,fmt.Errorf("unsupported configuration type for %s",field.Name)
        }
    }
    return value,nil
}

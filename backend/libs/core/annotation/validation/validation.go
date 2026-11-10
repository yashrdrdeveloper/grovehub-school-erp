// Package validation applies a small, explicit subset of declarative
// @Validator field-tag rules. Domain validators remain responsible for
// mapping technical failures to their service-owned error definitions.
package validation

import (
    "fmt"
    "reflect"
    "strconv"
    "strings"
    "unicode/utf8"
)

// Error reports the invalid field and rule only. It never includes an input
// value, password, token, or other potentially sensitive request content.
type Error struct {
    Field string
    Rule string
}
func (e Error) Error() string {return fmt.Sprintf("%s violates %s",e.Field,e.Rule)}

// Normalize returns a copy with safe declared transformations applied.
// Supports string and *string request fields, preserving nil PATCH fields.
// Rules: normalize:"trim"; validate:"required,notblank,maxrunes=100".
func Normalize[T any](input T) (T,error) {
    output:=input
    value:=reflect.ValueOf(&output).Elem()
    if value.Kind()!=reflect.Struct{return output,fmt.Errorf("validation target must be a struct")}
    typ:=value.Type()
    for i:=0;i<typ.NumField();i++{
        field:=typ.Field(i)
        normalizer:=field.Tag.Get("normalize")
        rules:=field.Tag.Get("validate")
        if normalizer==""&&rules==""{continue}
        if !value.Field(i).CanSet(){return output,fmt.Errorf("validation field %s must be exported",field.Name)}
        target:=value.Field(i)
        isPtr:=target.Kind()==reflect.Ptr
        if isPtr {
            if target.IsNil(){
                if hasRule(rules,"required"){return output,Error{field.Name,"required"}}
                continue
            }
            if target.Type().Elem().Kind()!=reflect.String{return output,fmt.Errorf("unsupported validation type for %s",field.Name)}
            pointer:=reflect.New(target.Type().Elem())
            pointer.Elem().Set(target.Elem())
            target.Set(pointer)
            target=pointer.Elem()
        }
        if target.Kind()!=reflect.String{return output,fmt.Errorf("unsupported validation type for %s",field.Name)}
        text:=target.String()
        if normalizer!=""{
            switch normalizer {
            case "trim":text=strings.TrimSpace(text)
            default:return output,fmt.Errorf("%s: unsupported normalization %s",field.Name,normalizer)
            }
        }
        target.SetString(text)
        for _,rule:=range strings.Split(rules,","){
            if rule==""{continue}
            switch {
            case rule=="required",rule=="notblank":
                if strings.TrimSpace(text)==""{return output,Error{field.Name,rule}}
            case strings.HasPrefix(rule,"maxrunes="):
                maximum,err:=strconv.Atoi(strings.TrimPrefix(rule,"maxrunes="))
                if err!=nil||maximum<0{return output,fmt.Errorf("%s: invalid maxrunes declaration",field.Name)}
                if utf8.RuneCountInString(text)>maximum{return output,Error{field.Name,"maxrunes"}}
            default:
                return output,fmt.Errorf("%s: unsupported validation rule %s",field.Name,rule)
            }
        }
    }
    return output,nil
}

func hasRule(rules, expected string)bool {
    for _,rule:=range strings.Split(rules,","){if rule==expected{return true}}
    return false
}

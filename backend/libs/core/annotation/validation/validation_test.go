package validation

import (
    "errors"
    "strings"
    "testing"
)

type input struct {
    Name string `normalize:"trim" validate:"notblank,maxrunes=3"`
    Optional *string `normalize:"trim" validate:"maxrunes=3"`
}

func TestNormalizeAndPreservePointerInput(t *testing.T) {
    original:=" abc "
    request:=input{Name:"  Ab ",Optional:&original}
    result,err:=Normalize(request)
    if err!=nil{t.Fatal(err)}
    if result.Name!="Ab"||*result.Optional!="abc"{t.Fatalf("bad normalization: %+v",result)}
    if original!=" abc "{t.Fatal("input pointer was mutated")}
}

func TestNilPatchField(t *testing.T) {
    result,err:=Normalize(input{Name:"Ada"})
    if err!=nil{t.Fatal(err)}
    if result.Optional!=nil{t.Fatal("nil patch field was created")}
}

func TestUnicodeRuneLength(t *testing.T) {
    _,err:=Normalize(input{Name:"नमस्ते"})
    var violation Error
    if !errors.As(err,&violation)||violation.Rule!="maxrunes"{t.Fatalf("expected rune-length failure, got %v",err)}
}

func TestNotBlank(t *testing.T) {
    _,err:=Normalize(input{Name:" \t "})
    var violation Error
    if !errors.As(err,&violation)||violation.Rule!="notblank"{t.Fatalf("expected blank failure, got %v",err)}
}

func TestFailClosedUnknownRule(t *testing.T) {
    type bad struct { Value string `validate:"unrecognized"` }
    _,err:=Normalize(bad{Value:"ok"})
    if err==nil||!strings.Contains(err.Error(),"unsupported validation"){t.Fatalf("expected unsupported-rule failure, got %v",err)}
}

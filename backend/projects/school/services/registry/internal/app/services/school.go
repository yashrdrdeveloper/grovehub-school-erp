package services

import (
 "context"
 "database/sql"
 "errors"
 "regexp"
 "strings"
 "unicode"
 "unicode/utf8"

 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/entity"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/repository"
)
var codePattern=regexp.MustCompile("^[a-z][a-z0-9_-]{1,63}$")
func ValidID(s string)bool{id,err:=uuid.Parse(s);return err==nil&&id!=uuid.Nil&&id.String()==s}
func Register(ctx context.Context,db *sql.DB,id,code,name string)error{
 name=strings.TrimSpace(name)
 if !ValidID(id)||!codePattern.MatchString(code)||!utf8.ValidString(name)||utf8.RuneCountInString(name)<1||utf8.RuneCountInString(name)>150||strings.ContainsFunc(name,unicode.IsControl){return errors.New("invalid school registration")}
 return repository.Register(ctx,db,entity.School{ID:id,Code:code,Name:name})
}

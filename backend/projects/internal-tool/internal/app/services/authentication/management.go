package authentication

import (
 "context"
 "crypto/rand"

 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
 operators "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/repository/operator"
)

func(s *Service) ListOperators(ctx context.Context,actor entity.Operator)([]operators.Summary,error){
 return s.operators.List(ctx,actor)
}
func(s *Service) CreateOperator(ctx context.Context,actor entity.Operator,email,password string,persona entity.Persona)(entity.Operator,error){
 if actor.Persona!=entity.L2{return entity.Operator{},operators.ErrForbidden}
 email,ok:=normalizeEmail(email)
 if !ok||len(password)<12||len(password)>256||(persona!=entity.L1&&persona!=entity.L2){return entity.Operator{},ErrInvalid}
 select{case s.hashing<-struct{}{}:defer func(){<-s.hashing}();default:return entity.Operator{},ErrLimited}
 salt:=make([]byte,16);if _,err:=rand.Read(salt);err!=nil{return entity.Operator{},err}
 op:=entity.Operator{ID:uuid.NewString(),Email:email,Persona:persona}
 err:=s.operators.Create(ctx,actor,entity.Credential{Operator:op,Salt:salt,Hash:hashPassword(password,salt)})
 return op,err
}
func(s *Service) SetOperatorStatus(ctx context.Context,actor entity.Operator,id,status string)error{
 v,err:=uuid.Parse(id)
 if err!=nil||v==uuid.Nil||v.String()!=id||(status!="ACTIVE"&&status!="DISABLED"){return ErrInvalid}
 return s.operators.SetStatus(ctx,actor,id,status)
}

func(s *Service) SearchOperators(ctx context.Context,actor entity.Operator,p entity.PageRequest)(entity.Page[operators.Summary],error){
 p,ok:=p.Normalize();if !ok{return entity.Page[operators.Summary]{},ErrInvalid}
 return s.operators.Search(ctx,actor,p)
}

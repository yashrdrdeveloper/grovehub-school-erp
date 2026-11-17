package impl
import (
 "context"
 "errors"
 "time"
 "grovehub/services/identity/internal/app/entity"
 contract "grovehub/services/identity/internal/app/services/credential"
 helpers "grovehub/services/identity/internal/app/services/credential/helpers"
 repo "grovehub/services/identity/internal/app/repository/credential"
)
type Service struct{repository repo.Repository}
func New(r repo.Repository)contract.Service{return &Service{repository:r}}
func validate(password string)error{if len(password)<12||len(password)>256{return errors.New("invalid password length")};return nil}
func(s *Service)Set(ctx context.Context,id,password string)error{
 if err:=validate(password);err!=nil{return err}
 _,err:=s.repository.Find(ctx,id);if err==nil{return errors.New("password already configured")}
 if !errors.Is(err,repo.ErrNotFound){return err}
 hash,err:=helpers.HashPassword(password);if err!=nil{return err}
 now:=time.Now().UTC()
 return s.repository.Create(ctx,entity.PasswordCredential{AccountID:id,PasswordHash:hash,Version:1,CreatedAt:now,UpdatedAt:now,PasswordChangedAt:now})
}
func(s *Service)Verify(ctx context.Context,id,password string)(bool,error){
 c,err:=s.repository.Find(ctx,id);if err!=nil{return false,err}
 return helpers.VerifyPassword(password,c.PasswordHash)
}
func(s *Service)Replace(ctx context.Context,id,password string)error{
 if err:=validate(password);err!=nil{return err}
 c,err:=s.repository.Find(ctx,id);if err!=nil{return err}
 hash,err:=helpers.HashPassword(password);if err!=nil{return err}
 now:=time.Now().UTC();c.Version++;c.PasswordHash=hash;c.PasswordChangedAt=now;c.UpdatedAt=now
 return s.repository.Update(ctx,c)
}

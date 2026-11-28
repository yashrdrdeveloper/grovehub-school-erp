package rbac_test

import (
 "context"
 "strings"
 "time"

 "github.com/google/uuid"
 . "github.com/onsi/gomega"
 "gorm.io/gorm"
 schoolgateway "grovehub/launcher/session"
 "grovehub/services/identity/internal/app/entity"
 sessionrepo "grovehub/services/identity/internal/app/repository/session"
 sessionimpl "grovehub/services/identity/internal/app/repository/session/impl"
)

func verifyTenantSessionReconciliation(ctx context.Context, db *gorm.DB) {
 tx:=db.Begin()
 Expect(tx.Error).NotTo(HaveOccurred())
 defer tx.Rollback()
 tenant,foreign:=uuid.NewString(),uuid.NewString()
 account,other:=uuid.NewString(),uuid.NewString()
 for _,row:=range []struct{id,tenant string}{{account,tenant},{other,foreign}}{
  Expect(tx.Exec("INSERT INTO identity.accounts(account_id,person_gid,tenant_id,email,status) VALUES (?,?,?,?,'ACTIVE')",row.id,uuid.NewString(),row.tenant,row.id+"@example.test").Error).NotTo(HaveOccurred())
 }
 now:=time.Now().UTC()
 create:=func(id string, generation int64, access,refresh string) entity.Session {
  item:=entity.Session{SessionID:uuid.NewString(),AccountID:id,SessionGeneration:generation,AccessHash:strings.Repeat(access,64),RefreshHash:strings.Repeat(refresh,64),AccessExpiresAt:now.Add(time.Hour),RefreshExpiresAt:now.Add(24*time.Hour),CreatedAt:now,UpdatedAt:now}
  Expect(tx.Create(&item).Error).NotTo(HaveOccurred())
  return item
 }
 legacy:=create(account,0,"a","b")
 stale:=create(account,1,"c","d")
 fresh:=create(account,2,"e","f")
 future:=create(account,3,"1","2")
 untouched:=create(other,1,"3","4")
 cleanup:=sessionimpl.NewTenantSessions(tx)
 sessions:=sessionimpl.New(tx)
 candidates,err:=cleanup.ListTenants(ctx,"",20)
 Expect(err).NotTo(HaveOccurred())
 Expect(candidates).To(ConsistOf(tenant,foreign))
 actorCtx:=schoolgateway.WithPrincipal(ctx,schoolgateway.Principal{AccountID:account,TenantID:tenant,SessionID:fresh.SessionID,SessionGeneration:2})
 listed,err:=sessions.List(actorCtx,account)
 Expect(err).NotTo(HaveOccurred())
 Expect(listed).To(HaveLen(1))
 Expect(listed[0].SessionID).To(Equal(fresh.SessionID))
 count,err:=cleanup.RevokeBefore(ctx,tenant,2)
 Expect(err).NotTo(HaveOccurred());Expect(count).To(Equal(int64(2)))
 count,err=cleanup.RevokeBefore(ctx,tenant,2)
 Expect(err).NotTo(HaveOccurred());Expect(count).To(BeZero())
 for _,item:=range []entity.Session{legacy,stale}{
  _,err:=sessions.FindAccess(ctx,item.AccessHash,now)
  Expect(err).To(MatchError(sessionrepo.ErrInvalidSession))
  _,err=sessions.FindRefresh(ctx,item.RefreshHash,now)
  Expect(err).To(MatchError(sessionrepo.ErrInvalidSession))
 }
 for _,item:=range []entity.Session{fresh,future,untouched}{
  _,err:=sessions.FindAccess(ctx,item.AccessHash,now)
  Expect(err).NotTo(HaveOccurred())
 }
}

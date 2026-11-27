package connection_test

import (
    "testing"
    "time"

    "grovehub/services/integration/internal/app/entity"
    impl "grovehub/services/integration/internal/app/mappers/connection/impl"
)

func TestConnectionMapperPreservesPublicFields(t *testing.T) {
    mapper:=impl.NewConnectionMapper()
    now:=time.Date(2026,time.November,15,0,0,0,0,time.UTC)
    email:="school@example.com"
    c:=entity.Connection{ID:"c",TenantID:"t",OwnerGID:"p",Provider:"GMAIL",ExternalID:&email,Status:entity.StatusActive,CreatedAt:now,UpdatedAt:now}
    result:=mapper.ToResponse(c)
    if result.ID!=c.ID || result.TenantID!=c.TenantID || result.OwnerGID!=c.OwnerGID || result.ExternalID==nil || *result.ExternalID!=email || result.Status!=c.Status || result.CreatedAt!=now {
        t.Fatalf("unexpected connection DTO: %+v",result)
    }
}

func TestConnectionTestMapperPreservesExecutionResult(t *testing.T) {
    mapper:=impl.NewConnectionMapper()
    now:=time.Date(2026,time.November,15,0,0,0,0,time.UTC)
    original:=entity.ConnectionTest{ID:"test",ConnectionID:"c",Direction:"OUTBOUND",Status:"SUCCESS",Detail:"accepted",StartedAt:now,CompletedAt:now}
    dto:=mapper.TestToResponse(original)
    if dto.ID!=original.ID || dto.ConnectionID!=original.ConnectionID || dto.Status!=original.Status || dto.Detail!=original.Detail || dto.CompletedAt!=now {
        t.Fatalf("unexpected connection test DTO: %+v",dto)
    }
}

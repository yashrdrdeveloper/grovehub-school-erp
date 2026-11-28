package connection_test

import (
    "errors"
    "testing"

    dto "grovehub/services/integration/internal/app/dto/connection"
    errorcode "grovehub/services/integration/internal/app/services/connection/helpers"
    impl "grovehub/services/integration/internal/app/validators/connection/impl"
)

func TestRejectInvalidConnectionIdentifiers(t *testing.T) {
    validator:=impl.NewConnectionValidator()
    for _,value:=range []string{"", "unknown", "00000000-0000-0000-0000-000000000000", "E93696CC-15F0-4E06-9577-0D5CFC2A83C8"} {
        if err:=validator.ValidatePair(value,"01a1b210-402c-7622-a33c-f82ceba26c7d");!errors.Is(err,errorcode.ErrInvalidConnection) {
            t.Fatalf("ID %q should be rejected: %v",value,err)
        }
    }
}

func TestCreateValidatesTenantAndOwnerWithoutProviderLookup(t *testing.T) {
    validator:=impl.NewConnectionValidator()
    tenant:="018f5d04-315c-75a9-b890-df0d02e787c1"
    owner:="01a1b210-402c-7622-a33c-f82ceba26c7d"
    if err:=validator.ValidateCreate(tenant,dto.CreateConnection{OwnerGID:owner,Provider:"GMAIL"});err!=nil {t.Fatal(err)}
    if err:=validator.ValidateCreate(tenant,dto.CreateConnection{OwnerGID:"bad",Provider:"GMAIL"});!errors.Is(err,errorcode.ErrInvalidConnection) {t.Fatalf("expected invalid owner: %v",err)}
}

func TestSchoolRegistryCanonicalUUID(t *testing.T) {
 validator:=impl.NewConnectionValidator()
 school:="11111111-1111-4111-8111-111111111111"
 if err:=validator.ValidatePair(school,"01a1b210-402c-7622-a33c-f82ceba26c7d");err!=nil{t.Fatal(err)}
 if err:=validator.ValidatePair(school,school);!errors.Is(err,errorcode.ErrInvalidConnection){t.Fatal("non-v7 resource accepted",err)}
}

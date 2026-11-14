package helpers_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    dto "grovehub/services/foundation/internal/app/dto/person"
    entity "grovehub/services/foundation/internal/app/entity"
    helper "grovehub/services/foundation/internal/app/services/person/operations/create/helpers"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    ids "grovehub/services/foundation/internal/tests/constants/person"
)

// @Testing suite=PersonCreateHelpers
func NewPersonSpecs() {
    Describe("NewPerson",func(){
        It("constructs a new active and version-one Person without persistence",func(){
            req:=dto.CreatePerson{PersonNames:dto.PersonNames{FirstName:ids.FirstName,LastName:ids.LastName}}
            p:=helper.NewPerson(ids.GID,req,fixtures.FixedTime())
            Expect(p.GID).To(Equal(ids.GID))
            Expect(p.Status).To(Equal(entity.PersonStatusActive))
            Expect(p.Version).To(Equal(int64(1)))
            Expect(p.MiddleName).To(BeNil())
            Expect(p.CreatedAt).To(Equal(fixtures.FixedTime()))
            Expect(p.UpdatedAt).To(Equal(fixtures.FixedTime()))
        })
        It("retains a specified middle name",func(){
            req:=dto.CreatePerson{PersonNames:dto.PersonNames{FirstName:ids.FirstName,MiddleName:"Kumar",LastName:ids.LastName}}
            p:=helper.NewPerson(ids.GID,req,fixtures.FixedTime())
            Expect(p.MiddleName).NotTo(BeNil())
            Expect(*p.MiddleName).To(Equal("Kumar"))
        })
    })
}

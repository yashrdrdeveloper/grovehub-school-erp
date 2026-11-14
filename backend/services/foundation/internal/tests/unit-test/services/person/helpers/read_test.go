package helpers_test

import (
    "context"
    "errors"
    "go.uber.org/mock/gomock"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    coreerrors "grovehub/libs/core/errors"
    errcode "grovehub/services/foundation/internal/app/constants/error"
    dto "grovehub/services/foundation/internal/app/dto/person"
    "grovehub/services/foundation/internal/app/entity"
    helpers "grovehub/services/foundation/internal/app/services/person/helpers"
    repoerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    ids "grovehub/services/foundation/internal/tests/constants/person"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
)

var persons *mocks.MockPersonRepository
var contacts *mocks.MockContactRepository
var projection *mocks.MockPersonMapper

// @BeforeEach
func setupReader() {
    controller:=gomock.NewController(GinkgoT())
    persons=mocks.NewMockPersonRepository(controller)
    contacts=mocks.NewMockContactRepository(controller)
    projection=mocks.NewMockPersonMapper(controller)
}

// @Testing suite=FoundationPersonHelpers
func ReadSpecs() {
    Describe("ReadPerson",func(){
        It("reads an aggregate with contact children",func(){
            p:=fixtures.Entity()
            c:=fixtures.Contact()
            persons.EXPECT().FindByGID(gomock.Any(),ids.GID).Return(p,nil)
            contacts.EXPECT().FindByPersonGID(gomock.Any(),ids.GID).Return([]entity.Contact{c},nil)
            found,children,err:=helpers.ReadPerson(context.Background(),ids.GID,persons,contacts,projection)
            Expect(err).NotTo(HaveOccurred())
            Expect(found.GID).To(Equal(ids.GID))
            Expect(children).To(Equal([]entity.Contact{c}))
        })
        It("returns an InvalidPersonGID without repository access",func(){
            projection.EXPECT().ToPersonErrorMetadata("invalid").Return(dto.PersonErrorMetadata{GID:"invalid"})
            _,_,err:=helpers.ReadPerson(context.Background(),"invalid",persons,contacts,projection)
            assertPersonReadError(err,errcode.InvalidPersonGID.Code())
        })
        It("maps missing persons to a domain NotFound error",func(){
            persons.EXPECT().FindByGID(gomock.Any(),ids.GID).Return(fixtures.Entity(),repoerrors.ErrPersonNotFound)
            projection.EXPECT().ToPersonErrorMetadata(ids.GID).Return(dto.PersonErrorMetadata{GID:ids.GID})
            _,_,err:=helpers.ReadPerson(context.Background(),ids.GID,persons,contacts,projection)
            assertPersonReadError(err,errcode.PersonNotFound.Code())
        })
        It("maps repository errors to PersonReadFailed",func(){
            persons.EXPECT().FindByGID(gomock.Any(),ids.GID).Return(fixtures.Entity(),errors.New("offline"))
            projection.EXPECT().ToPersonErrorMetadata(ids.GID).Return(dto.PersonErrorMetadata{GID:ids.GID})
            _,_,err:=helpers.ReadPerson(context.Background(),ids.GID,persons,contacts,projection)
            assertPersonReadError(err,errcode.PersonReadFailed.Code())
        })
        It("maps contact read errors to PersonReadFailed",func(){
            persons.EXPECT().FindByGID(gomock.Any(),ids.GID).Return(fixtures.Entity(),nil)
            contacts.EXPECT().FindByPersonGID(gomock.Any(),ids.GID).Return(nil,errors.New("offline"))
            projection.EXPECT().ToPersonErrorMetadata(ids.GID).Return(dto.PersonErrorMetadata{GID:ids.GID})
            _,_,err:=helpers.ReadPerson(context.Background(),ids.GID,persons,contacts,projection)
            assertPersonReadError(err,errcode.PersonReadFailed.Code())
        })
    })
}

func assertPersonReadError(err error,code coreerrors.Code) {
    var appErr *coreerrors.Error
    Expect(errors.As(err,&appErr)).To(BeTrue())
    Expect(appErr.Definition().Code()).To(Equal(code))
}

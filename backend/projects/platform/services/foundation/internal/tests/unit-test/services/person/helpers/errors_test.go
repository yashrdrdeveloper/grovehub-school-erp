package helpers_test

import (
    "errors"
    "fmt"

    . "github.com/onsi/ginkgo/v2"

    dto "grovehub/services/foundation/internal/app/dto/person"
    helpers "grovehub/services/foundation/internal/app/services/person/helpers"
    errcode "grovehub/services/foundation/internal/app/constants/error"
    repoerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
    ids "grovehub/services/foundation/internal/tests/constants/person"
)

// @Testing suite=FoundationPersonHelpers
func ErrorSpecs() {
    Describe("MapUpdateError",func(){
        It("maps missing persons to NotFound preserving the requested GID",func(){
            projection.EXPECT().ToPersonErrorMetadata(ids.GID).Return(dto.PersonErrorMetadata{GID:ids.GID})
            err:=helpers.MapUpdateError(ids.GID,fmt.Errorf("wrapped: %w",repoerrors.ErrPersonNotFound),projection)
            assertPersonReadError(err,errcode.PersonNotFound.Code())
        })
        It("maps duplicate contacts to a specific conflict without reading metadata",func(){
            err:=helpers.MapUpdateError(ids.GID,repoerrors.ErrContactAlreadyExists,projection)
            assertPersonReadError(err,errcode.PersonContactConflict.Code())
        })
        It("maps concurrent versions to the optimistic concurrency code",func(){
            projection.EXPECT().ToPersonErrorMetadata(ids.GID).Return(dto.PersonErrorMetadata{GID:ids.GID})
            err:=helpers.MapUpdateError(ids.GID,repoerrors.ErrPersonVersionConflict,projection)
            assertPersonReadError(err,errcode.PersonConcurrentModification.Code())
        })
        It("preserves unexpected backend failures as PersonUpdateFailed",func(){
            projection.EXPECT().ToPersonErrorMetadata(ids.GID).Return(dto.PersonErrorMetadata{GID:ids.GID})
            err:=helpers.MapUpdateError(ids.GID,errors.New("unknown write failure"),projection)
            assertPersonReadError(err,errcode.PersonUpdateFailed.Code())
        })
    })
}

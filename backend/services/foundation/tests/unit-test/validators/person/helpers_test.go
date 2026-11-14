package person_test

import (
    "errors"
    "testing"

    annotationvalidation "grovehub/libs/core/annotation/validation"
    coreerrors "grovehub/libs/core/errors"
    personerror "grovehub/services/foundation/constants/error"
    persondto "grovehub/services/foundation/dto/person"
    validatorhelpers "grovehub/services/foundation/validators/person/helpers"
)

func TestPersonDTOErrorMapping(t *testing.T) {
    cases := []struct {
        field string
        code string
    }{
        {"RequestBody.Contacts.Emails[0].Value", personerror.InvalidPersonContact.Code()},
        {"RequestBody.Contacts", personerror.InvalidPersonContact.Code()},
        {"RequestBody.FirstName", personerror.InvalidPersonRequest.Code()},
    }
    for _, tt := range cases {
        got := validatorhelpers.MapPersonDTOError(annotationvalidation.Error{Field: tt.field, Rule: "notblank"})
        var mapped *coreerrors.Error
        if !errors.As(got, &mapped) || mapped.Definition().Code() != tt.code {
            t.Errorf("%q: unexpected mapping %v", tt.field, got)
        }
    }
    original := errors.New("external error")
    if got := validatorhelpers.MapPersonDTOError(original); got != original {
        t.Errorf("non-DTO error must be preserved, got %v", got)
    }
}

func TestUpdateContactRulesKeepOmittedAndExplicitNull(t *testing.T) {
    for _, request := range []persondto.UpdatePersonRequest{
        {},
        {RequestBody: persondto.UpdatePerson{}},
    } {
        normalized, err := validatorhelpers.NormalizeUpdatePersonRequest(request)
        if err != nil || normalized.RequestBody.Contacts.Set {
            t.Fatalf("omitted contacts should remain omitted, got %+v, %v", normalized, err)
        }
    }
    explicitNull := persondto.UpdatePersonRequest{}
    explicitNull.RequestBody.Contacts.Set = true
    normalized, err := validatorhelpers.NormalizeUpdatePersonRequest(explicitNull)
    if err != nil || !normalized.RequestBody.Contacts.Set || normalized.RequestBody.Contacts.Value != nil {
        t.Fatalf("explicit null must be preserved, got %+v, %v", normalized, err)
    }
    if err := validatorhelpers.ValidateUpdatePersonRequest(normalized); err != nil {
        t.Fatalf("explicit null contacts must remain valid: %v", err)
    }
}

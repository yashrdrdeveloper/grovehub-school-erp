// Package person provides deterministic Foundation test fixtures.
package person

import (
    "time"

    coreclock "grovehub/libs/core/tests/helpers/clock"
    personconstants "grovehub/services/foundation/internal/tests/constants/person"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    "grovehub/services/foundation/internal/app/entity"
)

func FixedTime() time.Time { return coreclock.Fixed() }

func Entity() entity.Person {
    return entity.Person{
        GID:personconstants.GID,FirstName:personconstants.FirstName,
        LastName:personconstants.LastName,Status:entity.PersonStatusActive,
        Version:1,CreatedAt:FixedTime(),UpdatedAt:FixedTime(),
    }
}

func Contact() entity.Contact {
    return entity.Contact{PersonGID:personconstants.GID,Type:entity.ContactTypeEmail,
        Value:personconstants.Email,NormalizedValue:personconstants.Email,
        IsPrimary:true,CreatedAt:FixedTime(),UpdatedAt:FixedTime()}
}

func Response() persondto.Person {
    return persondto.Person{
        GID:personconstants.GID,FirstName:personconstants.FirstName,
        LastName:personconstants.LastName,Status:"ACTIVE",Version:1,
        Contacts:persondto.ResponseContacts{
            Emails:[]persondto.ResponseContactValue{{Value:personconstants.Email,IsPrimary:true}},
            Phones:[]persondto.ResponseContactValue{},
        },CreatedAt:FixedTime(),UpdatedAt:FixedTime(),
    }
}

func CreateRequest() persondto.CreatePersonRequest {
    return persondto.CreatePersonRequest{
        RequestBody:persondto.CreatePerson{
            PersonNames:persondto.PersonNames{FirstName:personconstants.FirstName,LastName:personconstants.LastName},
            Contacts:persondto.Contacts{Emails:[]persondto.ContactValue{{Value:personconstants.Email,IsPrimary:true}}},
        },
    }
}

func GetRequest() persondto.GetPersonByGIDRequest {
    return persondto.GetPersonByGIDRequest{Params:persondto.PersonByGID{GID:personconstants.GID}}
}

func UpdateRequest() persondto.UpdatePersonRequest {
    return persondto.UpdatePersonRequest{Params:persondto.PersonByGID{GID:personconstants.GID}}
}

func DeactivateRequest() persondto.DeactivatePersonRequest {
    return persondto.DeactivatePersonRequest{Params:persondto.PersonByGID{GID:personconstants.GID}}
}

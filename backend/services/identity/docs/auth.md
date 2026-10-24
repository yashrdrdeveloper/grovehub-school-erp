# Identity Auth Contract

## Scope

`AuthService` is the Identity entry-point boundary for public authentication flows. `Register` is the first implemented method. Future methods such as `Login`, `ForgotPassword`, `ResetPassword`, and `ChangePassword` belong on the same Auth service/controller boundary when implemented.

Current public endpoint:

```text
POST /auth/register
```

There is no public `/identity/accounts` endpoint and no standalone password CRUD endpoint.

## Auth structure

```text
AuthController
    ↓
AuthService
    ├── Register        implemented
    ├── Login           future
    ├── ForgotPassword  future
    ├── ResetPassword   future
    └── ChangePassword  future
```

Registration-specific errors remain registration-specific even though the owning service is Auth.

## Register request

```go
type Register struct {
    Gid         string  `json:"gid,omitempty"`
    FirstName   string  `json:"firstName"`
    MiddleName  *string `json:"middleName"`
    LastName    string  `json:"lastName"`
    Email       string  `json:"email"`
    PhoneNumber string  `json:"phoneNumber,omitempty"`
    Password    string  `json:"password"`
}
```

`firstName`, `middleName`, `lastName`, `email`, and `password` are mandatory JSON fields. `gid` and `phoneNumber` are optional.

`firstName` and `lastName` must be non-blank. `middleName` must be present but may be an empty string for a person with no middle name. The pointer exists only so the HTTP boundary can distinguish an omitted `middleName` field from an explicitly supplied empty value.

Foundation remains responsible for canonical name normalization and maximum-length validation when Auth patches the Person.

## Person resolution and name update

Register resolves a Foundation Person in this order:

```text
1. gid supplied?
      YES -> GetPersonByGID
             -> submitted email must belong to that Person
      NO  -> search by exact normalized email

2. Email search result
      one match -> reuse existing Person/GID
      no match  -> create Person/GID using submitted names + email + optional phone
      >1 match  -> fail safely

3. Always patch the resolved Person through Foundation UpdatePerson
      firstName
      middleName
      lastName

4. Create Identity Account

5. Set PasswordCredential
```

The name patch is performed for every successful resolution path: provided GID, existing email match, and newly created Person. This gives Register one consistent rule: the submitted names become the current canonical Foundation Person names before Identity Account creation.

Auth does not update the Foundation contact set during this name patch. Email is still the linkage proof/resolution key, and the optional Identity phone remains part of Account creation.

## Cross-service boundary

Auth depends on Foundation only through `foundation/client.Client` and currently uses:

```text
GetPersonByGID
SearchPersons
CreatePerson
UpdatePerson
```

Identity does not import Foundation repositories, service implementations, or tables.

## Compensation

Foundation Person updates are not rolled back if later Identity Account/password creation fails. Foundation is the canonical Person domain and the submitted names were explicitly accepted as the desired current names.

If Account creation succeeds but password credential creation fails, Auth compensates by permanently deleting the newly created Identity Account. It never deletes or deactivates the Foundation Person.

## Errors

```text
IDENTITY_REGISTRATION_INVALID        -> 400
IDENTITY_REGISTRATION_ALREADY_EXISTS -> 409
IDENTITY_REGISTRATION_FAILED         -> 500
```

## Register response

```text
accountId
personGid
email
phone (optional)
status
```

Register does not create a session or tokens yet.

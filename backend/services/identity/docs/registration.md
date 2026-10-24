# Identity Registration Contract

## Scope

Registration is the public Identity orchestration flow that resolves or creates the canonical Foundation Person, creates the internal Identity Account, and creates the password credential used by password login.

Public endpoint:

```text
POST /auth/register
```

There is no public `/identity/accounts` endpoint and no standalone password CRUD endpoint.

## Request contract

Registration owns a flat HTTP request contract:

```go
type Register struct {
    Gid         string `json:"gid,omitempty"`
    FirstName   string `json:"firstName"`
    MiddleName  string `json:"middleName,omitempty"`
    LastName    string `json:"lastName"`
    Email       string `json:"email"`
    PhoneNumber string `json:"phoneNumber,omitempty"`
    Password    string `json:"password"`
}
```

`gid` is optional. It is used when the caller already knows the canonical Foundation Person GID.

Example:

```json
{
  "firstName": "Yash",
  "middleName": "",
  "lastName": "Developer",
  "email": "yash@example.com",
  "phoneNumber": "+919876543210",
  "password": "StrongPassword123!"
}
```

Existing-person registration may also supply `gid`:

```json
{
  "gid": "019a1f4f-4abc-7def-8abc-1234567890ab",
  "firstName": "Yash",
  "middleName": "",
  "lastName": "Developer",
  "email": "yash@example.com",
  "phoneNumber": "+919876543210",
  "password": "StrongPassword123!"
}
```

## Registration validation

The Registration validator owns registration-level canonicalization:

```text
email -> normalized canonical email
phoneNumber -> normalized E.164 phone when present
gid -> validated/normalized UUIDv7 when present
password -> new-password policy
```

Passwords are not trimmed or normalized.

Person name validation remains Foundation-owned. This matters because when an existing Person is resolved by GID or email, Registration does not rewrite that Person's canonical name/contact record.

## Person resolution

Registration resolves the Foundation Person in this order:

```text
1. gid supplied?
      YES -> Foundation GetPersonByGID
             -> email must belong to that Person
      NO  -> continue

2. Search Foundation by exact contacts.email
      one match -> reuse that Person and GID
      no match  -> continue
      >1 match  -> fail safely (Foundation contact uniqueness should prevent this)

3. Create Foundation Person
      firstName + middleName + lastName
      primary email
      optional primary phone
      -> Foundation generates the new GID
```

Email is therefore the canonical automatic resolver when no GID is supplied.

Registration never creates a new Person merely because the request omitted a GID when the email already belongs to an existing Foundation Person.

If a GID is supplied but the submitted email is not one of that Person's email contacts, Registration returns `IDENTITY_REGISTRATION_INVALID` and does not create an Identity Account.

When Person creation races with another request, Foundation remains contact-aware/idempotent. Registration also verifies that any Person returned from the create path owns the submitted email before continuing, preventing accidental linkage through a conflicting phone contact.

## Identity creation

After the Person is resolved:

```text
Foundation Person.GID
        +
normalized registration email
        +
optional normalized phoneNumber
        ↓
AccountService.Create
        ↓
PasswordCredentialService.Set
        ↓
Registration response
```

The Identity Account email is always the normalized registration email. Phone is optional.

If an Identity Account already exists for the resolved Person GID, email, or phone, Registration returns the registration-level conflict error rather than creating another account.

## Cross-service boundary

Registration depends on Foundation only through `foundation/client.Client`.

It uses:

```text
GetPersonByGID
SearchPersons
CreatePerson
```

Identity does not import Foundation repositories, service implementations, or tables.

## Failure and compensation

Foundation Person and Identity Account are separate service responsibilities, so Registration does not use a cross-service database transaction.

If Person resolution/creation succeeds but Account creation fails, the Foundation Person remains unchanged.

If Account creation succeeds but password credential creation fails, Registration compensates by permanently deleting the newly created Identity Account.

Registration never deletes or deactivates the Foundation Person as compensation.

If compensation also fails, Registration returns `IDENTITY_REGISTRATION_FAILED` with the underlying causes retained internally.

## Error boundary

Public registration errors remain:

```text
IDENTITY_REGISTRATION_INVALID        -> 400
IDENTITY_REGISTRATION_ALREADY_EXISTS -> 409
IDENTITY_REGISTRATION_FAILED         -> 500
```

Foundation/Account/Password persistence details are not exposed as the public contract.

## Response

The response remains:

```text
accountId
personGid
email
phone (optional)
status
```

Registration does not create a session or access/refresh tokens yet.

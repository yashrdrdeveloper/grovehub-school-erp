# Foundation Persistence

Foundation persists the Person aggregate through GORM behind the repository contract.

```text
services/person.go
    ↓
repository.PersonRepository
    ↓
repository/impl.PersonRepository
    ↓
GORM
    ↓
PostgreSQL foundation schema
```

## Files

```text
repository/person.go
    public persistence contract and sentinel errors

repository/impl/person.go
    GORM operations, transactions, search execution, constraint mapping

repository/impl/person_model.go
    persistence-only GORM models for foundation.persons and foundation.person_contacts

repository/impl/person_mapper.go
    explicit mapping between persistence models and entity.Person/entity.Contact
```

`entity.Person` and `entity.Contact` must remain persistence-agnostic and must not contain GORM tags.

## Schema ownership

SQL migrations under `migrations/` are the source of truth. Foundation does not use GORM `AutoMigrate` at runtime.

Owned tables:

```text
foundation.persons
foundation.person_contacts
```

The repository keeps Person + contacts + `search_details` refresh atomic in one transaction.

## Search

The public/shared search contract is unchanged. The service validates fields/operators and produces the generic `coresearch.Query`. The repository implementation maps those approved fields to fixed PostgreSQL expressions through GORM.

No client-provided field is used directly as a SQL identifier.

## Errors

The PostgreSQL unique index:

```text
uq_person_contacts_type_normalized_value
```

maps to:

```text
repository.ErrPersonContactAlreadyExists
```

The database constraint remains the concurrency-safe source of truth for contact uniqueness.

# Foundation Person changed event

Topic: `foundation_person_changed`  
Type: `foundation.person.changed.v1`

Foundation emits one event after a successful state transition:

- `CREATED`: `old: null`, `new`: complete `dto/response.Person`.
- `UPDATED`: both full Person snapshots; no event for a no-op update.
- `DEACTIVATED`: both full Person snapshots; no separate UPDATED event and no event when already deactivated.
- A deletion operation is **not supported** by the current PersonService; `new: null` is reserved for a future `DELETED` operation.

Consumers decode `coreevents.Envelope[persons.Changed]`, then call `diff.Compare(envelope.Data.Old, envelope.Data.New)` from `libs/core/diff` to get flattened changes (`field`, `oldValue`, `newValue`), including nested JSON paths. `Compare` returns a list and an error.

The old `foundation_person_created` topic / `foundation.person.created.v1` contract is replaced; existing subscribers must migrate to the changed topic. Messages contain contact email/phone PII and must only be consumed by authorized services. Avoid logging raw payloads, and apply appropriate topic permissions and retention.

Publication is best-effort after persistence: an unavailable broker does not roll back the database operation, so events can be lost until an outbox/retry mechanism is implemented. A pre-update read provides the old snapshot; concurrent writers require careful consideration of ordering and snapshot consistency.

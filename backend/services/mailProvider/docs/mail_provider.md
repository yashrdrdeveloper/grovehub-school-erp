# GroveHub Mail Provider

## Scope

`mailProvider` is GroveHub's outbound email provider service. It owns domain ownership verification, sender registration, message persistence, and submission to the local Haraka MTA.

It does **not** own business communication concepts such as fee notices, admission messages, account invitations, templates, channel preference, or provider routing. Those belong to the future Communication service.

It is also not a mailbox product: there is no IMAP inbox, mailbox storage, or end-user mail client.

## Runtime

```text
GroveHub backend
    ↓ typed Client or /mail-provider HTTP contract
mailProvider service
    ↓ SMTP submission (wneessen/go-mail)
Haraka
    ↓ durable outbound queue + MX delivery + retry
recipient mail server
```

Haraka accepting SMTP submission means the message is `SUBMITTED`; it does not mean the recipient server has delivered it to the final mailbox. Delivery/bounce callbacks are a later phase.

The Go SMTP library is `github.com/wneessen/go-mail`.

## Domain ownership

A domain starts as `PENDING_VERIFICATION`.

Registering `school.example` returns a DNS TXT challenge:

```text
host:  _grovehub.school.example
value: grovehub-verification=<token>
```

`POST /mail-provider/domains/{domainId}/verify` resolves that TXT record. Only a verified domain can create senders or submit mail.

Statuses:

```text
PENDING_VERIFICATION
VERIFIED
DEACTIVATED
```

A deactivated domain cannot be re-verified through the verification endpoint. Re-activation, if required later, must be an explicit lifecycle operation.

## Sender identities

A sender belongs to exactly one verified domain.

Example:

```text
Domain: school.example
Sender: fees@school.example
```

The email domain must exactly match the registered domain. A sender cannot claim an unrelated Gmail, Microsoft, or third-party address.

Statuses:

```text
ACTIVE
DEACTIVATED
```

## Messages

A message stores its sender, recipients, subject, text/HTML content, submission status, and failure reason.

Statuses:

```text
QUEUED     persisted by GroveHub before SMTP submission
SUBMITTED  accepted by Haraka for its outbound path
FAILED     Haraka submission failed
```

Haraka owns SMTP retry after submission. GroveHub does not duplicate the MTA retry queue.

## HTTP surface

```text
POST /mail-provider/domains
GET  /mail-provider/domains/{domainId}
POST /mail-provider/domains/{domainId}/verify
POST /mail-provider/domains/{domainId}/deactivate

POST /mail-provider/senders
GET  /mail-provider/senders/{senderId}
POST /mail-provider/senders/{senderId}/deactivate

POST /mail-provider/messages
GET  /mail-provider/messages/{messageId}
```

The module also exposes a typed in-process client for future Communication service usage.

## Development runtime

Haraka lives under:

```text
deployments/docker/haraka/
deployments/docker/docker-compose/mail-provider.yaml
```

Local SMTP submission is bound only to `127.0.0.1:2525`. The relay ACL permits localhost/private container networks only.

`make dev-start` is the developer entry point. It detects the operating system and ensures the container runtime is available:

```text
macOS   -> Homebrew + Colima + Docker CLI/Compose
Windows -> Docker Desktop through winget/Chocolatey
Linux   -> Docker Engine through the official installer
```

Then it starts Haraka, runs all service migrations, starts the Go backend with Air, and starts Swagger.

## Production work still required

This first service slice establishes the provider boundary and local MTA transport. Before production internet delivery, deployment must add production sender-authentication and deliverability controls including DKIM signing/key rotation, SPF guidance, DMARC guidance, static public IP/rDNS, production TLS/auth policy, bounce/delivery event ingestion, suppression handling, rate limits, and reputation monitoring.

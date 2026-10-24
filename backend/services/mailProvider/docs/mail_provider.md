# GroveHub Mail Provider

## Scope

`mailProvider` is GroveHub's outbound email provider service. It owns domain ownership verification, sender registration, message persistence, DKIM signing material, and submission to the local Haraka MTA.

It does **not** own business communication concepts such as fee notices, admission messages, account invitations, templates, channel preference, or provider routing. Those belong to the future Communication service.

It is also not a mailbox product: there is no IMAP inbox, mailbox storage, or end-user mail client.

## Runtime

```text
GroveHub backend
    ↓ typed Client or /mail-provider HTTP contract
mailProvider service
    ↓ SMTP submission (wneessen/go-mail)
Haraka
    ↓ DKIM signing + durable outbound queue + MX delivery + retry
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

The domain can also be resolved idempotently by normalized name:

```text
GET /mail-provider/domains/by-name/{domainName}
```

Statuses:

```text
PENDING_VERIFICATION
VERIFIED
DEACTIVATED
```

A deactivated domain cannot be re-verified through the verification endpoint. Re-activation, if required later, must be an explicit lifecycle operation.

## Cloudflare domain readiness

For a Cloudflare-managed sending domain, GroveHub can prepare the DNS records and Haraka DKIM key automatically.

Create a Cloudflare API token scoped to the relevant zone with:

```text
Zone -> Zone -> Read
Zone -> DNS -> Edit
```

Put it in the repository-root `.env` file:

```text
CLOUDFLARE_API_TOKEN=<token>
# Optional; normally auto-discovered from DOMAIN.
CLOUDFLARE_ZONE_ID=
MAIL_PROVIDER_DKIM_SELECTOR=groovehub1
```

Then run:

```bash
make mail-domain-ready DOMAIN=school.example PUBLIC_IP=203.0.113.10
```

The command is safe to rerun. It:

```text
1. Starts/ensures Haraka.
2. Starts the backend when it is not already running.
3. Registers the domain or reuses the existing registration by name.
4. Generates/persists a 2048-bit RSA DKIM key inside the Haraka DKIM volume.
5. Creates/updates the GroveHub ownership TXT record in Cloudflare.
6. Creates/updates groovehub1._domainkey.<domain> with the DKIM public key.
7. Creates SPF when the domain does not already have one.
8. Preserves an existing SPF and tells the operator which ip4/ip6 mechanism to merge instead of creating an invalid second SPF record.
9. Creates a non-enforcing DMARC policy when no DMARC policy exists; an existing policy is preserved.
10. Restarts Haraka so the signing key is loaded and calls the GroveHub domain verification endpoint until the ownership record resolves.
```

The Cloudflare token is never written to source control. `.env`, `.state`, and `.runtime` are gitignored.

Cloudflare does not control reverse DNS/PTR for the server public IP. Production internet delivery still requires the VPS/cloud provider to set PTR/rDNS for the sending IP to the public mail hostname used by Haraka.

### Manual HTTP flow

Register the domain:

```bash
curl --request POST \
  --url http://localhost:8080/mail-provider/domains \
  --header 'Content-Type: application/json' \
  --data '{"name":"school.example"}'
```

Read it again by name:

```bash
curl --request GET \
  --url http://localhost:8080/mail-provider/domains/by-name/school.example
```

After publishing the returned verification TXT record:

```bash
curl --request POST \
  --url http://localhost:8080/mail-provider/domains/<domainId>/verify
```

To verify a Cloudflare API token itself:

```bash
curl --request GET \
  --url https://api.cloudflare.com/client/v4/user/tokens/verify \
  --header "Authorization: Bearer ${CLOUDFLARE_API_TOKEN}"
```

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
GET  /mail-provider/domains/by-name/{domainName}
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

Local SMTP submission is bound only to `127.0.0.1:2525`. The relay ACL permits localhost/private container networks only. Haraka's outbound queue and per-domain DKIM keys use separate persistent Docker volumes.

`make dev-start` is the developer entry point. It detects the operating system and ensures the container runtime is available:

```text
macOS   -> Homebrew + Colima + Docker CLI/Compose
Windows -> Docker Desktop through winget/Chocolatey
Linux   -> Docker Engine through the official installer
```

Then it starts Haraka, runs all service migrations, starts the Go backend with Air, and starts Swagger.

## Production work still required

The service now has domain verification, SPF/DMARC provisioning guidance, per-domain DKIM signing material, and Haraka outbound queueing. Production deployment still needs a static public IP, matching PTR/rDNS and public Haraka hostname, hardened SMTP/TLS/auth policy, bounce/delivery event ingestion, suppression handling, rate limits, abuse controls, key rotation, and reputation monitoring.

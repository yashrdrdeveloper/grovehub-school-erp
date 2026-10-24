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

## Cloudflare domain readiness

For a Cloudflare-managed sending domain, GroveHub can prepare its outbound-mail DNS identity and Haraka DKIM key automatically.

Create a Cloudflare API token scoped to the sending zone with:

```text
Zone -> Zone -> Read
Zone -> DNS -> Edit
```

Put only the token and zone ID in the repository-root `.env` file:

```text
CLOUDFLARE_API_TOKEN=<token>
CLOUDFLARE_ZONE_ID=<zone-id>
```

Optional overrides exist, but they are not required for the normal flow:

```text
MAIL_PROVIDER_DKIM_SELECTOR=groovehub1
MAIL_PROVIDER_PUBLIC_IP=
```

When `MAIL_PROVIDER_PUBLIC_IP` is empty, GroveHub discovers the public address of the machine running the command through Cloudflare.

Run exactly:

```bash
make mail-domain-ready
```

The command is safe to rerun. It:

```text
1. Verifies the Cloudflare API token.
2. Reads the zone directly from CLOUDFLARE_ZONE_ID and uses the zone name as the sending domain.
3. Discovers this machine's public IP unless MAIL_PROVIDER_PUBLIC_IP overrides it.
4. Persists the derived mail.<domain> Haraka hostname under .state/mailProvider/runtime.env.
5. Starts/ensures the OS-specific container runtime and Haraka.
6. Starts the GroveHub backend when it is not already running and applies Mail Provider migrations.
7. Registers the domain or reuses the existing Mail Provider domain by normalized name.
8. Generates/reuses a 2048-bit RSA DKIM private key in the persistent Haraka DKIM volume.
9. Creates/updates a DNS-only A/AAAA record for mail.<domain> pointing at the sending host.
10. Creates/updates GroveHub's ownership TXT record.
11. Creates/updates <selector>._domainkey.<domain> with the DKIM public key.
12. Creates SPF when absent, or merges GroveHub's ip4/ip6 mechanism into the single existing SPF record.
13. Creates a non-enforcing DMARC policy when absent; existing DMARC is preserved.
14. Restarts Haraka so its public hostname and DKIM material are active.
15. Calls the GroveHub verification endpoint until the ownership TXT record resolves.
16. Checks reverse DNS/PTR when possible and checks outbound TCP/25 when `nc` is available.
```

The command does **not** create or change MX records. GroveHub Mail Provider is outbound-only, so it must not take over a school's Google Workspace, Microsoft 365, or other inbound mail routing.

The Cloudflare token is never written to source control. `.env`, `.state`, and `.runtime` are gitignored.

Cloudflare cannot configure PTR/rDNS for a public IP owned by a VPS/cloud provider. For production delivery, the provider must set:

```text
<PUBLIC_IP> PTR -> mail.<domain>
mail.<domain> A/AAAA -> <PUBLIC_IP>
```

The hosting provider must also allow outbound TCP/25 for Haraka's direct MX delivery.

### Useful Cloudflare curls

Verify the API token:

```bash
curl --request GET \
  --url https://api.cloudflare.com/client/v4/user/tokens/verify \
  --header "Authorization: Bearer ${CLOUDFLARE_API_TOKEN}"
```

Read the zone that `mail-domain-ready` will use:

```bash
curl --request GET \
  --url "https://api.cloudflare.com/client/v4/zones/${CLOUDFLARE_ZONE_ID}" \
  --header "Authorization: Bearer ${CLOUDFLARE_API_TOKEN}"
```

List the resulting mail-related records:

```bash
curl --request GET \
  --url "https://api.cloudflare.com/client/v4/zones/${CLOUDFLARE_ZONE_ID}/dns_records?per_page=100" \
  --header "Authorization: Bearer ${CLOUDFLARE_API_TOKEN}"
```

### Manual GroveHub HTTP flow

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

The service now has domain verification, outbound DNS provisioning, persistent per-domain DKIM signing material, Haraka hostname alignment, and Haraka outbound queueing. Production still requires a static public IP, matching PTR/rDNS, outbound TCP/25, hardened SMTP/TLS/auth policy, bounce/delivery event ingestion, suppression handling, rate limits, abuse controls, key rotation, and reputation monitoring.

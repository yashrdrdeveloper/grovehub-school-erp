# RDR School public website

Public, photography-led school site built in `frontend/packages/landingPage`.

## Pages

`/`, `/about`, `/academics`, `/campus`, `/admissions`, `/gallery`, `/resources`, `/contact`, `/mandatory-public-disclosure`.

The existing `/login` and protected GrooveHub dashboard remain intact.

## Authentic RDR photography

The original school gallery contains recognizable student photographs. They are
historical images, not a representation of the current academic year.

Run from `frontend/`:

```bash
pnpm assets:school
```

This downloads eight original school JPEG files into `public/school-images/`.
The pages prefer local copies but fall back to the original school URLs until
the assets are migrated. Cropping is non-destructive, using CSS `object-fit`
and carefully selected focal positions. Do not artificially change students,
staff, certificates, buildings or facilities.

Original URLs and alt text are maintained in `content.ts`.
Obtain the school's rights/consent confirmation before production deployment,
particularly for identifiable minors.

## CBSE disclosures

`https://rdrschool.com/pdf/Mandatory.pdf` is the school's **18 June 2023**
SARAS 4.0 Appendix IX record, not a current certificate.

The public disclosure page explicitly labels it as historical, preserves
its original document references, and identifies required updates.
Before production:
- Provide the current self-attested Appendix IX as directed by CBSE.
- Verify affiliation, principal, teacher qualifications, staff counts, fees,
  board results, infrastructure, certificates and current academic calendar.
- Verify and repair each certificate hyperlink; links in the 2023 PDF are
  not guaranteed to be currently active.
- Resolve historical discrepancies (campus area and CBSE/HBSE admission policy).
- Confirm the current office email and phone number; the 2023 SARAS and
  website footer do not agree on the email spelling.

Current CBSE circulars include:
https://saras.cbse.gov.in/saras/Circulars/Circular01_2026639038229508601060.pdf
https://www.cbse.gov.in/cbsenew/documents/MPD_08052026.pdf

## Enquiry flow

The parent/guardian enquiry wizard asks for first, optional middle/last names,
phone or email (at least one required), enquiry type, and optional class/message.

**No fake CRM submission:** the current repository branch does not expose a
public Lead CRM endpoint. Step 2 prepares a real `mailto:` message and
prominently explains that the lead is not submitted until the user sends it.

Production integration is a separate backend task:
- Public, rate-limited and bot-protected tenant-specific lead endpoint
- Server-side normalization, validation, deduplication and idempotency
- Minimal data collection, consent metadata, retention and audit policies
- Lead CRM status, ownership, communication history and optional acknowledgement
- School staff review; enquiry must not automatically become a student record

## Architecture

`landingPage/` owns public page composition and school content.
Once a public API exists, place transport and business logic in
`packages/business/lead/`; do not call API endpoints from generic components.

Site CSS is isolated with `.rdr-*` classes to avoid overwriting the ERP shell
or authentication UI.

## Run

```bash
cd frontend
pnpm install
pnpm assets:school
pnpm build
pnpm dev
```

Known limitation: source website imagery has not been independently approved
for republishing; current disclosure documents and Lead CRM API are not available.

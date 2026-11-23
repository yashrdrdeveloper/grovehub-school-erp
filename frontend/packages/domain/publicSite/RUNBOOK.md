# RDR School public website

Public, photography-led school site built in `frontend/packages/domain/publicSite`.

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

Original URLs and alt text are maintained in `pages/landing/constants/school.content.ts`.
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

`packages/domain/publicSite/pages/landing/` owns public page composition, `constants/school.content.ts` holds school content, and `components/publicSheet/` owns the Radix dialog.
Once a public API exists, place transport and business logic in
`packages/business/lead/`; do not call API endpoints from generic components.

Page styles live under `pages/landing/styles/` and are isolated with `.rdr-*` classes to avoid overwriting the ERP shell
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

## Responsive navigation and accessibility

- The homepage prioritizes enquiry, telephone, directions, CBSE information,
  history, leadership, academics, campus and gallery.
- On mobile, long marketing sections are omitted **from the homepage view**.
  Public information remains available via dedicated URLs or information sheets.
- Public navigation opens **standalone addressable pages**, not a generic overlay. Only enquiry actions open the `PublicSheet` (Radix Dialog, focus-trapped, Escape and X close).
- Public pages have no shell-level public layout; each supported URL is independently addressable and suitable for SEO.
- `/mandatory-public-disclosure` intentionally remains a standalone HTML
  disclosure route; it does not open in a marketing panel.
- Director, principal, history and philosophy each have dedicated accessible
  URLs; in-site links open visual summaries. The legacy leadership titles
  are not treated as verified current office holders.
- No published future events = no fabricated upcoming-event date block.
  Automatic newsletter opt-in is hidden until a real service exists.
- Enquiry form fields use unique IDs so a sheet and background page never share
  label targets. The page behind the sheet is not scrollable; the modal body
  can scroll on small viewports to keep fields usable.

## SEO

After Vite builds, `scripts/prerender-public-pages.mjs` creates static,
crawlable HTML entries for all thirteen public URLs, including metadata,
canonical links, OpenGraph, Twitter cards, schema.org JSON-LD, a sitemap and
robots rules excluding the private dashboard/auth paths.

```bash
cd frontend
pnpm build
ls dist/about/director/index.html
ls dist/mandatory-public-disclosure/index.html
ls dist/sitemap.xml dist/robots.txt
```

This generates HTML snapshots and runs a client-side SPA over them. It is
**not full server-side React rendering**; it does however provide route-specific
metadata and informative HTML for direct links. Verify production hosting
serves the generated HTML paths instead of rewriting every request to the
single generic `dist/index.html`.

## Verification

GitHub CI runs dependency-boundary checks, TypeScript, Jest interaction tests,
a production build, and assertions that SEO files exist. The tests cover
navigation, enquiry-field validation, honest email handoff, and closing the
mobile navigation sheet.

## Current public-site decomposition

The route entry `pages/landing/LandingPage.tsx` is a page composer rather than a giant component. Chrome, reusable photo/text motifs, home and about content, academic/campus/gallery sections and resources/disclosures are separated into `components/` units. Enquiry form content is presentational, with state and orchestration in `hooks/useEnquiryForm` and `containers/enquiryFormContainer`; contact validation and mailto generation are pure helpers with colocated tests. The floating school information sheet splits story/photo/info panels from its Radix dialog shell.

The school site preserves original CSS classes until responsive browser review; Tailwind v4 and shadcn-style shared controls are available via `@Components`. The Enquiry form uses the shared native `InputControl` in unstyled migration mode so it does not unexpectedly change the public website.

## Revision 6 structure

`pages/landing/LandingPage.tsx` is presentational. `containers/landingContainer` renders the controller hook `hooks/useLandingPage`, which owns navigation, enquiry dialog state, scroll and SEO effects. Route-view registration lives in `registry/schoolPages.registry.tsx`; `PUBLIC_SITE_PATHS` is exported through the publicSite domain root to the shell, preventing duplicated route lists.

Static school information and photo references live in `constants/school.content.ts`. SEO in `constants/school.seo.ts`, panel paths in `constants/school.panels.ts`, site navigation in `constants/school.navigation.ts`, enquiry initial values and sheet copy in their corresponding component-local `constants/` folders. The old generic `config/` files have been deleted, as have unused testimonial/event constants that are not rendered.

The public site does **not** render a shared authenticated `Shell` or a separate public workspace layout. Existing CSS remains scoped while visual migration proceeds.

## Website copy and static SEO boundary

Public pages share the same `ShellContainer` i18n bootstrap as authentication and protected routes. Render-time human-readable headings, captions, links, alt text, enquiry labels, validation errors and accessibility names use bare `__()`, never an imported formatter or a new global registry. Translation parameter values are primitive named `{{value}}` replacements. School content is stored as source copy in domain constants; dynamic user input is never treated as a translation key. Retain real photos and accurate CBSE disclosures.

`scripts/prerender-public-pages.mjs` generates search-indexable HTML **without executing React effects** and must keep its static, escapable title/description/body metadata. Do not depend on `window.__` in that Node prerender script. The current resolver only formats English source text, not other languages; future public multilingual SEO requires per-locale static generation, canonical/hreflang mapping and locale-specific metadata. Preserve all 13 crawlable routes, robots.txt and sitemap.xml.

**Closed public-site variants:** `PhotoMotifs` exports `PHOTO_MOTIFS_VARIANT.HERO/GALLERY/BANNER` from its own `siteMedia` API. Public-site consumers use that typed constant instead of literal `variant="gallery"`. The global component-enum checker applies to production public-site pages as well as shared UI.

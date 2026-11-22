/**
 * Pre-render public route metadata and crawlable HTML snapshots after Vite
 * build. React replaces #root on hydrationless client startup; direct URLs
 * remain useful to search crawlers and visitors with JavaScript disabled.
 *
 * This complements (rather than replaces) full SSR. Do not expose ERP data.
 */
import { readFile, mkdir, writeFile } from 'node:fs/promises';
import { join, resolve } from 'node:path';

const origin = 'https://rdrschool.com';
const out = resolve('dist');
const routes = [
  {
    path: '/',
    title: 'RDR School | Rao Deena Ram Vidya Vihar, Haluhera',
    heading: 'Rao Deena Ram Vidya Vihar School',
    description:
      'Discover RDR School in Haluhera, Rewari, Haryana. School life, academics, facilities, CBSE information and admissions enquiries.',
    body: 'Serving the Haluhera community since 2006, RDR School focuses on meaningful education and student development.',
  },
  {
    path: '/about',
    title: 'About RDR School | History, Vision & Leadership',
    heading: 'About RDR School',
    description:
      'The story of Rao Deena Ram Vidya Vihar, founded in April 2006 in Haluhera, Rewari, and the school philosophy and leadership.',
    body: 'The school was founded by Shri Lal Singh Yadav (Retired HES-II) in memory of Late Shri Rao Deena Ram Ji. Explore the history and educational philosophy of the school.',
  },
  {
    path: '/about/history',
    title: 'School History | RDR School Haluhera',
    heading: 'Our School Story',
    description:
      'Read about the founding of Rao Deena Ram Vidya Vihar in April 2006 in Haluhera, Rewari.',
    body: 'Rao Deena Ram Vidya Vihar is named in memory of Late Shri Rao Deena Ram Ji and was founded by Shri Lal Singh Yadav (Retired HES-II) in April 2006.',
  },
  {
    path: '/about/director',
    title: "Director's Message | RDR School",
    heading: "Director's Message — Shri Lal Singh Yadav",
    description:
      'Discover the original school director’s message about holistic education, academic learning, arts and all-round child development.',
    body: 'The original school website shares the director’s commitment to nurturing children through academics, sports, arts, drama and co-curricular learning.',
  },
  {
    path: '/about/principal',
    title: "Principal's Message | RDR School",
    heading: "Principal's Message — Shri Arvind Yadav",
    description:
      'Read the historical principal’s message on education, discipline, school facilities, teachers and student development.',
    body: 'The principal’s historical message reflects on the school’s growth from modest beginnings and its commitment to quality education, values and experienced teaching.',
  },
  {
    path: '/about/philosophy',
    title: 'Our Philosophy & Vision | RDR School',
    heading: 'Our School Philosophy',
    description:
      'Explore RDR School’s vision of broadening minds, student confidence, independent thinking and compassion.',
    body: 'The school philosophy is guided by the idea “Be a lamp unto yourself”, emphasizing curiosity, personal responsibility and academic, social and emotional development.',
  },
  {
    path: '/academics',
    title: 'Academics & Learning | RDR School Rewari',
    heading: 'Academics at RDR School',
    description:
      'Learn about academic stages, school learning opportunities, laboratories, activities and student development at RDR School.',
    body: 'RDR School supports learning from foundational years through senior secondary with a focus on academic understanding and character.',
  },
  {
    path: '/campus',
    title: 'Campus & Facilities | RDR School Haluhera',
    heading: 'RDR School Campus and Facilities',
    description:
      'Explore RDR School facilities including science labs, computer education, library, sports and classrooms.',
    body: 'The school’s original website describes science and computer laboratories, library, sports, transport and other student facilities.',
  },
  {
    path: '/admissions',
    title: 'School Admissions & Enquiries | RDR School',
    heading: 'Admissions and Enquiries',
    description:
      'Ask RDR School about admission, fees, transport, school classes and campus visits. Contact the school office in Haluhera, Rewari.',
    body: 'Parents and guardians can enquire about admissions with a short form. Provide a guardian first name and at least one contact method: phone number or email address.',
  },
  {
    path: '/gallery',
    title: 'School Life & Photo Gallery | RDR School',
    heading: 'RDR School Photo Gallery',
    description:
      'View RDR School photographs of activities, performances, school life and past achievements.',
    body: 'Explore original RDR School photographs from cultural programmes, student events, achievements and community activities.',
  },
  {
    path: '/resources',
    title: 'School Resources | RDR School',
    heading: 'School Resources and Documents',
    description:
      'Find school information and access the available historical fee structure, calendar, results and other reference documents.',
    body: 'Historical school documents are available for reference. Contact the school office for current academic-year information.',
  },
  {
    path: '/contact',
    title: 'Contact RDR School | Haluhera, Rewari',
    heading: 'Contact RDR School',
    description:
      'Contact Rao Deena Ram Vidya Vihar in Haluhera, Rewari, Haryana. Call +91 90500 77550 or enquire about visiting the school.',
    body: 'Address: Village Haluhera, P.O. Musepur, District Rewari, Haryana 123401. Phone: +91 90500 77550.',
  },
  {
    path: '/mandatory-public-disclosure',
    title: 'CBSE Mandatory Public Disclosure | RDR School',
    heading: 'CBSE Mandatory Public Disclosure — Appendix IX',
    description:
      'View RDR School CBSE affiliation number 531923, school code 41925, historical Appendix IX information and links to school documents.',
    body: 'Rao Deena Ram Vidya Vihar Senior Secondary School, Haluhera, school code 41925, CBSE affiliation number 531923. The currently available legacy Appendix IX PDF is dated 18 June 2023 and requires confirmation before being presented as current.',
  },
];
const escapeHtml = (value) =>
  String(value).replace(
    /[&<>"']/g,
    (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[c],
  );
const base = await readFile(join(out, 'index.html'), 'utf8');
if (!base.includes('<div id="root"></div>')) throw new Error('Cannot find the Vite root marker');

for (const route of routes) {
  const canonical = origin + (route.path === '/' ? '/' : route.path);
  const crumb = route.path
    .split('/')
    .filter(Boolean)
    .map((part) => part.replaceAll('-', ' '))
    .join(' / ');
  const external =
    route.path === '/mandatory-public-disclosure'
      ? '<p><a href="https://rdrschool.com/pdf/Mandatory.pdf">Original archived Appendix IX PDF</a></p>'
      : '';
  const schema =
    route.path === '/'
      ? {
          '@context': 'https://schema.org',
          '@type': 'School',
          name: 'Rao Deena Ram Vidya Vihar Senior Secondary School',
          alternateName: 'RDR School',
          url: origin,
          telephone: '+91 90500 77550',
          address: {
            '@type': 'PostalAddress',
            streetAddress: 'Village Haluhera, P.O. Musepur',
            addressLocality: 'Rewari',
            addressRegion: 'Haryana',
            postalCode: '123401',
            addressCountry: 'IN',
          },
          sameAs: ['https://saras.cbse.gov.in/SARAS/AffiliatedList/AfflicationDetails/531923'],
        }
      : {
          '@context': 'https://schema.org',
          '@type': 'WebPage',
          name: route.title,
          url: canonical,
          description: route.description,
          isPartOf: { '@type': 'WebSite', url: origin, name: 'RDR School' },
        };
  const snapshot = `<div id="root"><main style="max-width:780px;margin:6rem auto;padding:1.5rem;font-family:system-ui,sans-serif;line-height:1.7;color:#1f2a44"><a href="/">RDR School</a><p>${escapeHtml(crumb || 'Welcome')}</p><h1>${escapeHtml(route.heading)}</h1><p>${escapeHtml(route.body)}</p>${external}<nav aria-label="Explore RDR School"><a href="/about">About the School</a> · <a href="/academics">Academics</a> · <a href="/campus">Facilities</a> · <a href="/admissions">Admissions</a> · <a href="/mandatory-public-disclosure">CBSE Disclosure</a> · <a href="/contact">Contact</a></nav></main></div>`;
  const head = `<meta name="robots" content="index,follow" />
    <link rel="canonical" href="${escapeHtml(canonical)}" />
    <meta property="og:type" content="website" />
    <meta property="og:site_name" content="RDR School" />
    <meta property="og:url" content="${escapeHtml(canonical)}" />
    <meta property="og:title" content="${escapeHtml(route.title)}" />
    <meta property="og:description" content="${escapeHtml(route.description)}" />
    <meta name="twitter:card" content="summary" />
    <script type="application/ld+json">${JSON.stringify(schema).replace(/</g, '\\u003c')}</script>`;
  let html = base
    .replace(/<title>[\s\S]*?<\/title>/, `<title>${escapeHtml(route.title)}</title>`)
    .replace(
      /<meta name="description"[^>]*>/,
      `<meta name="description" content="${escapeHtml(route.description)}" />`,
    )
    .replace('</head>', `    ${head}\n  </head>`)
    .replace('<div id="root"></div>', snapshot);
  const location = route.path === '/' ? out : join(out, route.path.slice(1));
  await mkdir(location, { recursive: true });
  await writeFile(join(location, 'index.html'), html, 'utf8');
}
const xml = `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${routes.map((r) => `  <url><loc>${escapeHtml(origin + (r.path === '/' ? '/' : r.path))}</loc></url>`).join('\n')}\n</urlset>\n`;
await writeFile(join(out, 'sitemap.xml'), xml, 'utf8');
await writeFile(
  join(out, 'robots.txt'),
  `User-agent: *\nAllow: /\nDisallow: /dashboard\nDisallow: /identity/\nDisallow: /foundation/\nDisallow: /integrations/\nDisallow: /login\nDisallow: /forgot-password\nDisallow: /reset-password\nSitemap: ${origin}/sitemap.xml\n`,
  'utf8',
);
console.log(`Prepared ${routes.length} crawlable RDR School routes, sitemap.xml, and robots.txt`);

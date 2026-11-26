import { readFile } from 'node:fs/promises';
import { resolveFrontendSource } from './lib/source-layout.mjs';

const fromProject = (filename) => readFile(resolveFrontendSource(filename), 'utf8');
const paths = [
  'packages/styles/index.css',
  'packages/styles/tokens/foundations.css',
  'packages/styles/themes/light.css',
  'packages/styles/themes/dark.css',
  'packages/styles/base.css',
  'app/styles.css',
  'app/main.tsx',
  'components.json',
  'app/theme/index.css',
  'app/theme/school-light.css',
  'app/theme/school-dark.css',
  'app/theme/compatibility.css',
];
const sources = await Promise.all(paths.map(fromProject));
const [entry, foundations, light, dark, base, legacy, main, shadcn, schoolEntry, schoolLight, schoolDark, compatibility] = sources;
const errors = [];

const definitions = (css) =>
  new Set([...css.matchAll(/--(gh-(?:color|shadow)-[\w-]+)\s*:/g)].map((match) => match[1]));
const lightTokens = definitions(light + schoolLight);
const darkTokens = definitions(dark + schoolDark);
for (const [shared, school] of [[light, schoolLight], [dark, schoolDark]]) {
  const sharedTokens = definitions(shared);
  for (const name of definitions(school)) {
    if (!name.startsWith('gh-color-public-') || sharedTokens.has(name)) {
      errors.push(`school-only token has invalid ownership: --${name}`);
    }
  }
  if (shared.includes('--gh-color-public-')) errors.push('public-site colors belong to school');
}
if (foundations.includes('--rdr-') || !compatibility.includes('--rdr-primary:')) {
  errors.push('school legacy aliases belong in the application theme');
}
for (const name of lightTokens) {
  if (!darkTokens.has(name)) errors.push(`missing dark theme token --${name}`);
}
for (const name of darkTokens) {
  if (!lightTokens.has(name)) errors.push(`missing light theme token --${name}`);
}
if (lightTokens.size < 20) errors.push('theme token set is unexpectedly small');
if (!foundations.includes('--gh-font-body:') || !foundations.includes('--gh-text-base:')) {
  errors.push('generic typography tokens must live under packages/styles/tokens');
}
if (!base.includes('var(--gh-font-body)')) {
  errors.push('global typography base must consume design tokens');
}
if (!entry.includes('tailwindcss/utilities.css') || !entry.includes('@theme inline')) {
  errors.push('Tailwind theme must be defined in @Styles/index.css');
}
if (!entry.includes('@source "../"')) {
  errors.push('Tailwind must scan the shared frontend packages');
}
if (/(^|\n)\s*:root(?:\[data-theme=[^\]]+\])?\s*\{/m.test(legacy)) {
  errors.push('legacy app/styles.css must not redefine global theme/root tokens');
}
if (/--(?:canvas|surface|accent|gh-color-[\w-]+)\s*:/.test(legacy)) {
  errors.push('legacy app/styles.css must not redeclare shared color tokens');
}
if ((main.match(/\.\/theme\/index\.css/g) ?? []).length !== 1 ||
    (schoolEntry.match(/@groovehub\/styles\/index\.css/g) ?? []).length !== 1) {
  errors.push('import the shared stylesheet exactly once through the school theme entry');
}
if (!schoolEntry.includes("@source '../../packages'")) errors.push('Tailwind must scan school packages');
if (JSON.parse(shadcn).tailwind.css !== 'app/theme/index.css') {
  errors.push('shadcn must point at the canonical @Styles entry CSS');
}
if (errors.length) {
  for (const error of errors) process.stderr.write(`${error}\n`);
  process.exitCode = 1;
} else {
  process.stdout.write(
    `Style ownership and ${lightTokens.size} light/dark semantic tokens verified.\n`,
  );
}

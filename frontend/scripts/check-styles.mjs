import { readFile } from 'node:fs/promises';
import path from 'node:path';

const fromProject = filename => readFile(path.join(process.cwd(), filename), 'utf8');
const paths = [
  'packages/styles/index.css',
  'packages/styles/tokens/foundations.css',
  'packages/styles/themes/light.css',
  'packages/styles/themes/dark.css',
  'packages/styles/base.css',
  'app/styles.css',
  'app/main.tsx',
  'components.json',
];
const sources = await Promise.all(paths.map(fromProject));
const [entry, foundations, light, dark, base, legacy, main, shadcn] = sources;
const errors = [];

const definitions = css => new Set(
  [...css.matchAll(/--(gh-(?:color|shadow)-[\w-]+)\s*:/g)].map(match => match[1]),
);
const lightTokens = definitions(light);
const darkTokens = definitions(dark);
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
if ((main.match(/@Styles\/index\.css/g) ?? []).length !== 1) {
  errors.push('import the shared style entry exactly once from app/main.tsx');
}
if (JSON.parse(shadcn).tailwind.css !== 'packages/styles/index.css') {
  errors.push('shadcn must point at the canonical @Styles entry CSS');
}
if (errors.length) {
  for (const error of errors) process.stderr.write(`${error}\n`);
  process.exitCode = 1;
} else {
  process.stdout.write(`Style ownership and ${lightTokens.size} light/dark semantic tokens verified.\n`);
}

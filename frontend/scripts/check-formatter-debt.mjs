import { readFileSync } from 'node:fs';

const expected = [
  '/packages/domain/publicSite/pages/landing/components/aboutContent/aboutContent.tsx',
  '/packages/domain/publicSite/pages/landing/components/enquiryForm/enquiryForm.tsx',
  '/packages/domain/publicSite/pages/landing/components/publicSheet/components/infoPanels/infoPanels.tsx',
  '/packages/domain/publicSite/pages/landing/components/resourceContent/resourceContent.tsx',
  '/packages/domain/publicSite/pages/landing/components/siteChrome/siteChrome.tsx',
].sort();

const lines = readFileSync(new URL('../.prettierignore', import.meta.url), 'utf8')
  .split('\n')
  .map((line) => line.trim())
  .filter((line) => line && !line.startsWith('#'))
  .sort();

if (lines.length !== expected.length || lines.some((line, index) => line !== expected[index])) {
  process.stderr.write(
    'Formatter exemptions must remain the seven verified oversized public-site units.\n',
  );
  process.exitCode = 1;
} else {
  process.stdout.write('Formatter debt locked to seven pre-existing public-site owner units.\n');
}

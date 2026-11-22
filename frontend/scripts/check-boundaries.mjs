import { readdir, readFile } from 'node:fs/promises';
import path from 'node:path';

const project = process.cwd();
const root = path.join(project, 'packages');
const problems = [];
const sourceExtensions = new Set(['.ts', '.tsx', '.js', '.jsx', '.mjs']);

const walk = async (dir) => {
  const results = [];
  for (const entry of await readdir(dir, { withFileTypes: true })) {
    const absolute = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      results.push(...(await walk(absolute)));
    } else if (sourceExtensions.has(path.extname(entry.name))) {
      results.push(absolute);
    }
  }
  return results;
};

const forbidden = {
  base: ['@Business', '@Components', '@Domain', '@Shell', '@Styles', '@Assets', '@Auth'],
  business: ['@Components', '@Domain', '@Shell', '@Styles', '@Assets', '@Auth'],
  components: ['@Business', '@Domain', '@Shell', '@Auth'],
  styles: ['@Base', '@Business', '@Components', '@Domain', '@Shell', '@Assets', '@Auth'],
  assets: ['@Base', '@Business', '@Components', '@Domain', '@Shell', '@Styles', '@Auth'],
  domain: ['@Shell', '@Auth/'],
  shell: ['@Auth'],
};

const files = await walk(root);
for (const filename of files) {
  const source = await readFile(filename, 'utf8');
  const relative = path.relative(root, filename);
  const category = relative.split(path.sep)[0];
  const imports = Array.from(
    source.matchAll(/(?:from\s*|import\s*\(|import\s*)['"]([^'"]+)['"]/g),
    (match) => match[1],
  );
  for (const name of imports) {
    if (
      (forbidden[category] ?? []).some((prefix) => name === prefix || name.startsWith(prefix + '/'))
    ) {
      problems.push(`${relative}: illegal dependency ${name}`);
    }
    if (
      category === 'shell' &&
      name.startsWith('@Business') &&
      name !== '@Business/app/routes' &&
      name !== '@Business/auth/routes'
    ) {
      problems.push(
        `${relative}: shell may only import pure @Business/app/routes or @Business/auth/routes, not ${name}`,
      );
    }
    if (category === 'shell' && /^@Domain\/[^/]+\//.test(name)) {
      problems.push(`${relative}: shell cannot deep-import a domain: ${name}`);
    }
    if (category === 'domain' && /^@Domain\/[^/]+\//.test(name)) {
      problems.push(`${relative}: cross-domain deep import: ${name}`);
    }
    if (name === '@LandingPage' || name.startsWith('@LandingPage/')) {
      problems.push(`${relative}: retired @LandingPage alias; import @Domain/publicSite instead`);
    }
    if (name === '@Auth' || name.startsWith('@Auth/')) {
      problems.push(`${relative}: retired @Auth alias; import @Domain/auth instead`);
    }
    if (category !== 'components' && name.startsWith('@radix-ui/')) {
      problems.push(`${relative}: Radix leaf components belong in @Components, not ${name}`);
    }
    if (
      category !== 'components' &&
      (name === 'lucide-react' || name.startsWith('lucide-react/'))
    ) {
      problems.push(`${relative}: direct lucide-react import; use @Components/icons`);
    }
    if (category !== 'base' && (name === 'lodash-es' || name.startsWith('lodash-es/'))) {
      problems.push(`${relative}: use @Base/utilities instead of direct Lodash imports`);
    }
  }
}

if (problems.length > 0) {
  problems.forEach((problem) => process.stderr.write(`${problem}\n`));
  process.exitCode = 1;
} else {
  process.stdout.write(`Architecture boundaries verified across ${files.length} source files.\n`);
}

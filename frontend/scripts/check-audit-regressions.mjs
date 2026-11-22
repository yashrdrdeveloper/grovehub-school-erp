import { execFileSync } from 'node:child_process';
import { readFileSync } from 'node:fs';
import path from 'node:path';

const baseline = JSON.parse(readFileSync(new URL('./audit-allowlist.json', import.meta.url), 'utf8'));
const output = execFileSync(
  'node',
  [path.join(process.cwd(), 'scripts/audit-guidelines.mjs')],
  { cwd: process.cwd(), encoding: 'utf8', maxBuffer: 8 * 1024 * 1024 },
);
const problems = [];
let remaining = 0;

for (const line of output.split('\n')) {
  if (!line.startsWith('GH_AUDIT|')) {
    continue;
  }

  const [, filename, observed] = line.split('|');
  if (!observed || observed === 'NO_PATTERN_FLAG') {
    continue;
  }

  const allowed = new Set(baseline.allowlisted[filename] ?? []);
  for (const issue of observed.split(',')) {
    const code = issue.split('@')[0];
    remaining += 1;

    if (!allowed.has(code)) {
      problems.push(`${filename}: newly introduced guideline violation ${code}`);
    }
  }
}

if (problems.length > 0) {
  problems.forEach((problem) => process.stderr.write(problem + '\n'));
  process.exitCode = 1;
} else {
  process.stdout.write(
    `Guideline regression gate passed; ${remaining} known legacy findings still require migration.\n`,
  );
}

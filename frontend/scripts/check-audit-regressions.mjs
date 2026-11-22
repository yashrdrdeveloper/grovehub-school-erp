import { execFileSync } from 'node:child_process';
import { readFileSync } from 'node:fs';
import path from 'node:path';

import { evaluateAudit } from './lib/audit-guard.mjs';

const baseline = JSON.parse(readFileSync(new URL('./audit-allowlist.json', import.meta.url), 'utf8'));
const output = execFileSync('node', [path.join(process.cwd(), 'scripts/audit-guidelines.mjs')], {
  cwd: process.cwd(),
  encoding: 'utf8',
  maxBuffer: 8 * 1024 * 1024,
});
const result = evaluateAudit(output, baseline);

if (result.failures.length > 0) {
  result.failures.forEach((failure) => process.stderr.write(failure + '\n'));
  process.exitCode = 1;
} else {
  process.stdout.write(
    'Guideline regression gate passed: ' +
      result.flaggedFiles +
      ' flagged files, ' +
      result.findings +
      ' legacy findings.\n',
  );
}


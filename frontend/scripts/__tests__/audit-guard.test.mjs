import assert from 'node:assert/strict';
import test from 'node:test';

import { evaluateAudit, parseAuditOutput } from '../lib/audit-guard.mjs';

const baseline = {
  allowlisted: { 'old.tsx': ['RAW_CLASS'] },
  maxByCode: { RAW_CLASS: 1 },
  maxFlaggedFiles: 1,
  maxFindings: 1,
};

const outputFor = (lines, counts, flagged = 1) => {
  return [
    ...lines,
    'GH_AUDIT_SUMMARY|' + JSON.stringify({ flagged, counts }),
  ].join('\n');
};

test('parseAuditOutput', () => {
  const output = outputFor(
    ['GH_AUDIT|old.tsx|RAW_CLASS@3', 'GH_AUDIT|new.tsx|NO_PATTERN_FLAG'],
    { RAW_CLASS: 1 },
  );

  assert.deepEqual([...parseAuditOutput(output).keys()], ['old.tsx']);
});

test('evaluateAudit', async (t) => {
  await t.test('accepts existing debt within its file and code ceilings', () => {
    const output = outputFor(['GH_AUDIT|old.tsx|RAW_CLASS@3'], { RAW_CLASS: 1 });
    assert.deepEqual(evaluateAudit(output, baseline).failures, []);
  });

  await t.test('rejects retired codes even when a stale allowlist contains them', () => {
    const output = outputFor(['GH_AUDIT|old.tsx|NATIVE_COLLECTION@3'], {
      NATIVE_COLLECTION: 1,
    });
    const result = evaluateAudit(output, {
      ...baseline,
      allowlisted: { 'old.tsx': ['RAW_CLASS', 'NATIVE_COLLECTION'] },
    });
    assert.ok(result.failures.some((failure) => failure.includes('retired or unapproved')));
  });

  await t.test('rejects a known issue introduced into another file', () => {
    const output = outputFor(['GH_AUDIT|new.tsx|RAW_CLASS@3'], { RAW_CLASS: 1 });
    assert.ok(evaluateAudit(output, baseline).failures.some((failure) => failure.includes('newly introduced')));
  });

  await t.test('rejects increases in total counts even for approved issue locations', () => {
    const output = outputFor(
      ['GH_AUDIT|old.tsx|RAW_CLASS@3', 'GH_AUDIT|other.tsx|RAW_CLASS@4'],
      { RAW_CLASS: 2 },
      2,
    );
    const result = evaluateAudit(output, {
      ...baseline,
      allowlisted: { 'old.tsx': ['RAW_CLASS'], 'other.tsx': ['RAW_CLASS'] },
    });
    assert.ok(result.failures.some((failure) => failure.includes('exceed ceiling')));
  });

  await t.test('rejects missing summaries rather than passing empty output', () => {
    const result = evaluateAudit('GH_AUDIT|old.tsx|RAW_CLASS@3', baseline);
    assert.ok(result.failures.some((failure) => failure.includes('missing GH_AUDIT_SUMMARY')));
  });
});

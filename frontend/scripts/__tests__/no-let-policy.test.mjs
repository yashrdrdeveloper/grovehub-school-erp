import assert from 'node:assert/strict';
import test from 'node:test';

import { findLetDeclarations } from '../lib/no-let-policy.mjs';

test('findLetDeclarations', async (t) => {
  await t.test('finds mutable declarations including for-loops', () => {
    const code = 'let count = 0; for (let i = 0; i < 2; i++) count += i;';
    const issues = findLetDeclarations('sample.ts', code);
    assert.deepEqual(
      issues.map((x) => x.variables[0]),
      ['count', 'i'],
    );
  });
  await t.test('accepts const owners and text containing the word let', () => {
    assert.deepEqual(
      findLetDeclarations('sample.tsx', 'const x = { value: 1 }; const message = "let";'),
      [],
    );
  });
});

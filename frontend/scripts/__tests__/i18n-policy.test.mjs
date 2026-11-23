import assert from 'node:assert/strict';
import test from 'node:test';

import { analyzeI18nSource, normalizeJsxText } from '../lib/i18n-policy.mjs';

test('normalizeJsxText', () => {
  assert.equal(normalizeJsxText('\n    Check your\n    email\n  '), 'Check your email');
  assert.equal(normalizeJsxText(' to '), ' to ');
  assert.equal(normalizeJsxText('\n  \n  '), '');
});

test('analyzeI18nSource', async (t) => {
  await t.test('detects visible raw JSX and attribute copy without technical values', () => {
    const source = "const Card = () => <div aria-live=\"polite\" data-testid=\"saved\" title=\"A title\">Check your email</div>;";
    const { rawCopy, issues } = analyzeI18nSource('Card.tsx', source);
    assert.equal(rawCopy.length, 2);
    assert.deepEqual(issues, []);
  });
  await t.test('does not flag global values when translated in render', () => {
    const source = "export const Card = () => <p>{__('Hi {{name}}', { name: 'Parent' })}</p>;";
    const { rawCopy, issues } = analyzeI18nSource('Card.tsx', source);
    assert.deepEqual(rawCopy, []);
    assert.deepEqual(issues, []);
  });
  await t.test('rejects direct formatter imports and unauthorized globals', () => {
    const source = "import { __ as fmt } from '@Base/i18n'; window.__ = fmt;";
    const { issues } = analyzeI18nSource('Card.tsx', source);
    assert.ok(issues.some((x) => x.code === 'IMPORTED_I18N'));
    assert.ok(issues.some((x) => x.code === 'GLOBAL_I18N_ACCESS'));
  });
  await t.test('exempts the designated initializer but not feature code', () => {
    const source = "import { __ as resolve } from '@Base/i18n'; window.__ = resolve;";
    const { issues } = analyzeI18nSource('useInitI18n.ts', source, {
      allowResolverImport: true,
      allowRuntimeGlobal: true,
    });
    assert.deepEqual(issues, []);
  });
  await t.test('rejects calling __ during module initialization', () => {
    const { issues } = analyzeI18nSource('copy.ts', "export const TITLE = __('Hello');");
    assert.equal(issues[0].code, 'EAGER_I18N');
  });
});

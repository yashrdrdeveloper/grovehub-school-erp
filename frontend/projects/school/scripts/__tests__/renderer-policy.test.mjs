import assert from 'node:assert/strict';
import test from 'node:test';

import { validateRendererSource } from '../lib/renderer-policy.mjs';

test('validateRendererSource', async (t) => {
  await t.test('accepts a stateless JSX renderer', () => {
    assert.deepEqual(
      validateRendererSource('safe.tsx', 'export const renderBadge = () => <Box />;'),
      [],
    );
  });
  await t.test('rejects hooks in renderers', () => {
    const errors = validateRendererSource(
      'bad.tsx',
      'export const renderBadge = () => useState(1);',
    );
    assert.ok(errors.some((error) => error.includes('hooks')));
  });
  await t.test('rejects browser and network IO', () => {
    const errors = validateRendererSource(
      'bad.tsx',
      'export const renderBadge = () => fetch(window.location.href);',
    );
    assert.ok(errors.some((error) => error.includes('globals')));
    assert.ok(errors.some((error) => error.includes('Network')));
  });
});

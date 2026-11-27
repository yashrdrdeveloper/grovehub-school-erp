import assert from 'node:assert/strict';
import test from 'node:test';

import { validateComponentEnums } from '../lib/component-enum-policy.mjs';

test('validateComponentEnums', async (t) => {
  await t.test('rejects raw variants and sizes for component-owned contracts', () => {
    const errors = validateComponentEnums('Demo.tsx', '<Button variant="primary" size="lg" />');
    assert.deepEqual(
      errors.map((problem) => problem.prop),
      ['variant', 'size'],
    );
  });

  await t.test('accepts imported constant expressions', () => {
    const errors = validateComponentEnums(
      'Demo.tsx',
      '<Button variant={BUTTON_VARIANT.PRIMARY} size={BUTTON_SIZE.LARGE} />',
    );
    assert.deepEqual(errors, []);
  });

  await t.test('does not treat DOM protocol options or arbitrary user text as enum props', () => {
    const errors = validateComponentEnums(
      'Demo.tsx',
      '<><input type="email" autoComplete="on" /><Box aria-live="polite" title="School" /></>',
    );
    assert.deepEqual(errors, []);
  });

  await t.test('rejects raw PhotoMotifs variants and Box tags', () => {
    const errors = validateComponentEnums(
      'Demo.tsx',
      '<><PhotoMotifs variant="gallery" /><Box as="h3" /></>',
    );
    assert.equal(errors.length, 2);
  });
});

import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { test } from 'node:test';

import { createSourceLayout } from '../lib/source-layout.mjs';

test('maps an extracted package to its original logical audit owner', () => {
  const directory = mkdtempSync(path.join(tmpdir(), 'groovehub-source-'));
  try {
    const app = path.join(directory, 'frontend');
    const shared = path.join(directory, 'libs/frontend/base');
    mkdirSync(path.join(app, 'packages'), { recursive: true });
    mkdirSync(shared, { recursive: true });
    const layout = createSourceLayout(app);
    assert.equal(layout.resolvePackageSource('base/index.ts'), path.join(shared, 'index.ts'));
    assert.equal(
      layout.packageRelativePath(path.join(shared, 'index.ts')),
      path.join('base', 'index.ts'),
    );
    assert.equal(layout.frontendRelativePath(path.join(shared, 'index.ts')), 'packages/base/index.ts');
    assert.equal(layout.resolveFrontendSource('packages/base/index.ts'), path.join(shared, 'index.ts'));
    assert.equal(layout.packageSourceRoots.length, 2);
    mkdirSync(path.join(app, 'packages/base'));
    assert.throws(() => createSourceLayout(app), /Duplicate package owner/);
  } finally {
    rmSync(directory, { recursive: true, force: true });
  }
});

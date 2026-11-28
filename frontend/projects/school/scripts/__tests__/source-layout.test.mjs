import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { test } from 'node:test';

import { createSourceLayout } from '../lib/source-layout.mjs';

test('maps an extracted package to its original logical audit owner', () => {
  const directory = mkdtempSync(path.join(tmpdir(), 'groovehub-source-'));
  try {
    const app = path.join(directory, 'frontend/projects/school');
    const shared = path.join(directory, 'frontend/libs/base');
    mkdirSync(path.join(app, 'packages'), { recursive: true });
    mkdirSync(shared, { recursive: true });
    const layout = createSourceLayout(app, ['base']);
    assert.equal(layout.resolvePackageSource('base/index.ts'), path.join(shared, 'index.ts'));
    assert.equal(
      layout.packageRelativePath(path.join(shared, 'index.ts')),
      path.join('base', 'index.ts'),
    );
    assert.equal(
      layout.frontendRelativePath(path.join(shared, 'index.ts')),
      'packages/base/index.ts',
    );
    assert.equal(
      layout.resolveFrontendSource('packages/base/index.ts'),
      path.join(shared, 'index.ts'),
    );
    assert.equal(layout.packageSourceRoots.length, 2);
    mkdirSync(path.join(app, 'packages/base'));
    assert.throws(() => createSourceLayout(app, ['base']), /Duplicate package owner/);
  } finally {
    rmSync(directory, { recursive: true, force: true });
  }
});

test('retains logical component ownership across a physical workspace move', () => {
  const directory = mkdtempSync(path.join(tmpdir(), 'groovehub-components-'));
  try {
    const app = path.join(directory, 'frontend/projects/school');
    const local = path.join(app, 'packages/components');
    const shared = path.join(directory, 'frontend/libs/components');
    mkdirSync(local, { recursive: true });
    const before = createSourceLayout(app, ['components']);
    assert.equal(before.resolvePackageSource('components/index.ts'), path.join(local, 'index.ts'));
    rmSync(local, { recursive: true });
    mkdirSync(shared, { recursive: true });
    const after = createSourceLayout(app, ['components']);
    assert.equal(after.resolvePackageSource('components/index.ts'), path.join(shared, 'index.ts'));
    assert.equal(
      after.frontendRelativePath(path.join(shared, 'index.ts')),
      'packages/components/index.ts',
    );
    assert.equal(
      after.packageRelativePath(path.join(shared, 'index.ts')),
      path.join('components', 'index.ts'),
    );
  } finally {
    rmSync(directory, { recursive: true, force: true });
  }
});

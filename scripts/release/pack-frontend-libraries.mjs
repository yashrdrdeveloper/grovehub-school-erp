// Pack tracked, source-distributed libraries. No tests, workspace references,
// application code or node_modules may appear in the release candidates.
import { execFileSync } from 'node:child_process';
import { mkdirSync, mkdtempSync, readFileSync, writeFileSync, copyFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';

const [destination, revision] = process.argv.slice(2);
if (!destination || !/^[0-9a-f]{40}$/.test(revision ?? '')) {
  throw new Error('Usage: pack-frontend-libraries.mjs DESTINATION COMMIT');
}
const root = execFileSync('git', ['rev-parse', '--show-toplevel'], { encoding: 'utf8' }).trim();
const output = path.resolve(destination);
const stage = mkdtempSync(path.join(tmpdir(), 'gh-libraries-'));
const libraries = ['base', 'styles', 'components'];
const manifests = new Map(libraries.map((name) => [
  name, JSON.parse(readFileSync(path.join(root, 'frontend/libs', name, 'package.json'), 'utf8')),
]));
const versions = new Map(libraries.map((name) => [
  '@groovehub/' + name, manifests.get(name).version + '-ci.' + revision.slice(0, 12),
]));

mkdirSync(output, { recursive: true });
try {
  for (const name of libraries) {
    const prefix = 'frontend/libs/' + name + '/';
    const files = execFileSync('git', ['ls-files', '-z', prefix], { cwd: root, encoding: 'utf8' }).split('\0').filter(Boolean);
    const packageRoot = path.join(stage, name);
    mkdirSync(packageRoot);
    for (const file of files) {
      const relative = file.slice(prefix.length);
      if (/(^|\/)(__tests__|tests|node_modules|dist|coverage)(\/|$)/.test(relative)) continue;
      if (/\.(test|spec)\.[^.]+$/.test(relative) || relative === 'package.json') continue;
      if (!/\.(ts|tsx|css|scss|svg|json|md)$/.test(relative) || /(^|\/)\./.test(relative)) continue;
      if (relative === 'tsconfig.json') continue;
      const target = path.join(packageRoot, relative);
      mkdirSync(path.dirname(target), { recursive: true });
      copyFileSync(path.join(root, file), target);
    }
    const original = manifests.get(name);
    const dependencies = Object.fromEntries(Object.entries(original.dependencies ?? {}).map(([dependency, version]) => {
      if (!version.startsWith('workspace:')) return [dependency, version];
      const pinned = versions.get(dependency);
      if (!pinned) throw new Error('Unknown workspace dependency: ' + dependency);
      return [dependency, pinned];
    }));
    const manifest = {
      name: original.name, version: versions.get(original.name), private: true,
      type: original.type, exports: original.exports, dependencies,
      peerDependencies: original.peerDependencies,
      repository: { type: 'git', url: 'https://github.com/yashrdrdeveloper/grovehub-school-erp', directory: prefix.slice(0,-1) },
    };
    writeFileSync(path.join(packageRoot, 'package.json'), JSON.stringify(manifest,null,2)+'\n');
    execFileSync('npm', ['pack', '--ignore-scripts', '--pack-destination', output], { cwd: packageRoot, stdio: 'inherit' });
  }
  writeFileSync(path.join(output, 'manifest.json'), JSON.stringify({
    commit: revision, format: 'source-packages', versions: Object.fromEntries(versions),
  }, null, 2)+'\n');
} finally {
  rmSync(stage, { recursive: true, force: true });
}

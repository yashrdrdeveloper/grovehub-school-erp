import { readFileSync, writeFileSync } from 'node:fs';
import path from 'node:path';

const [consumer, artifacts] = process.argv.slice(2);
if (!consumer || !artifacts) throw new Error('Usage: prepare-frontend-consumer.mjs CONSUMER ARTIFACTS');
const manifest = JSON.parse(readFileSync(path.join(artifacts, 'manifest.json'), 'utf8'));
const packagePath = path.join(consumer, 'package.json');
const pkg = JSON.parse(readFileSync(packagePath, 'utf8'));
const dependencies = Object.fromEntries(Object.entries(manifest.versions).map(([name, version]) => [
  name, 'file:' + path.relative(path.resolve(consumer), path.join(path.resolve(artifacts), name.replace('@','').replace('/','-') + '-' + version + '.tgz')),
]));
Object.assign(pkg.dependencies, dependencies);
pkg.pnpm = { ...(pkg.pnpm ?? {}), overrides: { ...(pkg.pnpm?.overrides ?? {}), ...dependencies } };
writeFileSync(packagePath, JSON.stringify(pkg,null,2)+'\n');
for (const file of ['vite.config.ts', 'tsconfig.json']) {
  const target = path.join(consumer,file);
  const source = readFileSync(target,'utf8').replaceAll('../../libs/', './node_modules/@groovehub/');
  writeFileSync(target,source);
}

const tsconfigPath = path.join(consumer, 'tsconfig.json');
const tsconfig = JSON.parse(readFileSync(tsconfigPath, 'utf8'));
const globals = './node_modules/@groovehub/components/types/i18n.global.d.ts';
tsconfig.include = [...new Set([...(tsconfig.include ?? []), globals])];
writeFileSync(tsconfigPath, JSON.stringify(tsconfig, null, 2) + '\n');

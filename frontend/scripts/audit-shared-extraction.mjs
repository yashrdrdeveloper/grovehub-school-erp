import { execFileSync } from 'node:child_process';
import { readFileSync } from 'node:fs';
import path from 'node:path';
import ts from 'typescript';

const repo = path.resolve(process.cwd(), '..');
const files = execFileSync('git', ['ls-files', '-z'], { cwd: repo, encoding: 'utf8' })
  .split('\0')
  .filter(Boolean);
const sourceFiles = files.filter((file) => /\.(?:tsx?|css|go)$/.test(file));
const candidates = [
  ['COMPONENTS', 'frontend/packages/components/', 'libs/frontend/components/'],
  ['CORE', 'backend/libs/core/', 'libs/backend/core/'],
];
const getOwner = (file) =>
  candidates.find(([, ...roots]) => roots.some((root) => file.startsWith(root)));
const report = (type, value) =>
  process.stdout.write('GH_EXTRACTION|' + type + '|' + JSON.stringify(value) + '\n');
const dependencies = new Map();
const findings = [];
const count = { components: 0, core: 0 };
const addDependency = (owner, name) => {
  const list = dependencies.get(owner) ?? new Set();
  list.add(name);
  dependencies.set(owner, list);
};
const reportFinding = (owner, file, kind, detail) => {
  findings.push({ owner, file, kind, detail });
};
const scanFrontend = (file, source, owner) => {
  if (file.endsWith('.css')) {
    for (const match of source.matchAll(/var\((--(?:rdr-|gh-color-public-)[\w-]+)/g)) {
      reportFinding(owner, file, 'SCHOOL_TOKEN', match[1]);
    }
    return;
  }
  const tree = ts.createSourceFile(file, source, ts.ScriptTarget.Latest, true);
  const visit = (node) => {
    if (ts.isImportDeclaration(node) || ts.isExportDeclaration(node)) {
      const specifier = node.moduleSpecifier;
      if (specifier && ts.isStringLiteral(specifier)) {
        const name = specifier.text;
        if (!name.startsWith('.')) addDependency(owner, name);
        if (/^@(?:Assets|Business|Domain|Shell|Auth)(?:\/|$)/.test(name)) {
          reportFinding(owner, file, 'APP_IMPORT', name);
        }
        if (name.startsWith('.')) {
          const resolved = path.posix.normalize(path.posix.join(path.posix.dirname(file), name));
          if (!getOwner(resolved)) reportFinding(owner, file, 'EXTERNAL_RELATIVE', resolved);
        }
      }
    }
    ts.forEachChild(node, visit);
  };
  visit(tree);
  if (/\b(?:localStorage|sessionStorage)\b/.test(source))
    reportFinding(owner, file, 'STORAGE', 'review namespace and app ownership');
  if (/\b(?:RDR|SchoolMonogram|School|groovehub\.theme)\b/.test(source))
    reportFinding(owner, file, 'BRANDING', 'review product marker');
};
const scanCore = (file, source, owner) => {
  const imports = [
    ...source.matchAll(
      /(?:^\s*import\s+(?:\w+\s+)?|^\s*(?:\w+\s+)?)"([\w.-]+(?:\/[\w.-]+)+)"\s*$/gm,
    ),
  ].map((match) => match[1]);
  for (const name of imports) {
    addDependency(owner, name);
    if (/^grovehub\/(?:services|launcher)(?:\/|$)/.test(name))
      reportFinding(owner, file, 'SCHOOL_IMPORT', name);
  }
  const rules = [
    ['ROOT_ASSUMPTION', /(?:libs\/core|filepath\.Join\([^\n]*"services"|grovehub\/)/],
    ['IDENTITY_CONTRACT', /(?:__Host-gh_|__Secure-gh_|gh_csrf|PersonGID|\bAdmin\b|TenantID)/],
    [
      'PRODUCT_REFERENCE',
      /(?:identity\.|foundation\.|integration\.|infrastructure\.|SCHOOL_ADMINISTRATION)/,
    ],
  ];
  for (const [kind, pattern] of rules) {
    if (pattern.test(source)) reportFinding(owner, file, kind, 'review source and test ownership');
  }
};
for (const file of sourceFiles) {
  const entry = getOwner(file);
  if (!entry) continue;
  const [owner] = entry;
  const source = readFileSync(path.join(repo, file), 'utf8');
  if (owner === 'COMPONENTS') {
    count.components += 1;
    scanFrontend(file, source, owner);
  } else {
    count.core += 1;
    scanCore(file, source, owner);
  }
}
for (const [owner, names] of dependencies)
  report('DEPENDENCIES', { owner, names: [...names].sort() });
for (const finding of findings) report('FINDING', finding);
report('SUMMARY', {
  ...count,
  findings: findings.length,
  note: 'Static inventory, not a completed semantic or runtime audit',
});

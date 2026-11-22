import { execFileSync } from 'node:child_process';
import { readFileSync } from 'node:fs';
import path from 'node:path';
import ts from 'typescript';

const root = process.cwd();
const files = execFileSync('git', ['ls-files', '--cached', '--', '.'], { cwd: root, encoding: 'utf8' })
  .split('\n').filter(Boolean).sort();
const nativeTags = new Set(['div', 'span', 'p', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
  'button', 'input', 'label', 'form', 'section', 'main', 'header', 'footer',
  'aside', 'nav', 'small', 'strong', 'a']);
const nativeOps = new Set(['map', 'flatMap', 'filter', 'find', 'includes', 'trim', 'keys', 'entries']);
const issues = new Map();
let sourceFiles = 0;
function flag(file, kind, line) {
  const found = issues.get(file) ?? new Map();
  if (!found.has(kind)) found.set(kind, line);
  issues.set(file, found);
}
for (const file of files) {
  const tsSource = /\.tsx?$/.test(file);
  const test = /\.test\.[jt]sx?$/.test(file);
  const feature = file.startsWith('packages/domain/') || file.startsWith('packages/shell/');
  const base = file.startsWith('packages/base/');
  const components = file.startsWith('packages/components/');
  const styles = file.startsWith('packages/styles/');
  if (/^packages\/domain\/(accounts|people|integrations|settings)\//.test(file)) flag(file, 'UNUSED_DOMAIN_REVIEW', 1);
  if (file.startsWith('packages/') && /\/config\//.test(file)) flag(file, 'CONFIG_FOLDER', 1);
  if (/\/helpers\/[^/]+\/(?:index\.ts|[^/]+\.ts|types\/|__tests__\/)/.test(file) && !base) flag(file, 'NESTED_HELPER', 1);
  if (/\/(?:helpers|constants|renderer)\/index\.ts$/.test(file)) flag(file, 'FORBIDDEN_BARREL', 1);
  if (test && !file.includes('/__tests__/')) flag(file, 'TEST_LOCATION', 1);
  if (/\.css$/.test(file) && file.startsWith('packages/domain/')) flag(file, 'LEGACY_DOMAIN_CSS', 1);
  if (file === 'app/styles.css') flag(file, 'LEGACY_GLOBAL_CSS', 1);
  if (file === 'packages/domain/auth/core/providers/AuthProvider.tsx') flag(file, 'FLAT_PROVIDER', 1);
  if (file === 'packages/domain/auth/core/guards/RecoveryFlowBoundary.tsx') flag(file, 'EFFECT_COMPONENT', 1);
  if (file === 'packages/shell/layout/AppLauncher.tsx') flag(file, 'MONOLITHIC_LAUNCHER', 1);
  if (!tsSource && !/\.css$/.test(file)) continue;
  const source = readFileSync(path.join(root, file), 'utf8');
  sourceFiles += 1;
  if (tsSource && !test && source.split('\n').length > 155) flag(file, 'LARGE_UNIT', 156);
  if (/\.css$/.test(file)) {
    if (source.length > 15000) flag(file, 'LARGE_STYLESHEET', 1);
    if (!styles && /#[0-9a-f]{3,8}\b|rgba?\(/i.test(source)) flag(file, 'CSS_COLOR_LITERAL', 1);
    continue;
  }
  const sf = ts.createSourceFile(file, source, ts.ScriptTarget.Latest, true,
    file.endsWith('.tsx') ? ts.ScriptKind.TSX : ts.ScriptKind.TS);
  const describes = [];
  const importedUnitNames = new Set();
  const line = node => sf.getLineAndCharacterOfPosition(node.getStart(sf)).line + 1;
  const visit = node => {
    if (ts.isConditionalExpression(node) && !test) flag(file, 'TERNARY', line(node));
    if (ts.isImportDeclaration(node) && ts.isStringLiteral(node.moduleSpecifier)) {
      const spec = node.moduleSpecifier.text;
      if (test && spec.startsWith('.') && node.importClause) {
        const binding = node.importClause.namedBindings;
        if (binding && ts.isNamedImports(binding)) {
          for (const item of binding.elements) importedUnitNames.add(item.name.text);
        }
        if (node.importClause.name) importedUnitNames.add(node.importClause.name.text);
      }
      if (spec === 'lucide-react' && !components) flag(file, 'DIRECT_LUCIDE', line(node));
      if ((spec === 'lodash-es' || spec.startsWith('lodash-es/')) && !base) flag(file, 'DIRECT_LODASH', line(node));
    }
    if (feature && !test && (ts.isJsxOpeningElement(node) || ts.isJsxSelfClosingElement(node))) {
      const tag = node.tagName.getText(sf);
      if (nativeTags.has(tag)) flag(file, 'RAW_NATIVE_HTML', line(node));
      for (const attr of node.attributes.properties) {
        if (!ts.isJsxAttribute(attr)) continue;
        const name = attr.name.getText(sf);
        if (!attr.initializer || !ts.isStringLiteral(attr.initializer)) continue;
        if (name === 'className') flag(file, 'RAW_CLASS', line(attr));
        if (name === 'dataTestId' || name === 'data-testid') flag(file, 'TEST_ID_LITERAL', line(attr));
        if ((name === 'to' || name === 'path') && attr.initializer.text.startsWith('/')) flag(file, 'ROUTE_LITERAL', line(attr));
        if (['fill', 'color', 'stroke'].includes(name) && /^#[0-9a-f]{3,8}$/i.test(attr.initializer.text)) flag(file, 'INLINE_COLOR', line(attr));
      }
    }
    if (feature && !test && ts.isCallExpression(node) && ts.isPropertyAccessExpression(node.expression)
      && nativeOps.has(node.expression.name.text)) flag(file, 'NATIVE_COLLECTION', line(node));
    if (test && ts.isCallExpression(node)) {
      const callee = node.expression.getText(sf);
      if (callee === 'describe') {
        const arg = node.arguments[0];
        describes.push(arg && ts.isStringLiteral(arg) ? arg.text : '<dynamic>');
        let parent = node.parent;
        while (parent) {
          if (ts.isCallExpression(parent) && parent.expression.getText(sf) === 'describe') {
            flag(file, 'NESTED_DESCRIBE', line(node)); break;
          }
          parent = parent.parent;
        }
      }
      if (callee === 'jest.mock' && node.arguments[0] && ts.isStringLiteral(node.arguments[0])
        && /helpers|utilities|lodash/.test(node.arguments[0].text)) flag(file, 'PURE_MOCK', line(node));
    }
    ts.forEachChild(node, visit);
  };
  visit(sf);
  if (test) {
    if (describes.length !== 1) flag(file, 'DESCRIBE_COUNT', 1);
    if (describes.length === 1 && !importedUnitNames.has(describes[0])) {
      flag(file, 'DESCRIBE_EXPORT_MISMATCH', 1);
    }
  }
}
const counts = {};
for (const file of files) {
  const found = Array.from(issues.get(file) ?? []).map(([name, line]) => name + '@' + line);
  for (const finding of found) {
    const name = finding.split('@')[0];
    counts[name] = (counts[name] ?? 0) + 1;
  }
  process.stdout.write('GH_AUDIT|' + file + '|' + (found.join(',') || 'NO_PATTERN_FLAG') + '\n');
}
process.stdout.write('GH_AUDIT_SUMMARY|' + JSON.stringify({
  files: files.length, sourceFiles, flagged: issues.size, counts,
}) + '\n');

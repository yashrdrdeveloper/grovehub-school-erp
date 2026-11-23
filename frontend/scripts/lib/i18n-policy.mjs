import ts from 'typescript';

/** UI copy lives at render time, whereas transport/HTML values retain native semantics. */
export const VISIBLE_TEXT_ATTRIBUTES = new Set([
  'title',
  'label',
  'description',
  'mobileTitle',
  'mobileDescription',
  'placeholder',
  'aria-label',
  'aria-description',
  'alt',
  'heading',
  'helperText',
  'loadingText',
  'buttonLabel',
  'emptyLabel',
  'caption',
  'message',
  'tooltip',
  'successMessage',
  'errorMessage',
  'hint',
  'secondaryLabel',
  'errorText',
  'helpText',
]);

const isVisible = (message) => /[A-Za-z\u00c0-\u024f]/.test(message);
const hasEntity = (message) => /&(?:[a-zA-Z][a-zA-Z0-9]+|#[xX]?[0-9a-fA-F]+);/.test(message);
const forbiddenGlobal = (node) =>
  ts.isPropertyAccessExpression(node) &&
  node.name.text === '__' &&
  (node.expression.getText() === 'window' || node.expression.getText() === 'globalThis');

/** JSX transforms normalize layout indentation; reproduce the same textual whitespace. */
export const normalizeJsxText = (input) => {
  const lines = input.replace(/\t/g, ' ').split(/\r\n|\r|\n/);
  const lastNonEmpty = lines.findLastIndex((line) => Boolean(line.trim()));
  if (lastNonEmpty < 0) return '';

  const normalized = lines.map((original, index) => {
    const withoutLeft = index === 0 ? original : original.replace(/^ +/, '');
    if (index === lines.length - 1) return withoutLeft;
    return withoutLeft.replace(/ +$/, '');
  });
  return normalized.reduce((result, line, index) => {
    if (!line) return result;
    return result + line + (index === lastNonEmpty ? '' : ' ');
  }, '');
};

export const analyzeI18nSource = (filename, source, options = {}) => {
  const file = ts.createSourceFile(
    filename,
    source,
    ts.ScriptTarget.Latest,
    true,
    ts.ScriptKind.TSX,
  );
  const issues = [];
  const rawCopy = [];
  const allowRuntimeGlobal = options.allowRuntimeGlobal === true;
  const allowResolverImport = options.allowResolverImport === true;
  const checkVisibleCopy = options.checkVisibleCopy !== false;
  const allowEagerCopy = options.allowEagerCopy === true;
  const add = (list, node, code, label) => {
    const loc = file.getLineAndCharacterOfPosition(node.getStart(file));
    list.push({ code, line: loc.line + 1, label });
  };
  const visit = (node, insideFunction = false) => {
    if (
      ts.isImportDeclaration(node) &&
      ts.isStringLiteral(node.moduleSpecifier) &&
      ['@Base', '@Base/i18n', '@Base/i18n/__'].includes(node.moduleSpecifier.text) &&
      !allowResolverImport
    ) {
      const clause = node.importClause;
      if (
        clause?.namedBindings &&
        ts.isNamedImports(clause.namedBindings) &&
        clause.namedBindings.elements.some(
          (name) => name.name.text === '__' || name.propertyName?.text === '__',
        )
      ) {
        add(issues, node, 'IMPORTED_I18N', 'Importing __ in a frontend owner is forbidden');
      }
    }
    if (forbiddenGlobal(node) && !allowRuntimeGlobal) {
      add(issues, node, 'GLOBAL_I18N_ACCESS', 'Only ShellContainer bootstrap installs window.__');
    }
    if (
      ts.isCallExpression(node) &&
      ts.isIdentifier(node.expression) &&
      node.expression.text === '__' &&
      !insideFunction &&
      !allowEagerCopy
    ) {
      add(
        issues,
        node,
        'EAGER_I18N',
        'Translate UI copy at render time, not during module initialization',
      );
    }
    if (ts.isCallExpression(node) && ts.isIdentifier(node.expression) &&
        node.expression.text === '__') {
      const first = node.arguments[0];
      if (!first || !(ts.isStringLiteral(first) || ts.isNoSubstitutionTemplateLiteral(first))) {
        add(issues, node, 'I18N_LITERAL_REQUIRED', '__() must receive a literal text string, not a variable or expression');
      }
    }
    if (checkVisibleCopy && ts.isJsxText(node)) {
      const message = normalizeJsxText(node.getText(file));
      if (isVisible(message) && !hasEntity(message)) {
        add(rawCopy, node, 'RAW_UI_COPY', message.slice(0, 100));
      }
    }
    if (
      checkVisibleCopy &&
      ts.isJsxAttribute(node) &&
      VISIBLE_TEXT_ATTRIBUTES.has(node.name.text) &&
      node.initializer &&
      ts.isStringLiteral(node.initializer) &&
      isVisible(node.initializer.text)
    ) {
      add(rawCopy, node, 'RAW_UI_COPY', node.initializer.text.slice(0, 100));
    }
    if (
      checkVisibleCopy &&
      ts.isJsxExpression(node) &&
      node.expression &&
      ts.isStringLiteral(node.expression) &&
      isVisible(node.expression.text)
    ) {
      add(rawCopy, node, 'RAW_UI_COPY', node.expression.text.slice(0, 100));
    }
    const nextInsideFunction = insideFunction || ts.isFunctionLike(node);
    ts.forEachChild(node, (child) => visit(child, nextInsideFunction));
  };
  ts.forEachChild(file, (child) => visit(child));
  return { issues, rawCopy };
};

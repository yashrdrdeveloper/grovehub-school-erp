import ts from 'typescript';

/** Pure syntax analysis: renderer execution must be hook-free and browser-IO-free. */
export const validateRendererSource = (filename, source) => {
  const tree = ts.createSourceFile(filename, source, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX);
  const violations = [];
  const forbiddenGlobals = new Set([
    'window', 'document', 'localStorage', 'sessionStorage', 'navigator', 'fetch',
  ]);
  const visit = (node) => {
    if (ts.isCallExpression(node) && ts.isIdentifier(node.expression)) {
      const invoked = node.expression.text;
      if (/^use[A-Z]/.test(invoked)) violations.push('React hooks are forbidden inside renderer functions: ' + invoked);
      if (invoked === 'fetch') violations.push('Network operations are forbidden in renderers');
    }
    if (ts.isIdentifier(node) && forbiddenGlobals.has(node.text)) {
      violations.push('Browser/network globals are forbidden in renderers: ' + node.text);
    }
    ts.forEachChild(node, visit);
  };
  visit(tree);
  return [...new Set(violations)];
};

import ts from 'typescript';

/** Immutable binding policy applies to production, scripts and unit tests. */
export const findLetDeclarations = (filename, source) => {
  const kind = filename.endsWith('.tsx') || filename.endsWith('.jsx')
    ? ts.ScriptKind.TSX
    : ts.ScriptKind.TS;
  const tree = ts.createSourceFile(filename, source, ts.ScriptTarget.Latest, true, kind);
  const results = [];
  const visit = (node) => {
    if (ts.isVariableDeclarationList(node) && (node.flags & ts.NodeFlags.Let)) {
      const line = tree.getLineAndCharacterOfPosition(node.getStart(tree)).line + 1;
      results.push({
        line,
        variables: node.declarations.map((declaration) => declaration.name.getText(tree)),
      });
    }
    ts.forEachChild(node, visit);
  };
  visit(tree);
  return results;
};

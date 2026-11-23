import ts from 'typescript';

/** Only actual closed component design props are controlled. HTML wire attributes are not. */
export const COMPONENT_CLOSED_PROPS = Object.freeze({
  Button: ['variant', 'size'],
  Box: ['as'],
  Text: ['as', 'tone', 'variant', 'size', 'weight'],
  Typography: ['as', 'tone', 'variant', 'size', 'weight'],
  Card: ['variant', 'padding'],
  Container: ['as', 'maxWidth'],
  Form: ['layout', 'spacing'],
  Input: ['controlSize'],
  InputControl: ['controlSize'],
  TextField: ['controlSize'],
  PasswordField: ['controlSize'],
  StatusPill: ['tone'],
  PhotoMotifs: ['variant'],
  LoginLayout: ['variant'],
});

/** Reports literal JSX closed-set props; callers must use the owner's runtime constants. */
export const validateComponentEnums = (filename, source) => {
  const sf = ts.createSourceFile(filename, source, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX);
  const problems = [];
  const visit = (node) => {
    if ((ts.isJsxOpeningElement(node) || ts.isJsxSelfClosingElement(node)) &&
        ts.isIdentifier(node.tagName)) {
      const closed = COMPONENT_CLOSED_PROPS[node.tagName.text];
      if (closed) {
        for (const attr of node.attributes.properties) {
          if (!ts.isJsxAttribute(attr) || !closed.includes(attr.name.text) ||
              !attr.initializer || !ts.isStringLiteral(attr.initializer)) continue;
          const line = sf.getLineAndCharacterOfPosition(attr.getStart(sf)).line + 1;
          problems.push({
            line,
            component: node.tagName.text,
            prop: attr.name.text,
            value: attr.initializer.text,
          });
        }
      }
    }
    ts.forEachChild(node, visit);
  };
  visit(sf);
  return problems;
};

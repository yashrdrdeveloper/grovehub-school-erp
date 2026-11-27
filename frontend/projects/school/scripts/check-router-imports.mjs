import { readdir, readFile } from 'node:fs/promises';
import path from 'node:path';
import ts from 'typescript';

/** Router infrastructure uses @Components/router; links use AppLink/AppNavLink. */
const violations = [];
const root = path.join(process.cwd(), 'packages');

const visit = async (folder) => {
  for (const entry of await readdir(folder, { withFileTypes: true })) {
    const file = path.join(folder, entry.name);

    if (entry.isDirectory()) {
      if (entry.name !== '__tests__') {
        await visit(file);
      }

      continue;
    }

    if (!/\.tsx?$/.test(entry.name) || /\.test\.tsx?$/.test(entry.name)) {
      continue;
    }

    const source = await readFile(file, 'utf8');
    const ast = ts.createSourceFile(file, source, ts.ScriptTarget.Latest, true);
    for (const statement of ast.statements) {
      if (!ts.isImportDeclaration(statement)) {
        continue;
      }

      if (!ts.isStringLiteral(statement.moduleSpecifier)) {
        continue;
      }

      const specifier = statement.moduleSpecifier.text;
      const relative = path.relative(process.cwd(), file);

      if (specifier === 'react-router-dom') {
        violations.push(relative + ': use the shared router boundary');
      }

      if (specifier !== '@Components/router') {
        continue;
      }

      const bindings = statement.importClause?.namedBindings;
      if (!bindings || !ts.isNamedImports(bindings)) {
        continue;
      }

      for (const binding of bindings.elements) {
        const imported = binding.propertyName?.text ?? binding.name.text;
        if (imported === 'Link' || imported === 'NavLink') {
          violations.push(relative + ': use AppLink or AppNavLink instead of router ' + imported);
        }
      }
    }
  }
};

await visit(path.join(root, 'domain'));
await visit(path.join(root, 'shell'));

if (violations.length > 0) {
  for (const source of violations) {
    process.stderr.write(source + '\\n');
  }

  process.exitCode = 1;
} else {
  process.stdout.write('Domain and shell navigation uses shared router and link adapters.\\n');
}

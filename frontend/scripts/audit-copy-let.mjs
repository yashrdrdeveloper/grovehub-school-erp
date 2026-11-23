import { readFile, readdir } from 'node:fs/promises';
import path from 'node:path';
import ts from 'typescript';

const results = { lets: [], translated: [], copy: [], eager: [] };
const root = process.cwd();
const files = [];
const walk = async (dir) => {
  for (const x of await readdir(dir, { withFileTypes: true })) {
    if (x.isDirectory()) {
      if (['node_modules', '.generated', 'dist'].includes(x.name)) continue;
      await walk(path.join(dir, x.name));
    } else if (/\.[mc]?[jt]sx?$/.test(x.name)) files.push(path.join(dir, x.name));
  }
};
for (const dir of ['app', 'packages', 'scripts']) await walk(path.join(root, dir));

for (const filename of files) {
  const rel = path.relative(root, filename);
  const src = await readFile(filename, 'utf8');
  const f = ts.createSourceFile(rel, src, ts.ScriptTarget.Latest, true, rel.endsWith('tsx') ? ts.ScriptKind.TSX : ts.ScriptKind.TS);
  const visit = (n) => {
    const line = f.getLineAndCharacterOfPosition(n.getStart(f)).line + 1;
    if (ts.isVariableDeclarationList(n) && (n.flags & ts.NodeFlags.Let)) {
      for (const v of n.declarations) results.lets.push({ f: rel, line, name: v.name.getText(f), value: (v.initializer?.getText(f) || '').slice(0,110) });
    }
    if (ts.isCallExpression(n) && n.expression.getText(f) === '__') {
      const first=n.arguments[0];
      results.translated.push({ f:rel,line, literal: Boolean(first && (ts.isStringLiteral(first)||ts.isNoSubstitutionTemplateLiteral(first))), arg: (first?.getText(f)||'').slice(0,130) });
    }
    if (ts.isVariableDeclaration(n) && /COPY|MESSAGES|LABELS|CONTENT|TEXT|ERROR/i.test(n.name.getText(f)) &&
        n.initializer && n.initializer.getText(f).includes("'")) {
      results.copy.push({f:rel,line,name:n.name.getText(f),body:n.initializer.getText(f).slice(0,180)});
    }
    ts.forEachChild(n, visit);
  };
  visit(f);
}
console.log('GH_COPY_LET_INVENTORY_START');
for (const [name,arr] of Object.entries(results)) {
 const counts = {};
 for (const x of arr) { const p=x.f; counts[p]=(counts[p]||0)+1; }
 console.log('GH_TOTAL_'+name.toUpperCase()+' '+arr.length+' files='+Object.keys(counts).length);
 console.log('GH_BY_FILE_'+name.toUpperCase()+' '+JSON.stringify(Object.entries(counts).sort((a,b)=>b[1]-a[1]).slice(0,100)));
 if (name==='lets') console.log('GH_LETS '+JSON.stringify(arr.slice(0,300)));
 if (name==='translated') console.log('GH_NONLITERAL '+JSON.stringify(arr.filter(x=>!x.literal).slice(0,250)));
 if (name==='copy') console.log('GH_COPY_OBJECTS '+JSON.stringify(arr.slice(0,100)));
}
console.log('GH_COPY_LET_INVENTORY_END');

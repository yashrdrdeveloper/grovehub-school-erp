const ts = require('typescript');

/** Convert the existing lodash-es ESM JS to CommonJS in Jest without mocking pure code. */
module.exports = {
  process(source, filename) {
    return {
      code: ts.transpileModule(source, {
        fileName: filename,
        compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2020, allowJs: true },
      }).outputText,
    };
  },
};

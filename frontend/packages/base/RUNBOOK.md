# Base

Pure technical functions only. Use `@Base/utilities/isEqual` and other public utility entrypoints. Dependencies must not point to business, components, shell, domain, or application code.

## Pure message resolution

`@Base/i18n` owns the pure, deterministic English source-message formatter `__` and typed `MessageParameters`. It accepts static message templates, replaces named `{{value}}` tokens (including zero, false and repeated occurrences) and throws when required values are missing. It does not evaluate HTML, access browser APIs, translate to another language or initialize React.

Only the shell's `useInitI18n` initializer and test bootstrap **import** this formatter. Feature components call the typed global `__()` after ShellContainer is ready, with no import. Base never writes `window.__`. Pure interpolation tests live in `i18n/__tests__/__.test.ts`.

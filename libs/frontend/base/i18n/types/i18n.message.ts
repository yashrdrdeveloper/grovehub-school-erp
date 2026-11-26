/** Values inserted as literal text by the shared source-message resolver. */
export type MessageParameter = string | number | boolean;

/** Read-only named interpolation values; no HTML, React nodes or nested objects. */
export type MessageParameters = Readonly<Record<string, MessageParameter>>;

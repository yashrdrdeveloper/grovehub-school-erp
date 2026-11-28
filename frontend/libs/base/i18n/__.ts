import type { MessageParameters } from './types/i18n.message';

/**
 * English source-message resolver.
 *
 * This is a pure interpolation foundation, not an active locale translator.
 * React safely escapes the returned string when rendering it as text.
 */
export const __ = (message: string, parameters?: MessageParameters): string =>
  message.replace(/\{\{\s*([A-Za-z_][A-Za-z0-9_]*)\s*\}\}/g, (_match, name: string) => {
    if (!parameters || !Object.prototype.hasOwnProperty.call(parameters, name)) {
      throw new Error('Missing message parameter: ' + name);
    }

    return String(parameters[name]);
  });

import type { MessageParameters } from '@Base/i18n';

/** One shell-installed browser UI message function; no feature-level imports. */
declare global {
  function __(message: string, parameters?: MessageParameters): string;

  interface Window {
    __: typeof __;
  }
}

export {};

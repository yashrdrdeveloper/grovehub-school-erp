import { __ as resolveMessage } from '@Base/i18n';

/**
 * Pre-install the single global message formatter before dynamic application imports
 * evaluate module-level translated COPY objects.
 * The ShellContainer hook invokes the same idempotent installer in its effect.
 */
export const installShellI18n = (): void => {
  if (typeof window === 'undefined') return;
  window.__ = resolveMessage;
  if (window.__ !== resolveMessage) {
    throw new Error('Unable to initialize school workspace messages');
  }
};

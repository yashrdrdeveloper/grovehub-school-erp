import { __ as resolveMessage } from '@Base/i18n';

export const installShellI18n = (): void => {
  window.__ = resolveMessage;
  if (window.__ !== resolveMessage) throw new Error('Unable to initialize operator workspace messages');
};

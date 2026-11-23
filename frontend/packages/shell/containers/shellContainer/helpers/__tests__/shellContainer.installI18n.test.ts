import { __ as resolveMessage } from '@Base/i18n';

import { installShellI18n } from '../shellContainer.installI18n';

describe('installShellI18n', () => {
  afterEach(() => {
    Object.defineProperty(window, '__', {
      configurable: true,
      writable: true,
      value: resolveMessage,
    });
  });

  it('installs the exact pure formatter and is idempotent', () => {
    Reflect.deleteProperty(window, '__');

    installShellI18n();
    installShellI18n();

    expect(window.__).toBe(resolveMessage);
    expect(__('Hello {{name}}', { name: 'Parent' })).toBe('Hello Parent');
  });

  it('rejects a locked window.__ registration', () => {
    Object.defineProperty(window, '__', {
      configurable: true,
      writable: false,
      value: undefined,
    });

    expect(installShellI18n).toThrow();
  });
});

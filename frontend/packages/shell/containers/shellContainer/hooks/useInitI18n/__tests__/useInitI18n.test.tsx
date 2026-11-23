import { StrictMode } from 'react';
import { renderHook } from '@testing-library/react';

import { __ as resolveMessage } from '@Base/i18n';

import { InitI18nStatus } from '../constants/useInitI18n.status';
import { useInitI18n } from '../useInitI18n';

describe('useInitI18n', () => {
  afterEach(() => {
    Object.defineProperty(window, '__', {
      configurable: true,
      writable: true,
      value: resolveMessage,
    });
  });

  it('installs the actual resolver before marking routes ready', () => {
    Reflect.deleteProperty(window, '__');

    const { result } = renderHook(useInitI18n);

    expect(result.current).toBe(InitI18nStatus.READY);
    expect(window.__).toBe(resolveMessage);
    expect(__('Hello {{name}}', { name: 'Parent' })).toBe('Hello Parent');
  });

  it('is idempotent when StrictMode replays effects', () => {
    Reflect.deleteProperty(window, '__');

    const { result, rerender } = renderHook(useInitI18n, { wrapper: StrictMode });
    expect(result.current).toBe(InitI18nStatus.READY);
    expect(window.__).toBe(resolveMessage);

    rerender();
    expect(window.__).toBe(resolveMessage);
    expect(__('Check your email')).toBe('Check your email');
  });

  it('returns FAILED if the window bridge cannot be installed', () => {
    Object.defineProperty(window, '__', {
      configurable: true,
      writable: false,
      value: undefined,
    });

    const { result } = renderHook(useInitI18n);

    expect(result.current).toBe(InitI18nStatus.FAILED);
  });
});

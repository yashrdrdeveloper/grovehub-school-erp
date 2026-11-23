import { StrictMode } from 'react';
import { renderHook } from '@testing-library/react';

import { __ as resolveMessage } from '@Base/i18n';

import { useInitI18n } from '../useInitI18n';

describe('useInitI18n', () => {
  afterEach(() => {
    window.__ = resolveMessage;
  });

  it('installs the actual resolver before marking routes ready', () => {
    Reflect.deleteProperty(window, '__');

    const { result } = renderHook(useInitI18n);

    expect(result.current).toBe(true);
    expect(window.__).toBe(resolveMessage);
    expect(__('Hello {{name}}', { name: 'Parent' })).toBe('Hello Parent');
  });

  it('is idempotent when StrictMode replays effects', () => {
    Reflect.deleteProperty(window, '__');

    const { result, rerender } = renderHook(useInitI18n, { wrapper: StrictMode });
    expect(result.current).toBe(true);
    expect(window.__).toBe(resolveMessage);

    rerender();
    expect(window.__).toBe(resolveMessage);
    expect(__('Check your email')).toBe('Check your email');
  });
});

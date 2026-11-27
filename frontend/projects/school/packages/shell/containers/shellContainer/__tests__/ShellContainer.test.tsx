import type { ReactNode } from 'react';
import { StrictMode } from 'react';
import { render, screen } from '@testing-library/react';

import { __ as resolveMessage } from '@Base/i18n';

import { ShellContainer } from '../ShellContainer';

const mockWasI18nReady = jest.fn();

jest.mock('@Domain/auth', () => ({
  AuthProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
}));

jest.mock('@Components/router', () => ({
  BrowserRouter: ({ children }: { children: ReactNode }) => <>{children}</>,
}));

jest.mock('@Components/components/ErrorBoundary', () => ({
  ErrorBoundary: ({ children }: { children: ReactNode }) => <>{children}</>,
}));

jest.mock('../../../components/appRoutes', () => ({
  AppRoutes: () => {
    mockWasI18nReady(typeof window.__ === 'function');
    return <p>{__('Shell ready')}</p>;
  },
}));

describe('ShellContainer', () => {
  afterEach(() => {
    Object.defineProperty(window, '__', {
      configurable: true,
      writable: true,
      value: resolveMessage,
    });
    mockWasI18nReady.mockClear();
  });

  it('never mounts public or protected routes before __ is registered', () => {
    Reflect.deleteProperty(window, '__');

    render(<ShellContainer />);

    expect(screen.getByText('Shell ready')).toBeTruthy();
    expect(mockWasI18nReady).toHaveBeenCalledWith(true);
    expect(mockWasI18nReady.mock.calls.every(([ready]) => ready)).toBe(true);
    expect(screen.queryByTestId('i18n-boot-status')).toBeNull();
  });

  it('blocks route mounting and shows accessible failure copy when global installation is denied', () => {
    Object.defineProperty(window, '__', {
      configurable: true,
      writable: false,
      value: undefined,
    });

    render(<ShellContainer />);

    expect(screen.getByRole('status').textContent).toContain('Unable to load RDR School.');
    expect(mockWasI18nReady).not.toHaveBeenCalled();
  });

  it('remains safe under StrictMode initialization replay', () => {
    Reflect.deleteProperty(window, '__');

    render(
      <StrictMode>
        <ShellContainer />
      </StrictMode>,
    );

    expect(screen.getByText('Shell ready')).toBeTruthy();
    expect(mockWasI18nReady.mock.calls.every(([ready]) => ready)).toBe(true);
    expect(__('Welcome {{name}}', { name: 'Parent' })).toBe('Welcome Parent');
  });
});

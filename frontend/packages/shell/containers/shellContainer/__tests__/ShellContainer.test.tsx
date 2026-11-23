import type { ReactNode } from 'react';
import { StrictMode } from 'react';
import { render, screen } from '@testing-library/react';
import { renderToString } from 'react-dom/server';

import { ShellContainer } from '../ShellContainer';

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
  AppRoutes: () => <p>{__('Shell ready')}</p>,
}));

describe('ShellContainer', () => {
  it('server-renders only untranslated boot status before the client effect can run', () => {
    const html = renderToString(<ShellContainer />);

    expect(html).toContain('Loading school workspace');
    expect(html).not.toContain('Shell ready');
  });

  it('mounts application routes only after the global function is installed', () => {
    Reflect.deleteProperty(window, '__');

    render(
      <StrictMode>
        <ShellContainer />
      </StrictMode>,
    );

    expect(screen.getByText('Shell ready')).toBeTruthy();
    expect(__('Welcome {{name}}', { name: 'Parent' })).toBe('Welcome Parent');
    expect(screen.queryByTestId('i18n-boot-status')).toBeNull();
  });
});

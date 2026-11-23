import { ErrorBoundary } from '@Components/components/ErrorBoundary';
import { QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from '@Components/router';

import { AuthProvider } from '@Domain/auth';

import { AppRoutes } from '../../components/appRoutes';
import { I18nBootstrapFallback } from './components/i18nBootstrapFallback';
import { useInitI18n, InitI18nStatus } from './hooks/useInitI18n';
import { createShellQueryClient } from './helpers/shellContainer.createQueryClient';
import type { ShellContainerComponent } from './types/shellContainer';

const queryClient = createShellQueryClient();

/** Provider composition only; route and workspace frame units own their markup. */
export const ShellContainer: ShellContainerComponent = () => {
  const i18nStatus = useInitI18n();

  if (i18nStatus !== InitI18nStatus.READY) {
    return <I18nBootstrapFallback status={i18nStatus} />;
  }

  return (
    <QueryClientProvider client={queryClient}>
      <ErrorBoundary>
        <BrowserRouter>
          <AuthProvider>
            <AppRoutes />
          </AuthProvider>
        </BrowserRouter>
      </ErrorBoundary>
    </QueryClientProvider>
  );
};

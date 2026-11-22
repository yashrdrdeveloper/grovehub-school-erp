import { ErrorBoundary } from '@Components/components/ErrorBoundary';
import { QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from '@Components/router';

import { AuthProvider } from '@Domain/auth';

import { AppRoutes } from '../../components/appRoutes';
import { createShellQueryClient } from './helpers/shellContainer.createQueryClient';
import type { ShellContainerComponent } from './types/shellContainer';

const queryClient = createShellQueryClient();

/** Provider composition only; route and workspace frame units own their markup. */
export const ShellContainer: ShellContainerComponent = () => (
  <QueryClientProvider client={queryClient}>
    <ErrorBoundary>
      <BrowserRouter>
      <AuthProvider>
        <AppRoutes/>
      </AuthProvider>
      </BrowserRouter>
    </ErrorBoundary>
  </QueryClientProvider>
);

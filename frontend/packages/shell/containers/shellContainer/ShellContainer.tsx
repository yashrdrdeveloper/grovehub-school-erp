import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';

import { AuthProvider, RecoveryFlowBoundary } from '@Domain/auth';

import { AppRoutes } from '../../components/appRoutes';
import type { ShellContainerComponent } from './types/shellContainer';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 30_000, retry: 1, refetchOnWindowFocus: false },
  },
});

/** Provider composition only; route and workspace frame units own their markup. */
export const ShellContainer: ShellContainerComponent = () => (
  <QueryClientProvider client={queryClient}>
    <BrowserRouter>
      <AuthProvider>
        <RecoveryFlowBoundary/>
        <AppRoutes/>
      </AuthProvider>
    </BrowserRouter>
  </QueryClientProvider>
);

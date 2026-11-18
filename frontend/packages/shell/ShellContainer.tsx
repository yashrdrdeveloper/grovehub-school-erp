import { useMemo } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';

import { AuthProvider, PublicOnly, RequireSession, useAuth } from '@Domain/auth';

import { Shell } from './Shell';
import { publicRoutes, workspaceRoutes } from './registry';

const client = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 30_000, retry: 1, refetchOnWindowFocus: false },
  },
});

const WorkspaceFrame = () => {
  const { account, signOut } = useAuth();
  const schoolLabel = useMemo(() => {
    const hostname = window.location.hostname;
    if (hostname.endsWith('.prod.groovehub.com')) {
      const first = hostname.split('.')[0] ?? 'school';
      return first.charAt(0).toUpperCase() + first.slice(1);
    }
    return 'Your school';
  }, []);
  return (
    <Shell
      accountLabel={account?.email ?? 'Account'}
      schoolLabel={schoolLabel}
      onLogout={() => void signOut()}
    />
  );
};

const AppRoutes = () => (
  <Routes>
    <Route element={<PublicOnly />}>
      {publicRoutes.map(({ path, element }) => <Route key={path} path={path} element={element} />)}
    </Route>
    <Route element={<RequireSession />}>
      <Route element={<WorkspaceFrame />}>
        {workspaceRoutes.map(({ path, element }) => <Route key={path} path={path} element={element} />)}
      </Route>
    </Route>
    <Route path="/" element={<Navigate to="/dashboard" replace />} />
    <Route path="*" element={<Navigate to="/dashboard" replace />} />
  </Routes>
);

export const ShellContainer = () => (
  <QueryClientProvider client={client}>
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  </QueryClientProvider>
);

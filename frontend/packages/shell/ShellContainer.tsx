import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';

import { AuthProvider, PublicOnly, RequireSession, useAuth } from '@Domain/auth';
import { LandingPage } from '@LandingPage';

import { Shell } from './Shell';
import { publicRoutes, workspaceRoutes } from './registry';

const client = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 30_000, retry: 1, refetchOnWindowFocus: false },
  },
});

const WorkspaceFrame = () => {
  const { account, signOut } = useAuth();
  const schoolLabel = 'RDR School';
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
    {['/', '/about', '/academics', '/campus', '/admissions', '/gallery', '/mandatory-public-disclosure', '/resources', '/contact'].map(path => (
      <Route key={path} path={path} element={<LandingPage />} />
    ))}
    <Route element={<PublicOnly />}>
      {publicRoutes.map(({ path, element }) => <Route key={path} path={path} element={element} />)}
    </Route>
    <Route element={<RequireSession />}>
      <Route element={<WorkspaceFrame />}>
        {workspaceRoutes.map(({ path, element }) => <Route key={path} path={path} element={element} />)}
      </Route>
    </Route>
    <Route path="*" element={<Navigate to="/" replace />} />
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

import { Box } from '@Components/components/Box';
import { useEffect } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { AUTH_ROUTE } from '@Business/auth/routes';
import { APP_ROUTE } from '@Business/app/routes';

import { useAuth } from '../providers/authProvider';

export const RequireSession = () => {
  const { account, loading, initialized, reload } = useAuth();
  const location = useLocation();

  // Session bootstrapping is required only by protected workspace routes.
  // Public login, set-password and reset-password pages never trigger it.
  useEffect(() => {
    if (!initialized) void reload();
  }, [initialized, reload]);

  if (loading || !initialized) return <Box as="div" className="screen-loading">Opening your workspace…</Box>;
  if (!account) return <Navigate to={AUTH_ROUTE.LOGIN} state={{ from: location.pathname }} replace />;
  return <Outlet />;
};

export const PublicOnly = () => {
  const { account, loading } = useAuth();
  if (loading) return <Box as="div" className="screen-loading">Checking your session…</Box>;
  if (account) return <Navigate to={APP_ROUTE.DASHBOARD} replace />;
  return <Outlet />;
};

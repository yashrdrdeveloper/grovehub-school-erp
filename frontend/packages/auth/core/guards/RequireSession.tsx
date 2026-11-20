import { useEffect } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';

import { useAuth } from '../providers/AuthProvider';

export const RequireSession = () => {
  const { account, loading, initialized, reload } = useAuth();
  const location = useLocation();

  // Session bootstrapping is required only by protected workspace routes.
  // Public login, set-password and reset-password pages never trigger it.
  useEffect(() => {
    if (!initialized) void reload();
  }, [initialized, reload]);

  if (loading || !initialized) return <div className="screen-loading">Opening your workspace…</div>;
  if (!account) return <Navigate to="/login" state={{ from: location.pathname }} replace />;
  return <Outlet />;
};

export const PublicOnly = () => {
  const { account, loading } = useAuth();
  if (loading) return <div className="screen-loading">Checking your session…</div>;
  return account ? <Navigate to="/dashboard" replace /> : <Outlet />;
};

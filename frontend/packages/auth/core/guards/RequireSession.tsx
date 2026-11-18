import { Navigate, Outlet, useLocation } from 'react-router-dom';

import { useAuth } from '../providers/AuthProvider';

export const RequireSession = () => {
  const { account, loading } = useAuth();
  const location = useLocation();

  if (loading) return <div className="screen-loading">Opening your workspace…</div>;
  if (!account) return <Navigate to="/login" state={{ from: location.pathname }} replace />;
  return <Outlet />;
};

export const PublicOnly = () => {
  const { account, loading } = useAuth();
  if (loading) return <div className="screen-loading">Checking your session…</div>;
  return account ? <Navigate to="/dashboard" replace /> : <Outlet />;
};

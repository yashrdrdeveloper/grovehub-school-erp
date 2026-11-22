import { useEffect } from 'react';

import { APP_ROUTE } from '@Business/app/routes';
import { AUTH_ROUTE } from '@Business/auth/routes';
import { Box, BOX_TAG } from '@Components/components/Box';
import { Navigate, Outlet, useLocation } from '@Components/router';

import { useAuth } from '../providers/authProvider';
import styles from './RequireSession.module.css';
import { REQUIRE_SESSION_TEST_IDS } from './constants/requireSession.testIds';

/** Session bootstrapping runs only inside the protected workspace boundary. */
export const RequireSession = () => {
  const { account, loading, initialized, reload } = useAuth();
  const location = useLocation();

  useEffect(() => {
    if (!initialized) {
      void reload();
    }
  }, [initialized, reload]);

  if (loading || !initialized) {
    return (
      <Box
        as={BOX_TAG.DIV}
        className={styles.loading}
        dataTestId={REQUIRE_SESSION_TEST_IDS.PRIVATE_LOADING}
        role="status"
        aria-live="polite"
      >
        Opening your workspace…
      </Box>
    );
  }

  if (!account) {
    return <Navigate to={AUTH_ROUTE.LOGIN} state={{ from: location.pathname }} replace />;
  }

  return <Outlet />;
};

/** Public-only pages never initiate an authenticated session or refresh request. */
export const PublicOnly = () => {
  const { account, loading } = useAuth();

  if (loading) {
    return (
      <Box
        as={BOX_TAG.DIV}
        className={styles.loading}
        dataTestId={REQUIRE_SESSION_TEST_IDS.PUBLIC_LOADING}
        role="status"
        aria-live="polite"
      >
        Checking your session…
      </Box>
    );
  }

  if (account) {
    return <Navigate to={APP_ROUTE.DASHBOARD} replace />;
  }

  return <Outlet />;
};

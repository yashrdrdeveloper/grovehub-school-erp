import { _map } from '@Base/utilities/map';
import { _filter } from '@Base/utilities/filter';
import type { ReactElement } from 'react';

import { Navigate } from '@Components/router';
import { AUTH_ROUTE_CONFIG, AUTH_ROUTE_IDS, getLoginRoute } from '@Business/auth/routes';
import { APP_ROUTE_CONFIG, APP_ROUTE_IDS, RouteAccess } from '@Business/app/routes';

import { ForgotPasswordPage, LoginV2Page, RecoveryPage } from '@Domain/auth';
import { DashboardPage } from '@Domain/dashboard';

/**
 * Route metadata is owned exclusively by Business.
 * Shell maps IDs to views without defining or duplicating browser paths.
 */
type AuthRouteKey = (typeof AUTH_ROUTE_IDS)[keyof typeof AUTH_ROUTE_IDS];
type AppRouteKey = (typeof APP_ROUTE_IDS)[keyof typeof APP_ROUTE_IDS];

const authViews: Record<AuthRouteKey, ReactElement> = {
  [AUTH_ROUTE_IDS.LOGIN]: <LoginV2Page />,
  [AUTH_ROUTE_IDS.LOGIN_LEGACY]: <Navigate to={getLoginRoute().path} replace />,
  [AUTH_ROUTE_IDS.FORGOT_PASSWORD]: <ForgotPasswordPage />,
  [AUTH_ROUTE_IDS.RESET_PASSWORD]: <RecoveryPage purpose="recovery" />,
  [AUTH_ROUTE_IDS.SET_PASSWORD]: <RecoveryPage purpose="invitation" />,
  [AUTH_ROUTE_IDS.SETUP_LEGACY]: <RecoveryPage purpose="invitation" />,
};

const workspaceViews: Record<AppRouteKey, ReactElement | null> = {
  [APP_ROUTE_IDS.HOME]: null,
  [APP_ROUTE_IDS.DASHBOARD]: <DashboardPage />,
};

export const publicRoutes = _map(AUTH_ROUTE_CONFIG, (route) => ({
  ...route,
  element: authViews[route.id],
}));

export const workspaceRoutes = _map(
  _filter(APP_ROUTE_CONFIG, (route) => route.type === RouteAccess.PROTECTED),
  (route) => ({
    ...route,
    element: workspaceViews[route.id],
  }),
);

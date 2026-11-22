import { Navigate } from 'react-router-dom';

import { AUTH_ROUTE } from '@Business/auth/routes';
import { APP_ROUTE } from '@Business/app/routes';

import {
  ForgotPasswordPage,
  LoginV2Page,
  RecoveryPage,
} from '@Domain/auth';
import { DashboardPage } from '@Domain/dashboard';

// Phase 0.2: identity entrypoints and the dashboard shell only.
// Business modules are deliberately not routable until implemented.
export const publicRoutes = [
  { path: AUTH_ROUTE.LOGIN, element: <LoginV2Page /> },
  // Preserve old bookmarks without maintaining a second login screen.
  { path: AUTH_ROUTE.LOGIN_LEGACY, element: <Navigate to={AUTH_ROUTE.LOGIN} replace /> },
  { path: AUTH_ROUTE.FORGOT_PASSWORD, element: <ForgotPasswordPage /> },
  { path: AUTH_ROUTE.RESET_PASSWORD, element: <RecoveryPage purpose="recovery" /> },
  { path: AUTH_ROUTE.SET_PASSWORD, element: <RecoveryPage purpose="invitation" /> },
  // Preserve links issued by the existing Identity invitation service.
  { path: AUTH_ROUTE.SETUP_LEGACY, element: <RecoveryPage purpose="invitation" /> },
] as const;

export const workspaceRoutes = [
  { path: APP_ROUTE.DASHBOARD, element: <DashboardPage /> },
] as const;

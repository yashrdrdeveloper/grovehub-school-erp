import { Navigate } from 'react-router-dom';

import {
  ForgotPasswordPage,
  LoginV2Page,
  RecoveryPage,
} from '@Domain/auth';
import { DashboardPage } from '@Domain/dashboard';

// Phase 0.2: identity entrypoints and the dashboard shell only.
// Business modules are deliberately not routable until implemented.
export const publicRoutes = [
  { path: '/login', element: <LoginV2Page /> },
  // Preserve old bookmarks without maintaining a second login screen.
  { path: '/loginv2', element: <Navigate to="/login" replace /> },
  { path: '/forgot-password', element: <ForgotPasswordPage /> },
  { path: '/reset-password', element: <RecoveryPage purpose="recovery" /> },
  { path: '/set-password', element: <RecoveryPage purpose="invitation" /> },
  // Preserve links issued by the existing Identity invitation service.
  { path: '/setup', element: <RecoveryPage purpose="invitation" /> },
] as const;

export const workspaceRoutes = [
  { path: '/dashboard', element: <DashboardPage /> },
] as const;

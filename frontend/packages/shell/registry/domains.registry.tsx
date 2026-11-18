import {
  ForgotPasswordPage,
  LoginPage,
  RecoveryPage,
} from '@Domain/auth';
import { DashboardPage } from '@Domain/dashboard';

// Phase 0.2: identity entrypoints and the dashboard shell only.
// Business modules are deliberately not routable until implemented.
export const publicRoutes = [
  { path: '/login', element: <LoginPage /> },
  { path: '/forgot-password', element: <ForgotPasswordPage /> },
  { path: '/reset-password', element: <RecoveryPage purpose="recovery" /> },
  { path: '/set-password', element: <RecoveryPage purpose="invitation" /> },
  // Preserve links issued by the existing Identity invitation service.
  { path: '/setup', element: <RecoveryPage purpose="invitation" /> },
] as const;

export const workspaceRoutes = [
  { path: '/dashboard', element: <DashboardPage /> },
] as const;

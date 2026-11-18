import {
  ForgotPasswordPage,
  LoginPage,
  RecoveryPage,
} from '@Domain/auth';
import { AccountsPage } from '@Domain/accounts';
import { DashboardPage } from '@Domain/dashboard';
import { IntegrationsPage } from '@Domain/integrations';
import { SecurityPage } from '@Domain/settings';

// The shell imports only domain public entrypoints.
// Domain modules own their screens, transport and behavior.
export const publicRoutes = [
  { path: '/login', element: <LoginPage /> },
  { path: '/forgot-password', element: <ForgotPasswordPage /> },
  { path: '/reset-password', element: <RecoveryPage purpose="recovery" /> },
  { path: '/setup', element: <RecoveryPage purpose="invitation" /> },
] as const;

export const workspaceRoutes = [
  { path: '/dashboard', element: <DashboardPage /> },
  { path: '/accounts', element: <AccountsPage /> },
  { path: '/integrations', element: <IntegrationsPage /> },
  { path: '/settings/security', element: <SecurityPage /> },
] as const;

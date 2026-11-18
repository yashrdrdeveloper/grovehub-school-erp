import {
  ForgotPasswordPage,
  LoginPage,
  RecoveryPage,
} from '@Domain/auth';
import { AccountDetailsPage, AccountsPage } from '@Domain/accounts';
import { DashboardPage } from '@Domain/dashboard';
import { PeoplePage } from '@Domain/people';
import { ConnectionDetailsPage, IntegrationsPage } from '@Domain/integrations';
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
  { path: '/people', element: <PeoplePage /> },
  { path: '/accounts', element: <AccountsPage /> },
  { path: '/accounts/:accountId', element: <AccountDetailsPage /> },
  { path: '/integrations', element: <IntegrationsPage /> },
  { path: '/integrations/:connectionId', element: <ConnectionDetailsPage /> },
  { path: '/settings/security', element: <SecurityPage /> },
] as const;

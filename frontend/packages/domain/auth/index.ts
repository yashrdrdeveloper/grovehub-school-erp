// Intentional compatibility facade: pages live under packages/auth/pages.
// Shell consumes only this domain root, never feature internals.
export {
  AuthProvider,
  useAuth,
  RequireSession,
  PublicOnly,
  LoginPage,
  ForgotPasswordPage,
  RecoveryPage,
} from '@Auth';

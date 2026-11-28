/** Integration transport endpoints are never browser route paths. */
export const INTEGRATION_ENDPOINT = Object.freeze({
  PUBLIC_PROVIDERS: '/integrations/public/providers',
  PROTECTED_CONNECTIONS: '/integrations/protected/connections',
} as const);

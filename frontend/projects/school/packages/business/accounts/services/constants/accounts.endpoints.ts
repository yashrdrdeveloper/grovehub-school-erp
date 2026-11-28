/** These are Identity transport endpoints, not browser navigation routes. */
export const ACCOUNTS_ENDPOINT = Object.freeze({
  PROTECTED: '/identity/protected/accounts',
  PUBLIC_SETUP_STATUS: '/identity/public/accounts/setup/status',
  PUBLIC_SETUP: '/identity/public/accounts/setup',
} as const);

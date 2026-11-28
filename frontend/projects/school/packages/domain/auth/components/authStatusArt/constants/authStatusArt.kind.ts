/** Closed set of reusable authentication status illustrations. */
export const AUTH_STATUS_KIND = Object.freeze({
  EMAIL: 'email',
  SUCCESS: 'success',
  EXPIRED: 'expired',
  MISSING: 'missing',
  CHECKING: 'checking',
  UNAVAILABLE: 'unavailable',
  ATTEMPTS_EXHAUSTED: 'attempts-exhausted',
} as const);

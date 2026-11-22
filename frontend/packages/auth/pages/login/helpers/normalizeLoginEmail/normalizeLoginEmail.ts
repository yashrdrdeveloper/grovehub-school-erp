import type { NormalizeLoginEmail } from './types/normalizeLoginEmail';

/** Deliberately trim only: preserve the original email casing for identity matching. */
export const normalizeLoginEmail: NormalizeLoginEmail = (email) => email.trim();

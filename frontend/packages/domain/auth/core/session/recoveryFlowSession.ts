import { AUTH_ROUTE } from '@Business/auth/routes';

/** A short-lived single-tab continuation for one password-link flow.
 * Raw bearer proofs must never enter web storage, URLs beyond first use, or logs.
 * Only a non-secret marker for the HttpOnly grant cookie may survive a refresh.
 */
export type RecoveryPurpose = 'invitation' | 'recovery';
export type RecoveryFlow =
  | { purpose: RecoveryPurpose; stage: 'proof'; proof: string; expiresAt: number }
  | { purpose: RecoveryPurpose; stage: 'grant'; expiresAt: number };

const KEY = 'groovehub.identity.password-flow.v1';
const MAX_LIFETIME_MS = 10 * 60 * 1000;

// A failed exchange may retry while this JS context survives; page reload
// discards the bearer proof. This intentionally favors secrecy over replay.
const sessionState = {
  pendingProofState: null as Extract<RecoveryFlow, { stage: 'proof' }> | null,
};

export const recoveryPurposeForPath = (path: string): RecoveryPurpose | null => {
  if (path === AUTH_ROUTE.SET_PASSWORD || path === AUTH_ROUTE.SETUP_LEGACY) return 'invitation';
  if (path === AUTH_ROUTE.RESET_PASSWORD) return 'recovery';
  return null;
};

export const clearRecoveryFlow = (): void => {
  sessionState.pendingProofState = null;
  try {
    window.sessionStorage.removeItem(KEY);
  } catch {
    /* Storage may be disabled. */
  }
};

const active = (expiresAt: number): boolean =>
  Number.isFinite(expiresAt) && expiresAt > Date.now() && expiresAt <= Date.now() + MAX_LIFETIME_MS;

export const readRecoveryFlow = (purpose: RecoveryPurpose): RecoveryFlow | null => {
  if (sessionState.pendingProofState) {
    if (
      sessionState.pendingProofState.purpose !== purpose ||
      !active(sessionState.pendingProofState.expiresAt)
    ) {
      clearRecoveryFlow();
      return null;
    }
    return sessionState.pendingProofState;
  }
  try {
    const raw = window.sessionStorage.getItem(KEY);
    if (!raw) return null;
    const item: unknown = JSON.parse(raw);
    if (!item || typeof item !== 'object') {
      clearRecoveryFlow();
      return null;
    }
    const marker = item as Partial<RecoveryFlow>;
    // Legacy stored proof records are invalidated and never replayed.
    if (
      marker.purpose !== purpose ||
      marker.stage !== 'grant' ||
      typeof marker.expiresAt !== 'number' ||
      !active(marker.expiresAt)
    ) {
      clearRecoveryFlow();
      return null;
    }
    return { purpose, stage: 'grant', expiresAt: marker.expiresAt };
  } catch {
    clearRecoveryFlow();
    return null;
  }
};

export const saveRecoveryProof = (purpose: RecoveryPurpose, proof: string): void => {
  clearRecoveryFlow();
  sessionState.pendingProofState = {
    purpose,
    stage: 'proof',
    proof,
    expiresAt: Date.now() + MAX_LIFETIME_MS,
  };
};

export const saveRecoveryGrant = (purpose: RecoveryPurpose, expiresIn: number): void => {
  sessionState.pendingProofState = null;
  const lifetime = Math.max(0, Math.min(MAX_LIFETIME_MS, Math.floor(expiresIn * 1000)));
  if (lifetime === 0) {
    clearRecoveryFlow();
    return;
  }
  try {
    window.sessionStorage.setItem(
      KEY,
      JSON.stringify({
        purpose,
        stage: 'grant',
        expiresAt: Date.now() + lifetime,
      }),
    );
  } catch {
    /* Cookie-backed flow remains active in memory. */
  }
};

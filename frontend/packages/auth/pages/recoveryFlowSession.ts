/** A short-lived, tab-local continuation for one password-link flow.
 * Raw proof is removed once exchanged; only the HttpOnly cookie can authorize setup.
 */
export type RecoveryPurpose = 'invitation' | 'recovery';
export type RecoveryFlow =
  | { purpose: RecoveryPurpose; stage: 'proof'; proof: string; expiresAt: number }
  | { purpose: RecoveryPurpose; stage: 'grant'; expiresAt: number };

const KEY = 'groovehub.identity.password-flow.v1';
const MAX_LIFETIME_MS = 10 * 60 * 1000;

export const recoveryPurposeForPath = (path: string): RecoveryPurpose | null => {
  if (path === '/set-password' || path === '/setup') return 'invitation';
  if (path === '/reset-password') return 'recovery';
  return null;
};

export const clearRecoveryFlow = (): void => {
  try { window.sessionStorage.removeItem(KEY); } catch { /* Private mode can disable storage. */ }
};

export const readRecoveryFlow = (purpose: RecoveryPurpose): RecoveryFlow | null => {
  try {
    const raw = window.sessionStorage.getItem(KEY);
    if (!raw) return null;
    const item: unknown = JSON.parse(raw);
    if (!item || typeof item !== 'object') { clearRecoveryFlow(); return null; }
    const flow = item as Partial<RecoveryFlow>;
    if (flow.purpose !== purpose || !Number.isFinite(flow.expiresAt) || !flow.expiresAt ||
        flow.expiresAt <= Date.now() || flow.expiresAt > Date.now() + MAX_LIFETIME_MS ||
        (flow.stage !== 'grant' && (flow.stage !== 'proof' || typeof flow.proof !== 'string' || !flow.proof.trim()))) {
      clearRecoveryFlow();
      return null;
    }
    return flow as RecoveryFlow;
  } catch { clearRecoveryFlow(); return null; }
};

const save = (flow: RecoveryFlow): void => {
  try { window.sessionStorage.setItem(KEY, JSON.stringify(flow)); } catch { /* In-memory flow still works. */ }
};

export const saveRecoveryProof = (purpose: RecoveryPurpose, proof: string): void => {
  save({ purpose, stage: 'proof', proof, expiresAt: Date.now() + MAX_LIFETIME_MS });
};

export const saveRecoveryGrant = (purpose: RecoveryPurpose, expiresIn: number): void => {
  const lifetime = Math.max(0, Math.min(MAX_LIFETIME_MS, Math.floor(expiresIn * 1000)));
  if (lifetime === 0) { clearRecoveryFlow(); return; }
  save({ purpose, stage: 'grant', expiresAt: Date.now() + lifetime });
};

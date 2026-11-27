import { AlertCircle, Clock3, Link2Off, PartyPopper, ShieldCheck } from '@Components/icons';
import type { LucideIcon } from '@Components/icons';
import type { AuthStatusKind } from '../../../components/authStatusArt';
import type { RecoveryViewState } from '../types';

export interface RecoveryArtworkDescriptor {
  readonly Icon: LucideIcon;
  readonly statusKind: AuthStatusKind;
  readonly visualClass: 'success' | 'checking' | 'problem';
}

/** Exhaustive descriptors replace cascaded render-time branch expressions. */
export const RECOVERY_ARTWORK: Readonly<Record<RecoveryViewState, RecoveryArtworkDescriptor>> = {
  missing: { Icon: Link2Off, statusKind: 'missing', visualClass: 'problem' },
  checking: { Icon: ShieldCheck, statusKind: 'checking', visualClass: 'checking' },
  expired: { Icon: Clock3, statusKind: 'expired', visualClass: 'problem' },
  unavailable: { Icon: Clock3, statusKind: 'unavailable', visualClass: 'problem' },
  ready: { Icon: ShieldCheck, statusKind: 'checking', visualClass: 'checking' },
  saving: { Icon: ShieldCheck, statusKind: 'checking', visualClass: 'checking' },
  success: { Icon: PartyPopper, statusKind: 'success', visualClass: 'success' },
  'attempts-exhausted': {
    Icon: AlertCircle,
    statusKind: 'attempts-exhausted',
    visualClass: 'problem',
  },
};

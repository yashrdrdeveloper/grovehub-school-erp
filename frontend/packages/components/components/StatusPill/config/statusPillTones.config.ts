import type { StatusPillTone } from '../types/statusPill.props';

export const STATUS_PILL_TONES: Readonly<Record<string, StatusPillTone>> = Object.freeze({
  ACTIVE: 'positive',
  PENDING_SETUP: 'notice',
  SUSPENDED: 'warning',
  DEACTIVATED: 'neutral',
});

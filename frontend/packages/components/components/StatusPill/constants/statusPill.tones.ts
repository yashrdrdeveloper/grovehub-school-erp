import { STATUS_PILL_TONE } from './statusPill.tone';
import type { StatusPillTone } from '../types/statusPill.props';

export const STATUS_PILL_TONES: Readonly<Record<string, StatusPillTone>> = Object.freeze({
  ACTIVE: STATUS_PILL_TONE.POSITIVE,
  PENDING_SETUP: STATUS_PILL_TONE.NOTICE,
  SUSPENDED: STATUS_PILL_TONE.WARNING,
  DEACTIVATED: STATUS_PILL_TONE.NEUTRAL,
});

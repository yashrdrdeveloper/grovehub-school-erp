import { STATUS_PILL_TONE } from '../constants/statusPill.tone';
import { STATUS_PILL_TONES } from '../constants/statusPill.tones';
import type { GetStatusPillTone } from './types/statusPill.getTone';

export const getStatusPillTone: GetStatusPillTone = (status) =>
  STATUS_PILL_TONES[status] ?? STATUS_PILL_TONE.NEUTRAL;

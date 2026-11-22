import { STATUS_PILL_TONES } from '../../config/statusPillTones.config';
import type { GetStatusPillTone } from './types/getStatusPillTone';

export const getStatusPillTone: GetStatusPillTone = status =>
  STATUS_PILL_TONES[status] ?? 'neutral';

import { PHONE_VIEWPORT } from '../../../../config';
import type { GetMobileSnapshot } from './types/getMobileSnapshot';

/** The viewport query is isolated from the React hook. */
export const getMobileSnapshot: GetMobileSnapshot = () =>
  typeof window !== 'undefined'
  && typeof window.matchMedia === 'function'
  && window.matchMedia(PHONE_VIEWPORT).matches;

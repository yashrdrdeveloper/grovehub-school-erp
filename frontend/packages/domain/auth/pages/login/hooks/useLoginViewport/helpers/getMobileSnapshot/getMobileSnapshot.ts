import { STYLE_MEDIA_QUERY } from '@Styles/tokens/breakpoints';
import type { GetMobileSnapshot } from './types/getMobileSnapshot';

/** The viewport query is isolated from the React hook. */
export const getMobileSnapshot: GetMobileSnapshot = () =>
  typeof window !== 'undefined'
  && typeof window.matchMedia === 'function'
  && window.matchMedia(STYLE_MEDIA_QUERY.PHONE).matches;

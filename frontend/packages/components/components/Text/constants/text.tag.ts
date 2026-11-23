import { TYPOGRAPHY_TAG } from '../../Typography/constants/typography.tag';

/** Text's semantic tag subset, reusing the Typography contract. */
export const TEXT_TAG = Object.freeze({
  SPAN: TYPOGRAPHY_TAG.SPAN,
  P: TYPOGRAPHY_TAG.P,
  SMALL: TYPOGRAPHY_TAG.SMALL,
  STRONG: TYPOGRAPHY_TAG.STRONG,
} as const);

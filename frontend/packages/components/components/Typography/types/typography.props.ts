import type { TYPOGRAPHY_WEIGHT } from '../constants/typography.weight';
import type { TYPOGRAPHY_SIZE } from '../constants/typography.size';
import type { TYPOGRAPHY_TONE } from '../constants/typography.tone';
import type { TYPOGRAPHY_VARIANT } from '../constants/typography.variant';
import type { TYPOGRAPHY_TAG } from '../constants/typography.tag';
import type { HTMLAttributes, ReactElement } from 'react';
export type TypographyTag = typeof TYPOGRAPHY_TAG[keyof typeof TYPOGRAPHY_TAG];
export type TypographyVariant = typeof TYPOGRAPHY_VARIANT[keyof typeof TYPOGRAPHY_VARIANT];
export type TypographyTone = typeof TYPOGRAPHY_TONE[keyof typeof TYPOGRAPHY_TONE];
export type TypographySize = typeof TYPOGRAPHY_SIZE[keyof typeof TYPOGRAPHY_SIZE];
export type TypographyWeight = typeof TYPOGRAPHY_WEIGHT[keyof typeof TYPOGRAPHY_WEIGHT];
export interface TypographyProps extends HTMLAttributes<HTMLElement> {
  readonly as?: TypographyTag;
  readonly variant?: TypographyVariant;
  readonly tone?: TypographyTone;
  readonly size?: TypographySize;
  readonly weight?: TypographyWeight;
  readonly dataTestId?: string;
}
export type TypographyComponent = (props: TypographyProps) => ReactElement;

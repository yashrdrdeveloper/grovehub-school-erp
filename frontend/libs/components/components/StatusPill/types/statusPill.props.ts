import type { STATUS_PILL_TONE } from '../constants/statusPill.tone';
import type { ReactElement } from 'react';
export interface StatusPillProps {
  readonly value: string;
  readonly dataTestId?: string;
  /** Text and semantic color without a filled or outlined badge. */
  readonly plain?: boolean;
}
export type StatusPillTone = (typeof STATUS_PILL_TONE)[keyof typeof STATUS_PILL_TONE];
export type StatusPillComponent = (props: StatusPillProps) => ReactElement;

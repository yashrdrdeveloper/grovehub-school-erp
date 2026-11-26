import type { BOX_TAG } from '../constants/box.tag';
import type {
  ForwardRefExoticComponent,
  HTMLAttributes,
  PropsWithoutRef,
  RefAttributes,
} from 'react';

export type BoxTag = (typeof BOX_TAG)[keyof typeof BOX_TAG];
export interface BoxProps extends HTMLAttributes<HTMLElement> {
  readonly as?: BoxTag;
  readonly dataTestId?: string;
}
export type BoxComponent = ForwardRefExoticComponent<
  PropsWithoutRef<BoxProps> & RefAttributes<HTMLElement>
>;

import type { CONTAINER_TAG } from '../constants/container.tag';
import type { CONTAINER_WIDTH } from '../constants/container.width';
import type { HTMLAttributes, ReactElement } from 'react';
export type ContainerTag = typeof CONTAINER_TAG[keyof typeof CONTAINER_TAG];
export type ContainerWidth = typeof CONTAINER_WIDTH[keyof typeof CONTAINER_WIDTH];
export interface ContainerProps extends HTMLAttributes<HTMLElement> {
  readonly as?: ContainerTag;
  readonly maxWidth?: ContainerWidth;
  readonly padded?: boolean;
  readonly dataTestId?: string;
}
export type ContainerComponent = (props: ContainerProps) => ReactElement;

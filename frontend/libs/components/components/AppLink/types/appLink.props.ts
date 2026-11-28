import type { ForwardRefExoticComponent, PropsWithoutRef, RefAttributes } from 'react';
import type { LinkProps } from 'react-router-dom';
export interface AppLinkProps extends LinkProps {
  readonly dataTestId?: string;
}
export type AppLinkComponent = ForwardRefExoticComponent<
  PropsWithoutRef<AppLinkProps> & RefAttributes<HTMLAnchorElement>
>;

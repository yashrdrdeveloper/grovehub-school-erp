import type { ForwardRefExoticComponent, PropsWithoutRef, RefAttributes } from 'react';
import type { NavLinkProps } from 'react-router-dom';

export interface AppNavLinkProps extends NavLinkProps {
  readonly dataTestId?: string;
}

export type AppNavLinkComponent = ForwardRefExoticComponent<
  PropsWithoutRef<AppNavLinkProps> & RefAttributes<HTMLAnchorElement>
>;

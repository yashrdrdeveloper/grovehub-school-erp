import type { ReactElement, ReactNode } from 'react';
import type { LoginLayoutVariant } from '../constants/loginLayout.variant';
export interface LoginLayoutProps {
  readonly children: ReactNode;
  readonly title: string;
  readonly description: string;
  readonly mobileDescription?: string;
  readonly minimalMobile?: boolean;
  readonly mobileTitle?: string;
  readonly variant?: LoginLayoutVariant;
  readonly eyebrow?: string;
}
export type LoginLayoutComponent = (props: LoginLayoutProps) => ReactElement;

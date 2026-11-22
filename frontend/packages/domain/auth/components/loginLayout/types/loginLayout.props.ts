import type { ReactElement, ReactNode } from 'react';
export interface LoginLayoutProps {
  readonly children: ReactNode;
  readonly title: string;
  readonly description: string;
  readonly mobileDescription?: string;
  readonly minimalMobile?: boolean;
  readonly mobileTitle?: string;
  readonly variant?: 'login' | 'forgot' | 'recovery';
  readonly eyebrow?: string;
}
export type LoginLayoutComponent = (props: LoginLayoutProps) => ReactElement;

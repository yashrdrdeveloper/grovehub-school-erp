import type { LoginLayoutProps } from '../../types/loginLayout.props';
export type GetLoginLayoutClassName = (
  variant: NonNullable<LoginLayoutProps['variant']>,
  minimalMobile: boolean,
) => string;

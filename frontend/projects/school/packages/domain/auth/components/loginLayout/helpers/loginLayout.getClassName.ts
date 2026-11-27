import type { GetLoginLayoutClassName } from './types/loginLayout.getClassName';

export const getLoginLayoutClassName: GetLoginLayoutClassName = (variant, minimalMobile) => {
  return (
    'rdr-login-v2 rdr-login-v2--' +
    variant +
    ((minimalMobile && ' rdr-login-v2--mobile-minimal') || '')
  );
};

import type { GetLoginLayoutClassName } from './types/loginLayout.getClassName';

export const getLoginLayoutClassName: GetLoginLayoutClassName = (variant, minimalMobile) => {
  let result = 'rdr-login-v2 rdr-login-v2--' + variant;
  if (minimalMobile) result += ' rdr-login-v2--mobile-minimal';
  return result;
};

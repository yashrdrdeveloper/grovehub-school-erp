/** Existing auth layout classes; preserve responsive CSS until owner-module migration. */
export const LOGIN_LAYOUT_LEGACY_CLASSES = Object.freeze({
  VISUAL: 'rdr-login-v2__visual',
  MOBILE_HEADER: 'rdr-login-v2__mobile-header',
  MOBILE_IDENTITY: 'rdr-login-v2__mobile-identity',
  MOBILE_SCHOOL_NAME: 'rdr-login-v2__mobile-school-name',
  VISUAL_BRAND: 'rdr-login-v2__visual-brand',
  VISUAL_COPY: 'rdr-login-v2__visual-copy',
  VISUAL_KICKER: 'rdr-login-v2__visual-kicker',
  VISUAL_FOOTER: 'rdr-login-v2__visual-footer',
  CLOSING: 'rdr-login-v2__closing',
  FOOTER: 'rdr-o1-auth-footer',
} as const);

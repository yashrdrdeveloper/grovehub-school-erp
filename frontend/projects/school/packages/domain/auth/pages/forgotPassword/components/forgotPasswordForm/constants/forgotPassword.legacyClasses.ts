/** Transitional CSS contract: keep existing auth visuals until the scoped CSS migration. */
export const FORGOT_PASSWORD_LEGACY_CLASSES = Object.freeze({
  NOTICE_INFO: 'rdr-auth-flow__notice rdr-auth-flow__notice--info',
  NOTICE_SUCCESS: 'rdr-auth-flow__notice rdr-auth-flow__notice--success',
  NOTICE_ICON: 'rdr-auth-flow__notice-icon',
  NOTICE_COPY: 'rdr-auth-flow__notice-copy',
  SUCCESS: 'rdr-auth-flow__mobile-success',
  SUBMIT: 'rdr-login-v2__submit',
  RESEND: 'rdr-auth-flow__resend',
} as const);

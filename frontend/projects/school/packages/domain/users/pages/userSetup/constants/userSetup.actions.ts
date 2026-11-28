export const UserAction = Object.freeze({
  CREATE: 'CREATE',
  EDIT: 'EDIT',
  VIEW: 'VIEW',
  ACTIVATE: 'ACTIVATE',
  SUSPEND: 'SUSPEND',
  DEACTIVATE: 'DEACTIVATE',
  RESEND: 'RESEND',
} as const);

export type UserAction = (typeof UserAction)[keyof typeof UserAction];

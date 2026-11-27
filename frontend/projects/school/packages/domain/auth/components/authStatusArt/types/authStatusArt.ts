import type { AUTH_STATUS_KIND } from '../constants/authStatusArt.kind';

export type AuthStatusKind = (typeof AUTH_STATUS_KIND)[keyof typeof AUTH_STATUS_KIND];

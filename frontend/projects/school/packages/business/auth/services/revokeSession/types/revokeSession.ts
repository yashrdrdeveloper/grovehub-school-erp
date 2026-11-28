import type { RevokeSessionRequest } from './revokeSession.request';
import type { RevokeSessionResponse } from './revokeSession.response';

export type RevokeSession = (
  sessionId: RevokeSessionRequest['sessionId'],
) => Promise<RevokeSessionResponse>;

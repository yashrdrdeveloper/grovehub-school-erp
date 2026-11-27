/** Session summary returned by the identity service. */
export interface SessionEntity {
  sessionId: string;
  current: boolean;
  createdAt: string;
  lastActiveAt?: string;
  expiresAt: string;
}

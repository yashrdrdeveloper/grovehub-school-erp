export type { AccountEntity as Account } from './baseEntities';
export type { SessionEntity as SessionSummary } from './baseEntities';

export interface LoginResult {
  accountId: string;
  personGid: string;
  sessionId: string;
  expiresIn: number;
  requiresSetup: boolean;
}

export interface RefreshResult {
  sessionId: string;
  expiresIn: number;
  refreshed: boolean;
}

export interface StatusResult {
  status: string;
  sessionsRevoked?: boolean;
}

export interface GrantResult {
  verified: boolean;
  purpose: string;
  expiresIn: number;
}

export interface SetupResult {
  accountId: string;
  personGid: string;
  status: string;
  setupCompleted: boolean;
}

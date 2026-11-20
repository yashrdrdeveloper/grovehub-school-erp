export interface Account {
  accountId: string;
  personGid: string;
  username: string | null;
  email: string;
  status: string;
  version: number;
  createdAt: string;
  updatedAt: string;
}

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

export interface SessionSummary {
  sessionId: string;
  current: boolean;
  createdAt: string;
  lastActiveAt?: string;
  expiresAt: string;
}

export interface GrantResult {
  verified: boolean;
  purpose: string;
  expiresIn: number;
}

export interface SetupResult {
  accountId: string;
  personGid: string;
  username: string | null;
  status: string;
  setupCompleted: boolean;
}

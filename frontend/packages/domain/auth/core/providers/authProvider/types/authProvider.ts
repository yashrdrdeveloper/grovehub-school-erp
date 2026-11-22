import type { ReactNode } from 'react';
import type { Account } from '@Business/auth';

export interface AuthProviderProps {
  readonly children: ReactNode;
}
export interface SessionContextValue {
  readonly account: Account | null;
  readonly loading: boolean;
  readonly initialized: boolean;
  readonly signIn: (identifier: string, password: string) => Promise<void>;
  readonly signOut: () => Promise<void>;
  readonly reload: () => Promise<void>;
}

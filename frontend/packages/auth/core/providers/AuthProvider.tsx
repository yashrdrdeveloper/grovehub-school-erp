import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';

import { authenticationApi, type Account } from '@Business/auth';

interface SessionContextValue {
  account: Account | null;
  loading: boolean;
  signIn: (identifier: string, password: string) => Promise<void>;
  signOut: () => Promise<void>;
  reload: () => Promise<void>;
}

const SessionContext = createContext<SessionContextValue | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [account, setAccount] = useState<Account | null>(null);
  const [loading, setLoading] = useState(true);

  const reload = useCallback(async () => {
    try {
      setAccount(await authenticationApi.me());
    } catch {
      try {
        await authenticationApi.refresh();
        setAccount(await authenticationApi.me());
      } catch {
        setAccount(null);
      }
    }
  }, []);

  useEffect(() => {
    let active = true;
    const bootstrap = async () => {
      await reload();
      if (active) setLoading(false);
    };
    void bootstrap();
    return () => { active = false; };
  }, [reload]);

  const signIn = useCallback(async (identifier: string, password: string) => {
    await authenticationApi.login(identifier, password);
    await reload();
  }, [reload]);

  const signOut = useCallback(async () => {
    try {
      await authenticationApi.logout();
    } finally {
      setAccount(null);
    }
  }, []);

  const value = useMemo(
    () => ({ account, loading, signIn, signOut, reload }),
    [account, loading, signIn, signOut, reload],
  );

  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;
};

export const useAuth = (): SessionContextValue => {
  const context = useContext(SessionContext);
  if (!context) throw new Error('AuthProvider is missing');
  return context;
};

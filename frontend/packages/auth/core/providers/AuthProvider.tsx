import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from 'react';

import { authenticationApi, type Account } from '@Business/auth';

interface SessionContextValue {
  account: Account | null;
  loading: boolean;
  initialized: boolean;
  signIn: (identifier: string, password: string) => Promise<void>;
  signOut: () => Promise<void>;
  reload: () => Promise<void>;
}

const SessionContext = createContext<SessionContextValue | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [account, setAccount] = useState<Account | null>(null);
  const [loading, setLoading] = useState(false);
  const [initialized, setInitialized] = useState(false);
  const pendingReload = useRef<Promise<void> | null>(null);

  // Public routes do not bootstrap an authenticated session. The protected
  // route guard invokes reload when a workspace route actually needs one.
  const reload = useCallback((): Promise<void> => {
    if (pendingReload.current) return pendingReload.current;
    setLoading(true);
    const task = (async () => {
      try {
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
      } finally {
        setInitialized(true);
        setLoading(false);
        pendingReload.current = null;
      }
    })();
    pendingReload.current = task;
    return task;
  }, []);

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
    () => ({ account, loading, initialized, signIn, signOut, reload }),
    [account, loading, initialized, signIn, signOut, reload],
  );

  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;
};

export const useAuth = (): SessionContextValue => {
  const context = useContext(SessionContext);
  if (!context) throw new Error('AuthProvider is missing');
  return context;
};

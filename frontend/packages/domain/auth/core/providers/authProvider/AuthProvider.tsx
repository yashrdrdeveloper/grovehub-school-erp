import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
  useState,
} from 'react';

import { authenticationApi } from '@Business/auth';
import type { AuthProviderProps, SessionContextValue } from './types/authProvider';

const SessionContext = createContext<SessionContextValue | null>(null);

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [account, setAccount] = useState<SessionContextValue['account']>(null);
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

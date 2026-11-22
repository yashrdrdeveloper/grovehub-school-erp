import { useCallback } from 'react';
import { useAuth } from '@Domain/auth';

import { Shell } from '../../Shell';
import type { ProtectedLayoutComponent } from './types/protectedLayout';

/** Persistent authenticated workspace chrome. Only the nested <Outlet/> changes on child routes. */
export const ProtectedLayout: ProtectedLayoutComponent = () => {
  const { account, signOut } = useAuth();
  const handleLogout = useCallback(() => { void signOut(); }, [signOut]);
  return (
    <Shell accountLabel={account?.email ?? 'Account'} schoolLabel="RDR School"
      onLogout={handleLogout}/>
  );
};

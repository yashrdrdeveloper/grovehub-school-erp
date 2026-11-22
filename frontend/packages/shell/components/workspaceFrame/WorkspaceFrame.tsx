import { useCallback } from 'react';
import { useAuth } from '@Domain/auth';

import { Shell } from '../../Shell';
import type { WorkspaceFrameComponent } from './types/workspaceFrame';

export const WorkspaceFrame: WorkspaceFrameComponent = () => {
  const { account, signOut } = useAuth();
  const handleLogout = useCallback(() => { void signOut(); }, [signOut]);
  return (
    <Shell accountLabel={account?.email ?? 'Account'} schoolLabel="RDR School"
      onLogout={handleLogout}/>
  );
};

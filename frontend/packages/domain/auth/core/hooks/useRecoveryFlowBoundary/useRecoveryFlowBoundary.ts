import { useEffect } from 'react';
import { useLocation } from '@Components/router';
import { clearRecoveryFlow, readRecoveryFlow, recoveryPurposeForPath } from '../../session/recoveryFlowSession';

/** Keep an in-memory proof only while navigating inside the same recovery flow. */
export const useRecoveryFlowBoundary = (): void => {
  const { pathname } = useLocation();
  useEffect(() => {
    const purpose = recoveryPurposeForPath(pathname);
    if (!purpose) {
      clearRecoveryFlow();
      return;
    }
    readRecoveryFlow(purpose);
  }, [pathname]);
};

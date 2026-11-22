import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';

import { clearRecoveryFlow, readRecoveryFlow, recoveryPurposeForPath } from '../session/recoveryFlowSession';

/** Remove password-link state as soon as navigation leaves its specific flow. */
export const RecoveryFlowBoundary = () => {
  const { pathname } = useLocation();
  useEffect(() => {
    const purpose = recoveryPurposeForPath(pathname);
    if (!purpose) {
      clearRecoveryFlow();
      return;
    }
    // readRecoveryFlow deliberately clears a record belonging to the other
    // purpose (e.g. set-password -> reset-password).
    readRecoveryFlow(purpose);
  }, [pathname]);
  return null;
};

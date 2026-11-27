import { RecoveryPasswordForm } from '../components/recoveryPasswordForm';
import { RecoveryState } from '../components/recoveryState';

import type { RenderRecoveryContent } from './types/recovery.content';

/** Form/state JSX selection only; recovery hooks retain all lifecycle ownership. */
export const renderRecoveryContent: RenderRecoveryContent = ({ props, formVisible }) => {
  if (formVisible) {
    return <RecoveryPasswordForm {...props} />;
  }

  return <RecoveryState {...props} />;
};

import { LoginLayout } from '../../components/loginLayout';

import { RecoveryPasswordForm } from './components/recoveryPasswordForm';
import { RecoveryState } from './components/recoveryState';
import type { RecoveryComponent } from './types';

/** Stateless recovery/activation composition, one form in the DOM. */
export const Recovery: RecoveryComponent = (props) => {
  const formVisible = props.view === 'ready' || props.view === 'saving';
  const mobileTitle = props.invitation ? 'Set password' : 'Reset password';
  const description = 'Choose a strong password to protect your school account.';

  return (
    <LoginLayout variant="recovery" title="Set a new password."
      description={description} mobileDescription={formVisible ? description : ''}
      minimalMobile mobileTitle={mobileTitle} eyebrow="SECURE YOUR ACCOUNT">
      <div className="rdr-recovery">
        {formVisible ? <RecoveryPasswordForm {...props}/> : <RecoveryState {...props}/>}
      </div>
    </LoginLayout>
  );
};

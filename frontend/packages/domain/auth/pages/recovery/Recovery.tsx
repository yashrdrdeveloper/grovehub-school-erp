import { Box } from '@Components/components/Box';
import type { ReactNode } from 'react';
import { LoginLayout } from '../../components/loginLayout';

import { RecoveryPasswordForm } from './components/recoveryPasswordForm';
import { RecoveryState } from './components/recoveryState';
import type { RecoveryComponent } from './types';

/** Stateless recovery/activation composition, one form in the DOM. */
export const Recovery: RecoveryComponent = (props) => {
  const formVisible = props.view === 'ready' || props.view === 'saving';
  let mobileTitle = 'Reset password';
  if (props.invitation) mobileTitle = 'Set password';
  const description = 'Choose a strong password to protect your school account.';
  let mobileDescription = '';
  if (formVisible) mobileDescription = description;
  let content: ReactNode = <RecoveryState {...props}/>;
  if (formVisible) content = <RecoveryPasswordForm {...props}/>;

  return (
    <LoginLayout variant="recovery" title="Set a new password."
      description={description} mobileDescription={mobileDescription}
      minimalMobile mobileTitle={mobileTitle} eyebrow="SECURE YOUR ACCOUNT">
      <Box as="div" className="rdr-recovery">
        {content}
      </Box>
    </LoginLayout>
  );
};

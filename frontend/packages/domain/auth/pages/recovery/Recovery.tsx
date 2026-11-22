import { Box, BOX_TAG } from '@Components/components/Box';
import { LoginLayout } from '../../components/loginLayout';
import { LoginLayoutVariant } from '../../components/loginLayout/constants/loginLayout.variant';

import { renderRecoveryContent } from './renderer/recovery.content';
import type { RecoveryComponent } from './types';

/** Stateless recovery/activation composition, one form in the DOM. */
export const Recovery: RecoveryComponent = (props) => {
  const formVisible = props.view === 'ready' || props.view === 'saving';
  let mobileTitle = 'Reset password';
  if (props.invitation) mobileTitle = 'Set password';
  const description = 'Choose a strong password to protect your school account.';
  let mobileDescription = '';
  if (formVisible) mobileDescription = description;
  const content = renderRecoveryContent({ props, formVisible });

  return (
    <LoginLayout
      variant={LoginLayoutVariant.RECOVERY}
      title="Set a new password."
      description={description}
      mobileDescription={mobileDescription}
      minimalMobile
      mobileTitle={mobileTitle}
      eyebrow="SECURE YOUR ACCOUNT"
    >
      <Box as={BOX_TAG.DIV} className="rdr-recovery">
        {content}
      </Box>
    </LoginLayout>
  );
};

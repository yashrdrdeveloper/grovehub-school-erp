import { Box, BOX_TAG } from '@Components/components/Box';
import { LoginLayout } from '../../components/loginLayout';
import { LoginLayoutVariant } from '../../components/loginLayout/constants/loginLayout.variant';

import { renderRecoveryContent } from './renderer/recovery.content';
import { getRecoveryLayoutCopy } from './helpers/recovery.getLayoutCopy';
import type { RecoveryComponent } from './types';

/** Stateless recovery/activation composition, one form in the DOM. */
export const Recovery: RecoveryComponent = (props) => {
  const formVisible = props.view === 'ready' || props.view === 'saving';
  const { mobileTitle, description, mobileDescription } = getRecoveryLayoutCopy({
    invitation: props.invitation,
    formVisible,
  });
  const content = renderRecoveryContent({ props, formVisible });

  return (
    <LoginLayout
      variant={LoginLayoutVariant.RECOVERY}
      title={__('Set a new password.')}
      description={description}
      mobileDescription={mobileDescription}
      minimalMobile
      mobileTitle={mobileTitle}
      eyebrow={__('SECURE YOUR ACCOUNT')}
    >
      <Box as={BOX_TAG.DIV} className="rdr-recovery">
        {content}
      </Box>
    </LoginLayout>
  );
};

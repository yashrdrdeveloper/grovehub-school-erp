import { LoginLayout } from '../../components/loginLayout';
import { LoginLayoutVariant } from '../../components/loginLayout/constants/loginLayout.variant';

import { ForgotPasswordForm } from './components/forgotPasswordForm';
import type { ForgotPasswordComponent } from './types';

/** Stateless recovery page composition. */
export const ForgotPassword: ForgotPasswordComponent = (props) => (
  <LoginLayout
    variant={LoginLayoutVariant.FORGOT}
    title={__('Forgot your password?')}
    description={__('Enter your school email address to request a reset link.')}
    mobileDescription={__('Enter your school email address to request a reset link.')}
    minimalMobile
    mobileTitle={__('Forgot your password?')}
    eyebrow="ACCOUNT RECOVERY"
  >
    <ForgotPasswordForm {...props} />
  </LoginLayout>
);

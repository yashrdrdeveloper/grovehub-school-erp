import { LoginLayout, LoginLayoutVariant } from '../../components/loginLayout';

import { ForgotPasswordForm } from './components/forgotPasswordForm';
import type { ForgotPasswordComponent } from './types';

/** Stateless recovery page composition. */
export const ForgotPassword: ForgotPasswordComponent = (props) => (
  <LoginLayout
    variant={LoginLayoutVariant.FORGOT}
    title="Forgot your password?"
    description="Enter your school email address to request a reset link."
    mobileDescription="Enter your school email address to request a reset link."
    minimalMobile
    mobileTitle="Forgot your password?"
    eyebrow="ACCOUNT RECOVERY"
  >
    <ForgotPasswordForm {...props} />
  </LoginLayout>
);

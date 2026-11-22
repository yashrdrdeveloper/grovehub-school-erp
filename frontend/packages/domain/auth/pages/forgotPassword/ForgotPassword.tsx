import { LoginLayout } from '../LoginLayout';

import { ForgotPasswordForm } from './components/forgotPasswordForm';
import type { ForgotPasswordComponent } from './types';

/** Stateless recovery page composition. */
export const ForgotPassword: ForgotPasswordComponent = (props) => (
  <LoginLayout variant="forgot" title="Forgot your password?"
    description="Enter your school email address to request a reset link."
    mobileDescription="Enter your school email address to request a reset link."
    minimalMobile mobileTitle="Forgot your password?" eyebrow="ACCOUNT RECOVERY">
    <ForgotPasswordForm {...props}/>
  </LoginLayout>
);

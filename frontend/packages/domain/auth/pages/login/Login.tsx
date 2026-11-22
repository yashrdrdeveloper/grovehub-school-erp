import { LoginLayout } from '../LoginLayout';

import { LOGIN_COPY } from './config';
import { DesktopLogin } from './templates/desktopLogin';
import { MobileLogin } from './templates/mobileLogin';
import type { LoginPageComponent } from './types';

/** Stateless page composition. Only one active form is rendered per viewport. */
export const Login: LoginPageComponent = ({ controls, isMobile }) => {
  const FormTemplate = isMobile ? MobileLogin : DesktopLogin;

  return (
    <LoginLayout
      title={LOGIN_COPY.desktopTitle}
      description={LOGIN_COPY.desktopDescription}
      mobileTitle={LOGIN_COPY.mobileTitle}
      mobileDescription={LOGIN_COPY.mobileDescription}
      minimalMobile>
      <FormTemplate {...controls}/>
    </LoginLayout>
  );
};

import { LoginLayout, LoginLayoutVariant } from '../../components/loginLayout';

import { LOGIN_COPY } from './constants/login.copy';
import { DesktopLogin } from './templates/desktopLogin';
import { MobileLogin } from './templates/mobileLogin';
import type { LoginPageComponent } from './types';

/** Stateless page composition. Only one active form is rendered per viewport. */
export const Login: LoginPageComponent = ({ controls, isMobile }) => {
  let FormTemplate = DesktopLogin;
  if (isMobile) FormTemplate = MobileLogin;

  return (
    <LoginLayout
      variant={LoginLayoutVariant.LOGIN}
      title={LOGIN_COPY.desktopTitle}
      description={LOGIN_COPY.desktopDescription}
      mobileTitle={LOGIN_COPY.mobileTitle}
      mobileDescription={LOGIN_COPY.mobileDescription}
      minimalMobile
    >
      <FormTemplate {...controls} />
    </LoginLayout>
  );
};

import { LoginLayout } from '../../components/loginLayout';
import { LoginLayoutVariant } from '../../components/loginLayout/constants/loginLayout.variant';

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
      title={__(LOGIN_COPY.desktopTitle)}
      description={__(LOGIN_COPY.desktopDescription)}
      mobileTitle={__(LOGIN_COPY.mobileTitle)}
      mobileDescription={__(LOGIN_COPY.mobileDescription)}
      minimalMobile
    >
      <FormTemplate {...controls} />
    </LoginLayout>
  );
};

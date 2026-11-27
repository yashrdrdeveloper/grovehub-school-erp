import { LoginLayout } from '../../components/loginLayout';
import { LoginLayoutVariant } from '../../components/loginLayout/constants/loginLayout.variant';

import { LOGIN_COPY } from './constants/login.copy';
import { getLoginFormTemplate } from './renderer/login.formTemplate';
import type { LoginPageComponent } from './types';

/** Stateless page composition. Only one active form is rendered per viewport. */
export const Login: LoginPageComponent = ({ controls, isMobile }) => {
  const FormTemplate = getLoginFormTemplate(isMobile);

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

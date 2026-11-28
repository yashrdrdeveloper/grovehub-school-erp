import { AppLink } from '@Components/components/AppLink';

import { AUTH_ROUTE } from '@Business/auth/routes';
import { LOGIN_COPY } from '../../constants/login.copy';
import { LoginForm } from '../../components/LoginForm';
import type { DesktopLoginComponent } from './types/desktopLogin';

/** Desktop composition: recovery link beside the password label. */
export const DesktopLogin: DesktopLoginComponent = (props) => (
  <LoginForm
    {...props}
    passwordLabelAction={
      <AppLink className="rdr-login-v2__forgot-desktop" to={AUTH_ROUTE.FORGOT_PASSWORD}>
        {LOGIN_COPY.forgotPassword}
      </AppLink>
    }
  />
);

import { AppLink } from '@Components/components/AppLink';

import { AUTH_ROUTE } from '@Business/auth/routes';
import { LOGIN_COPY } from '../../constants/login.copy';
import { LoginForm } from '../../components/LoginForm';
import type { MobileLoginComponent } from './types/mobileLogin';

/** Mobile composition: recovery link below the password. */
export const MobileLogin: MobileLoginComponent = (props) => (
  <LoginForm
    {...props}
    afterFields={
      <AppLink className="rdr-login-v2__forgot-mobile" to={AUTH_ROUTE.FORGOT_PASSWORD}>
        {LOGIN_COPY.forgotPassword}
      </AppLink>
    }
  />
);

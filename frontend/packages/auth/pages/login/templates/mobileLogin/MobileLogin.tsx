import { Link } from 'react-router-dom';

import { LOGIN_COPY, LOGIN_ROUTES } from '../../config';
import { LoginForm } from '../../components/LoginForm';
import type { LoginProps } from '../../types';

/** Mobile-specific composition: recovery action below the credential fields. */
export const MobileLogin = (props: LoginProps) => (
  <LoginForm {...props}
    afterFields={<Link className="rdr-login-v2__forgot-mobile" to={LOGIN_ROUTES.forgotPassword}>
      {LOGIN_COPY.forgotPassword}
    </Link>}/>
);

import { Link } from 'react-router-dom';

import { LOGIN_COPY, LOGIN_ROUTES } from '../../config';
import { LoginForm } from '../../components/LoginForm';
import type { MobileLoginComponent } from './types/mobileLogin';

/** Mobile composition: recovery link between email and password. */
export const MobileLogin: MobileLoginComponent = (props) => (
  <LoginForm {...props}
    betweenFields={<Link className="rdr-login-v2__forgot-mobile" to={LOGIN_ROUTES.forgotPassword}>
      {LOGIN_COPY.forgotPassword}
    </Link>}/>
);

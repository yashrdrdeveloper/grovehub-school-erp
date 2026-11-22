import { Link } from 'react-router-dom';

import { LOGIN_COPY, LOGIN_ROUTES } from '../../config';
import { LoginForm } from '../../components/LoginForm';
import type { DesktopLoginComponent } from './types/desktopLogin';

/** Desktop composition: recovery link beside the password label. */
export const DesktopLogin: DesktopLoginComponent = (props) => (
  <LoginForm {...props}
    passwordLabelAction={<Link className="rdr-login-v2__forgot-desktop" to={LOGIN_ROUTES.forgotPassword}>
      {LOGIN_COPY.forgotPassword}
    </Link>}/>
);

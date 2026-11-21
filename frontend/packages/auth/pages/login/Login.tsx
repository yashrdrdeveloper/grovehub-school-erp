import { LoginLayout } from '../LoginLayout';

import { LOGIN_COPY } from './config';
import { DesktopLogin } from './templates/desktopLogin';
import { MobileLogin } from './templates/mobileLogin';
import type { LoginPageProps } from './types';

/** Stateless page composition. LoginContainer owns state and side effects. */
export const Login = ({ controls, isMobile }: LoginPageProps) => (
  <LoginLayout
    title={LOGIN_COPY.desktopTitle}
    description={LOGIN_COPY.desktopDescription}
    mobileTitle={LOGIN_COPY.mobileTitle}
    mobileDescription={LOGIN_COPY.mobileDescription}
    minimalMobile>
    {isMobile ? <MobileLogin {...controls}/> : <DesktopLogin {...controls}/>}
  </LoginLayout>
);

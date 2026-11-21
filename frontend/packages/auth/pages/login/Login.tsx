import { LoginLayout } from '../LoginLayout';

import { LOGIN_COPY } from './config';
import { useLogin } from './hooks/useLogin';
import { useLoginViewport } from './hooks/useLoginViewport';
import { DesktopLogin } from './templates/desktopLogin';
import { MobileLogin } from './templates/mobileLogin';

/** Page container: only binds the login use case to its active presentation. */
export const Login = () => {
  const login = useLogin();
  const isMobile = useLoginViewport();

  return (
    <LoginLayout
      title={LOGIN_COPY.desktopTitle}
      description={LOGIN_COPY.desktopDescription}
      mobileTitle={LOGIN_COPY.mobileTitle}
      mobileDescription={LOGIN_COPY.mobileDescription}
      minimalMobile>
      {isMobile ? <MobileLogin {...login}/> : <DesktopLogin {...login}/>}
    </LoginLayout>
  );
};

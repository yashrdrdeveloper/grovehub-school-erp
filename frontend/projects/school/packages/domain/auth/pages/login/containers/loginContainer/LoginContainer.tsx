import { Login } from '../../Login';
import { useLogin } from '../../hooks/useLogin';
import { useLoginViewport } from '../../hooks/useLoginViewport';

import type { LoginContainer as LoginContainerContract } from './types/loginContainer';

/** Bind the login use case to its responsive, presentational page. */
export const LoginContainer: LoginContainerContract = () => {
  const controls = useLogin();
  const isMobile = useLoginViewport();

  return <Login controls={controls} isMobile={isMobile} />;
};

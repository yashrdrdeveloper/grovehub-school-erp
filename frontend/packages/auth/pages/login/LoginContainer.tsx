import { Login } from './Login';
import { useLogin } from './hooks/useLogin';
import { useLoginViewport } from './hooks/useLoginViewport';

/** Stateful login entry point: bind auth use case and responsive presentation. */
export const LoginContainer = () => {
  const controls = useLogin();
  const isMobile = useLoginViewport();

  return <Login controls={controls} isMobile={isMobile}/>;
};

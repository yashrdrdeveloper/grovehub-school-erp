import { useId } from 'react';

import { ForgotPassword } from '../../ForgotPassword';
import { useForgotPassword } from '../../hooks/useForgotPassword';
import type { ForgotPasswordContainerComponent } from './types/forgotPasswordContainer';

export const ForgotPasswordContainer: ForgotPasswordContainerComponent = () => {
  const id = useId();
  const formId = useId();
  const controller = useForgotPassword();
  return <ForgotPassword id={id} formId={formId} {...controller}/>;
};

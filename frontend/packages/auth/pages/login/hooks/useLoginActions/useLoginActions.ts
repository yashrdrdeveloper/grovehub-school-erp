import { useCallback } from 'react';

import type { LoginFormModel } from '../useLoginForm';
import type { UseLoginActions } from './types/useLoginActions';

/** Owns named field callbacks; no field-level orchestration in JSX. */
export const useLoginActions: UseLoginActions = (form: LoginFormModel) => {
  const { setEmail, setPassword, setPasswordVisible, setFieldErrors, setErrorMessage } = form;

  const onEmailChange = useCallback((value: LoginFormModel['email']) => {
    setEmail(value);
    setFieldErrors(current => current.email ? { ...current, email: undefined } : current);
    setErrorMessage('');
  }, [setEmail, setFieldErrors, setErrorMessage]);

  const onPasswordChange = useCallback((value: LoginFormModel['password']) => {
    setPassword(value);
    setFieldErrors(current => current.password ? { ...current, password: undefined } : current);
    setErrorMessage('');
  }, [setPassword, setFieldErrors, setErrorMessage]);

  const onTogglePassword = useCallback(() => {
    setPasswordVisible(current => !current);
  }, [setPasswordVisible]);

  return { onEmailChange, onPasswordChange, onTogglePassword };
};

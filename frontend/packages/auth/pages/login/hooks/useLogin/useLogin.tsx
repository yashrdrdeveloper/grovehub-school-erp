import { useCallback, useId, useRef, useState, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';

import { useAuth } from '../../../../core/providers/AuthProvider';
import { LOGIN_ROUTES } from '../../config';
import { getLoginErrorMessage, validateLoginCredentials } from '../../helpers';
import type { LoginFieldErrors, LoginProps } from '../../types';

/** Login use-case controller: state, validation, session creation and redirect. */
export const useLogin = (): LoginProps => {
  const fieldId = useId();
  const navigate = useNavigate();
  const { signIn } = useAuth();
  const submittingRef = useRef(false);

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isPasswordVisible, setPasswordVisible] = useState(false);
  const [isSubmitting, setSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [fieldErrors, setFieldErrors] = useState<LoginFieldErrors>({});

  const handleEmailChange = useCallback((value: string) => {
    setEmail(value);
    setFieldErrors(current => current.email ? { ...current, email: undefined } : current);
    setErrorMessage('');
  }, []);

  const handlePasswordChange = useCallback((value: string) => {
    setPassword(value);
    setFieldErrors(current => current.password ? { ...current, password: undefined } : current);
    setErrorMessage('');
  }, []);

  const handleTogglePassword = useCallback(() => {
    setPasswordVisible(current => !current);
  }, []);

  const handleSubmit = useCallback(async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (submittingRef.current) return;

    const result = validateLoginCredentials({ email, password });
    setFieldErrors(result.fieldErrors);
    setErrorMessage('');
    if (Object.keys(result.fieldErrors).length > 0) return;

    submittingRef.current = true;
    setSubmitting(true);
    try {
      await signIn(result.normalizedEmail, password);
      navigate(LOGIN_ROUTES.success, { replace: true });
    } catch (failure: unknown) {
      setErrorMessage(getLoginErrorMessage(failure));
    } finally {
      submittingRef.current = false;
      setSubmitting(false);
    }
  }, [email, password, signIn, navigate]);

  return {
    fieldId,
    state: { email, password, isPasswordVisible, isSubmitting, errorMessage, fieldErrors },
    onEmailChange: handleEmailChange,
    onPasswordChange: handlePasswordChange,
    onTogglePassword: handleTogglePassword,
    onSubmit: event => { void handleSubmit(event); },
  };
};

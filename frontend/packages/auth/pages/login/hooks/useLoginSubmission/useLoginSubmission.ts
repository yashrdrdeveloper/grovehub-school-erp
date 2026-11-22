import { useCallback, useRef, useState, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';

import { useAuth } from '../../../../core/providers/AuthProvider';
import { LOGIN_ROUTES } from '../../config';
import { getLoginErrorMessage, validateLoginCredentials } from '../../helpers';
import type { LoginFormModel } from '../useLoginForm';
import type { UseLoginSubmission } from './types/useLoginSubmission';

/** Validates, serializes submission, maps failures and redirects on success. */
export const useLoginSubmission: UseLoginSubmission = (form: LoginFormModel) => {
  const { email, password, setFieldErrors, setErrorMessage } = form;
  const { signIn } = useAuth();
  const navigate = useNavigate();
  const submittingRef = useRef(false);
  const [isSubmitting, setSubmitting] = useState(false);

  const submitCredentials = useCallback(async (): Promise<void> => {
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
  }, [email, password, signIn, navigate, setFieldErrors, setErrorMessage]);

  const onSubmit = useCallback((event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    void submitCredentials();
  }, [submitCredentials]);

  return { isSubmitting, onSubmit };
};

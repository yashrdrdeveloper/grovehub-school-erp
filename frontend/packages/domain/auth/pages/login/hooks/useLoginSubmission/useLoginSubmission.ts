import { useCallback, useRef, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';

import { useAuth } from '../../../../core/providers/AuthProvider';
import { APP_ROUTE } from '@Business/app/routes';
import { _keys } from '@Base/utilities/keys';
import { getLoginErrorMessage } from '../../helpers/login.getErrorMessage';
import { validateLoginCredentials } from '../../helpers/login.validateCredentials';
import type { LoginFormModel } from '../useLoginForm';
import type { UseLoginSubmission } from './types/useLoginSubmission';

/** Validates, serializes submission, maps failures and redirects on success. */
export const useLoginSubmission: UseLoginSubmission = (form: LoginFormModel) => {
  const { email, password, setFieldErrors, setErrorMessage, setSubmitting } = form;
  const { signIn } = useAuth();
  const navigate = useNavigate();
  const submittingRef = useRef(false);

  const submitCredentials = useCallback(async (): Promise<void> => {
    if (submittingRef.current) return;

    const result = validateLoginCredentials({ email, password });
    setFieldErrors(result.fieldErrors);
    setErrorMessage('');
    if (_keys(result.fieldErrors).length > 0) return;

    submittingRef.current = true;
    setSubmitting(true);
    try {
      await signIn(result.normalizedEmail, password);
      navigate(APP_ROUTE.DASHBOARD, { replace: true });
    } catch (failure: unknown) {
      setErrorMessage(getLoginErrorMessage(failure));
    } finally {
      submittingRef.current = false;
      setSubmitting(false);
    }
  }, [email, password, signIn, navigate, setFieldErrors, setErrorMessage, setSubmitting]);

  const onSubmit = useCallback((event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    void submitCredentials();
  }, [submitCredentials]);

  return { onSubmit };
};

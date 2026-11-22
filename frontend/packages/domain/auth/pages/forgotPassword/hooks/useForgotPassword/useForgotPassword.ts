import { useCallback, useRef, useState, type FormEvent } from 'react';

import { authenticationApi } from '@Business/auth';

import { getForgotPasswordViewModel } from '../../helpers/getForgotPasswordViewModel';
import { validateForgotPasswordEmail } from '../../helpers/validateForgotPasswordEmail';
import { useForgotPasswordCountdown } from '../useForgotPasswordCountdown';
import type { UseForgotPassword } from './types/useForgotPassword';

/** Owns recovery request orchestration, input state and named actions. */
export const useForgotPassword: UseForgotPassword = () => {
  const [email, setEmail] = useState('');
  const [lastRequestedEmail, setLastRequestedEmail] = useState('');
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [emailError, setEmailError] = useState('');
  const requestInFlightRef = useRef(false);
  const { secondsLeft, startCountdown, isCoolingDown } = useForgotPasswordCountdown();

  const onEmailChange = useCallback((value: string) => {
    setEmail(value);
    setEmailError('');
    setError('');
  }, []);

  const onUseAnotherEmail = useCallback(() => {
    setLastRequestedEmail('');
    setEmail('');
    setEmailError('');
    setError('');
  }, []);

  const submitRequest = useCallback(async (): Promise<void> => {
    if (busy || requestInFlightRef.current || isCoolingDown()) return;
    const result = validateForgotPasswordEmail(email);
    setError('');
    setEmailError(result.error);
    if (result.error) return;

    requestInFlightRef.current = true;
    setBusy(true);
    try {
      await authenticationApi.forgotPassword(result.normalizedEmail);
      setLastRequestedEmail(result.normalizedEmail);
      startCountdown();
    } catch {
      setError('Unable to send a reset link right now. Please try again.');
    } finally {
      requestInFlightRef.current = false;
      setBusy(false);
    }
  }, [busy, email, isCoolingDown, startCountdown]);

  const onSubmit = useCallback((event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    void submitRequest();
  }, [submitRequest]);

  const state = { email, lastRequestedEmail, busy, secondsLeft, error, emailError };
  const view = getForgotPasswordViewModel(state);
  return { state, view, onSubmit, onEmailChange, onUseAnotherEmail };
};

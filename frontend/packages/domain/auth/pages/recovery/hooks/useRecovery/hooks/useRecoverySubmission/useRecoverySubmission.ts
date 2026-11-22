import { useCallback, useEffect, useRef, type FormEvent } from 'react';
import { authenticationApi } from '@Business/auth';
import { accountsApi } from '@Business/accounts';
import { _keys } from '@Base/utilities/keys';
import { clearRecoveryFlow } from '../../../../../../core/session/recoveryFlowSession';
import { validateRecoveryPassword } from '../../../../helpers/recovery.validatePassword';
import { classifyRecoveryFailure } from '../../../../helpers/recovery.classifyFailure';
import type { UseRecoverySubmission } from './types/useRecoverySubmission';

/** Serializes password changes; server classifications alone consume retry budget. */
export const useRecoverySubmission: UseRecoverySubmission = (purpose, form, link) => {
  const submittingRef = useRef(false);
  const flowGeneration = useRef(0);
  const invitation = purpose === 'invitation';
  const { password, confirmation, setFieldErrors, setError, clearCredentials } = form;
  const { view, setView } = link;

  useEffect(() => () => { flowGeneration.current += 1; }, [purpose]);

  const submitPassword = useCallback(async (): Promise<void> => {
    if (view !== 'ready' || submittingRef.current) return;
    const next = validateRecoveryPassword({ password, confirmation });
    setFieldErrors(next);
    if (_keys(next).length > 0) return;
    setError('');
    submittingRef.current = true;
    setView('saving');
    const generation = flowGeneration.current;
    try {
      if (invitation) await accountsApi.setup(password);
      else await authenticationApi.resetPassword(password);
      if (generation !== flowGeneration.current) return;
      clearCredentials();
      clearRecoveryFlow();
      setView('success');
    } catch (failure) {
      if (generation !== flowGeneration.current) return;
      const resolution = classifyRecoveryFailure(failure);
      clearCredentials();
      if (resolution === 'attempts-exhausted') {
        clearRecoveryFlow();
        setView('attempts-exhausted');
      } else if (resolution === 'expired') {
        clearRecoveryFlow();
        setView('expired');
      } else {
        // A transient request failure does not prove that an attempt was consumed.
        if (invitation) setError('Unable to activate your account right now. Please try again.');
        else setError('Unable to reset your password. Please try again.');
        setView('ready');
      }
    } finally {
      submittingRef.current = false;
    }
  }, [view, password, confirmation, invitation, setFieldErrors, setError, setView, clearCredentials]);

  const onSubmit = useCallback((event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    void submitPassword();
  }, [submitPassword]);

  return { onSubmit };
};

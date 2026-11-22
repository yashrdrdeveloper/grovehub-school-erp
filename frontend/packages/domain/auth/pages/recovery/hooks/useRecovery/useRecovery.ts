import { useCallback, useEffect, useId, useRef, useState, type FormEvent } from 'react';

import { authenticationApi } from '@Business/auth';
import { accountsApi } from '@Business/accounts';
import { ApiError } from '@Business/api';

import { clearRecoveryFlow, readRecoveryFlow, saveRecoveryGrant, saveRecoveryProof } from '../../../../core/session/recoveryFlowSession';
import { validateRecoveryPassword } from '../../helpers/validateRecoveryPassword';
import { classifyRecoveryFailure } from '../../helpers/classifyRecoveryFailure';
import type { RecoveryFieldErrors, RecoveryViewState } from '../../types';
import type { UseRecovery } from './types/useRecovery';

/** Controller for the shared reset/activation flow; presentation never calls APIs. */
export const useRecovery: UseRecovery = (purpose) => {
  const invitation = purpose === 'invitation';
  const id = useId();
  const proofRef = useRef<string | null>(null);
  const flowSequence = useRef(0);
  const submittingRef = useRef(false);
  const [view, setView] = useState<RecoveryViewState>('checking');
  const [slow, setSlow] = useState(false);
  const [password, setPassword] = useState('');
  const [confirmation, setConfirmation] = useState('');
  const [visible, setVisible] = useState(false);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState<RecoveryFieldErrors>({});

  const verify = useCallback(async (proof: string, sequence: number): Promise<void> => {
    setView('checking');
    try {
      if (invitation) await authenticationApi.validateInvitation(proof);
      else await authenticationApi.validateRecovery(proof);
      if (sequence !== flowSequence.current) return;
      const grant = invitation
        ? await authenticationApi.exchangeInvitation(proof)
        : await authenticationApi.exchangeRecovery(proof);
      if (sequence !== flowSequence.current) return;
      saveRecoveryGrant(purpose, grant.expiresIn);
      proofRef.current = null;
      setView('ready');
    } catch (failure) {
      if (sequence !== flowSequence.current) return;
      const expired = failure instanceof ApiError &&
        (failure.code === 'IDENTITY_TOKEN_INVALID' || [401, 403, 404, 410].includes(failure.status));
      if (expired) { clearRecoveryFlow(); proofRef.current = null; }
      setView(expired ? 'expired' : 'unavailable');
    }
  }, [invitation, purpose]);

  const verifyGrant = useCallback(async (sequence: number): Promise<void> => {
    setView('checking');
    try {
      if (invitation) await accountsApi.setupStatus();
      else await authenticationApi.resetGrantStatus();
      if (sequence === flowSequence.current) setView('ready');
    } catch (failure) {
      if (sequence !== flowSequence.current) return;
      if (failure instanceof ApiError && failure.code === 'IDENTITY_SETUP_ATTEMPTS_EXCEEDED') {
        clearRecoveryFlow();
        setView('attempts-exhausted');
      } else if (failure instanceof ApiError && [401, 403, 404, 410].includes(failure.status)) {
        clearRecoveryFlow();
        setView('expired');
      } else {
        setView('unavailable');
      }
    }
  }, [invitation]);

  useEffect(() => {
    const handleLink = () => {
      const url = new URL(window.location.href);
      const incoming = new URLSearchParams(url.hash.slice(1)).get('proof')?.trim();
      if (incoming) {
        const sequence = ++flowSequence.current;
        saveRecoveryProof(purpose, incoming);
        proofRef.current = incoming;
        setPassword('');
        setConfirmation('');
        setError('');
        setFieldErrors({});
        window.history.replaceState(window.history.state, '', url.pathname + url.search);
        void verify(incoming, sequence);
        return;
      }

      const saved = readRecoveryFlow(purpose);
      if (!saved) { setView('missing'); return; }
      const sequence = ++flowSequence.current;
      if (saved.stage === 'proof') {
        proofRef.current = saved.proof;
        void verify(saved.proof, sequence);
        return;
      }
      proofRef.current = null;
      void verifyGrant(sequence);
    };

    handleLink();
    window.addEventListener('hashchange', handleLink);
    return () => {
      ++flowSequence.current;
      window.removeEventListener('hashchange', handleLink);
    };
  }, [purpose, verify, verifyGrant]);

  useEffect(() => {
    if (view !== 'checking') { setSlow(false); return; }
    const timer = window.setTimeout(() => setSlow(true), 5000);
    return () => window.clearTimeout(timer);
  }, [view]);

  const submitPassword = useCallback(async (): Promise<void> => {
    if (view !== 'ready' || submittingRef.current) return;
    const next = validateRecoveryPassword({ password, confirmation });
    setFieldErrors(next);
    if (Object.keys(next).length > 0) return;
    setError('');
    submittingRef.current = true;
    setView('saving');
    try {
      if (invitation) await accountsApi.setup(password);
      else await authenticationApi.resetPassword(password);
      setPassword('');
      setConfirmation('');
      clearRecoveryFlow();
      setView('success');
    } catch (failure) {
      const resolution = classifyRecoveryFailure(failure);
      setPassword('');
      setConfirmation('');
      if (resolution === 'attempts-exhausted') {
        clearRecoveryFlow();
        setView('attempts-exhausted');
      } else if (resolution === 'expired') {
        clearRecoveryFlow();
        setView('expired');
      } else {
        // The frontend must not guess whether a failed request cost a server attempt.
        setError(invitation
          ? 'Unable to activate your account right now. Please try again.'
          : 'Unable to reset your password. Please try again.');
        setView('ready');
      }
    } finally {
      submittingRef.current = false;
    }
  }, [view, password, confirmation, invitation]);

  const onSubmit = useCallback((event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    void submitPassword();
  }, [submitPassword]);
  const onPasswordChange = useCallback((value: string) => {
    setPassword(value);
    setFieldErrors(current => ({ ...current, password: undefined }));
  }, []);
  const onConfirmationChange = useCallback((value: string) => {
    setConfirmation(value);
    setFieldErrors(current => ({ ...current, confirmation: undefined }));
  }, []);
  const onToggleVisibility = useCallback(() => setVisible(current => !current), []);
  const onRetry = useCallback(() => {
    const sequence = ++flowSequence.current;
    if (proofRef.current) { void verify(proofRef.current, sequence); return; }
    const saved = readRecoveryFlow(purpose);
    if (saved?.stage === 'grant') void verifyGrant(sequence);
  }, [purpose, verify, verifyGrant]);

  return {
    purpose, invitation, id, view, slow, password, confirmation, visible,
    error, fieldErrors, onPasswordChange, onConfirmationChange,
    onToggleVisibility, onSubmit, onRetry,
  };
};

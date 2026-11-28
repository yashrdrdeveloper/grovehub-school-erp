import { _includes } from '@Base/utilities/includes';
import { _trim } from '@Base/utilities/trim';
import { useCallback, useEffect, useRef, useState } from 'react';
import { authenticationApi } from '@Business/auth';
import { accountsApi } from '@Business/accounts';
import { ApiError } from '@Business/api';
import {
  clearRecoveryFlow,
  readRecoveryFlow,
  saveRecoveryGrant,
  saveRecoveryProof,
} from '../../../../../../core/session/recoveryFlowSession';
import type { RecoveryViewState } from '../../../../types';
import type { UseRecoveryLink } from './types/useRecoveryLink';

/** Handles proof exchange, grant revalidation, URL cleanup and stale responses. */
export const useRecoveryLink: UseRecoveryLink = (purpose, clearForm) => {
  const invitation = purpose === 'invitation';
  const proofRef = useRef<string | null>(null);
  const flowSequence = useRef(0);
  const [view, setView] = useState<RecoveryViewState>('checking');
  const [slow, setSlow] = useState(false);

  const verify = useCallback(
    async (proof: string, sequence: number): Promise<void> => {
      setView('checking');
      try {
        if (invitation) await authenticationApi.validateInvitation(proof);
        else await authenticationApi.validateRecovery(proof);
        if (sequence !== flowSequence.current) return;
        const exchangeGrant = async () => {
          if (invitation) return authenticationApi.exchangeInvitation(proof);
          return authenticationApi.exchangeRecovery(proof);
        };
        const grant = await exchangeGrant();
        if (sequence !== flowSequence.current) return;
        saveRecoveryGrant(purpose, grant.expiresIn);
        proofRef.current = null;
        setView('ready');
      } catch (failure) {
        if (sequence !== flowSequence.current) return;
        const expired =
          failure instanceof ApiError &&
          (failure.code === 'IDENTITY_TOKEN_INVALID' ||
            _includes([401, 403, 404, 410], failure.status));
        if (expired) {
          clearRecoveryFlow();
          proofRef.current = null;
          setView('expired');
        } else {
          setView('unavailable');
        }
      }
    },
    [invitation, purpose],
  );

  const verifyGrant = useCallback(
    async (sequence: number): Promise<void> => {
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
        } else if (failure instanceof ApiError && _includes([401, 403, 404, 410], failure.status)) {
          clearRecoveryFlow();
          setView('expired');
        } else {
          setView('unavailable');
        }
      }
    },
    [invitation],
  );

  useEffect(() => {
    const handleLink = () => {
      const url = new URL(window.location.href);
      const incoming = _trim(new URLSearchParams(url.hash.slice(1)).get('proof') ?? '');
      if (incoming) {
        const sequence = ++flowSequence.current;
        saveRecoveryProof(purpose, incoming);
        proofRef.current = incoming;
        clearForm();
        window.history.replaceState(window.history.state, '', url.pathname + url.search);
        void verify(incoming, sequence);
        return;
      }
      const saved = readRecoveryFlow(purpose);
      if (!saved) {
        setView('missing');
        return;
      }
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
  }, [purpose, clearForm, verify, verifyGrant]);

  useEffect(() => {
    if (view !== 'checking') {
      setSlow(false);
      return;
    }
    const timer = window.setTimeout(() => setSlow(true), 5000);
    return () => window.clearTimeout(timer);
  }, [view]);

  const onRetry = useCallback(() => {
    const sequence = ++flowSequence.current;
    if (proofRef.current) {
      void verify(proofRef.current, sequence);
      return;
    }
    const saved = readRecoveryFlow(purpose);
    if (saved?.stage === 'grant') void verifyGrant(sequence);
  }, [purpose, verify, verifyGrant]);

  return { view, slow, setView, onRetry };
};

import type { RecoveryArtworkCopy, RecoveryPurpose, RecoveryViewState } from '../types';
import type { GetRecoveryArtworkCopy } from './types/recovery.getArtworkCopy';

/** Single source of truth for the invitation/reset state text. */
const COPY: Record<RecoveryPurpose, Record<RecoveryViewState, RecoveryArtworkCopy>> = {
  invitation: {
    missing: {
      mobileStatusTitle: 'This link looks incomplete',
      title: 'This link looks incomplete',
      message: 'Please contact your RDR School administrator for a fresh invitation.',
    },
    checking: {
      mobileStatusTitle: 'Checking your link',
      title: 'Checking your secure link',
      message: 'We’re verifying your invitation and preparing your workspace.',
    },
    expired: {
      mobileStatusTitle: 'This link has expired',
      title: 'This link is no longer valid',
      message:
        'It may have expired or already been used. Ask your school administrator to regenerate the invitation.',
    },
    unavailable: {
      mobileStatusTitle: 'Unable to check your link',
      title: 'We couldn’t check your link',
      message: 'A connection issue prevented verification. Your link may still be valid.',
    },
    ready: { mobileStatusTitle: '', title: '', message: '' },
    saving: { mobileStatusTitle: '', title: '', message: '' },
    success: {
      mobileStatusTitle: 'Account ready',
      title: 'Welcome to RDR School!',
      message: 'Your account is ready. A new chapter begins here.',
    },
    'attempts-exhausted': {
      mobileStatusTitle: 'Setup unavailable',
      title: 'We couldn’t finish your setup',
      message: 'Please contact your RDR School administrator to regenerate your invitation.',
    },
  },
  recovery: {
    missing: {
      mobileStatusTitle: 'This link looks incomplete',
      title: 'This link looks incomplete',
      message: 'Please request another password reset link.',
    },
    checking: {
      mobileStatusTitle: 'Checking your link',
      title: 'Checking your secure link',
      message: 'We’re verifying your password reset link securely.',
    },
    expired: {
      mobileStatusTitle: 'This link has expired',
      title: 'This link is no longer valid',
      message: 'Your reset link has expired or has already been used. Please request a new one.',
    },
    unavailable: {
      mobileStatusTitle: 'Unable to check your link',
      title: 'We couldn’t check your link',
      message: 'A connection issue prevented verification. Your link may still be valid.',
    },
    ready: { mobileStatusTitle: '', title: '', message: '' },
    saving: { mobileStatusTitle: '', title: '', message: '' },
    success: {
      mobileStatusTitle: 'Password updated',
      title: 'Password updated!',
      message: 'Your password was reset. You can now sign in securely.',
    },
    'attempts-exhausted': {
      mobileStatusTitle: 'Setup unavailable',
      title: 'We couldn’t finish your setup',
      message: 'Please contact your RDR School administrator to regenerate your invitation.',
    },
  },
};

export const getRecoveryArtworkCopy: GetRecoveryArtworkCopy = ({ state, purpose }) =>
  COPY[purpose][state];

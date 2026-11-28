import type { RecoveryArtworkCopy, RecoveryPurpose, RecoveryViewState } from '../types';

/** Static human-facing copy translated once during shell-prebootstrapped module evaluation. */
export const RECOVERY_ARTWORK_COPY: Record<
  RecoveryPurpose,
  Record<RecoveryViewState, RecoveryArtworkCopy>
> = {
  invitation: {
    missing: {
      mobileStatusTitle: __('This link looks incomplete'),
      title: __('This link looks incomplete'),
      message: __('Please contact your RDR School administrator for a fresh invitation.'),
    },
    checking: {
      mobileStatusTitle: __('Checking your link'),
      title: __('Checking your secure link'),
      message: __('We’re verifying your invitation and preparing your workspace.'),
    },
    expired: {
      mobileStatusTitle: __('This link has expired'),
      title: __('This link is no longer valid'),
      message: __(
        'It may have expired or already been used. Ask your school administrator to regenerate the invitation.',
      ),
    },
    unavailable: {
      mobileStatusTitle: __('Unable to check your link'),
      title: __('We couldn’t check your link'),
      message: __('A connection issue prevented verification. Your link may still be valid.'),
    },
    ready: { mobileStatusTitle: '', title: '', message: '' },
    saving: { mobileStatusTitle: '', title: '', message: '' },
    success: {
      mobileStatusTitle: __('Account ready'),
      title: __('Welcome to RDR School!'),
      message: __('Your account is ready. A new chapter begins here.'),
    },
    'attempts-exhausted': {
      mobileStatusTitle: __('Setup unavailable'),
      title: __('We couldn’t finish your setup'),
      message: __('Please contact your RDR School administrator to regenerate your invitation.'),
    },
  },
  recovery: {
    missing: {
      mobileStatusTitle: __('This link looks incomplete'),
      title: __('This link looks incomplete'),
      message: __('Please request another password reset link.'),
    },
    checking: {
      mobileStatusTitle: __('Checking your link'),
      title: __('Checking your secure link'),
      message: __('We’re verifying your password reset link securely.'),
    },
    expired: {
      mobileStatusTitle: __('This link has expired'),
      title: __('This link is no longer valid'),
      message: __(
        'Your reset link has expired or has already been used. Please request a new one.',
      ),
    },
    unavailable: {
      mobileStatusTitle: __('Unable to check your link'),
      title: __('We couldn’t check your link'),
      message: __('A connection issue prevented verification. Your link may still be valid.'),
    },
    ready: { mobileStatusTitle: '', title: '', message: '' },
    saving: { mobileStatusTitle: '', title: '', message: '' },
    success: {
      mobileStatusTitle: __('Password updated'),
      title: __('Password updated!'),
      message: __('Your password was reset. You can now sign in securely.'),
    },
    'attempts-exhausted': {
      mobileStatusTitle: __('Setup unavailable'),
      title: __('We couldn’t finish your setup'),
      message: __('Please contact your RDR School administrator to regenerate your invitation.'),
    },
  },
};

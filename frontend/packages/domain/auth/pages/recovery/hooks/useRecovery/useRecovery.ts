import { useId } from 'react';
import { useRecoveryForm } from './hooks/useRecoveryForm';
import { useRecoveryLink } from './hooks/useRecoveryLink';
import { useRecoverySubmission } from './hooks/useRecoverySubmission';
import type { UseRecovery } from './types/useRecovery';

/** Compose independently owned link, form and submission state for recovery. */
export const useRecovery: UseRecovery = (purpose) => {
  const id = useId();
  const invitation = purpose === 'invitation';
  const form = useRecoveryForm();
  const link = useRecoveryLink(purpose, form.clearForm);
  const { onSubmit } = useRecoverySubmission(purpose, form, link);

  return {
    purpose, invitation, id,
    view: link.view, slow: link.slow, onRetry: link.onRetry,
    password: form.password, confirmation: form.confirmation, visible: form.visible,
    error: form.error, fieldErrors: form.fieldErrors,
    onPasswordChange: form.onPasswordChange,
    onConfirmationChange: form.onConfirmationChange,
    onToggleVisibility: form.onToggleVisibility,
    onSubmit,
  };
};

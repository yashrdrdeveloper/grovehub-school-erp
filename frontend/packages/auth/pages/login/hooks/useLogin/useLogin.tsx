import { useId } from 'react';

import { useLoginActions } from '../useLoginActions';
import { useLoginForm } from '../useLoginForm';
import { useLoginSubmission } from '../useLoginSubmission';
import type { UseLogin } from './types/useLogin';

/** Assembles cohesive login units without owning duplicate state. */
export const useLogin: UseLogin = () => {
  const fieldId = useId();
  const form = useLoginForm();
  const submission = useLoginSubmission(form);
  const actions = useLoginActions(form);

  return {
    fieldId,
    state: {
      email: form.email,
      password: form.password,
      isPasswordVisible: form.isPasswordVisible,
      isSubmitting: submission.isSubmitting,
      errorMessage: form.errorMessage,
      fieldErrors: form.fieldErrors,
    },
    ...actions,
    onSubmit: submission.onSubmit,
  };
};

import { useCallback, useState } from 'react';
import type { RecoveryFieldErrors } from '../../../../types';
import type { UseRecoveryForm } from './types/useRecoveryForm';

/** Single owner of password fields, visibility and validation feedback. */
export const useRecoveryForm: UseRecoveryForm = () => {
  const [password, setPassword] = useState('');
  const [confirmation, setConfirmation] = useState('');
  const [visible, setVisible] = useState(false);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState<RecoveryFieldErrors>({});

  const clearCredentials = useCallback(() => {
    setPassword('');
    setConfirmation('');
  }, []);
  const clearForm = useCallback(() => {
    clearCredentials();
    setFieldErrors({});
    setError('');
    setVisible(false);
  }, [clearCredentials]);
  const onPasswordChange = useCallback((value: string) => {
    setPassword(value);
    setFieldErrors((current) => ({ ...current, password: undefined }));
  }, []);
  const onConfirmationChange = useCallback((value: string) => {
    setConfirmation(value);
    setFieldErrors((current) => ({ ...current, confirmation: undefined }));
  }, []);
  const onToggleVisibility = useCallback(() => setVisible((current) => !current), []);

  return {
    password,
    confirmation,
    visible,
    error,
    fieldErrors,
    setError,
    setFieldErrors,
    clearCredentials,
    clearForm,
    onPasswordChange,
    onConfirmationChange,
    onToggleVisibility,
  };
};

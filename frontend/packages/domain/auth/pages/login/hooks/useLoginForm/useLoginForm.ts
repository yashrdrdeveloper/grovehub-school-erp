import { useState } from 'react';

import type { LoginFieldErrors, LoginState } from '../../types';
import type { UseLoginForm } from './types/useLoginForm';

/** Single owner of all login form state, including async submit status. */
export const useLoginForm: UseLoginForm = () => {
  const [email, setEmail] = useState<LoginState['email']>('');
  const [password, setPassword] = useState<LoginState['password']>('');
  const [isPasswordVisible, setPasswordVisible] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [fieldErrors, setFieldErrors] = useState<LoginFieldErrors>({});
  const [isSubmitting, setSubmitting] = useState(false);

  return {
    email, password, isPasswordVisible, isSubmitting, errorMessage, fieldErrors,
    setEmail, setPassword, setPasswordVisible, setSubmitting, setErrorMessage, setFieldErrors,
  };
};

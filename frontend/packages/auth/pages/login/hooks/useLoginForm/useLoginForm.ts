import { useState } from 'react';

import type { LoginFieldErrors, LoginState } from '../../types';
import type { UseLoginForm } from './types/useLoginForm';

/** Owns only the editable login form state and validation feedback. */
export const useLoginForm: UseLoginForm = () => {
  const [email, setEmail] = useState<LoginState['email']>('');
  const [password, setPassword] = useState<LoginState['password']>('');
  const [isPasswordVisible, setPasswordVisible] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [fieldErrors, setFieldErrors] = useState<LoginFieldErrors>({});

  return {
    email, password, isPasswordVisible, errorMessage, fieldErrors,
    setEmail, setPassword, setPasswordVisible, setErrorMessage, setFieldErrors,
  };
};

import { useState, useRef, type FormEvent } from 'react';

import { request, errorText } from '../../../../../business/api/client';
import type { Operator } from '../../../../../business/api/types';
import type { LoginPageProps } from '../types/login.props';

export function useLoginSubmission({ onLogin }: LoginPageProps) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const submitting = useRef(false);
  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (submitting.current) return;
    const form = event.currentTarget;
    const data = new FormData(form);
    submitting.current = true;
    setBusy(true);
    setError('');
    try {
      const operator = await request<Operator>('/internal/public/auth/login', 'POST', {
        email: data.get('email'),
        password: data.get('password'),
      });
      form.reset();
      onLogin(operator);
    } catch (failure) {
      setError(errorText(failure));
    } finally {
      submitting.current = false;
      setBusy(false);
    }
  }
  return { busy, error, onSubmit };
}

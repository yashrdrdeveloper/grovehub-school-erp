import { act, renderHook, waitFor } from '@testing-library/react';
import type { FormEvent, ReactNode } from 'react';
import { MemoryRouter } from 'react-router-dom';

import { useLogin } from '../useLogin';

const signIn = jest.fn<Promise<void>, [string, string]>();

jest.mock('../../../../../core/providers/AuthProvider', () => ({
  useAuth: () => ({ signIn }),
}));

const wrapper = ({ children }: { children: ReactNode }) => (
  <MemoryRouter>{children}</MemoryRouter>
);

const submitEvent = () => ({
  preventDefault: jest.fn(),
} as unknown as FormEvent<HTMLFormElement>);

describe('useLogin', () => {
  beforeEach(() => {
    signIn.mockReset();
  });

  it('keeps validation local and never calls auth with invalid input', () => {
    const { result } = renderHook(() => useLogin(), { wrapper });
    act(() => result.current.onSubmit(submitEvent()));
    expect(result.current.state.fieldErrors.email).toBe('Please enter your email address.');
    expect(result.current.state.fieldErrors.password).toBe('Please enter your password.');
    expect(signIn).not.toHaveBeenCalled();
  });

  it('normalizes email, preserves password, and prevents concurrent requests', async () => {
    let resolveSignIn: (() => void) | undefined;
    signIn.mockImplementation(() => new Promise<void>(resolve => { resolveSignIn = resolve; }));

    const { result } = renderHook(() => useLogin(), { wrapper });
    act(() => {
      result.current.onEmailChange(' user@school.edu ');
      result.current.onPasswordChange(' test-secret ');
    });
    act(() => {
      result.current.onSubmit(submitEvent());
      result.current.onSubmit(submitEvent());
    });
    expect(signIn).toHaveBeenCalledTimes(1);
    expect(signIn).toHaveBeenCalledWith('user@school.edu', ' test-secret ');
    expect(result.current.state.isSubmitting).toBe(true);
    await act(async () => { resolveSignIn?.(); });
    await waitFor(() => expect(result.current.state.isSubmitting).toBe(false));
  });

  it('makes errors retryable without revealing raw authentication failures', async () => {
    signIn.mockRejectedValueOnce(new Error('Invalid credentials'));
    const { result } = renderHook(() => useLogin(), { wrapper });
    act(() => {
      result.current.onEmailChange('user@school.edu');
      result.current.onPasswordChange('wrong');
    });
    act(() => result.current.onSubmit(submitEvent()));
    await waitFor(() => expect(result.current.state.errorMessage).toContain('Incorrect email or password'));
    act(() => result.current.onPasswordChange('new-secret'));
    expect(result.current.state.errorMessage).toBe('');
  });
});

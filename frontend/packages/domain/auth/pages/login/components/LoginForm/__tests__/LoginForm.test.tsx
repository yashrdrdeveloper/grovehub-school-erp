import { fireEvent, render, screen } from '@testing-library/react';
import type { FormEvent } from 'react';
import { MemoryRouter } from 'react-router-dom';

import { LoginForm } from '../LoginForm';
import { MobileLogin } from '../../../templates/mobileLogin/MobileLogin';
import type { LoginFormProps } from '../../../types';

const onEmailChange = jest.fn();
const onPasswordChange = jest.fn();
const onTogglePassword = jest.fn();
const onSubmit = jest.fn((event: FormEvent<HTMLFormElement>) => event.preventDefault());

const props: LoginFormProps = {
  fieldId: 'login',
  state: {
    email: '',
    password: '',
    isPasswordVisible: false,
    isSubmitting: false,
    errorMessage: '',
    fieldErrors: {},
  },
  onEmailChange,
  onPasswordChange,
  onTogglePassword,
  onSubmit,
};

describe('LoginForm', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders accessible fields and delegates interaction to callbacks', () => {
    render(<LoginForm {...props}/>);
    fireEvent.change(screen.getByRole('textbox', { name: 'Email address' }), {
      target: { value: 'user@school.edu' },
    });
    fireEvent.change(screen.getByLabelText('Password'), {
      target: { value: 'a-password' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Show password' }));
    fireEvent.click(screen.getByRole('button', { name: 'Sign in' }));
    expect(onEmailChange).toHaveBeenCalledWith('user@school.edu');
    expect(onPasswordChange).toHaveBeenCalledWith('a-password');
    expect(onTogglePassword).toHaveBeenCalledTimes(1);
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });

  it('keeps mobile forgot password between email and password and shows the complete API error', () => {
    const message = 'Unable to sign in right now. Please try again.';
    const { container } = render(
      <MemoryRouter>
        <MobileLogin {...props} state={{ ...props.state, errorMessage: message }}/>
      </MemoryRouter>,
    );
    const emailField = container.querySelector('#login-email')?.closest('.rdr-login-v2__field');
    const passwordField = container.querySelector('#login-password')?.closest('.rdr-login-v2__field');
    const forgot = screen.getByRole('link', { name: 'Forgot password?' });
    expect(emailField).not.toBeNull();
    expect(passwordField).not.toBeNull();
    expect(Boolean(emailField!.compareDocumentPosition(forgot) & Node.DOCUMENT_POSITION_FOLLOWING)).toBe(true);
    expect(Boolean(forgot.compareDocumentPosition(passwordField!) & Node.DOCUMENT_POSITION_FOLLOWING)).toBe(true);
    expect(screen.getByRole('alert').textContent).toContain(message);
  });

  it('renders errors and disables submit without owning any state', () => {
    const state = {
      ...props.state,
      isSubmitting: true,
      errorMessage: 'Unable to sign in right now. Please try again.',
      fieldErrors: { email: 'Enter a valid email address.' },
    };
    render(<LoginForm {...props} state={state}/>);
    expect(screen.getByRole('alert').textContent).toContain('Unable to sign in');
    expect(screen.getByText('Enter a valid email address.')).toBeTruthy();
    const fieldErrors=document.querySelectorAll('.rdr-login-v2__field-error--reserved');
    expect(fieldErrors).toHaveLength(2);
    expect(Array.from(fieldErrors).every(node => node.getAttribute('aria-live') === 'polite')).toBe(true);
    expect((screen.getByRole('button', { name: 'Signing in…' }) as HTMLButtonElement).disabled).toBe(true);
  });
});

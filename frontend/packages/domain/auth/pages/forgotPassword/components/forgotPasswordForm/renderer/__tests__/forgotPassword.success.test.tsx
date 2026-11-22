import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { renderForgotPasswordSuccess } from '../forgotPassword.success';

import type { ForgotPasswordProps } from '../../../../types/forgotPassword.props';

const props: Pick<ForgotPasswordProps, 'formId' | 'state' | 'view' | 'onUseAnotherEmail'> = {
  formId: 'forgot-submit',
  state: {
    email: 'parent@example.com',
    lastRequestedEmail: 'parent@example.com',
    busy: false,
    secondsLeft: 30,
    error: '',
    emailError: '',
  },
  view: {
    cooldown: true,
    noticeSuccess: true,
    noticeTitle: 'Request received!',
    noticeMessage: 'Check your inbox.',
    buttonLabel: 'Resend in 30s',
    resendLabel: 'Resend in 30s',
  },
  onUseAnotherEmail: () => {},
};

describe('renderForgotPasswordSuccess', () => {
  it('preserves its accessible heading, email, login link and disabled resend', () => {
    render(
      <MemoryRouter>
        {renderForgotPasswordSuccess(props)}
      </MemoryRouter>,
    );

    expect(screen.getByRole('heading', { level: 3, name: 'Check your email' })).toBeTruthy();
    expect(screen.getByText('parent@example.com')).toBeTruthy();
    expect(screen.getByRole('link', { name: 'Back to sign in' }).getAttribute('href')).toBe('/login');
    expect((screen.getByRole('button', { name: 'Resend in 30s' }) as HTMLButtonElement).disabled).toBe(true);
  });
});

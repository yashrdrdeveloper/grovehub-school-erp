import { render, screen } from '@testing-library/react';

import { renderForgotPasswordNotice } from '../forgotPassword.notice';

import type { ForgotPasswordViewModel } from '../../../../types/forgotPassword.viewModel';

const info: ForgotPasswordViewModel = {
  cooldown: false,
  noticeSuccess: false,
  noticeTitle: 'Secure password reset',
  noticeMessage: 'Enter your email address.',
  buttonLabel: 'Send reset link',
  resendLabel: 'Resend link',
};

describe('renderForgotPasswordNotice', () => {
  it('renders an accessible information notice', () => {
    render(<>{renderForgotPasswordNotice(info)}</>);

    expect(screen.getByRole('status').textContent).toContain('Secure password reset');
    expect(screen.getByTestId('forgot-password-notice').className).toContain('--info');
  });

  it('renders the success notice with matching title and status', () => {
    render(
      <>
        {renderForgotPasswordNotice({
          ...info,
          noticeSuccess: true,
          noticeTitle: 'Request received!',
        })}
      </>,
    );

    expect(screen.getByRole('status').textContent).toContain('Request received!');
    expect(screen.getByTestId('forgot-password-notice').className).toContain('--success');
  });
});

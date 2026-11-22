import { render, screen } from '@testing-library/react';

import { Mail } from '@Components/icons';

import { CredentialField } from '../CredentialField';

describe('CredentialField', () => {
  it('retains the caller and validation descriptions using the shared form helper', () => {
    render(
      <CredentialField
        Icon={Mail}
        id="guardian-email"
        label="Email"
        error="Please enter a valid email"
        dataTestId="guardian-email"
        aria-describedby="email-hint"
      />,
    );

    const input = screen.getByRole('textbox', { name: 'Email' });

    expect(input.getAttribute('aria-describedby')).toBe('email-hint guardian-email-error');
    expect(screen.getByText('Please enter a valid email').id).toBe('guardian-email-error');
    expect(screen.getByTestId('guardian-email-wrapper')).toBeTruthy();
  });

  it('keeps a reserved error row without falsely marking a valid field', () => {
    render(<CredentialField Icon={Mail} id="student-email" label="Email" />);

    const input = screen.getByRole('textbox', { name: 'Email' });
    expect(input.getAttribute('aria-invalid')).toBe('false');
    expect(input.getAttribute('aria-describedby')).toBeNull();
    expect(document.getElementById('student-email-error')).toBeTruthy();
  });
});

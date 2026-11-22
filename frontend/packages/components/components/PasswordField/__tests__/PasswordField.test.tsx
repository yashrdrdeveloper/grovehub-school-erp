import { fireEvent, render, screen } from '@testing-library/react';

import { PasswordField } from '../PasswordField';

describe('PasswordField', () => {
  it('preserves password masking and toggles visibility via the shared Button', () => {
    render(<PasswordField label="Password" id="password" />);
    const input = screen.getByLabelText('Password') as HTMLInputElement;

    expect(input.type).toBe('password');
    fireEvent.click(screen.getByRole('button', { name: 'Show password' }));
    expect(input.type).toBe('text');
    fireEvent.click(screen.getByRole('button', { name: 'Hide password' }));
    expect(input.type).toBe('password');
  });

  it('keeps descriptions and disables its visibility control when disabled', () => {
    render(
      <PasswordField
        label="New password"
        id="new-password"
        error="Password too short"
        aria-describedby="instructions"
        disabled
      />,
    );

    const input = screen.getByLabelText('New password');
    expect(input.getAttribute('aria-describedby')).toBe('instructions new-password-description');
    expect(screen.getByRole('button', { name: 'Show password' }).hasAttribute('disabled')).toBe(
      true,
    );
  });
});

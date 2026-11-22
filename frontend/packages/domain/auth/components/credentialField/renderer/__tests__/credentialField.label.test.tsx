import { render, screen } from '@testing-library/react';

import { renderCredentialFieldLabel } from '../credentialField.label';

describe('renderCredentialFieldLabel', () => {
  it('links a standalone label to its control', () => {
    render(<>{renderCredentialFieldLabel({ id: 'email-field', label: 'Email' })}</>);

    expect(screen.getByText('Email').getAttribute('for')).toBe('email-field');
    expect(document.querySelector('.rdr-login-v2__password-label')).toBeNull();
  });

  it('groups the label with an optional action', () => {
    render(
      <>
        {renderCredentialFieldLabel({
          id: 'password-field',
          label: 'Password',
          groupLabel: true,
          labelAction: <span>Help</span>,
        })}
      </>,
    );

    const label = screen.getByText('Password');
    expect(label.getAttribute('for')).toBe('password-field');
    expect(label.parentElement?.className).toContain('rdr-login-v2__password-label');
    expect(screen.getByText('Help')).toBeTruthy();
  });
});

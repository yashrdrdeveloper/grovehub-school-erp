import { render, screen } from '@testing-library/react';

import { Input } from '../Input';

describe('Input', () => {
  it('preserves caller-provided and generated accessibility descriptions', () => {
    render(
      <Input
        label="School Email"
        id="school-email"
        dataTestId="email-control"
        aria-describedby="external-description"
        hint="Use the registered email"
      />,
    );

    const input = screen.getByRole('textbox', { name: 'School Email' });
    expect(input.getAttribute('data-testid')).toBe('email-control');
    expect(input.getAttribute('aria-describedby')).toBe(
      'external-description school-email-description',
    );
    expect(screen.getByText('Use the registered email').id).toBe('school-email-description');
  });

  it('uses Lodash-normalized label selectors when neither ID is provided', () => {
    render(<Input label="Parent's Full Name" />);
    expect(screen.getByRole('textbox', { name: "Parent's Full Name" }).dataset.testid).toBe(
      'parents-full-name',
    );
  });
});

import { render, screen } from '@testing-library/react';

import { FormField } from '../FormField';

describe('FormField', () => {
  it('renders a native linked label and accessible validation feedback', () => {
    render(
      <FormField
        label="School email"
        htmlFor="email"
        dataTestId="form-email"
        error="Please enter a valid email"
        required
      >
        <input id="email" />
      </FormField>,
    );

    expect(screen.getByLabelText(/School email/)).toBeTruthy();
    expect(screen.getByRole('alert').id).toBe('email-description');
    expect(screen.getByTestId('form-email-label').tagName).toBe('LABEL');
    expect(screen.getByTestId('form-email-field')).toBeTruthy();
  });
});

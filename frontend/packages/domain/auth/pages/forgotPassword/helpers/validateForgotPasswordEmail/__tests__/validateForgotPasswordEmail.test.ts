import { validateForgotPasswordEmail } from '../validateForgotPasswordEmail';

describe('validateForgotPasswordEmail', () => {
  it('normalizes a valid email', () => {
    expect(validateForgotPasswordEmail('  Parent@school.edu ')).toEqual({
      normalizedEmail: 'Parent@school.edu', error: '',
    });
  });
  it('requires an email', () => {
    expect(validateForgotPasswordEmail(' ').error).toBe('Please enter your email address.');
  });
  it('rejects malformed addresses', () => {
    expect(validateForgotPasswordEmail('invalid').error).toBe('Enter a valid email address.');
  });
});

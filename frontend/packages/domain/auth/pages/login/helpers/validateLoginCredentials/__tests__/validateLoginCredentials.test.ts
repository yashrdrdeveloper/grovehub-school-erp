import { validateLoginCredentials } from '../validateLoginCredentials';

describe('login credential validation', () => {
  it('normalizes email without changing the original password', () => {
    expect(validateLoginCredentials({
      email: '  person@school.edu  ', password: '  secret  ',
    })).toEqual({ normalizedEmail: 'person@school.edu', fieldErrors: {} });
  });

  it('uses configured required-field messages', () => {
    expect(validateLoginCredentials({ email: ' ', password: '' }).fieldErrors).toEqual({
      email: 'Please enter your email address.',
      password: 'Please enter your password.',
    });
  });

  it('rejects malformed email without rejecting a provided password', () => {
    expect(validateLoginCredentials({ email: 'bad-address', password: 'valid' }).fieldErrors).toEqual({
      email: 'Enter a valid email address.',
    });
  });
});

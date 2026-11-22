import { validateRecoveryPassword } from '../validateRecoveryPassword';

describe('validateRecoveryPassword', () => {
  it('requires both fields', () => {
    expect(validateRecoveryPassword({ password: '', confirmation: '' })).toEqual({
      password: 'Please enter a password.', confirmation: 'Please confirm your password.',
    });
  });
  it('enforces the minimum length without calling the API', () => {
    expect(validateRecoveryPassword({ password: 'short', confirmation: 'short' })).toEqual({
      password: 'Use at least 12 characters.',
    });
  });
  it('requires matching confirmation', () => {
    expect(validateRecoveryPassword({ password: 'a-strong-password-123', confirmation: 'different' })).toEqual({
      confirmation: 'Passwords do not match.',
    });
  });
});

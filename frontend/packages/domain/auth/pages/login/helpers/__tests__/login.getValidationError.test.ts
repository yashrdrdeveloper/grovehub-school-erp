import { LoginField } from '../../constants/login.fields';
import { getLoginValidationError } from '../login.getValidationError';

describe('getLoginValidationError', () => {
  it('uses field-specific required text in field order', () => {
    expect(getLoginValidationError(LoginField.EMAIL, '')).toBe('Please enter your email address.');
    expect(getLoginValidationError(LoginField.PASSWORD, '')).toBe('Please enter your password.');
  });
  it('validates email formatting but not password formatting', () => {
    expect(getLoginValidationError(LoginField.EMAIL, 'bad-address')).toBe('Enter a valid email address.');
    expect(getLoginValidationError(LoginField.PASSWORD, 'bad-address')).toBeUndefined();
  });
});

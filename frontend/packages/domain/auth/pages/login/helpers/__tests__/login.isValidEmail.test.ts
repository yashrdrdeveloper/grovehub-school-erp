import { isValidLoginEmail } from '../login.isValidEmail';

describe('isValidLoginEmail', () => {
  it('accepts email addresses with a domain and local part', () => {
    expect(isValidLoginEmail('Person@School.edu')).toBe(true);
  });
  it('rejects malformed addresses', () => {
    expect(isValidLoginEmail('not-email')).toBe(false);
    expect(isValidLoginEmail('')).toBe(false);
  });
});

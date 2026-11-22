import { normalizeLoginEmail } from '../normalizeLoginEmail';

describe('normalizeLoginEmail', () => {
  it('trims whitespace without changing the casing', () => {
    expect(normalizeLoginEmail('  Person@School.edu  ')).toBe('Person@School.edu');
  });
});

import { _toLower } from '../toLower';

describe('_toLower', () => {
  it('converts backend-style status strings to lower-case strings without changing punctuation', () => {
    expect(_toLower('PENDING_SETUP')).toBe('pending_setup');
    expect(_toLower('ACTIVE')).toBe('active');
  });
});

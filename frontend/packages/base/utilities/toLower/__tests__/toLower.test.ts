import { _toLower } from '../toLower';

describe('_toLower', () => {
  it('converts backend-style status strings to readable lower-case labels', () => {
    expect(_toLower('PENDING_SETUP')).toBe('pending setup');
    expect(_toLower('ACTIVE')).toBe('active');
  });
});

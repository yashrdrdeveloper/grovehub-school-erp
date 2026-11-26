import { getStatusPillDisplayLabel } from '../statusPill.getDisplayLabel';

describe('getStatusPillDisplayLabel', () => {
  it('converts backend status identifiers into readable text', () => {
    expect(getStatusPillDisplayLabel('PENDING_SETUP')).toBe('pending setup');
    expect(getStatusPillDisplayLabel('SUSPENDED')).toBe('suspended');
  });
});

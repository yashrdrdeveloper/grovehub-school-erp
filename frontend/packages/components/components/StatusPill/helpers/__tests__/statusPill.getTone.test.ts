import { getStatusPillTone } from '../statusPill.getTone';

describe('getStatusPillTone', () => {
  it('maps recognized statuses to their visual treatment', () => {
    expect(getStatusPillTone('ACTIVE')).toBe('positive');
    expect(getStatusPillTone('PENDING_SETUP')).toBe('notice');
  });
  it('uses neutral for future or unknown statuses', () => {
    expect(getStatusPillTone('NEW_STATUS')).toBe('neutral');
  });
});

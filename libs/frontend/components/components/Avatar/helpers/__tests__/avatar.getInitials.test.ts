import { getAvatarInitials } from '../avatar.getInitials';

describe('getAvatarInitials', () => {
  it('uses the first and last names', () => {
    expect(getAvatarInitials(' Rahul Kumar Sharma ')).toBe('RS');
  });

  it('supports single and empty names', () => {
    expect(getAvatarInitials('Rahul')).toBe('RA');
    expect(getAvatarInitials('')).toBe('');
  });
});

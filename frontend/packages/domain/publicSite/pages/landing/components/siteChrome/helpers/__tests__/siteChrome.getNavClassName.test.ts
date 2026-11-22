import { getSiteNavClassName } from '../siteChrome.getNavClassName';

describe('getSiteNavClassName', () => {
  it('marks only the currently active route', () => {
    expect(getSiteNavClassName('/campus', '/campus')).toBe('active');
    expect(getSiteNavClassName('/campus', '/gallery')).toBe('');
  });
});

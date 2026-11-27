import { getPublicSheetStoryKey } from '../publicSheet.getStoryKey';
describe('getPublicSheetStoryKey', () => {
  it('preserves actual story panel and defaults non-story routes', () => {
    expect(getPublicSheetStoryKey('director')).toBe('director');
    expect(getPublicSheetStoryKey('philosophy')).toBe('philosophy');
    expect(getPublicSheetStoryKey('enquiry')).toBe('history');
    expect(getPublicSheetStoryKey(null)).toBe('history');
  });
});

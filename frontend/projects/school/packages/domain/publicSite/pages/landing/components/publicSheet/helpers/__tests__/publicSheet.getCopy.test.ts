import { getPublicSheetCopy } from '../publicSheet.getCopy';
describe('getPublicSheetCopy', () => {
  it('resolves panel labels and the empty-view fallback', () => {
    expect(getPublicSheetCopy('about').title).toBe('Discover RDR School');
    expect(getPublicSheetCopy(null).title).toBe('RDR School information');
  });
});

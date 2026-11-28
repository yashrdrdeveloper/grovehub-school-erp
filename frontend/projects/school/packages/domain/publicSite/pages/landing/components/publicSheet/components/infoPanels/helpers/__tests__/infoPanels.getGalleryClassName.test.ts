import { getInfoGalleryClassName } from '../infoPanels.getGalleryClassName';

describe('getInfoGalleryClassName', () => {
  it('features only the first gallery preview', () => {
    expect(getInfoGalleryClassName(0)).toBe('featured');
    expect(getInfoGalleryClassName(1)).toBe('');
  });
});

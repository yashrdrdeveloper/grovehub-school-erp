import { getGalleryPhotoClassName } from '../campusContent.getGalleryPhotoClassName';

describe('getGalleryPhotoClassName', () => {
  it('features the first and every fifth archived photo', () => {
    expect(getGalleryPhotoClassName(0)).toBe('rdr-gallery-grid__wide');
    expect(getGalleryPhotoClassName(5)).toBe('rdr-gallery-grid__wide');
    expect(getGalleryPhotoClassName(1)).toBe('');
  });
});

import type { GetGalleryPhotoClassName } from './types/campusContent.getGalleryPhotoClassName';

export const getGalleryPhotoClassName: GetGalleryPhotoClassName = (index) => {
  if (index % 5 === 0) return 'rdr-gallery-grid__wide';
  return '';
};

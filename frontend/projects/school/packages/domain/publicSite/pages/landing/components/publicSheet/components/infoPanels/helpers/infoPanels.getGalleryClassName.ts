import type { GetInfoGalleryClassName } from './types/infoPanels.getGalleryClassName';
export const getInfoGalleryClassName: GetInfoGalleryClassName = (index) => {
  if (index === 0) return 'featured';
  return '';
};

/** Closed school-art motif variants, exported by the owning siteMedia component unit. */
export const PHOTO_MOTIFS_VARIANT = Object.freeze({
  HERO: 'hero',
  GALLERY: 'gallery',
  BANNER: 'banner',
} as const);
export type PhotoMotifsVariant = (typeof PHOTO_MOTIFS_VARIANT)[keyof typeof PHOTO_MOTIFS_VARIANT];

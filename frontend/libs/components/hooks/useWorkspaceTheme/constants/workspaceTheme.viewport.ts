/** Theme selection is always resolved through an exhaustive viewport strategy. */
export enum ThemeViewport {
  MOBILE = 'MOBILE',
  DESKTOP = 'DESKTOP',
}

export const THEME_VIEWPORT_FROM_MEDIA: Readonly<Record<string, ThemeViewport>> = {
  true: ThemeViewport.MOBILE,
  false: ThemeViewport.DESKTOP,
};

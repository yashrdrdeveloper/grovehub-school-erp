import type { GetLoginLayoutHeading } from './types/loginLayout.getHeading';

/** Display one heading string while retaining a single accessible H1 across viewports. */
export const getLoginLayoutHeading: GetLoginLayoutHeading = (isMobile, title, mobileTitle) => {
  if (isMobile) return mobileTitle ?? title;
  return title;
};

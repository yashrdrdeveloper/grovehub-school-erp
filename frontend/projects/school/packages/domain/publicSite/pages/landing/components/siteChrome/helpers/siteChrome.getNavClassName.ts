import type { GetSiteNavClassName } from './types/siteChrome.getNavClassName';
export const getSiteNavClassName: GetSiteNavClassName = (currentPath, targetPath) => {
  if (currentPath === targetPath) return 'active';
  return '';
};

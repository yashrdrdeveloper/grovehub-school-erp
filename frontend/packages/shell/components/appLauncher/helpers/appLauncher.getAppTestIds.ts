import type { GetAppTestIds } from './types/appLauncher.getAppTestIds';

export const getAppTestIds: GetAppTestIds = (path) => {
  const key = path.replace(/^\/+/, '').replace(/[^a-zA-Z0-9-]+/g, '-').toLowerCase();
  return { card: 'catalog-app-' + key, pin: 'catalog-pin-' + key };
};

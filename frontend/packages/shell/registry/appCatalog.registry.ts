import { Building2 } from '@Components/icons';
import { APP_ROUTE } from '@Business/app/routes';

/** Only implemented protected domains belong in the catalog; no speculative pages. */
export type CatalogApp = {
  readonly name: string;
  readonly description: string;
  readonly path: string;
  readonly icon: typeof Building2;
  readonly category: string;
};

export const appCatalog: CatalogApp[] = [
  { name: 'Dashboard', description: 'Overview of your school', path: APP_ROUTE.DASHBOARD, icon: Building2, category: 'Workspace' },
];

export const defaultPinnedApps = [APP_ROUTE.DASHBOARD];

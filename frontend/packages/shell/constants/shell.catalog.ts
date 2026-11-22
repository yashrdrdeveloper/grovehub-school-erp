import { APP_ROUTE } from '@Business/app/routes';
import { Building2 } from '@Components/icons';

export interface CatalogApp {
  readonly name: string;
  readonly description: string;
  readonly path: string;
  readonly icon: typeof Building2;
  readonly category: string;
}

/** Register implemented protected domains only. No planned or imaginary pages. */
export const appCatalog: CatalogApp[] = [
  {
    name: 'Dashboard',
    description: 'Overview of your school',
    path: APP_ROUTE.DASHBOARD,
    icon: Building2,
    category: 'Workspace',
  },
];

export const defaultPinnedApps = [APP_ROUTE.DASHBOARD];

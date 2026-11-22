import { Building2 } from 'lucide-react';

/** Only implemented protected domains belong in the catalog; no speculative pages. */
export type CatalogApp = {
  readonly name: string;
  readonly description: string;
  readonly path: string;
  readonly icon: typeof Building2;
  readonly category: string;
};

export const appCatalog: CatalogApp[] = [
  { name: 'Dashboard', description: 'Overview of your school', path: '/dashboard', icon: Building2, category: 'Workspace' },
];

export const defaultPinnedApps = ['/dashboard'];

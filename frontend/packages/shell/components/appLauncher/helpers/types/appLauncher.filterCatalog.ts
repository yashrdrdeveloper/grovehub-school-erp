import type { CatalogApp } from '../../../../registry/appCatalog.registry';
export interface CatalogResult {
  readonly matching: CatalogApp[];
  readonly pinnedApps: CatalogApp[];
  readonly otherApps: CatalogApp[];
}
export type FilterCatalog = (search: string, pinned: readonly string[]) => CatalogResult;

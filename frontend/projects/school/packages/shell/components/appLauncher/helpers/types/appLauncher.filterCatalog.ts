import type { CatalogApp } from '../../../../constants/shell.catalog';
export interface CatalogResult {
  readonly matching: CatalogApp[];
  readonly pinnedApps: CatalogApp[];
  readonly otherApps: CatalogApp[];
}
export type FilterCatalog = (search: string, pinned: readonly string[]) => CatalogResult;

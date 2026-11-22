import type { ChangeEventHandler } from 'react';
import type { CatalogApp } from '../../../registry/appCatalog.registry';

export interface AppLauncherProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly pinned: string[];
  readonly onTogglePin: (path: string) => void;
  readonly sidebarCollapsed: boolean;
}
export interface AppLauncherModel {
  readonly filter: string;
  readonly matching: CatalogApp[];
  readonly pinnedApps: CatalogApp[];
  readonly otherApps: CatalogApp[];
  readonly onFilterChange: ChangeEventHandler<HTMLInputElement>;
}
export interface AppLauncherViewProps extends AppLauncherProps {
  readonly model: AppLauncherModel;
}

import type { ReactElement } from 'react';
import type { CatalogApp } from '../../../registry/appCatalog.registry';
export interface RenderAppLauncherAppProps {
  readonly app: CatalogApp;
  readonly isPinned: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onTogglePin: (path: string) => void;
}
export type RenderAppLauncherApp = (props: RenderAppLauncherAppProps) => ReactElement;

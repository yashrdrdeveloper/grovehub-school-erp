import type { ReactElement } from 'react';
import type { CatalogApp } from '../../../registry/appCatalog.registry';
import type { ShellViewModel } from '../../../hooks/useShell';

export interface PinnedNavigationProps {
  readonly pinned: ShellViewModel['pinned'];
  readonly closeSidebar: ShellViewModel['closeSidebar'];
  readonly togglePin: ShellViewModel['togglePin'];
}

export interface PinnedAppEntry {
  readonly app: CatalogApp;
  readonly path: string;
}

export interface PinnedAppProps extends PinnedNavigationProps {
  readonly entry: PinnedAppEntry;
}

export type PinnedNavigationComponent = (props: PinnedNavigationProps) => ReactElement;

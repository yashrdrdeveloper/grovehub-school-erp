import type { ShellStateModel } from '../../useShellState';

export interface ShellActions {
  readonly openSidebar: () => void;
  readonly closeSidebar: () => void;
  readonly openCatalog: () => void;
  readonly openAccountPreferences: () => void;
  readonly toggleQuickSettings: () => void;
  readonly toggleSidebarCollapsed: () => void;
  readonly togglePin: (path: string) => void;
}

export type UseShellActions = (state: ShellStateModel) => ShellActions;

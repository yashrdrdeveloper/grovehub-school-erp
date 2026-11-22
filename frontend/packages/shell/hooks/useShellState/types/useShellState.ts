import type { Dispatch, SetStateAction } from 'react';

export interface ShellStateModel {
  readonly sidebarOpen: boolean;
  readonly sidebarCollapsed: boolean;
  readonly quickSettings: boolean;
  readonly launcherOpen: boolean;
  readonly pinned: string[];
  readonly setSidebarOpen: Dispatch<SetStateAction<boolean>>;
  readonly setSidebarCollapsed: Dispatch<SetStateAction<boolean>>;
  readonly setQuickSettings: Dispatch<SetStateAction<boolean>>;
  readonly setLauncherOpen: Dispatch<SetStateAction<boolean>>;
  readonly setPinned: Dispatch<SetStateAction<string[]>>;
}

export type UseShellState = () => ShellStateModel;

import type { Dispatch, SetStateAction } from 'react';
import type { ThemePreference } from '@Components/hooks/useWorkspaceTheme';

export interface ShellViewModel {
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
  readonly openSidebar: () => void;
  readonly closeSidebar: () => void;
  readonly openCatalog: () => void;
  readonly openAccountPreferences: () => void;
  readonly toggleQuickSettings: () => void;
  readonly toggleSidebarCollapsed: () => void;
  readonly togglePin: (path: string) => void;
  readonly preference: ThemePreference;
  readonly setPreference: (next: ThemePreference) => void;
}
export type UseShell = () => ShellViewModel;

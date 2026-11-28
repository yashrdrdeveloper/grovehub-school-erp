import { useCallback, useEffect, useState } from 'react';

import { useBoolean } from '@Components/hooks/useBoolean';
import { useWorkspaceTheme } from '@Components/hooks/useWorkspaceTheme';
import { SCHOOL_THEME_STORAGE_KEY } from '@Domain/preferences';

import {
  readPinnedApps,
  readSidebarCollapsed,
  writePinnedApps,
  writeSidebarCollapsed,
} from './helpers/shell.preferences';
import { togglePinnedPath } from './helpers/shell.togglePinnedPath';

import type { UseShell } from './types/useShell';

/** The sole owner of workspace chrome state; persistence only runs after state commits. */
export const useShell: UseShell = () => {
  const sidebar = useBoolean();
  const quick = useBoolean();
  const launcher = useBoolean();

  const [sidebarCollapsed, setSidebarCollapsed] = useState(readSidebarCollapsed);
  const [pinned, setPinned] = useState<string[]>(readPinnedApps);
  const theme = useWorkspaceTheme(SCHOOL_THEME_STORAGE_KEY);

  useEffect(() => {
    writeSidebarCollapsed(sidebarCollapsed);
  }, [sidebarCollapsed]);

  useEffect(() => {
    writePinnedApps(pinned);
  }, [pinned]);

  const openCatalog = useCallback(() => {
    launcher.setTrue();
    sidebar.setFalse();
  }, [launcher.setTrue, sidebar.setFalse]);

  const openAccountPreferences = useCallback(() => {
    quick.setTrue();
    sidebar.setTrue();
  }, [quick.setTrue, sidebar.setTrue]);

  const toggleSidebarCollapsed = useCallback(() => {
    setSidebarCollapsed((value) => !value);
  }, []);

  const togglePin = useCallback((path: string) => {
    setPinned((previous) => togglePinnedPath(previous, path));
  }, []);

  return {
    sidebarOpen: sidebar.value,
    sidebarCollapsed,
    quickSettings: quick.value,
    launcherOpen: launcher.value,
    pinned,
    setSidebarOpen: sidebar.setValue,
    setSidebarCollapsed,
    setQuickSettings: quick.setValue,
    setLauncherOpen: launcher.setValue,
    setPinned,
    openSidebar: sidebar.setTrue,
    closeSidebar: sidebar.setFalse,
    openCatalog,
    openAccountPreferences,
    toggleQuickSettings: quick.toggle,
    toggleSidebarCollapsed,
    togglePin,
    ...theme,
  };
};

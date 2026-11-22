import { useCallback } from 'react';

import type { ShellStateModel } from '../useShellState';
import type { UseShellActions } from './types/useShellActions';

/** Named, stable actions; browser persistence never appears in Shell JSX. */
export const useShellActions: UseShellActions = (state: ShellStateModel) => {
  const { setSidebarOpen, setSidebarCollapsed, setQuickSettings, setLauncherOpen, setPinned } = state;

  const openSidebar = useCallback(() => setSidebarOpen(true), [setSidebarOpen]);
  const closeSidebar = useCallback(() => setSidebarOpen(false), [setSidebarOpen]);
  const openCatalog = useCallback(() => {
    setLauncherOpen(true);
    setSidebarOpen(false);
  }, [setLauncherOpen, setSidebarOpen]);
  const openAccountPreferences = useCallback(() => {
    setQuickSettings(true);
    setSidebarOpen(true);
  }, [setQuickSettings, setSidebarOpen]);
  const toggleQuickSettings = useCallback(() => {
    setQuickSettings(previous => !previous);
  }, [setQuickSettings]);
  const toggleSidebarCollapsed = useCallback(() => {
    setSidebarCollapsed(previous => !previous);
  }, [setSidebarCollapsed]);
  const togglePin = useCallback((path: string) => {
    setPinned(previous => {
      if (previous.includes(path)) return previous.filter(item => item !== path);
      return [...previous, path];
    });
  }, [setPinned]);

  return {
    openSidebar, closeSidebar, openCatalog, openAccountPreferences,
    toggleQuickSettings, toggleSidebarCollapsed, togglePin,
  };
};

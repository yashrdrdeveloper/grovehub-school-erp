import { useCallback } from 'react';

import { writePinnedApps, writeSidebarCollapsed } from '../../storage';
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
    setSidebarCollapsed(previous => {
      const next = !previous;
      writeSidebarCollapsed(next);
      return next;
    });
  }, [setSidebarCollapsed]);
  const togglePin = useCallback((path: string) => {
    setPinned(previous => {
      const next = previous.includes(path)
        ? previous.filter(item => item !== path)
        : [...previous, path];
      writePinnedApps(next);
      return next;
    });
  }, [setPinned]);

  return {
    openSidebar, closeSidebar, openCatalog, openAccountPreferences,
    toggleQuickSettings, toggleSidebarCollapsed, togglePin,
  };
};

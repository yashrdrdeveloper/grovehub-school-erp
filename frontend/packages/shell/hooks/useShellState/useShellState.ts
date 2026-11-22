import { useEffect, useState } from 'react';

import { readPinnedApps, readSidebarCollapsed, writePinnedApps, writeSidebarCollapsed } from '../../storage';
import type { UseShellState } from './types/useShellState';

/** Owns the workspace chrome's state; persistence belongs to storage/actions. */
export const useShellState: UseShellState = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(readSidebarCollapsed);
  const [quickSettings, setQuickSettings] = useState(false);
  const [launcherOpen, setLauncherOpen] = useState(false);
  const [pinned, setPinned] = useState<string[]>(readPinnedApps);

  // Persist committed state, not speculative functional state updates.
  useEffect(() => { writeSidebarCollapsed(sidebarCollapsed); }, [sidebarCollapsed]);
  useEffect(() => { writePinnedApps(pinned); }, [pinned]);

  return {
    sidebarOpen, sidebarCollapsed, quickSettings, launcherOpen, pinned,
    setSidebarOpen, setSidebarCollapsed, setQuickSettings, setLauncherOpen, setPinned,
  };
};

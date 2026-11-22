import { appCatalog, defaultPinnedApps } from '../registry/appCatalog.registry';
import type {
  ReadPinnedApps, WritePinnedApps, ReadSidebarCollapsed, WriteSidebarCollapsed,
} from './types/workspacePreferences';

const PINNED_KEY = 'groovehub.workspace.pinned.v1';
const COLLAPSED_KEY = 'groovehub.workspace.sidebar.collapsed.v1';

/** Browser state is optional: private browsing and SSR must not break navigation. */
export const readPinnedApps: ReadPinnedApps = () => {
  try {
    const raw = window.localStorage.getItem(PINNED_KEY);
    if (!raw) return [...defaultPinnedApps];
    const parsed: unknown = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [...defaultPinnedApps];
    const known = new Set(appCatalog.filter(app => app.path).map(app => app.path));
    return [...new Set(parsed.filter((item): item is string =>
      typeof item === 'string' && known.has(item)))];
  } catch {
    return [...defaultPinnedApps];
  }
};

export const writePinnedApps: WritePinnedApps = (paths) => {
  try { window.localStorage.setItem(PINNED_KEY, JSON.stringify(paths)); }
  catch { /* Browser storage is optional. */ }
};

export const readSidebarCollapsed: ReadSidebarCollapsed = () => {
  try { return window.localStorage.getItem(COLLAPSED_KEY) === 'true'; }
  catch { return false; }
};

export const writeSidebarCollapsed: WriteSidebarCollapsed = (collapsed) => {
  try { window.localStorage.setItem(COLLAPSED_KEY, String(collapsed)); }
  catch { /* Browser storage is optional. */ }
};

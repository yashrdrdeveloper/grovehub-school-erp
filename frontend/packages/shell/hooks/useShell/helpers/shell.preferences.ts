import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';
import { appCatalog, defaultPinnedApps } from '../../../constants/shell.catalog';
import { SHELL_STORAGE_KEYS } from '../constants/shell.storageKeys';

import type {
  ReadPinnedApps, WritePinnedApps, ReadSidebarCollapsed, WriteSidebarCollapsed,
} from './types/shell.preferences';



/** Browser state is optional: private browsing and SSR must not break navigation. */
export const readPinnedApps: ReadPinnedApps = () => {
  try {
    const raw = window.localStorage.getItem(SHELL_STORAGE_KEYS.PINNED);
    if (!raw) return [...defaultPinnedApps];
    const parsed: unknown = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [...defaultPinnedApps];
    const known = new Set(_map(_filter(appCatalog, (app) => Boolean(app.path)), (app) => app.path));
    const valid = _filter(parsed as unknown[], (item): item is string =>
      typeof item === 'string' && known.has(item));
    return [...new Set(valid)];
  } catch {
    return [...defaultPinnedApps];
  }
};

export const writePinnedApps: WritePinnedApps = (paths) => {
  try { window.localStorage.setItem(SHELL_STORAGE_KEYS.PINNED, JSON.stringify(paths)); }
  catch { /* Browser storage is optional. */ }
};

export const readSidebarCollapsed: ReadSidebarCollapsed = () => {
  try { return window.localStorage.getItem(SHELL_STORAGE_KEYS.COLLAPSED) === 'true'; }
  catch { return false; }
};

export const writeSidebarCollapsed: WriteSidebarCollapsed = (collapsed) => {
  try { window.localStorage.setItem(SHELL_STORAGE_KEYS.COLLAPSED, String(collapsed)); }
  catch { /* Browser storage is optional. */ }
};

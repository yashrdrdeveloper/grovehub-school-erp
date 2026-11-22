import { _filter } from '@Base/utilities/filter';
import { _trim } from '@Base/utilities/trim';
import { _includes } from '@Base/utilities/includes';
import { appCatalog } from '../../../constants/shell.catalog';
import type { FilterCatalog } from './types/appLauncher.filterCatalog';

/** Search only implemented apps and derive disjoint pinned/available groups. */
export const filterCatalog: FilterCatalog = (search, pinned) => {
  const query = _trim(search).toLowerCase();
  const matching = _filter(appCatalog, app => {
    const haystack = (app.name + ' ' + app.description).toLowerCase();
    return haystack.indexOf(query) !== -1;
  });
  const pinnedApps = _filter(matching, app => _includes(pinned, app.path));
  const otherApps = _filter(matching, app => !_includes(pinned, app.path));
  return { matching, pinnedApps, otherApps };
};

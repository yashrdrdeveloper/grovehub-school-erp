import { filterCatalog } from '../appLauncher.filterCatalog';

describe('filterCatalog', () => {
  it('shows implemented dashboard and preserves its pinned status', () => {
    const result = filterCatalog(' DASHBOARD ', ['/dashboard']);
    expect(result.matching.map(item => item.path)).toEqual(['/dashboard']);
    expect(result.pinnedApps).toHaveLength(1);
    expect(result.otherApps).toHaveLength(0);
  });
  it('filters without mutating the catalog or pinned selection', () => {
    expect(filterCatalog('no such app', ['/dashboard']).matching).toEqual([]);
    expect(filterCatalog('', []).otherApps[0]?.path).toBe('/dashboard');
  });
});

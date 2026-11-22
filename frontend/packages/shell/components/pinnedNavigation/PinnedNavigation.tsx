import { appCatalog } from '../../registry/appCatalog.registry';
import { PinnedApp } from './PinnedApp';
import type { PinnedNavigationComponent, PinnedAppEntry } from './types/pinnedNavigation';

/** Pure sidebar navigation: no pin state or storage writes. */
export const PinnedNavigation: PinnedNavigationComponent = (props) => {
  const entries: PinnedAppEntry[] = props.pinned.flatMap(path => {
    const app = appCatalog.find(item => item.path === path);
    return app && app.path ? [{ app, path }] : [];
  });
  const nodes = entries.map(entry => <PinnedApp key={entry.path} entry={entry} {...props}/>);

  return (
    <nav aria-label="Pinned applications" className="sidebar__nav">
      {nodes.length > 0 ? nodes : <p className="sidebar__empty">Pin apps from the catalog to see them here.</p>}
    </nav>
  );
};

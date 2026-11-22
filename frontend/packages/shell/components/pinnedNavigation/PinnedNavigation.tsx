import { _find } from '@Base/utilities/find';
import { _flatMap } from '@Base/utilities/flatMap';
import { _map } from '@Base/utilities/map';
import { Box } from '@Components/components/Box';
import type { ReactNode } from 'react';
import { appCatalog } from '../../constants/shell.catalog';
import { PinnedApp } from './PinnedApp';
import type { PinnedNavigationComponent, PinnedAppEntry } from './types/pinnedNavigation';

/** Pure sidebar navigation: no pin state or storage writes. */
export const PinnedNavigation: PinnedNavigationComponent = (props) => {
  const entries: PinnedAppEntry[] = _flatMap(props.pinned, path => {
    const app = _find(appCatalog, item => item.path === path);
    if (app && app.path) return [{ app, path }];
    return [];
  });
  const nodes = _map(entries, entry => <PinnedApp key={entry.path} entry={entry} {...props}/>);
  let content: ReactNode = <Box as="p" className="sidebar__empty">Pin apps from the catalog to see them here.</Box>;
  if (nodes.length > 0) content = nodes;

  return (
    <Box as="nav" aria-label="Pinned applications" className="sidebar__nav">
      {content}
    </Box>
  );
};

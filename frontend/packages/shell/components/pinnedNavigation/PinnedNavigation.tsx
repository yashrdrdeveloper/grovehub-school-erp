import type { ReactNode } from 'react';

import { _find } from '@Base/utilities/find';
import { _flatMap } from '@Base/utilities/flatMap';
import { _map } from '@Base/utilities/map';
import { Box, BOX_TAG } from '@Components/components/Box';

import { appCatalog } from '../../constants/shell.catalog';
import styles from './PinnedNavigation.module.css';
import { PinnedApp } from './PinnedApp';

import type { PinnedNavigationComponent, PinnedAppEntry } from './types/pinnedNavigation';

/** Pure, owner-styled pinned navigation: storage and pin mutations live in useShell. */
export const PinnedNavigation: PinnedNavigationComponent = (props) => {
  const entries: PinnedAppEntry[] = _flatMap(props.pinned, (path) => {
    const app = _find(appCatalog, (item) => item.path === path);

    if (app && app.path) {
      return [{ app, path }];
    }

    return [];
  });

  const nodes = _map(entries, (entry) => <PinnedApp key={entry.path} entry={entry} {...props} />);

  const content: ReactNode = (nodes.length > 0 && nodes) || (
    <Box as={BOX_TAG.P} className={styles.empty}>
      {__('Pin apps from the catalog to see them here.')}
    </Box>
  );

  return (
    <Box as={BOX_TAG.NAV} aria-label={__('Pinned applications')} className={styles.nav}>
      {content}
    </Box>
  );
};

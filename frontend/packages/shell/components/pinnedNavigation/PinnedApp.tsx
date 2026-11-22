import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { PinOff } from '@Components/icons';
import { AppNavLink } from '@Components/components/AppNavLink';

import styles from './PinnedNavigation.module.css';
import { getPinnedAppClassName } from './helpers/pinnedNavigation.getAppClassName';

import type { PinnedAppProps } from './types/pinnedNavigation';

/** One accessible pinned application and its independent unpin action. */
export const PinnedApp = ({ entry, closeSidebar, togglePin }: PinnedAppProps) => {
  const Icon = entry.app.icon;

  const handleTogglePin = () => {
    togglePin(entry.path);
  };

  return (
    <Box as="div" className={styles.row}>
      <AppNavLink
        to={entry.path}
        title={entry.app.name}
        aria-label={entry.app.name}
        onClick={closeSidebar}
        className={getPinnedAppClassName}
      >
        <Icon size={19} strokeWidth={1.8} aria-hidden="true" />
        <Box as="span" className={styles.label}>
          {entry.app.name}
        </Box>
      </AppNavLink>

      <Button
        unstyled
        type="button"
        className={styles.unpin}
        title={`Unpin ${entry.app.name}`}
        aria-label={`Unpin ${entry.app.name}`}
        onClick={handleTogglePin}
      >
        <PinOff size={14} aria-hidden="true" />
      </Button>
    </Box>
  );
};

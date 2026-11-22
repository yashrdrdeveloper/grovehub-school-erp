import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { PinOff } from '@Components/icons';
import { NavLink } from 'react-router-dom';

import { getShellNavClassName } from '../../helpers/shell.getNavClassName';
import type { PinnedAppProps } from './types/pinnedNavigation';

/** One pinned entry, with a named pin action and route navigation handler. */
export const PinnedApp = ({ entry, closeSidebar, togglePin }: PinnedAppProps) => {
  const Icon = entry.app.icon;
  const handleTogglePin = () => togglePin(entry.path);
  return (
    <Box as="div" className="sidebar__pinned-row">
      <NavLink to={entry.path} title={entry.app.name} aria-label={entry.app.name}
        onClick={closeSidebar}
        className={getShellNavClassName}>
        <Icon size={19} strokeWidth={1.8}/><Box as="span" className="sidebar__label">{entry.app.name}</Box>
      </NavLink>
      <Button unstyled type="button" className="sidebar__unpin" title={`Unpin ${entry.app.name}`}
        aria-label={`Unpin ${entry.app.name}`} onClick={handleTogglePin}>
        <PinOff size={14}/>
      </Button>
    </Box>
  );
};

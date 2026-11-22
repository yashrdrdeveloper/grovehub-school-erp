import { PinOff } from 'lucide-react';
import { NavLink } from 'react-router-dom';

import type { PinnedAppProps } from './types/pinnedNavigation';

/** One pinned entry, with a named pin action and route navigation handler. */
export const PinnedApp = ({ entry, closeSidebar, togglePin }: PinnedAppProps) => {
  const Icon = entry.app.icon;
  const handleTogglePin = () => togglePin(entry.path);
  return (
    <div className="sidebar__pinned-row">
      <NavLink to={entry.path} title={entry.app.name} aria-label={entry.app.name}
        onClick={closeSidebar}
        className={({ isActive }) => `nav-item${isActive ? ' nav-item--active' : ''}`}>
        <Icon size={19} strokeWidth={1.8}/><span className="sidebar__label">{entry.app.name}</span>
      </NavLink>
      <button type="button" className="sidebar__unpin" title={`Unpin ${entry.app.name}`}
        aria-label={`Unpin ${entry.app.name}`} onClick={handleTogglePin}>
        <PinOff size={14}/>
      </button>
    </div>
  );
};

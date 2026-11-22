import { ChevronsLeft, ChevronsRight, Grid2X2, Settings2, X } from '@Components/icons';
import { NavLink } from 'react-router-dom';

import { SchoolMonogram } from '@Components';
import { APP_ROUTE } from '@Business/app/routes';

import { PinnedNavigation } from '../pinnedNavigation';
import { ThemeChoices } from '../themeChoices';
import type { SidebarComponent } from './types/sidebar';

/** Stateless sidebar. Interaction and persistence come from the shell model. */
export const Sidebar: SidebarComponent = ({ accountLabel, schoolLabel, onLogout, model }) => {
  const collapsedLabel = model.sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar';
  const CollapseIcon = model.sidebarCollapsed ? ChevronsRight : ChevronsLeft;

  return (
    <aside className={`sidebar${model.sidebarOpen ? ' sidebar--open' : ''}`} aria-label="School workspace sidebar">
      <div className="sidebar__header">
        <button type="button" className="icon-button app-switcher-button" title="Open app catalog"
          aria-label="Open app catalog" onClick={model.openCatalog}><Grid2X2 size={19}/></button>
        <NavLink to={APP_ROUTE.DASHBOARD} className="sidebar__brand" onClick={model.closeSidebar} aria-label="GrooveHub home">
          <span className="brand-mark brand-mark--small">G</span><span className="sidebar__brand-name">groovehub<span className="brand-period">.</span></span>
        </NavLink>
        <button type="button" className="icon-button sidebar__close" aria-label="Close navigation"
          onClick={model.closeSidebar}><X size={20}/></button>
      </div>
      <div className="sidebar__school" title={schoolLabel}>
        <SchoolMonogram/>
        <span className="sidebar__school-details"><strong>{schoolLabel}</strong><small>School operating system</small></span>
      </div>
      <div className="sidebar__section-label">PINNED APPS</div>
      <PinnedNavigation pinned={model.pinned} closeSidebar={model.closeSidebar} togglePin={model.togglePin}/>
      <button type="button" className="sidebar__catalog" title="Browse all applications" onClick={model.openCatalog}>
        <Grid2X2 size={18}/><span className="sidebar__label">Browse all apps</span>
      </button>
      <div className="sidebar__bottom">
        <button type="button" className="sidebar__collapse" onClick={model.toggleSidebarCollapsed}
          title={model.sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'} aria-label={collapsedLabel}>
          <CollapseIcon size={17}/>
          <span className="sidebar__label">{model.sidebarCollapsed ? 'Expand' : 'Collapse sidebar'}</span>
        </button>
        <button type="button" className="sidebar__settings" onClick={model.toggleQuickSettings}
          aria-expanded={model.quickSettings} title="Account and preferences">
          <span className="profile-avatar">{accountLabel.charAt(0).toUpperCase()}</span>
          <span className="sidebar__profile-details"><strong>{accountLabel}</strong><small>Account & preferences</small></span>
          <Settings2 size={17} className="sidebar__settings-icon"/>
        </button>
        {model.quickSettings && (
          <div className="settings-popover">
            <strong>Appearance</strong>
            <ThemeChoices preference={model.preference} onPreferenceChange={model.setPreference}/>
            <button type="button" onClick={onLogout}>Sign out</button>
          </div>
        )}
      </div>
    </aside>
  );
};

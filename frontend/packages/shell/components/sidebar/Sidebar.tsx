import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { ChevronsLeft, ChevronsRight, Grid2X2, Settings2, X } from '@Components/icons';
import { NavLink } from 'react-router-dom';

import { SchoolMonogram } from '@Components';
import { APP_ROUTE } from '@Business/app/routes';

import { PinnedNavigation } from '../pinnedNavigation';
import { ThemeChoices } from '../themeChoices';
import type { SidebarComponent } from './types/sidebar';

/** Stateless sidebar. Interaction and persistence come from the shell model. */
export const Sidebar: SidebarComponent = ({ accountLabel, schoolLabel, onLogout, model }) => {
  let collapsedLabel = 'Collapse sidebar';
  let collapsedShortLabel = 'Collapse sidebar';
  let CollapseIcon = ChevronsLeft;
  let sidebarClass = 'sidebar';
  if (model.sidebarCollapsed) {
    collapsedLabel = 'Expand sidebar';
    collapsedShortLabel = 'Expand';
    CollapseIcon = ChevronsRight;
  }
  if (model.sidebarOpen) sidebarClass += ' sidebar--open';

  return (
    <Box as="aside" className={sidebarClass} aria-label="School workspace sidebar">
      <Box as="div" className="sidebar__header">
        <Button unstyled type="button" className="icon-button app-switcher-button" title="Open app catalog"
          aria-label="Open app catalog" onClick={model.openCatalog}><Grid2X2 size={19}/></Button>
        <NavLink to={APP_ROUTE.DASHBOARD} className="sidebar__brand" onClick={model.closeSidebar} aria-label="GrooveHub home">
          <Box as="span" className="brand-mark brand-mark--small">G</Box><Box as="span" className="sidebar__brand-name">groovehub<Box as="span" className="brand-period">.</Box></Box>
        </NavLink>
        <Button unstyled type="button" className="icon-button sidebar__close" aria-label="Close navigation"
          onClick={model.closeSidebar}><X size={20}/></Button>
      </Box>
      <Box as="div" className="sidebar__school" title={schoolLabel}>
        <SchoolMonogram/>
        <Box as="span" className="sidebar__school-details"><Box as="strong">{schoolLabel}</Box><Box as="small">School operating system</Box></Box>
      </Box>
      <Box as="div" className="sidebar__section-label">PINNED APPS</Box>
      <PinnedNavigation pinned={model.pinned} closeSidebar={model.closeSidebar} togglePin={model.togglePin}/>
      <Button unstyled type="button" className="sidebar__catalog" title="Browse all applications" onClick={model.openCatalog}>
        <Grid2X2 size={18}/><Box as="span" className="sidebar__label">Browse all apps</Box>
      </Button>
      <Box as="div" className="sidebar__bottom">
        <Button unstyled type="button" className="sidebar__collapse" onClick={model.toggleSidebarCollapsed}
          title={collapsedLabel} aria-label={collapsedLabel}>
          <CollapseIcon size={17}/>
          <Box as="span" className="sidebar__label">{collapsedShortLabel}</Box>
        </Button>
        <Button unstyled type="button" className="sidebar__settings" onClick={model.toggleQuickSettings}
          aria-expanded={model.quickSettings} title="Account and preferences">
          <Box as="span" className="profile-avatar">{accountLabel.charAt(0).toUpperCase()}</Box>
          <Box as="span" className="sidebar__profile-details"><Box as="strong">{accountLabel}</Box><Box as="small">Account & preferences</Box></Box>
          <Settings2 size={17} className="sidebar__settings-icon"/>
        </Button>
        {model.quickSettings && (
          <Box as="div" className="settings-popover">
            <Box as="strong">Appearance</Box>
            <ThemeChoices preference={model.preference} onPreferenceChange={model.setPreference}/>
            <Button unstyled type="button" onClick={onLogout}>Sign out</Button>
          </Box>
        )}
      </Box>
    </Box>
  );
};

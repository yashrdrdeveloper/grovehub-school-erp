import { Button } from '@Components/components/Button';
import { Box, BOX_TAG } from '@Components/components/Box';
import { ChevronsLeft, ChevronsRight, Grid2X2, Settings2, X } from '@Components/icons';
import { AppNavLink } from '@Components/components/AppNavLink';

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
    <Box as={BOX_TAG.ASIDE} className={sidebarClass} aria-label="School workspace sidebar">
      <Box as={BOX_TAG.DIV} className="sidebar__header">
        <Button
          unstyled
          type="button"
          className="icon-button app-switcher-button"
          title="Open app catalog"
          aria-label="Open app catalog"
          onClick={model.openCatalog}
        >
          <Grid2X2 size={19} />
        </Button>
        <AppNavLink
          to={APP_ROUTE.DASHBOARD}
          className="sidebar__brand"
          onClick={model.closeSidebar}
          aria-label="GrooveHub home"
        >
          <Box as={BOX_TAG.SPAN} className="brand-mark brand-mark--small">
            G
          </Box>
          <Box as={BOX_TAG.SPAN} className="sidebar__brand-name">
            groovehub
            <Box as={BOX_TAG.SPAN} className="brand-period">
              .
            </Box>
          </Box>
        </AppNavLink>
        <Button
          unstyled
          type="button"
          className="icon-button sidebar__close"
          aria-label="Close navigation"
          onClick={model.closeSidebar}
        >
          <X size={20} />
        </Button>
      </Box>
      <Box as={BOX_TAG.DIV} className="sidebar__school" title={schoolLabel}>
        <SchoolMonogram />
        <Box as={BOX_TAG.SPAN} className="sidebar__school-details">
          <Box as={BOX_TAG.STRONG}>{schoolLabel}</Box>
          <Box as={BOX_TAG.SMALL}>School operating system</Box>
        </Box>
      </Box>
      <Box as={BOX_TAG.DIV} className="sidebar__section-label">
        PINNED APPS
      </Box>
      <PinnedNavigation
        pinned={model.pinned}
        closeSidebar={model.closeSidebar}
        togglePin={model.togglePin}
      />
      <Button
        unstyled
        type="button"
        className="sidebar__catalog"
        title="Browse all applications"
        onClick={model.openCatalog}
      >
        <Grid2X2 size={18} />
        <Box as={BOX_TAG.SPAN} className="sidebar__label">
          Browse all apps
        </Box>
      </Button>
      <Box as={BOX_TAG.DIV} className="sidebar__bottom">
        <Button
          unstyled
          type="button"
          className="sidebar__collapse"
          onClick={model.toggleSidebarCollapsed}
          title={collapsedLabel}
          aria-label={collapsedLabel}
        >
          <CollapseIcon size={17} />
          <Box as={BOX_TAG.SPAN} className="sidebar__label">
            {collapsedShortLabel}
          </Box>
        </Button>
        <Button
          unstyled
          type="button"
          className="sidebar__settings"
          onClick={model.toggleQuickSettings}
          aria-expanded={model.quickSettings}
          title="Account and preferences"
        >
          <Box as={BOX_TAG.SPAN} className="profile-avatar">
            {accountLabel.charAt(0).toUpperCase()}
          </Box>
          <Box as={BOX_TAG.SPAN} className="sidebar__profile-details">
            <Box as={BOX_TAG.STRONG}>{accountLabel}</Box>
            <Box as={BOX_TAG.SMALL}>Account & preferences</Box>
          </Box>
          <Settings2 size={17} className="sidebar__settings-icon" />
        </Button>
        {model.quickSettings && (
          <Box as={BOX_TAG.DIV} className="settings-popover">
            <Box as={BOX_TAG.STRONG}>Appearance</Box>
            <ThemeChoices preference={model.preference} onPreferenceChange={model.setPreference} />
            <Button unstyled type="button" onClick={onLogout}>
              Sign out
            </Button>
          </Box>
        )}
      </Box>
    </Box>
  );
};

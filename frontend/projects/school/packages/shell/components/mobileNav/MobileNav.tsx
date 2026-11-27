import { APP_ROUTE } from '@Business/app/routes';
import { Box, BOX_TAG } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { Grid2X2, LayoutDashboard, Menu } from '@Components/icons';
import { AppNavLink } from '@Components/components/AppNavLink';

import styles from './MobileNav.module.css';
import { getMobileNavClassName } from './helpers/mobileNav.getClassName';

import type { MobileNavComponent } from './types/mobileNav';

/** Responsive workspace navigation using owner-scoped CSS Module styles. */
export const MobileNav: MobileNavComponent = ({ model }) => (
  <Box as={BOX_TAG.NAV} className={styles.root} aria-label={__('Mobile navigation')}>
    <AppNavLink to={APP_ROUTE.DASHBOARD} className={getMobileNavClassName}>
      <LayoutDashboard size={19} aria-hidden="true" />
      <Box as={BOX_TAG.SPAN}>{__('Home')}</Box>
    </AppNavLink>
    <Button unstyled type="button" className={styles.item} onClick={model.openCatalog}>
      <Grid2X2 size={19} aria-hidden="true" />
      <Box as={BOX_TAG.SPAN}>{__('Apps')}</Box>
    </Button>
    <Button
      unstyled
      type="button"
      id="workspace-mobile-menu"
      className={styles.item}
      aria-controls="workspace-sidebar"
      aria-expanded={model.sidebarOpen}
      onClick={model.openSidebar}
    >
      <Menu size={19} aria-hidden="true" />
      <Box as={BOX_TAG.SPAN}>{__('Menu')}</Box>
    </Button>
  </Box>
);

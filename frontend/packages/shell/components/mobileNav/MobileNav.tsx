import { APP_ROUTE } from '@Business/app/routes';
import { Box } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { Grid2X2, LayoutDashboard, Menu } from '@Components/icons';
import { AppNavLink } from '@Components/components/AppNavLink';

import styles from './MobileNav.module.css';
import { getMobileNavClassName } from './helpers/mobileNav.getClassName';

import type { MobileNavComponent } from './types/mobileNav';

/** Responsive workspace navigation using owner-scoped CSS Module styles. */
export const MobileNav: MobileNavComponent = ({ model }) => (
  <Box as="nav" className={styles.root} aria-label="Mobile navigation">
    <AppNavLink to={APP_ROUTE.DASHBOARD} className={getMobileNavClassName}>
      <LayoutDashboard size={19} aria-hidden="true" />
      <Box as="span">Home</Box>
    </AppNavLink>
    <Button unstyled type="button" className={styles.item} onClick={model.openCatalog}>
      <Grid2X2 size={19} aria-hidden="true" />
      <Box as="span">Apps</Box>
    </Button>
    <Button unstyled type="button" className={styles.item} onClick={model.openSidebar}>
      <Menu size={19} aria-hidden="true" />
      <Box as="span">Menu</Box>
    </Button>
  </Box>
);

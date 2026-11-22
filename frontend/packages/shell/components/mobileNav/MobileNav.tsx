import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { Grid2X2, LayoutDashboard, Menu } from '@Components/icons';
import { NavLink } from 'react-router-dom';
import { APP_ROUTE } from '@Business/app/routes';
import { getShellNavClassName } from '../../helpers/shell.getNavClassName';
import type { MobileNavComponent } from './types/mobileNav';

/** Same mobile navigation semantics with no inline action callbacks. */
export const MobileNav: MobileNavComponent = ({ model }) => (
  <Box as="nav" className="mobile-nav" aria-label="Mobile navigation">
    <NavLink to={APP_ROUTE.DASHBOARD} className={getShellNavClassName}>
      <LayoutDashboard size={19}/><Box as="span">Home</Box>
    </NavLink>
    <Button unstyled type="button" className="nav-item" onClick={model.openCatalog}>
      <Grid2X2 size={19}/><Box as="span">Apps</Box>
    </Button>
    <Button unstyled type="button" className="nav-item" onClick={model.openSidebar}>
      <Menu size={19}/><Box as="span">Menu</Box>
    </Button>
  </Box>
);

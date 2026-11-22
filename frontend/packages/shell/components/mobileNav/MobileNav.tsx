import { Grid2X2, LayoutDashboard, Menu } from '@Components/icons';
import { NavLink } from 'react-router-dom';
import type { MobileNavComponent } from './types/mobileNav';

const homeClassName = ({ isActive }: { isActive: boolean }) =>
  `nav-item${isActive ? ' nav-item--active' : ''}`;

/** Same mobile navigation semantics with no inline action callbacks. */
export const MobileNav: MobileNavComponent = ({ model }) => (
  <nav className="mobile-nav" aria-label="Mobile navigation">
    <NavLink to="/dashboard" className={homeClassName}>
      <LayoutDashboard size={19}/><span>Home</span>
    </NavLink>
    <button type="button" className="nav-item" onClick={model.openCatalog}>
      <Grid2X2 size={19}/><span>Apps</span>
    </button>
    <button type="button" className="nav-item" onClick={model.openSidebar}>
      <Menu size={19}/><span>Menu</span>
    </button>
  </nav>
);

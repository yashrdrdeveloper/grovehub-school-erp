import { Box } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { Outlet } from '@Components/router';

import { AppLauncherContainer } from '../../components/appLauncher';
import { MobileNav } from '../../components/mobileNav';
import { Sidebar } from '../../components/sidebar';
import { Topbar } from '../../components/topbar';
import { SHELL_TEST_IDS } from '../../constants/shell.testIds';
import { getShellClassName } from '../../helpers/shell.getClassName';

import type { WorkspaceShellComponent } from './types/workspaceShell.props';

/** Pure workspace template: state, persistence and provider logic belong in containers/hooks. */
export const WorkspaceShell: WorkspaceShellComponent = ({ model, ...props }) => (
  <Box
    as="div"
    className={getShellClassName(model.sidebarCollapsed)}
    dataTestId={SHELL_TEST_IDS.ROOT}
  >
    <Sidebar {...props} model={model} />

    {model.sidebarOpen && (
      <Button
        unstyled
        type="button"
        className="sidebar-backdrop"
        aria-label="Dismiss navigation"
        onClick={model.closeSidebar}
      />
    )}

    <Box as="div" className="shell__content">
      <Topbar schoolLabel={props.schoolLabel} accountLabel={props.accountLabel} model={model} />
      <Box as="main" className="shell__main">
        <Outlet />
      </Box>
      <MobileNav model={model} />
    </Box>

    <AppLauncherContainer
      open={model.launcherOpen}
      onOpenChange={model.setLauncherOpen}
      pinned={model.pinned}
      onTogglePin={model.togglePin}
      sidebarCollapsed={model.sidebarCollapsed && !model.sidebarOpen}
    />
  </Box>
);

import { Outlet } from 'react-router-dom';
import { SHELL_TEST_IDS } from './constants/shell.testIds';

import { AppLauncherContainer } from './components/appLauncher';
import { Sidebar } from './components/sidebar';
import { Topbar } from './components/topbar';
import { MobileNav } from './components/mobileNav';
import { useShell } from './hooks/useShell';
import { getShellClassName } from './helpers/shell.getClassName';
import type { ShellComponent } from './types/shell.props';

export type { ShellProps } from './types/shell.props';

/** Stable protected layout. React Router replaces only the nested Outlet when changing domain routes. */
export const Shell: ShellComponent = (props) => {
  const model = useShell();
  return (
    <div className={getShellClassName(model.sidebarCollapsed)} data-testid={SHELL_TEST_IDS.ROOT}>
      <Sidebar {...props} model={model}/>
      {model.sidebarOpen && <button type="button" className="sidebar-backdrop"
        aria-label="Dismiss navigation" onClick={model.closeSidebar}/>}
      <div className="shell__content">
        <Topbar schoolLabel={props.schoolLabel} accountLabel={props.accountLabel} model={model}/>
        <main className="shell__main"><Outlet/></main>
        <MobileNav model={model}/>
      </div>
      <AppLauncherContainer open={model.launcherOpen} onOpenChange={model.setLauncherOpen}
        pinned={model.pinned} onTogglePin={model.togglePin}
        sidebarCollapsed={model.sidebarCollapsed && !model.sidebarOpen}/>
    </div>
  );
};

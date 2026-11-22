import { WorkspaceShell } from './templates/workspaceShell';

import type { ShellComponent } from './types/shell.props';

/** Shell is only the presentation root; its container owns state and side effects. */
export const Shell: ShellComponent = (props) => <WorkspaceShell {...props} />;

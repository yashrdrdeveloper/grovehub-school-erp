import { Shell } from '../../Shell';
import { useShell } from '../../hooks/useShell';

import type { WorkspaceShellContainerComponent } from './types/workspaceShellContainer';

/** Persistent stateful workspace owner, independent of provider bootstrap. */
export const WorkspaceShellContainer: WorkspaceShellContainerComponent = (props) => {
  const model = useShell();

  return <Shell {...props} model={model} />;
};

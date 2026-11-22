import { useWorkspaceTheme } from '@Components/hooks/theme';

import { useShellActions } from '../useShellActions';
import { useShellState } from '../useShellState';
import type { UseShell } from './types/useShell';

/** Thin composer for shell state, named actions and the shared theme controller. */
export const useShell: UseShell = () => {
  const state = useShellState();
  const actions = useShellActions(state);
  const theme = useWorkspaceTheme();

  return { ...state, ...actions, ...theme };
};

import type { ThemePreference } from '@Components/hooks/theme';
import type { ShellActions } from '../../useShellActions';
import type { ShellStateModel } from '../../useShellState';

export type ShellViewModel = ShellStateModel & ShellActions & {
  readonly preference: ThemePreference;
  readonly setPreference: (next: ThemePreference) => void;
};

export type UseShell = () => ShellViewModel;

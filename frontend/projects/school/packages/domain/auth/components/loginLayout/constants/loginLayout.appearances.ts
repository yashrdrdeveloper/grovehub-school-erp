import { Monitor, Moon, Sun } from '@Components/icons';
import type { LucideIcon } from '@Components/icons';
import { ThemePreference } from '@Components/hooks/useWorkspaceTheme';

export interface LoginLayoutAppearance {
  readonly value: ThemePreference;
  readonly label: string;
  readonly Icon: LucideIcon;
}
export const LOGIN_LAYOUT_APPEARANCES: readonly LoginLayoutAppearance[] = [
  { value: ThemePreference.SYSTEM, label: __('System'), Icon: Monitor },
  { value: ThemePreference.LIGHT, label: __('Light'), Icon: Sun },
  { value: ThemePreference.DARK, label: __('Dark'), Icon: Moon },
];

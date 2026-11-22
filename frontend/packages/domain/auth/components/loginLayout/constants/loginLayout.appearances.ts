import { Monitor, Moon, Sun } from '@Components/icons';
import type { LucideIcon } from '@Components/icons';
import type { ThemePreference } from '@Components/hooks/theme';

export interface LoginLayoutAppearance {
  readonly value: ThemePreference;
  readonly label: string;
  readonly Icon: LucideIcon;
}
export const LOGIN_LAYOUT_APPEARANCES: readonly LoginLayoutAppearance[] = [
  { value: 'system', label: 'System', Icon: Monitor },
  { value: 'light', label: 'Light', Icon: Sun },
  { value: 'dark', label: 'Dark', Icon: Moon },
];

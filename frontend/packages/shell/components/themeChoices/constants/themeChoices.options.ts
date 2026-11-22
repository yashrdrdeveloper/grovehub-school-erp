import { Monitor, Moon, Sun } from '@Components/icons';
import { ThemePreference } from '@Components/hooks/useWorkspaceTheme';

import type { LucideIcon } from '@Components/icons';

export interface ThemeChoiceOption {
  readonly value: ThemePreference;
  readonly label: string;
  readonly Icon: LucideIcon;
}

export const THEME_CHOICES: readonly ThemeChoiceOption[] = [
  { value: ThemePreference.SYSTEM, label: 'System', Icon: Monitor },
  { value: ThemePreference.LIGHT, label: 'Light', Icon: Sun },
  { value: ThemePreference.DARK, label: 'Dark', Icon: Moon },
];

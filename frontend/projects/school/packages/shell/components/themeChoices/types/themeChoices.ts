import type { ReactElement } from 'react';
import type { ThemePreference } from '@Components/hooks/useWorkspaceTheme';

export interface ThemeChoicesProps {
  readonly compact?: boolean;
  readonly preference: ThemePreference;
  readonly onPreferenceChange: (next: ThemePreference) => void;
}

export type ThemeChoicesComponent = (props: ThemeChoicesProps) => ReactElement;

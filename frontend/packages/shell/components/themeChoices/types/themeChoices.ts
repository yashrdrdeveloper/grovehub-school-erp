import type { ReactElement } from 'react';
import type { ThemePreference } from '@Components/hooks/theme';

export interface ThemeChoicesProps {
  readonly preference: ThemePreference;
  readonly onPreferenceChange: (next: ThemePreference) => void;
}

export type ThemeChoicesComponent = (props: ThemeChoicesProps) => ReactElement;

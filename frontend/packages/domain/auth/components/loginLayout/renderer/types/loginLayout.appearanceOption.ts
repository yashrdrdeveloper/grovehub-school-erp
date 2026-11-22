import type { ReactElement } from 'react';
import type { ThemePreference } from '@Components/hooks/useWorkspaceTheme';
import type { LoginLayoutAppearance } from '../../constants/loginLayout.appearances';
export interface LoginLayoutAppearanceOptionProps extends LoginLayoutAppearance {
  readonly preference: ThemePreference;
  readonly setPreference: (next: ThemePreference) => void;
}
export type RenderLoginLayoutAppearanceOption = (
  props: LoginLayoutAppearanceOptionProps,
) => ReactElement;

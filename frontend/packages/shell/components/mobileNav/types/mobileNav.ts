import type { ReactElement } from 'react';
import type { ShellViewModel } from '../../../hooks/useShell';

export interface MobileNavProps {
  readonly model: ShellViewModel;
}

export type MobileNavComponent = (props: MobileNavProps) => ReactElement;

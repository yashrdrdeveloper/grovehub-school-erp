import type { ReactElement } from 'react';

import type { ShellViewModel } from '../hooks/useShell';

export interface ShellProps {
  readonly accountLabel: string;
  readonly schoolLabel: string;
  readonly onLogout: () => void;
}

export interface ShellViewProps extends ShellProps {
  readonly model: ShellViewModel;
}

export type ShellComponent = (props: ShellViewProps) => ReactElement;

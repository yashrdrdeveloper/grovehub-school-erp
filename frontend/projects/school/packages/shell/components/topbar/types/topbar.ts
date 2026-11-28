import type { ReactElement } from 'react';
import type { ShellProps } from '../../../types/shell.props';
import type { ShellViewModel } from '../../../hooks/useShell';

export interface TopbarProps {
  readonly schoolLabel: ShellProps['schoolLabel'];
  readonly model: ShellViewModel;
}

export type TopbarComponent = (props: TopbarProps) => ReactElement;

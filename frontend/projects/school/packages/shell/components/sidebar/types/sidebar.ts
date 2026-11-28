import type { ReactElement } from 'react';
import type { ShellProps } from '../../../types/shell.props';
import type { ShellViewModel } from '../../../hooks/useShell';

export interface SidebarProps extends ShellProps {
  readonly model: ShellViewModel;
}

export type SidebarComponent = (props: SidebarProps) => ReactElement;

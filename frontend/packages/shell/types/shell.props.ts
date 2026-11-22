import type { ReactElement } from 'react';

export interface ShellProps {
  readonly accountLabel: string;
  readonly schoolLabel: string;
  readonly onLogout: () => void;
}

export type ShellComponent = (props: ShellProps) => ReactElement;

import type { TableDensity } from '../../Table/constants/table.type';

export interface ColumnManagerOption {
  readonly id: string;
  readonly label: string;
  readonly visible: boolean;
  readonly disabled?: boolean;
}

export interface ColumnManagerProps {
  readonly options: readonly ColumnManagerOption[];
  readonly density: (typeof TableDensity)[keyof typeof TableDensity];
  readonly onToggle: (id: string) => void;
  readonly onDensityChange: (density: (typeof TableDensity)[keyof typeof TableDensity]) => void;
  readonly onReset: () => void;
  readonly dataTestId?: string;
}

import type { TableDensity } from '../../Table/constants/table.type';
import type { ColumnPin } from '../../Table/types/table.column';

export interface ColumnManagerOption {
  readonly id: string;
  readonly label: string;
  readonly visible: boolean;
  readonly disabled?: boolean;
  readonly width?: number;
  readonly pin?: ColumnPin;
  readonly resizable?: boolean;
  readonly reorderable?: boolean;
}

export interface ColumnManagerProps {
  readonly options: readonly ColumnManagerOption[];
  readonly density: (typeof TableDensity)[keyof typeof TableDensity];
  readonly onToggle: (id: string) => void;
  readonly onDensityChange: (density: (typeof TableDensity)[keyof typeof TableDensity]) => void;
  readonly maxPinnedColumns?: number;
  readonly onReorder?: (ids: readonly string[]) => void;
  readonly onMove?: (id: string, offset: number) => void;
  readonly onResize?: (id: string, width: number) => void;
  readonly onPin?: (id: string, pin: ColumnPin) => void;
  readonly onReset: () => void;
  readonly dataTestId?: string;
  readonly mobileMenu?: boolean;
  readonly mobileList?: boolean;
}

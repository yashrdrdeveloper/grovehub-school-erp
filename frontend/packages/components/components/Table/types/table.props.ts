import type { LucideIcon } from '@Components/icons';
import type { FilterCondition, FiltersProps } from '../../Filters';
import type { SortCondition, SortProps } from '../../Sort';
import type { QuickFiltersProps } from '../../QuickFilters';
import type { PaginationAction, PaginationProps } from '../../Pagination';
import type { TableType, TableDataMode, TableDensity } from '../constants/table.type';
import type { TableColumn } from './table.column';

export interface TableActions<T> {
  readonly SEARCH?: (value: string) => void;
  readonly FILTER?: (value: readonly FilterCondition[]) => void;
  readonly SORT?: (value: readonly SortCondition[]) => void;
  readonly QUICK_FILTER?: (value: readonly string[]) => void;
  readonly PAGINATION?: (action: PaginationAction) => void;
  readonly HEADER_ACTION?: (id: string) => void;
  readonly ROW_ACTION?: (action: { id: string; row: T }) => void;
  readonly ROW_CLICK?: (row: T) => void;
  readonly SELECTION_CHANGE?: (ids: readonly string[]) => void;
  readonly REFRESH?: () => void;
  readonly RETRY?: () => void;
}

export interface TableHeaderAction {
  readonly id: string;
  readonly label: string;
  readonly icon?: LucideIcon;
  readonly disabled?: boolean;
}

export interface TableHeaderProps {
  readonly title?: string;
  readonly description?: string;
  readonly leftConfig?: readonly TableHeaderAction[];
  readonly rightConfig?: readonly TableHeaderAction[];
}

export interface TableRowAction<T> extends TableHeaderAction {
  readonly isVisible?: (row: T) => boolean;
  readonly destructive?: boolean;
}

export interface TableSearchProps {
  readonly value: string;
  readonly placeholder?: string;
  readonly debounceMs?: number;
}

type WithoutChange<T> = T extends unknown ? Omit<T, 'onChange'> : never;

export interface TableProps<T> {
  readonly id: string;
  readonly type: (typeof TableType)[keyof typeof TableType];
  readonly dataMode: (typeof TableDataMode)[keyof typeof TableDataMode];
  readonly data: readonly T[];
  readonly columns: readonly TableColumn<T>[];
  readonly getRowId: (row: T) => string;
  readonly actions?: TableActions<T>;
  readonly headerProps?: TableHeaderProps;
  readonly searchProps?: TableSearchProps;
  readonly filterProps?: Omit<FiltersProps, 'onChange'>;
  readonly sortProps?: Omit<SortProps, 'onChange'>;
  readonly quickFilterProps?: Omit<QuickFiltersProps, 'onChange'>;
  readonly paginationProps?: WithoutChange<PaginationProps>;
  readonly rowActions?: readonly TableRowAction<T>[];
  readonly loading?: boolean;
  readonly refreshing?: boolean;
  readonly error?: string;
  readonly emptyTitle?: string;
  readonly dataTestId?: string;
  readonly density?: (typeof TableDensity)[keyof typeof TableDensity];
}

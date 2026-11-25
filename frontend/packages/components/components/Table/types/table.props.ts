import type { ReactNode } from 'react';
import type { LucideIcon } from '@Components/icons';
import type { FilterCondition, FiltersProps } from '../../Filters';
import type { SortCondition, SortProps } from '../../Sort';
import type { QuickFiltersProps } from '../../QuickFilters';
import type { PaginationAction, PaginationProps } from '../../Pagination';
import type { TableType, TableDataMode, TableDensity } from '../constants/table.type';
import type { ColumnPin, TableColumn } from './table.column';

export type Density = (typeof TableDensity)[keyof typeof TableDensity];

export interface TablePreferences {
  readonly hidden: readonly string[];
  readonly order: readonly string[];
  readonly widths: Readonly<Record<string, number>>;
  readonly pins: Readonly<Record<string, ColumnPin>>;
  readonly density: Density;
}

export interface TablePreferencesProps {
  /** Caller must scope by tenant and account. No business data is persisted. */
  readonly storageKey?: string;
  readonly value?: TablePreferences;
  readonly onChange?: (value: TablePreferences) => void;
}

export interface TableActions<T> {
  readonly SEARCH?: (value: string) => void;
  readonly FILTER?: (value: readonly FilterCondition[]) => void;
  readonly SORT?: (value: readonly SortCondition[]) => void;
  readonly QUICK_FILTER?: (value: readonly string[]) => void;
  readonly PAGINATION?: (action: PaginationAction) => void;
  readonly HEADER_ACTION?: (id: string) => void;
  readonly SUBHEADER_ACTION?: (id: string) => void;
  readonly FOOTER_ACTION?: (id: string) => void;
  readonly ROW_ACTION?: (action: { id: string; row: T }) => void;
  readonly ROW_CLICK?: (row: T) => void;
  readonly ROW_DOUBLE_CLICK?: (row: T) => void;
  readonly ROW_EXPAND?: (action: { id: string; expanded: boolean }) => void;
  readonly GROUP_CHANGE?: (action: { id: string; expanded: boolean }) => void;
  readonly CELL_CLICK?: (action: { row: T; columnId: string }) => void;
  readonly CELL_EDIT?: (action: { row: T; columnId: string; value: unknown }) => void;
  readonly SELECTION_CHANGE?: (ids: readonly string[]) => void;
  readonly BULK_ACTION?: (action: { id: string; rowIds: readonly string[] }) => void;
  readonly COLUMN_VISIBILITY?: (hidden: readonly string[]) => void;
  readonly COLUMN_REORDER?: (order: readonly string[]) => void;
  readonly COLUMN_RESIZE?: (action: { id: string; width: number }) => void;
  readonly COLUMN_PIN?: (action: { id: string; pin: ColumnPin }) => void;
  readonly DENSITY_CHANGE?: (density: Density) => void;
  readonly REFRESH?: () => void;
  readonly RETRY?: () => void;
  readonly EXPORT?: () => void;
  readonly RESET?: () => void;
  readonly LOAD_MORE?: () => void;
  readonly FULLSCREEN_CHANGE?: (value: boolean) => void;
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
  /** Explicit searchable fields for server data; never inferred from visible columns. */
  readonly fields?: readonly string[];
}

export interface TableSelectionProps<T> {
  readonly value: readonly string[];
  readonly canSelect?: (row: T) => boolean;
}

export interface TableExpansionProps<T> {
  readonly value?: readonly string[];
  readonly render: (row: T) => ReactNode;
  readonly canExpand?: (row: T) => boolean;
}

export interface TableGroupingProps<T> {
  readonly getGroupId: (row: T) => string;
  readonly getGroupLabel: (row: T) => string;
  readonly collapsed?: readonly string[];
}

export interface TableTreeProps<T> {
  readonly getChildren: (row: T) => readonly T[];
  readonly expanded?: readonly string[];
}

export interface TableVirtualizationProps {
  readonly height?: number;
  readonly rowHeight?: number;
  readonly overscan?: number;
}

type WithoutChange<T> = T extends unknown ? Omit<T, 'onChange'> : never;

type TableStrategy<T> =
  | { readonly type: typeof TableType.DEFAULT; readonly virtualizationProps?: TableVirtualizationProps }
  | { readonly type: typeof TableType.VIRTUALIZED; readonly virtualizationProps?: TableVirtualizationProps }
  | { readonly type: typeof TableType.EXPANDABLE; readonly expansionProps: TableExpansionProps<T> }
  | { readonly type: typeof TableType.GROUPED; readonly groupingProps: TableGroupingProps<T>; readonly virtualizationProps?: TableVirtualizationProps }
  | { readonly type: typeof TableType.TREE; readonly treeProps: TableTreeProps<T>; readonly virtualizationProps?: TableVirtualizationProps };

type TableCommonProps<T> = {
  readonly id: string;
  readonly dataMode: (typeof TableDataMode)[keyof typeof TableDataMode];
  readonly data: readonly T[];
  readonly totalRecords?: number;
  readonly columns: readonly TableColumn<T>[];
  readonly getRowId: (row: T) => string;
  readonly actions?: TableActions<T>;
  readonly headerProps?: TableHeaderProps;
  readonly subheaderProps?: TableHeaderProps;
  readonly footerProps?: TableHeaderProps;
  readonly searchProps?: TableSearchProps;
  readonly filterProps?: Omit<FiltersProps, 'onChange'>;
  readonly sortProps?: Omit<SortProps, 'onChange'>;
  readonly quickFilterProps?: Omit<QuickFiltersProps, 'onChange'>;
  readonly matchesQuickFilter?: (row: T, values: readonly string[]) => boolean;
  readonly paginationProps?: WithoutChange<PaginationProps>;
  readonly rowActions?: readonly TableRowAction<T>[];
  readonly selectionProps?: TableSelectionProps<T>;
  readonly bulkActions?: readonly TableHeaderAction[];
  readonly fillHeight?: boolean;
  readonly maxPinnedColumns?: number;
  readonly preferencesProps?: TablePreferencesProps;
  readonly loading?: boolean;
  readonly refreshing?: boolean;
  readonly error?: string;
  readonly emptyTitle?: string;
  readonly dataTestId?: string;
  readonly density?: Density;
};

export type TableProps<T> = TableCommonProps<T> & TableStrategy<T>;

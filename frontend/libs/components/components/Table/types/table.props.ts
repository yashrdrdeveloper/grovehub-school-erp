import type { ReactNode } from 'react';
import type { LucideIcon } from '@Components/icons';
import type { FiltersProps } from '../../Filters';
import type { SortProps } from '../../Sort';
import type { QuickFiltersProps } from '../../QuickFilters';
import type { PaginationProps } from '../../Pagination';
import type { TableType, TableDataMode } from '../constants/table.type';
import type { TableColumn } from './table.column';
import type { Density, TablePreferencesProps } from './table.preferences';
import type { TableActions } from './table.actions';

export type { Density, TablePreferences, TablePreferencesProps } from './table.preferences';
export type { TableActions } from './table.actions';

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
  readonly mobilePlaceholder?: string;
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
  | {
      readonly type: typeof TableType.DEFAULT;
      readonly virtualizationProps?: TableVirtualizationProps;
    }
  | {
      readonly type: typeof TableType.VIRTUALIZED;
      readonly virtualizationProps?: TableVirtualizationProps;
    }
  | { readonly type: typeof TableType.EXPANDABLE; readonly expansionProps: TableExpansionProps<T> }
  | {
      readonly type: typeof TableType.GROUPED;
      readonly groupingProps: TableGroupingProps<T>;
      readonly virtualizationProps?: TableVirtualizationProps;
    }
  | {
      readonly type: typeof TableType.TREE;
      readonly treeProps: TableTreeProps<T>;
      readonly virtualizationProps?: TableVirtualizationProps;
    };

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
  readonly mobileRowProps?: {
    readonly render: (row: T) => ReactNode;
    readonly onClick?: (row: T) => void;
    readonly rowActions?: readonly TableRowAction<T>[];
  };
  /** Server-owned cursor accumulation; desktop pagination is independent. */
  readonly mobileInfiniteProps?: {
    readonly hasNext: boolean;
    readonly loading: boolean;
    readonly error?: boolean;
    readonly resetKey: string;
  };
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

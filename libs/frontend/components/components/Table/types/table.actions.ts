import type { FilterCondition } from '../../Filters';
import type { SortCondition } from '../../Sort';
import type { PaginationAction } from '../../Pagination';
import type { ColumnPin } from './table.column';
import type { Density } from './table.preferences';

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

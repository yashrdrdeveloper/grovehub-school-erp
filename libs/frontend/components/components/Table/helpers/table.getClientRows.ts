import { FilterOperator } from '../../Filters';
import { SortDirection } from '../../Sort';
import { PaginationType } from '../../Pagination';
import { TableDataMode } from '../constants/table.type';
import type { TableProps } from '../types/table.props';
import type { FilterCondition } from '../../Filters';

const normalize = (value: unknown): string => String(value ?? '').toLocaleLowerCase();

const match = (value: unknown, condition: FilterCondition): boolean => {
  const text = normalize(value);
  const values = condition.values.map(normalize);
  const first = values[0] ?? '';
  const second = values[1] ?? '';
  const present = value !== null && value !== undefined && value !== '';
  const compare = (target: string): number => {
    if (typeof value === 'number' && target !== '' && Number.isFinite(Number(target))) {
      return value - Number(target);
    }
    return text.localeCompare(target);
  };
  switch (condition.operator) {
    case FilterOperator.EQUALS:
      return text === first;
    case FilterOperator.NOT_EQUALS:
      return text !== first;
    case FilterOperator.IN:
      return values.includes(text);
    case FilterOperator.NOT_IN:
      return !values.includes(text);
    case FilterOperator.EXISTS:
      return present;
    case FilterOperator.NOT_EXISTS:
      return !present;
    case FilterOperator.CONTAINS:
      return text.includes(first);
    case FilterOperator.STARTS_WITH:
      return text.startsWith(first);
    case FilterOperator.ENDS_WITH:
      return text.endsWith(first);
    case FilterOperator.GREATER_THAN:
      return present && compare(first) > 0;
    case FilterOperator.LESS_THAN:
      return present && compare(first) < 0;
    case FilterOperator.BETWEEN:
      return present && compare(first) >= 0 && compare(second) <= 0;
    default:
      return false;
  }
};

export const getTableClientRows = <T>(props: TableProps<T>) => {
  if (props.dataMode === TableDataMode.SERVER)
    return { rows: props.data, total: props.data.length };
  const search = normalize(props.searchProps?.value).trim();
  const columns = new Map(props.columns.map((column) => [column.id, column]));
  const filtered = props.data.filter((row) => {
    if (
      search &&
      !props.columns.some(
        (column) => column.searchable !== false && normalize(column.accessor(row)).includes(search),
      )
    )
      return false;
    if (
      props.matchesQuickFilter &&
      !props.matchesQuickFilter(row, props.quickFilterProps?.value ?? [])
    )
      return false;
    return (props.filterProps?.value ?? []).every((condition) => {
      const column = columns.get(condition.field);
      return Boolean(column && match(column.accessor(row), condition));
    });
  });
  const rows = [...filtered].sort((left, right) => {
    for (const rule of props.sortProps?.value ?? []) {
      const column = props.columns.find((entry) => (entry.sortField ?? entry.id) === rule.field);
      if (!column) continue;
      const a = column.accessor(left);
      const b = column.accessor(right);
      const getComparison = () => {
        if (column.compare) return column.compare(left, right);
        if (typeof a === 'number' && typeof b === 'number') return a - b;
        return normalize(a).localeCompare(normalize(b), undefined, { numeric: true });
      };
      const order = getComparison();
      if (!order) continue;
      if (rule.direction === SortDirection.DESC) return -order;
      return order;
    }
    return 0;
  });
  if (props.paginationProps?.type === PaginationType.OFFSET) {
    const { currentPage, pageSize } = props.paginationProps;
    const start = Math.max(0, currentPage - 1) * pageSize;
    return { rows: rows.slice(start, start + pageSize), total: rows.length };
  }
  return { rows, total: rows.length };
};

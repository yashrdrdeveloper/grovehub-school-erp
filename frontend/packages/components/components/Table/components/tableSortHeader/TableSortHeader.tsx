import { Button, BUTTON_VARIANT } from '../../../Button';
import { SortDirection } from '../../../Sort';

import type { TableSortHeaderProps } from './types/tableSortHeader.props';

export const TableSortHeader = <T,>({ column, props }: TableSortHeaderProps<T>) => {
  if (!column.sortable || !props.actions?.SORT) return <>{column.header}</>;
  const field = column.sortField ?? column.id;
  const rules = props.sortProps?.value ?? [];
  const current = rules.find((entry) => entry.field === field);
  const getDirection = () => {
    if (current?.direction === SortDirection.ASC) return SortDirection.DESC;
    return SortDirection.ASC;
  };
  const getLabel = () => {
    if (current?.direction === SortDirection.ASC) return ' ↑';
    if (current?.direction === SortDirection.DESC) return ' ↓';
    return ' ↕';
  };
  const change = () => {
    const next = { field, direction: getDirection() };
    if (current) {
      props.actions?.SORT?.(rules.map((rule) => {
        if (rule.field === field) return next;
        return rule;
      }));
      return;
    }
    props.actions?.SORT?.([next, ...rules].slice(0, props.sortProps?.maxFields ?? rules.length + 1));
  };
  return <Button variant={BUTTON_VARIANT.GHOST} onClick={change}>{column.header}{getLabel()}</Button>;
};

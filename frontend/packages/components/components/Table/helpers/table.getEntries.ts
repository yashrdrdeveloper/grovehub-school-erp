import { TableType } from '../constants/table.type';
import type { TableProps } from '../types/table.props';
import type { TableEntry } from '../types/table.row';

export const getTableEntries = <T>(
  props: TableProps<T>,
  rows: readonly T[],
  expanded: readonly string[],
  collapsed: readonly string[],
): TableEntry<T>[] => {
  const result: TableEntry<T>[] = [];
  const visit = (row: T, depth: number, ancestors: ReadonlySet<string>) => {
    const id = props.getRowId(row);
    if (ancestors.has(id)) return;
    if (props.type === TableType.TREE) {
      const children = props.treeProps.getChildren(row);
      result.push({ kind: 'ROW', id, row, depth, expandable: children.length > 0 });
      if (expanded.includes(id)) {
        const next = new Set([...ancestors, id]);
        for (const child of children) visit(child, depth + 1, next);
      }
      return;
    }
    if (props.type === TableType.EXPANDABLE) {
      const expandable = props.expansionProps.canExpand?.(row) ?? true;
      result.push({ kind: 'ROW', id, row, depth, expandable });
      if (expandable && expanded.includes(id)) result.push({ kind: 'DETAIL', id, row });
      return;
    }
    result.push({ kind: 'ROW', id, row, depth, expandable: false });
  };
  if (props.type === TableType.GROUPED) {
    const groups = new Map<string, { label: string; rows: T[] }>();
    for (const row of rows) {
      const id = props.groupingProps.getGroupId(row);
      const group = groups.get(id);
      if (group) {
        group.rows.push(row);
      } else {
        groups.set(id, { label: props.groupingProps.getGroupLabel(row), rows: [row] });
      }
    }
    for (const [id, group] of groups) {
      const open = !collapsed.includes(id);
      result.push({
        kind: 'GROUP',
        id,
        label: group.label,
        count: group.rows.length,
        expanded: open,
      });
      if (open) for (const row of group.rows) visit(row, 0, new Set());
    }
    return result;
  }
  for (const row of rows) visit(row, 0, new Set());
  return result;
};

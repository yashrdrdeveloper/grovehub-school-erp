import { useState } from 'react';
import { _filter } from '@Base/utilities/filter';

import { TableType } from '../../constants/table.type';
import { getTableClientRows } from '../../helpers/table.getClientRows';
import { getTableEntries } from '../../helpers/table.getEntries';
import { useTablePreferences } from '../useTablePreferences';
import type { TableProps } from '../../types/table.props';

export const useTable = <T>(props: TableProps<T>) => {
  const { actions } = props;
  const columnPreferences = useTablePreferences(props);
  const [localExpanded, setLocalExpanded] = useState<readonly string[]>([]);
  const [localCollapsed, setLocalCollapsed] = useState<readonly string[]>([]);
  const [fullscreen, setFullscreen] = useState(false);
  const getExpanded = () => {
    if (props.type === TableType.TREE) return props.treeProps.expanded ?? localExpanded;
    if (props.type === TableType.EXPANDABLE) return props.expansionProps.value ?? localExpanded;
    return localExpanded;
  };
  const expanded = getExpanded();
  const getCollapsed = () => {
    if (props.type === TableType.GROUPED) return props.groupingProps.collapsed ?? localCollapsed;
    return localCollapsed;
  };
  const collapsed = getCollapsed();
  const client = getTableClientRows(props);
  const entries = getTableEntries(props, client.rows, expanded, collapsed);
  const selected = props.selectionProps?.value ?? [];
  const selectable = entries.flatMap((entry) => {
    if (entry.kind !== 'ROW') return [];
    if (props.selectionProps?.canSelect && !props.selectionProps.canSelect(entry.row)) return [];
    return [entry.id];
  });
  const toggleExpanded = (id: string) => {
    const next = _filter(expanded, (entry) => entry !== id);
    if (!expanded.includes(id)) next.push(id);
    setLocalExpanded(next);
    actions?.ROW_EXPAND?.({ id, expanded: next.includes(id) });
  };
  const toggleGroup = (id: string) => {
    const next = _filter(collapsed, (entry) => entry !== id);
    if (!collapsed.includes(id)) next.push(id);
    setLocalCollapsed(next);
    actions?.GROUP_CHANGE?.({ id, expanded: !next.includes(id) });
  };
  const selectRow = (id: string) => {
    if (!selectable.includes(id)) return;
    const next = _filter(selected, (entry) => entry !== id);
    if (!selected.includes(id)) next.push(id);
    actions?.SELECTION_CHANGE?.(next);
  };
  const selectPage = () => {
    if (selectable.every((id) => selected.includes(id))) {
      actions?.SELECTION_CHANGE?.(_filter(selected, (id) => !selectable.includes(id)));
      return;
    }
    actions?.SELECTION_CHANGE?.(Array.from(new Set([...selected, ...selectable])));
  };
  const changeFullscreen = (value: boolean) => {
    setFullscreen(value);
    actions?.FULLSCREEN_CHANGE?.(value);
  };
  return {
    ...columnPreferences,
    entries,
    total: client.total,
    expanded,
    collapsed,
    selected,
    selectable,
    toggleExpanded,
    toggleGroup,
    selectRow,
    selectPage,
    fullscreen,
    changeFullscreen,
  };
};

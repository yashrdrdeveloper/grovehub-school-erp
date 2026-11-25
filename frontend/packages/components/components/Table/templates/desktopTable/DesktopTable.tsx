import { Fragment, useEffect, useRef, useState } from 'react';
import { _map } from '@Base/utilities/map';

import { Box } from '../../../Box';
import { Button, BUTTON_VARIANT } from '../../../Button';

import { TableCell } from '../../components/tableCell';
import { TableRowActions } from '../../components/tableRowActions';
import { TableSelection } from '../../components/tableSelection';
import { TableSortHeader } from '../../components/tableSortHeader';
import { ColumnResizeHandle } from '../../components/columnResizeHandle';
import { TABLE_DIMENSIONS } from '../../constants/table.dimensions';
import { getTableFillWidth } from '../../helpers/table.getFillWidth';
import { getTableStretchLayout } from '../../helpers/table.getStretchLayout';
import { TablePin, TableType } from '../../constants/table.type';
import { getVirtualWindow } from '../../helpers/table.getVirtualWindow';
import { getTableCellStyle } from '../../helpers/table.getCellStyle';

import type { DesktopTableProps } from './types/desktopTable.props';
import styles from '../../Table.module.css';

export const DesktopTable = <T,>(props: DesktopTableProps<T>) => {
  const { controls } = props;
  const [scrollTop, setScrollTop] = useState(0);
  const frame = useRef<HTMLElement>(null);
  const [available, setAvailable] = useState(0);
  useEffect(() => {
    const element = frame.current;
    if (!element) return;
    const measure = () => setAvailable(element.clientWidth);
    measure();
    if (typeof ResizeObserver === 'undefined') return;
    const observer = new ResizeObserver(measure);
    observer.observe(element);
    return () => observer.disconnect();
  }, []);
  const virtual = props.type === TableType.VIRTUALIZED || ('virtualizationProps' in props && Boolean(props.virtualizationProps));
  const settings = ('virtualizationProps' in props && props.virtualizationProps) || {};
  const height = Math.max(100, settings.height ?? 480);
  const rowHeight = Math.max(44, settings.rowHeight ?? 64);
  const range = getVirtualWindow(controls.entries.length, scrollTop, height, rowHeight, settings.overscan ?? 6);
  const getEntries = () => {
    if (virtual) return controls.entries.slice(range.start, range.end);
    return controls.entries;
  };
  const selection = Boolean(props.selectionProps);
  const rowActions = Boolean(props.rowActions?.length);
  const fixedWidth = Number(selection) * TABLE_DIMENSIONS.selection + Number(rowActions) * TABLE_DIMENSIONS.actions;
  const layout = getTableStretchLayout(controls.layout, available, fixedWidth, controls.preferences.widths);
  const columnsWidth = layout.reduce((sum, item) => sum + item.width, 0);
  const fill = getTableFillWidth(available, columnsWidth, fixedWidth);
  const filler = fill > 0;
  const firstRight = layout.findIndex((entry) => entry.pin === TablePin.RIGHT);
  const fillBefore = (index: number) => filler && firstRight === index;
  const fillAtEnd = filler && firstRight < 0;

  const span = Math.max(1, layout.length + Number(selection) + Number(rowActions) + Number(filler));
  const pageChecked = Boolean(controls.selectable.length) && controls.selectable.every((id) => controls.selected.includes(id));
  const someChecked = controls.selectable.some((id) => controls.selected.includes(id));
  const tableWidth = columnsWidth + fixedWidth + fill;
  const frameStyle = () => {
    if (virtual) return { height };
    return undefined;
  };
  const rowStyle = () => {
    if (virtual) return { height: rowHeight };
    return undefined;
  };
  return (
    <Box ref={frame} className={styles.scroll} style={frameStyle()} onScroll={(event) => setScrollTop(event.currentTarget.scrollTop)}>
      <table className={styles.table} data-density={controls.preferences.density} data-virtual={virtual} style={{ width: tableWidth }} aria-busy={props.refreshing} aria-label={props.headerProps?.title ?? __('Table')}>
        <colgroup>
          {selection && <col style={{ width: TABLE_DIMENSIONS.selection }} />}
          {_map(layout, (entry, index) => (
            <Fragment key={entry.column.id}>
              {fillBefore(index) && <col style={{ width: fill }} />}
              <col style={{ width: entry.width }} />
            </Fragment>
          ))}
          {fillAtEnd && <col style={{ width: fill }} />}
          {rowActions && <col style={{ width: TABLE_DIMENSIONS.actions }} />}
        </colgroup>
        <thead>
          <tr>
            {selection && <th scope="col" className={styles.selection}><TableSelection label={__('Select visible rows')} checked={pageChecked} indeterminate={someChecked && !pageChecked} disabled={!controls.selectable.length} onChange={controls.selectPage} /></th>}
            {_map(layout, (entry, index) => (
              <Fragment key={entry.column.id}>
              {fillBefore(index) && <th aria-hidden="true" className={styles.filler} />}
              <th scope="col" style={{ ...getTableCellStyle(entry, selection, rowActions), position: 'sticky', top: 0, zIndex: 3 + Number(entry.pin !== TablePin.NONE) }} data-pinned={entry.pin !== TablePin.NONE}>
                <TableSortHeader column={entry.column} props={props} />
                {entry.column.resizable !== false && <ColumnResizeHandle label={__('Resize {{column}}', { column: entry.column.header })} width={entry.width} min={entry.column.minWidth ?? 80} max={entry.column.maxWidth ?? 1200} onResize={(width) => controls.resize(entry.column.id, width)} />}
              </th>
              </Fragment>
            ))}
            {fillAtEnd && <th aria-hidden="true" className={styles.filler} style={{ width: fill }} />}
            {rowActions && <th className={styles.actionColumn} scope="col" aria-label={__('Actions')} />}
          </tr>
        </thead>
        <tbody>
          {virtual && range.before > 0 && <tr aria-hidden="true"><td colSpan={span} className={styles.spacer} style={{ height: range.before }} /></tr>}
          {_map(getEntries(), (entry) => {
            if (entry.kind === 'GROUP') return (
              <tr key={`group:${entry.id}`} style={rowStyle()}>
                <td colSpan={span} className={styles.group}>
                  <Button variant={BUTTON_VARIANT.GHOST} aria-expanded={entry.expanded} onClick={() => controls.toggleGroup(entry.id)}>
                    {entry.expanded && '▾ '}{!entry.expanded && '▸ '}{entry.label} ({entry.count})
                  </Button>
                </td>
              </tr>
            );
            if (entry.kind === 'DETAIL' && props.type === TableType.EXPANDABLE) return <tr key={`detail:${entry.id}`}><td colSpan={span}>{props.expansionProps.render(entry.row)}</td></tr>;
            if (entry.kind !== 'ROW') return null;
            const row = entry.row;
            return (
              <tr key={`row:${entry.id}`} style={rowStyle()} data-clickable={Boolean(props.actions?.ROW_CLICK)} tabIndex={props.actions?.ROW_CLICK && 0} onClick={(event) => {
                if (!(event.target instanceof Element) || !event.currentTarget.contains(event.target) || event.target.closest('button, input, a, select, textarea, [role="button"], [role="checkbox"]')) return;
                props.actions?.ROW_CLICK?.(row);
              }} onKeyDown={(event) => {
                if (event.target !== event.currentTarget || !props.actions?.ROW_CLICK) return;
                if (event.key === 'Enter' || event.key === ' ') { event.preventDefault(); props.actions.ROW_CLICK(row); }
              }} onDoubleClick={() => props.actions?.ROW_DOUBLE_CLICK?.(row)} data-selected={controls.selected.includes(entry.id)}>
                {selection && <td className={styles.selection}><TableSelection label={__('Select row {{id}}', { id: entry.id })} checked={controls.selected.includes(entry.id)} disabled={!controls.selectable.includes(entry.id)} onChange={() => controls.selectRow(entry.id)} /></td>}
                {_map(layout, (layout, index) => (
                  <Fragment key={layout.column.id}>
                  {fillBefore(index) && <td aria-hidden="true" className={styles.filler} />}
                  <td style={getTableCellStyle(layout, selection, rowActions)} data-pinned={layout.pin !== TablePin.NONE}>
                    <Box className={styles.cell} style={{ paddingInlineStart: Number(index === 0) * entry.depth * 16 }}>
                      {index === 0 && entry.expandable && (
                        <Button variant={BUTTON_VARIANT.GHOST} aria-expanded={controls.expanded.includes(entry.id)} aria-label={__('Toggle row details')} onClick={() => controls.toggleExpanded(entry.id)}>
                          {controls.expanded.includes(entry.id) && '▾'}{!controls.expanded.includes(entry.id) && '▸'}
                        </Button>
                      )}
                      <TableCell row={row} column={layout.column} actions={props.actions} />
                    </Box>
                  </td>
                  </Fragment>
                ))}
                {fillAtEnd && <td aria-hidden="true" className={styles.filler} />}
                {rowActions && <td className={styles.actionColumn}><TableRowActions row={row} props={props} /></td>}
              </tr>
            );
          })}
          {virtual && range.after > 0 && <tr aria-hidden="true"><td colSpan={span} className={styles.spacer} style={{ height: range.after }} /></tr>}
        </tbody>
      </table>
    </Box>
  );
};

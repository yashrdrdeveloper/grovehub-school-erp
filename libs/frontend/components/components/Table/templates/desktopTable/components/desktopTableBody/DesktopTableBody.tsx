import { Fragment } from 'react';
import { _map } from '@Base/utilities/map';
import { Box } from '../../../../../Box';
import { Button, BUTTON_VARIANT } from '../../../../../Button';

import { TableCell } from '../../../../components/tableCell';
import { TableRowActions } from '../../../../components/tableRowActions';
import { TableSelection } from '../../../../components/tableSelection';
import { TablePin, TableType } from '../../../../constants/table.type';
import { getTablePinnedBoundary } from '../../../../helpers/table.getPinnedBoundary';
import { getTableCellStyle } from '../../../../helpers/table.getCellStyle';
import styles from '../../../../Table.module.css';
import type { DesktopTableBodyProps } from './types/desktopTableBody.props';

export const DesktopTableBody = <T,>({
  props,
  entries,
  layout,
  selection,
  rowActions,
  fillBefore,
  fillAtEnd,
  span,
  virtual,
  range,
  rowStyle,
}: DesktopTableBodyProps<T>) => {
  const { controls } = props;
  return (
    <tbody>
      {virtual && range.before > 0 && (
        <tr aria-hidden="true">
          <td colSpan={span} className={styles.spacer} style={{ height: range.before }} />
        </tr>
      )}
      {_map(entries, (entry) => {
        if (entry.kind === 'GROUP')
          return (
            <tr key={`group:${entry.id}`} style={rowStyle()}>
              <td colSpan={span} className={styles.group}>
                <Button
                  variant={BUTTON_VARIANT.GHOST}
                  aria-expanded={entry.expanded}
                  onClick={() => controls.toggleGroup(entry.id)}
                >
                  {entry.expanded && '▾ '}
                  {!entry.expanded && '▸ '}
                  {entry.label} ({entry.count})
                </Button>
              </td>
            </tr>
          );
        if (entry.kind === 'DETAIL' && props.type === TableType.EXPANDABLE)
          return (
            <tr key={`detail:${entry.id}`}>
              <td colSpan={span}>{props.expansionProps.render(entry.row)}</td>
            </tr>
          );
        if (entry.kind !== 'ROW') return null;
        const row = entry.row;
        return (
          <tr
            key={`row:${entry.id}`}
            style={rowStyle()}
            data-clickable={Boolean(props.actions?.ROW_CLICK)}
            tabIndex={props.actions?.ROW_CLICK && 0}
            onClick={(event) => {
              if (
                !(event.target instanceof Element) ||
                !event.currentTarget.contains(event.target) ||
                event.target.closest(
                  'button, input, a, select, textarea, [role="button"], [role="checkbox"]',
                )
              )
                return;
              props.actions?.ROW_CLICK?.(row);
            }}
            onKeyDown={(event) => {
              if (event.target !== event.currentTarget || !props.actions?.ROW_CLICK) return;
              if (event.key === 'Enter' || event.key === ' ') {
                event.preventDefault();
                props.actions.ROW_CLICK(row);
              }
            }}
            onDoubleClick={() => props.actions?.ROW_DOUBLE_CLICK?.(row)}
            data-selected={controls.selected.includes(entry.id)}
          >
            {selection && (
              <td className={styles.selection}>
                <TableSelection
                  label={__('Select row {{id}}', { id: entry.id })}
                  checked={controls.selected.includes(entry.id)}
                  disabled={!controls.selectable.includes(entry.id)}
                  onChange={() => controls.selectRow(entry.id)}
                />
              </td>
            )}
            {_map(layout, (columnLayout, index) => (
              <Fragment key={columnLayout.column.id}>
                {fillBefore(index) && <td aria-hidden="true" className={styles.filler} />}
                <td
                  style={getTableCellStyle(columnLayout, selection, rowActions)}
                  data-pinned={columnLayout.pin !== TablePin.NONE}
                  data-pin-boundary={getTablePinnedBoundary(layout, index)}
                >
                  <Box
                    className={styles.cell}
                    style={{ paddingInlineStart: Number(index === 0) * entry.depth * 16 }}
                  >
                    {index === 0 && entry.expandable && (
                      <Button
                        variant={BUTTON_VARIANT.GHOST}
                        aria-expanded={controls.expanded.includes(entry.id)}
                        aria-label={__('Toggle row details')}
                        onClick={() => controls.toggleExpanded(entry.id)}
                      >
                        {controls.expanded.includes(entry.id) && '▾'}
                        {!controls.expanded.includes(entry.id) && '▸'}
                      </Button>
                    )}
                    <TableCell row={row} column={columnLayout.column} actions={props.actions} />
                  </Box>
                </td>
              </Fragment>
            ))}
            {fillAtEnd && <td aria-hidden="true" className={styles.filler} />}
            {rowActions && (
              <td className={styles.actionColumn}>
                <TableRowActions row={row} props={props} />
              </td>
            )}
          </tr>
        );
      })}
      {virtual && range.after > 0 && (
        <tr aria-hidden="true">
          <td colSpan={span} className={styles.spacer} style={{ height: range.after }} />
        </tr>
      )}
    </tbody>
  );
};

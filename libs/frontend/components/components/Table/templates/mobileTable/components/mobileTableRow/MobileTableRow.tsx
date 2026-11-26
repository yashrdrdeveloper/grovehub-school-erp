import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';
import { Box, BOX_TAG } from '../../../../../Box';
import { Button, BUTTON_VARIANT } from '../../../../../Button';

import { TableCell } from '../../../../components/tableCell';
import { TableRowActions } from '../../../../components/tableRowActions';
import { TableSelection } from '../../../../components/tableSelection';
import styles from '../../../../Table.module.css';
import type { MobileTableRowProps } from './types/mobileTableRow.props';

export const MobileTableRow = <T,>({ props, entry }: MobileTableRowProps<T>) => {
  const { controls } = props;
  const rowClick = props.mobileRowProps?.onClick ?? props.actions?.ROW_CLICK;
  const rowActions = props.mobileRowProps?.rowActions ?? props.rowActions;
  const visible = _filter(controls.layout, (entry) => entry.column.mobileVisible !== false);
  const getVisible = () => {
    if (visible.length) return visible;
    return controls.layout.slice(0, 1);
  };
  return (
    <Box
      as={BOX_TAG.ARTICLE}
      role="listitem"
      className={styles.mobileCard}
      data-compact-row={Boolean(props.mobileRowProps)}
      data-clickable={Boolean(rowClick)}
      tabIndex={rowClick && 0}
      onClick={(event) => {
        if (
          !(event.target instanceof Element) ||
          !event.currentTarget.contains(event.target) ||
          event.target.closest(
            'button, input, a, select, textarea, [role="button"], [role="checkbox"]',
          )
        )
          return;
        rowClick?.(entry.row);
      }}
      onKeyDown={(event) => {
        if (event.target !== event.currentTarget || !rowClick) return;
        if (event.key === 'Enter' || event.key === ' ') {
          event.preventDefault();
          rowClick(entry.row);
        }
      }}
      data-selected={controls.selected.includes(entry.id)}
    >
      {props.mobileRowProps && (
        <Box className={styles.mobileRowContent}>{props.mobileRowProps.render(entry.row)}</Box>
      )}
      {Boolean(props.selectionProps || entry.expandable || rowActions?.length) && (
        <Box className={styles.mobileActions}>
          {props.selectionProps && (
            <TableSelection
              label={__('Select row {{id}}', { id: entry.id })}
              checked={controls.selected.includes(entry.id)}
              disabled={!controls.selectable.includes(entry.id)}
              onChange={() => controls.selectRow(entry.id)}
            />
          )}
          {entry.expandable && (
            <Button
              variant={BUTTON_VARIANT.GHOST}
              aria-expanded={controls.expanded.includes(entry.id)}
              onClick={() => controls.toggleExpanded(entry.id)}
            >
              {__('Details')}
            </Button>
          )}
          <TableRowActions row={entry.row} props={{ ...props, rowActions }} />
        </Box>
      )}
      {!props.mobileRowProps &&
        _map(getVisible(), (layout) => (
          <Box key={layout.column.id} className={styles.mobileField}>
            <Box className={styles.mobileLabel}>{layout.column.header}</Box>
            <TableCell row={entry.row} column={layout.column} actions={props.actions} />
          </Box>
        ))}
    </Box>
  );
};

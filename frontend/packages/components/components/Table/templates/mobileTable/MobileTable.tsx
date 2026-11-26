import { useState } from 'react';
import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';

import { Box, BOX_TAG } from '../../../Box';
import { Button, BUTTON_VARIANT } from '../../../Button';

import { TableCell } from '../../components/tableCell';
import { TableRowActions } from '../../components/tableRowActions';
import { TableSelection } from '../../components/tableSelection';
import { MobileInfiniteEnd } from '../../components/mobileInfiniteEnd';
import { TableType } from '../../constants/table.type';

import type { MobileTableProps } from './types/mobileTable.props';
import styles from '../../Table.module.css';

export const MobileTable = <T,>(props: MobileTableProps<T>) => {
  const { controls } = props;
  const rowClick = props.mobileRowProps?.onClick ?? props.actions?.ROW_CLICK;
  const rowActions = props.mobileRowProps?.rowActions ?? props.rowActions;
  const [limit, setLimit] = useState(40);
  const getEntries = () => {
    if (props.mobileInfiniteProps) return controls.entries;
    return controls.entries.slice(0, limit);
  };
  const entries = getEntries();
  const visible = _filter(controls.layout, (entry) => entry.column.mobileVisible !== false);
  const getVisible = () => {
    if (visible.length) return visible;
    return controls.layout.slice(0, 1);
  };
  const pageChecked = Boolean(controls.selectable.length) && controls.selectable.every((id) => controls.selected.includes(id));
  const someChecked = controls.selectable.some((id) => controls.selected.includes(id));
  return (
    <Box className={styles.mobileList} data-density={controls.preferences.density} role="list" aria-label={props.headerProps?.title ?? __('Table')}>
      {props.selectionProps && <Box className={styles.mobileActions}><TableSelection label={__('Select visible rows')} checked={pageChecked} indeterminate={someChecked && !pageChecked} disabled={!controls.selectable.length} onChange={controls.selectPage} /><Box>{__('Select page')}</Box></Box>}
      {_map(entries, (entry) => {
        if (entry.kind === 'GROUP') return (
          <Button key={`group:${entry.id}`} variant={BUTTON_VARIANT.OUTLINE} aria-expanded={entry.expanded} onClick={() => controls.toggleGroup(entry.id)}>{entry.expanded && '▾ '}{!entry.expanded && '▸ '}{entry.label} ({entry.count})</Button>
        );
        if (entry.kind === 'DETAIL' && props.type === TableType.EXPANDABLE) return <Box key={`detail:${entry.id}`} className={styles.mobileCard}>{props.expansionProps.render(entry.row)}</Box>;
        if (entry.kind !== 'ROW') return null;
        return (
          <Box key={`row:${entry.id}`} as={BOX_TAG.ARTICLE} role="listitem" className={styles.mobileCard} data-compact-row={Boolean(props.mobileRowProps)} data-clickable={Boolean(rowClick)} tabIndex={rowClick && 0} onClick={(event) => {
            if (!(event.target instanceof Element) || !event.currentTarget.contains(event.target) || event.target.closest('button, input, a, select, textarea, [role="button"], [role="checkbox"]')) return;
            rowClick?.(entry.row);
          }} onKeyDown={(event) => {
            if (event.target !== event.currentTarget || !rowClick) return;
            if (event.key === 'Enter' || event.key === ' ') { event.preventDefault(); rowClick(entry.row); }
          }} data-selected={controls.selected.includes(entry.id)}>
            {props.mobileRowProps && <Box className={styles.mobileRowContent}>{props.mobileRowProps.render(entry.row)}</Box>}
            {(props.selectionProps || entry.expandable || rowActions?.length) && <Box className={styles.mobileActions}>
              {props.selectionProps && <TableSelection label={__('Select row {{id}}', { id: entry.id })} checked={controls.selected.includes(entry.id)} disabled={!controls.selectable.includes(entry.id)} onChange={() => controls.selectRow(entry.id)} />}
              {entry.expandable && <Button variant={BUTTON_VARIANT.GHOST} aria-expanded={controls.expanded.includes(entry.id)} onClick={() => controls.toggleExpanded(entry.id)}>{__('Details')}</Button>}
              <TableRowActions row={entry.row} props={{ ...props, rowActions }} />
            </Box>}
            {!props.mobileRowProps && _map(getVisible(), (layout) => (
              <Box key={layout.column.id} className={styles.mobileField}>
                <Box className={styles.mobileLabel}>{layout.column.header}</Box>
                <TableCell row={entry.row} column={layout.column} actions={props.actions} />
              </Box>
            ))}
          </Box>
        );
      })}
      {!props.mobileInfiniteProps && controls.entries.length > limit && <Button variant={BUTTON_VARIANT.OUTLINE} onClick={() => setLimit((current) => current + 40)}>{__('Show more rows')}</Button>}
      {props.mobileInfiniteProps && <MobileInfiniteEnd {...props.mobileInfiniteProps} loaded={props.data.length} total={props.totalRecords} onLoadMore={props.actions?.LOAD_MORE} />}
    </Box>
  );
};

import { useState } from 'react';
import { _map } from '@Base/utilities/map';

import { Box } from '../../../Box';
import { Button, BUTTON_VARIANT } from '../../../Button';

import { TableSelection } from '../../components/tableSelection';
import { MobileInfiniteEnd } from '../../components/mobileInfiniteEnd';
import { TableType } from '../../constants/table.type';

import { MobileTableRow } from './components/mobileTableRow';
import type { MobileTableProps } from './types/mobileTable.props';
import styles from '../../Table.module.css';

export const MobileTable = <T,>(props: MobileTableProps<T>) => {
  const { controls } = props;
  const [limit, setLimit] = useState(40);
  const getEntries = () => {
    if (props.mobileInfiniteProps) return controls.entries;
    return controls.entries.slice(0, limit);
  };
  const entries = getEntries();
  const pageChecked =
    Boolean(controls.selectable.length) &&
    controls.selectable.every((id) => controls.selected.includes(id));
  const someChecked = controls.selectable.some((id) => controls.selected.includes(id));
  return (
    <Box
      className={styles.mobileList}
      data-density={controls.preferences.density}
      role="list"
      aria-label={props.headerProps?.title ?? __('Table')}
    >
      {props.selectionProps && (
        <Box className={styles.mobileActions}>
          <TableSelection
            label={__('Select visible rows')}
            checked={pageChecked}
            indeterminate={someChecked && !pageChecked}
            disabled={!controls.selectable.length}
            onChange={controls.selectPage}
          />
          <Box>{__('Select page')}</Box>
        </Box>
      )}
      {_map(entries, (entry) => {
        if (entry.kind === 'GROUP')
          return (
            <Button
              key={`group:${entry.id}`}
              variant={BUTTON_VARIANT.OUTLINE}
              aria-expanded={entry.expanded}
              onClick={() => controls.toggleGroup(entry.id)}
            >
              {entry.expanded && '▾ '}
              {!entry.expanded && '▸ '}
              {entry.label} ({entry.count})
            </Button>
          );
        if (entry.kind === 'DETAIL' && props.type === TableType.EXPANDABLE)
          return (
            <Box key={`detail:${entry.id}`} className={styles.mobileCard}>
              {props.expansionProps.render(entry.row)}
            </Box>
          );
        if (entry.kind !== 'ROW') return null;
        return <MobileTableRow key={`row:${entry.id}`} props={props} entry={entry} />;
      })}
      {!props.mobileInfiniteProps && controls.entries.length > limit && (
        <Button
          variant={BUTTON_VARIANT.OUTLINE}
          onClick={() => setLimit((current) => current + 40)}
        >
          {__('Show more rows')}
        </Button>
      )}
      {props.mobileInfiniteProps && (
        <MobileInfiniteEnd {...props.mobileInfiniteProps} onLoadMore={props.actions?.LOAD_MORE} />
      )}
    </Box>
  );
};

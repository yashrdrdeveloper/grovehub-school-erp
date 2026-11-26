import { Fragment } from 'react';
import { _map } from '@Base/utilities/map';

import { TableSelection } from '../../../../components/tableSelection';
import { TableSortHeader } from '../../../../components/tableSortHeader';
import { ColumnResizeHandle } from '../../../../components/columnResizeHandle';
import { TablePin } from '../../../../constants/table.type';
import { getTablePinnedBoundary } from '../../../../helpers/table.getPinnedBoundary';
import { getTableCellStyle } from '../../../../helpers/table.getCellStyle';
import styles from '../../../../Table.module.css';
import type { DesktopTableHeadProps } from './types/desktopTableHead.props';

export const DesktopTableHead = <T,>({
  props,
  layout,
  selection,
  rowActions,
  fill,
  fillBefore,
  fillAtEnd,
}: DesktopTableHeadProps<T>) => {
  const { controls } = props;
  const pageChecked =
    Boolean(controls.selectable.length) &&
    controls.selectable.every((id) => controls.selected.includes(id));
  const someChecked = controls.selectable.some((id) => controls.selected.includes(id));
  return (
    <thead>
      <tr>
        {selection && (
          <th scope="col" className={styles.selection}>
            <TableSelection
              label={__('Select visible rows')}
              checked={pageChecked}
              indeterminate={someChecked && !pageChecked}
              disabled={!controls.selectable.length}
              onChange={controls.selectPage}
            />
          </th>
        )}
        {_map(layout, (entry, index) => (
          <Fragment key={entry.column.id}>
            {fillBefore(index) && <th aria-hidden="true" className={styles.filler} />}
            <th
              scope="col"
              style={{
                ...getTableCellStyle(entry, selection, rowActions),
                position: 'sticky',
                top: 0,
                zIndex: 3 + Number(entry.pin !== TablePin.NONE),
              }}
              data-pinned={entry.pin !== TablePin.NONE}
              data-pin-boundary={getTablePinnedBoundary(layout, index)}
            >
              <TableSortHeader column={entry.column} props={props} />
              {entry.column.resizable !== false && (
                <ColumnResizeHandle
                  label={__('Resize {{column}}', { column: entry.column.header })}
                  width={entry.width}
                  min={entry.column.minWidth ?? 80}
                  max={entry.column.maxWidth ?? 1200}
                  onResize={(width) => controls.resize(entry.column.id, width)}
                />
              )}
            </th>
          </Fragment>
        ))}
        {fillAtEnd && (
          <th aria-hidden="true" className={styles.filler} style={{ width: fill }} />
        )}
        {rowActions && (
          <th className={styles.actionColumn} scope="col" aria-label={__('Actions')} />
        )}
      </tr>
    </thead>
  );
};

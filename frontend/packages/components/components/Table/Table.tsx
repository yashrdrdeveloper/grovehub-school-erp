import { _map } from '@Base/utilities/map';
import { _filter } from '@Base/utilities/filter';
import { RefreshCw } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { ColumnManager } from '../ColumnManager';
import { DebouncedInput } from '../DebouncedInput';
import { Dropdown } from '../Dropdown';
import { Filters } from '../Filters';
import { Sort, SortDirection } from '../Sort';
import { QuickFilters } from '../QuickFilters';
import { Pagination } from '../Pagination';

import { TableHeader } from './components/tableHeader';
import styles from './Table.module.css';
import type { TableProps } from './types/table.props';
import type { TableColumn } from './types/table.column';
import type { useTable } from './hooks/useTable';

export const TableView = <T,>(props: TableProps<T> & { readonly controls: ReturnType<typeof useTable<T>> }) => {
  const { actions, controls, data, columns, headerProps, searchProps, filterProps, sortProps, quickFilterProps, paginationProps, rowActions } = props;
  const visible = controls.visibleColumns;
  const renderCell = (column: TableColumn<T>, row: T) => {
    if (column.renderer) return column.renderer(row);
    return String(column.accessor(row) ?? '—');
  };
  const sortColumn = (column: TableColumn<T>) => {
    const field = column.sortField ?? column.id;
    const current = sortProps?.value.find((entry) => entry.field === field);
    if (current?.direction === SortDirection.ASC) {
      actions?.SORT?.([{ field, direction: SortDirection.DESC }]);
      return;
    }
    actions?.SORT?.([{ field, direction: SortDirection.ASC }]);
  };
  const sortLabel = (column: TableColumn<T>) => {
    const current = sortProps?.value.find((entry) => entry.field === (column.sortField ?? column.id));
    if (current?.direction === SortDirection.ASC) return ' ↑';
    if (current?.direction === SortDirection.DESC) return ' ↓';
    return ' ↕';
  };
  return (
    <Box className={styles.root} dataTestId={props.dataTestId ?? props.id}>
      {headerProps && <TableHeader {...headerProps} onAction={actions?.HEADER_ACTION} />}
      <Box className={styles.panel}>
        <Box className={styles.toolbar}>
          {searchProps && <Box className={styles.search}><DebouncedInput label={__('Search')} value={searchProps.value} placeholder={searchProps.placeholder} debounceMs={searchProps.debounceMs} onChange={(value) => actions?.SEARCH?.(value)} /></Box>}
          {filterProps && <Filters {...filterProps} onChange={(value) => actions?.FILTER?.(value)} />}
          {sortProps && <Sort {...sortProps} onChange={(value) => actions?.SORT?.(value)} />}
          <ColumnManager options={_map(columns, (column) => ({ id: column.id, label: column.header, visible: !controls.hidden.includes(column.id), disabled: column.hideable === false }))} density={controls.rowDensity} onToggle={controls.toggle} onDensityChange={controls.setRowDensity} onReset={controls.reset} />
          {actions?.REFRESH && <ButtonWithIcon icon={RefreshCw} variant={BUTTON_VARIANT.OUTLINE} aria-label={__('Refresh')} onClick={actions.REFRESH} disabled={props.refreshing} />}
        </Box>
        {quickFilterProps && <QuickFilters {...quickFilterProps} onChange={(value) => actions?.QUICK_FILTER?.(value)} />}
        {props.error && <Box role="alert" className={styles.state}>{props.error}{actions?.RETRY && <Button onClick={actions.RETRY}>{__('Retry')}</Button>}</Box>}
        {props.loading && <Box role="status" className={styles.state}>{__('Loading…')}</Box>}
        {!props.loading && !props.error && !data.length && <Box className={styles.state}>{props.emptyTitle ?? __('No results found')}</Box>}
        {!!data.length && (
          <Box className={styles.scroll}>
            <table className={styles.table} data-density={controls.rowDensity} aria-busy={props.refreshing}>
              <thead>
                <tr>
                  {_map(visible, (column) => (
                    <th key={column.id} scope="col">
                      {column.sortable && <Button variant={BUTTON_VARIANT.GHOST} onClick={() => sortColumn(column)}>{column.header}{sortLabel(column)}</Button>}
                      {!column.sortable && column.header}
                    </th>
                  ))}
                  {!!rowActions?.length && <th scope="col">{__('Actions')}</th>}
                </tr>
              </thead>
              <tbody>
                {_map(data, (row) => (
                  <tr key={props.getRowId(row)}>
                    {_map(visible, (column) => (
                      <td key={column.id} data-label={column.header}>{renderCell(column, row)}</td>
                    ))}
                    {!!rowActions?.length && (
                      <td data-label={__('Actions')}>
                        <Dropdown label={__('Actions')} items={_filter(rowActions, (action) => !action.isVisible || action.isVisible(row))} onSelect={(id) => actions?.ROW_ACTION?.({ id, row })} />
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </Box>
        )}
        {paginationProps && <Pagination {...paginationProps} onChange={(action) => actions?.PAGINATION?.(action)} />}
      </Box>
    </Box>
  );
};

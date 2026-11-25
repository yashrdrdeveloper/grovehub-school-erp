import { _map } from '@Base/utilities/map';
import { RefreshCw, Search } from '@Components/icons';

import { Icon } from '../Icon';
import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { ColumnManager } from '../ColumnManager';
import { DebouncedInput } from '../DebouncedInput';
import { Filters } from '../Filters';
import { Sort } from '../Sort';
import { QuickFilters } from '../QuickFilters';
import { Pagination, PaginationType } from '../Pagination';

import { TableHeader } from './components/tableHeader';
import { DesktopTable } from './templates/desktopTable';
import { MobileTable } from './templates/mobileTable';
import { useTableViewport } from './hooks/useTableViewport';
import { getTableColumnWidth } from './helpers/table.preferences';
import { TableDataMode } from './constants/table.type';
import styles from './Table.module.css';
import type { TableViewProps } from './types/table.view';

export const TableView = <T,>(props: TableViewProps<T>) => {
  const { actions, controls, columns, headerProps, subheaderProps, footerProps, searchProps, filterProps, sortProps, quickFilterProps, paginationProps } = props;
  const mobile = useTableViewport();
  const getPagination = () => {
    if (paginationProps?.type === PaginationType.OFFSET && props.dataMode === TableDataMode.CLIENT) return { ...paginationProps, totalRecords: controls.total };
    return paginationProps;
  };
  const pagination = getPagination();
  const order = Array.from(new Set([...controls.preferences.order, ...columns.map((column) => column.id)]));
  const options = order.flatMap((id) => {
    const column = columns.find((entry) => entry.id === id);
    if (!column) return [];
    return [{
      id, label: column.header, visible: controls.layout.some((entry) => entry.column.id === id),
      disabled: column.hideable === false, width: getTableColumnWidth(column, controls.preferences.widths[id]),
      pin: controls.preferences.pins[id] ?? column.pin, resizable: column.resizable, reorderable: column.reorderable,
    }];
  });
  return (
    <Box className={styles.root} dataTestId={props.dataTestId ?? props.id} data-fullscreen={controls.fullscreen} data-fill-height={props.fillHeight}>
      {headerProps && <TableHeader {...headerProps} onAction={actions?.HEADER_ACTION} />}
      <Box className={styles.panel}>
        <Box className={styles.toolbar}>
          {searchProps && <Box className={styles.search}><Box className={styles.searchIcon} aria-hidden="true"><Icon icon={Search} /></Box><DebouncedInput type="search" label={__('Search')} value={searchProps.value} placeholder={searchProps.placeholder} debounceMs={searchProps.debounceMs} onChange={(value) => actions?.SEARCH?.(value)} /></Box>}
          {filterProps && <Filters {...filterProps} onChange={(value) => actions?.FILTER?.(value)} />}
          {sortProps && <Sort {...sortProps} onChange={(value) => actions?.SORT?.(value)} />}
          <ColumnManager options={options} density={controls.preferences.density} onToggle={controls.toggle} onDensityChange={controls.setRowDensity} onReset={controls.reset} onReorder={controls.reorder} onPin={controls.pin} maxPinnedColumns={props.maxPinnedColumns} />
          {actions?.REFRESH && <ButtonWithIcon icon={RefreshCw} variant={BUTTON_VARIANT.OUTLINE} aria-label={__('Refresh')} onClick={actions.REFRESH} disabled={props.refreshing} />}
          {actions?.EXPORT && <Button variant={BUTTON_VARIANT.OUTLINE} onClick={actions.EXPORT}>{__('Export')}</Button>}
          {actions?.FULLSCREEN_CHANGE && <Button variant={BUTTON_VARIANT.OUTLINE} aria-pressed={controls.fullscreen} onClick={() => controls.changeFullscreen(!controls.fullscreen)}>{!controls.fullscreen && __('Fullscreen')}{controls.fullscreen && __('Exit fullscreen')}</Button>}
        </Box>
        {subheaderProps && <TableHeader {...subheaderProps} onAction={actions?.SUBHEADER_ACTION} />}
        {quickFilterProps && <QuickFilters {...quickFilterProps} onChange={(value) => actions?.QUICK_FILTER?.(value)} />}
        {!!controls.selected.length && props.bulkActions && <Box className={styles.toolbar}>
          <Box>{__('{{count}} selected', { count: controls.selected.length })}</Box>
          {_map(props.bulkActions, (item) => <Button key={item.id} disabled={item.disabled} variant={BUTTON_VARIANT.OUTLINE} onClick={() => actions?.BULK_ACTION?.({ id: item.id, rowIds: controls.selected })}>{item.label}</Button>)}
          <Button variant={BUTTON_VARIANT.GHOST} onClick={() => actions?.SELECTION_CHANGE?.([])}>{__('Clear selection')}</Button>
        </Box>}
        <Box className={styles.body}>
        {props.error && <Box role="alert" className={styles.state}>{props.error}{actions?.RETRY && <Button onClick={actions.RETRY}>{__('Retry')}</Button>}</Box>}
        {props.loading && <Box role="status" className={styles.state}>{__('Loading…')}</Box>}
        {!props.loading && !props.error && !controls.entries.length && <Box className={styles.state}>{props.emptyTitle ?? __('No results found')}</Box>}
        {!!controls.entries.length && !mobile && <DesktopTable {...props} />}
        {!!controls.entries.length && mobile && <MobileTable {...props} />}
        </Box>
        <Box className={styles.footer}>
        {pagination && <Pagination {...pagination} onChange={(action) => {
          if (action.type === 'LOAD_MORE' && actions?.LOAD_MORE) {
            actions.LOAD_MORE();
            return;
          }
          actions?.PAGINATION?.(action);
        }} />}
        {footerProps && <TableHeader {...footerProps} onAction={actions?.FOOTER_ACTION} />}
        </Box>
      </Box>
    </Box>
  );
};

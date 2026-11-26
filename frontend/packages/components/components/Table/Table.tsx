import { useState } from 'react';
import { getFormattedNumber } from '@Base/localization/getFormattedNumber';
import { _map } from '@Base/utilities/map';
import { Ellipsis, RefreshCw } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { ColumnManager } from '../ColumnManager';
import { EmptyState } from '../EmptyState';
import { Filters } from '../Filters';
import { Tooltip } from '../Tooltip';
import { Sort } from '../Sort';
import { Sheet } from '../Sheet';
import { Select } from '../Select';
import { QuickFilters } from '../QuickFilters';
import { Pagination, PaginationType } from '../Pagination';

import { TableSearch } from './components/tableSearch';
import { TableOptions } from './components/tableOptions';
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
  const [mobileOptionsOpen, setMobileOptionsOpen] = useState(false);
  const infinite = mobile && props.mobileInfiniteProps;
  const getPagination = () => {
    if (paginationProps?.type === PaginationType.OFFSET && props.dataMode === TableDataMode.CLIENT) return { ...paginationProps, totalRecords: controls.total };
    return paginationProps;
  };
  const pagination = getPagination();
  const getTotal = () => {
    if (pagination?.type === PaginationType.OFFSET) return pagination.totalRecords;
    if (props.dataMode === TableDataMode.CLIENT) return controls.total;
    return props.totalRecords;
  };
  const total = getTotal();
  const getSummary = () => {
    if (total !== undefined) return __('Total records: {{count}}', { count: getFormattedNumber(total) });
    return __('Records on this page: {{count}}', { count: getFormattedNumber(props.data.length) });
  };
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
        {!mobile && <Box className={styles.toolbar}>
          <Box className={styles.toolbarLeft}>
            {filterProps && <Filters {...filterProps} iconOnly onChange={(value) => actions?.FILTER?.(value)} />}
            {sortProps && <Sort {...sortProps} iconOnly onChange={(value) => actions?.SORT?.(value)} />}
          </Box>
          <Box className={styles.toolbarRight}>
            {searchProps && <TableSearch {...searchProps} onChange={(value) => actions?.SEARCH?.(value)} />}
            {actions?.REFRESH && (
              <Tooltip label={props.refreshing ? __('Refreshing…') : __('Refresh')} focusable={Boolean(props.refreshing)}><ButtonWithIcon icon={RefreshCw} iconSize={16} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} aria-label={__('Refresh')} onClick={actions.REFRESH} disabled={props.refreshing} /></Tooltip>
            )}
            <ColumnManager options={options} density={controls.preferences.density} onToggle={controls.toggle} onDensityChange={controls.setRowDensity} onReset={controls.reset} onReorder={controls.reorder} onPin={controls.pin} maxPinnedColumns={props.maxPinnedColumns} />
            {(actions?.EXPORT || actions?.FULLSCREEN_CHANGE) && (
              <TableOptions>
                {actions?.EXPORT && <Button variant={BUTTON_VARIANT.GHOST} onClick={actions.EXPORT}>{__('Export')}</Button>}
                {actions?.FULLSCREEN_CHANGE && <Button variant={BUTTON_VARIANT.GHOST} aria-pressed={controls.fullscreen} onClick={() => controls.changeFullscreen(!controls.fullscreen)}>{!controls.fullscreen && __('Fullscreen')}{controls.fullscreen && __('Exit fullscreen')}</Button>}
              </TableOptions>
            )}
          </Box>
        </Box>}
        {mobile && <Box className={styles.mobileToolbar}>
          {searchProps && <TableSearch {...searchProps} placeholder={searchProps.mobilePlaceholder ?? searchProps.placeholder} fields={undefined} onChange={(value) => actions?.SEARCH?.(value)} />}
          <ButtonWithIcon icon={Ellipsis} iconSize={18} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} aria-label={__('Table options')} aria-expanded={mobileOptionsOpen} onClick={() => setMobileOptionsOpen(true)} />
          <Sheet open={mobileOptionsOpen} onOpenChange={setMobileOptionsOpen} title={__('Table options')}>
            <Box className={styles.mobileOptions}>
              {filterProps && <Filters {...filterProps} onChange={(value) => { actions?.FILTER?.(value); setMobileOptionsOpen(false); }} />}
              {sortProps && <Sort {...sortProps} onChange={(value) => { actions?.SORT?.(value); setMobileOptionsOpen(false); }} />}
              <ColumnManager options={options} mobileMenu mobileList={Boolean(props.mobileRowProps)} density={controls.preferences.density} onToggle={controls.toggle} onDensityChange={controls.setRowDensity} onReset={controls.reset} onReorder={controls.reorder} onPin={controls.pin} maxPinnedColumns={props.maxPinnedColumns} />
              {pagination && <Select<number> label={__('Records per load')} value={pagination.pageSize} options={_map(pagination.pageSizeOptions ?? [10, 25, 50, 100], (size) => ({ value: size, label: String(size) }))} searchable={false} clearable={false} onChange={(pageSize) => { if (pageSize) { actions?.PAGINATION?.({ type: 'PAGE_SIZE_CHANGE', pageSize }); setMobileOptionsOpen(false); } }} />}
              {actions?.REFRESH && <ButtonWithIcon icon={RefreshCw} variant={BUTTON_VARIANT.OUTLINE} disabled={props.refreshing} onClick={() => { setMobileOptionsOpen(false); actions.REFRESH?.(); }}>{__('Refresh')}</ButtonWithIcon>}
              {actions?.EXPORT && <Button variant={BUTTON_VARIANT.OUTLINE} onClick={() => { setMobileOptionsOpen(false); actions.EXPORT?.(); }}>{__('Export')}</Button>}
            </Box>
          </Sheet>
        </Box>}
        {subheaderProps && <TableHeader {...subheaderProps} onAction={actions?.SUBHEADER_ACTION} />}
        {quickFilterProps && <QuickFilters {...quickFilterProps} onChange={(value) => actions?.QUICK_FILTER?.(value)} />}
        {!!controls.selected.length && props.bulkActions && <Box className={styles.toolbar}>
          <Box>{__('{{count}} selected', { count: controls.selected.length })}</Box>
          {_map(props.bulkActions, (item) => <Button key={item.id} disabled={item.disabled} variant={BUTTON_VARIANT.OUTLINE} onClick={() => actions?.BULK_ACTION?.({ id: item.id, rowIds: controls.selected })}>{item.label}</Button>)}
          <Button variant={BUTTON_VARIANT.GHOST} onClick={() => actions?.SELECTION_CHANGE?.([])}>{__('Clear selection')}</Button>
        </Box>}
        <Box className={styles.body} key={infinite ? infinite.resetKey : 'table-body'}>
        {props.error && <Box role="alert" className={styles.state}>{props.error}{actions?.RETRY && <Button onClick={actions.RETRY}>{__('Retry')}</Button>}</Box>}
        {props.loading && <Box role="status" className={styles.state}>{__('Loading…')}</Box>}
        {!props.loading && !props.error && !controls.entries.length && <Box className={styles.state}><EmptyState title={props.emptyTitle ?? __('No results found')} description={__('No records match the current view. Try changing your search or filters.')} /></Box>}
        {!!controls.entries.length && !mobile && <DesktopTable {...props} />}
        {!!controls.entries.length && mobile && <MobileTable {...props} />}
        </Box>
        {!infinite && <Box className={styles.footer}>
        <Box className={styles.footerRow}>
          <Box className={styles.summary} aria-live="polite">
            {!props.loading && !props.error && getSummary()}
          </Box>
          {pagination && <Pagination {...pagination} onChange={(action) => {
            if (action.type === 'LOAD_MORE' && actions?.LOAD_MORE) {
              actions.LOAD_MORE();
              return;
            }
            actions?.PAGINATION?.(action);
          }} />}
        </Box>
        {footerProps && <TableHeader {...footerProps} onAction={actions?.FOOTER_ACTION} />}
        </Box>}
      </Box>
    </Box>
  );
};

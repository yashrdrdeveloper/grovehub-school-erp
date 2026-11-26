import { getFormattedNumber } from '@Base/localization/getFormattedNumber';
import { _map } from '@Base/utilities/map';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { EmptyState } from '../EmptyState';
import { QuickFilters } from '../QuickFilters';
import { Pagination, PaginationType } from '../Pagination';

import { DesktopTableToolbar } from './components/desktopTableToolbar';
import { MobileTableToolbar } from './components/mobileTableToolbar';
import { TableHeader } from './components/tableHeader';
import { DesktopTable } from './templates/desktopTable';
import { MobileTable } from './templates/mobileTable';
import { useTableViewport } from './hooks/useTableViewport';
import { getTableColumnOptions } from './helpers/table.getColumnOptions';
import { TableDataMode } from './constants/table.type';
import styles from './Table.module.css';
import type { TableViewProps } from './types/table.view';

export const TableView = <T,>(props: TableViewProps<T>) => {
  const {
    actions,
    controls,
    headerProps,
    subheaderProps,
    footerProps,
    quickFilterProps,
    paginationProps,
  } = props;
  const mobile = useTableViewport();
  const infinite = mobile && props.mobileInfiniteProps;
  const getPagination = () => {
    if (paginationProps?.type === PaginationType.OFFSET && props.dataMode === TableDataMode.CLIENT)
      return { ...paginationProps, totalRecords: controls.total };
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
    if (total !== undefined)
      return __('Total records: {{count}}', { count: getFormattedNumber(total) });
    return __('Records on this page: {{count}}', { count: getFormattedNumber(props.data.length) });
  };
  const options = getTableColumnOptions(props);
  return (
    <Box
      className={styles.root}
      dataTestId={props.dataTestId ?? props.id}
      data-fullscreen={controls.fullscreen}
      data-fill-height={props.fillHeight}
    >
      {headerProps && <TableHeader {...headerProps} onAction={actions?.HEADER_ACTION} />}
      <Box className={styles.panel} data-mobile-list={Boolean(mobile && props.mobileRowProps)}>
        {!mobile && <DesktopTableToolbar {...props} options={options} />}
        {mobile && <MobileTableToolbar {...props} options={options} />}
        {subheaderProps && <TableHeader {...subheaderProps} onAction={actions?.SUBHEADER_ACTION} />}
        {quickFilterProps && (
          <QuickFilters
            {...quickFilterProps}
            onChange={(value) => actions?.QUICK_FILTER?.(value)}
          />
        )}
        {!!controls.selected.length && props.bulkActions && (
          <Box className={styles.toolbar}>
            <Box>{__('{{count}} selected', { count: controls.selected.length })}</Box>
            {_map(props.bulkActions, (item) => (
              <Button
                key={item.id}
                disabled={item.disabled}
                variant={BUTTON_VARIANT.OUTLINE}
                onClick={() => actions?.BULK_ACTION?.({ id: item.id, rowIds: controls.selected })}
              >
                {item.label}
              </Button>
            ))}
            <Button variant={BUTTON_VARIANT.GHOST} onClick={() => actions?.SELECTION_CHANGE?.([])}>
              {__('Clear selection')}
            </Button>
          </Box>
        )}
        <Box className={styles.body} key={(infinite && infinite.resetKey) || 'table-body'}>
          {props.error && (
            <Box role="alert" className={styles.state}>
              {props.error}
              {actions?.RETRY && <Button onClick={actions.RETRY}>{__('Retry')}</Button>}
            </Box>
          )}
          {props.loading && (
            <Box role="status" className={styles.state}>
              {__('Loading…')}
            </Box>
          )}
          {!props.loading && !props.error && !controls.entries.length && (
            <Box className={styles.state}>
              <EmptyState
                title={props.emptyTitle ?? __('No results found')}
                description={__(
                  'No records match the current view. Try changing your search or filters.',
                )}
              />
            </Box>
          )}
          {!!controls.entries.length && !mobile && <DesktopTable {...props} />}
          {!!controls.entries.length && mobile && <MobileTable {...props} />}
        </Box>
        {!infinite && (
          <Box className={styles.footer}>
            <Box className={styles.footerRow}>
              <Box className={styles.summary} aria-live="polite">
                {!props.loading && !props.error && getSummary()}
              </Box>
              {pagination && (
                <Pagination
                  {...pagination}
                  onChange={(action) => {
                    if (action.type === 'LOAD_MORE' && actions?.LOAD_MORE) {
                      actions.LOAD_MORE();
                      return;
                    }
                    actions?.PAGINATION?.(action);
                  }}
                />
              )}
            </Box>
            {footerProps && <TableHeader {...footerProps} onAction={actions?.FOOTER_ACTION} />}
          </Box>
        )}
      </Box>
    </Box>
  );
};

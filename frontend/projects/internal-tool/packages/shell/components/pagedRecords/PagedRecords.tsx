import { useMemo } from 'react';
import { PaginationType } from '@Components/components/Pagination';
import { Table, TableDataMode, TableType } from '@Components/components/Table';

import { getRecordColumns } from './helpers/pagedRecords.columns';
import { usePagedRecords } from './hooks/usePagedRecords';
import type { PagedRecordsProps } from './types/pagedRecords.props';

export function PagedRecords<T>({
  resource,
  revision,
  filters,
  columns,
  rowKey,
  onSelect,
  title,
}: PagedRecordsProps<T>) {
  const records = usePagedRecords<T>(resource, revision, filters);
  const tableColumns = useMemo(() => getRecordColumns(columns), [columns]);

  return (
    <Table<T>
      id={resource + '-records'}
      dataTestId={resource + '-records'}
      type={TableType.DEFAULT}
      dataMode={TableDataMode.SERVER}
      data={records.data}
      getRowId={(record) => String(rowKey(record))}
      columns={tableColumns}
      headerProps={{ title }}
      loading={records.loading}
      error={records.error}
      actions={{
        ROW_CLICK: onSelect,
        ROW_ACTION: ({ row }) => onSelect?.(row),
        REFRESH: records.refresh,
        RETRY: records.refresh,
        PAGINATION: records.paginate,
      }}
      rowActions={onSelect && [{ id: 'open', label: __('Open') }]}
      paginationProps={{
        type: PaginationType.HYBRID,
        pageSize: records.rowSize,
        currentPage: records.pageNumber,
        hasNext: records.hasMore,
        hasPrevious: records.pageNumber > 1,
        loading: records.loading,
        dataTestId: resource + '-pagination',
      }}
    />
  );
}

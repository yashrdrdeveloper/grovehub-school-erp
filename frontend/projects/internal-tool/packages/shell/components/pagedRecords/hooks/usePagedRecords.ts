import { useEffect, useState } from 'react';
import type { PaginationAction } from '@Components/components/Pagination';

import { errorText, search } from '../../../../business/api/client';
import type { Page } from '../../../../business/api/types';

export function usePagedRecords<T>(resource: string, revision: number) {
  const [pageNumber, setPageNumber] = useState(1);
  const [rowSize, setRowSize] = useState(25);
  const [refresh, setRefresh] = useState(0);
  const [page, setPage] = useState<Page<T> | null>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const controller = new AbortController();
    setLoading(true);
    setError('');
    setPage(null);
    void search<T>(resource, pageNumber, rowSize, controller.signal)
      .then((value) => {
        if (controller.signal.aborted) return;
        if (!value.records.length && pageNumber > 1) {
          setPageNumber((current) => Math.max(1, current - 1));
          return;
        }
        setPage(value);
      })
      .catch((failure: unknown) => {
        if (!controller.signal.aborted) setError(errorText(failure));
      })
      .finally(() => {
        if (!controller.signal.aborted) setLoading(false);
      });
    return () => controller.abort();
  }, [resource, pageNumber, rowSize, refresh, revision]);

  function paginate(action: PaginationAction) {
    if (loading) return;
    switch (action.type) {
      case 'NEXT':
        if (page?.hasMore) setPageNumber((current) => current + 1);
        return;
      case 'PREVIOUS':
        setPageNumber((current) => Math.max(1, current - 1));
        return;
      case 'PAGE_SIZE_CHANGE':
        setRowSize(action.pageSize);
        setPageNumber(1);
        return;
      case 'PAGE_CHANGE':
        if (action.page >= 1 && action.page <= 10000) setPageNumber(action.page);
        return;
      default:
        return;
    }
  }

  return {
    data: page?.records ?? [],
    hasMore: Boolean(page?.hasMore) && pageNumber < 10000,
    pageNumber,
    rowSize,
    loading,
    error,
    paginate,
    refresh: () => setRefresh((current) => current + 1),
  };
}

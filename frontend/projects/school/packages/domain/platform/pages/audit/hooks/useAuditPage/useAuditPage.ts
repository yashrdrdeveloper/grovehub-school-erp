import { useState } from 'react';
import { useQuery, useInfiniteQuery } from '@tanstack/react-query';
import { rbacOperations } from '@Business/domain/users/rbac/services/operations';
import { useAuth } from '@Domain/auth';
import { useMobileViewport } from '@Components/hooks/useMobileViewport';
import { PaginationType } from '@Components/components/Pagination';
import type { PaginationAction } from '@Components/components/Pagination';
import type { TableActions } from '@Components/components/Table';
import type { FilterCondition } from '@Components/components/Filters';
import type { SortCondition } from '@Components/components/Sort';
import { auditApi } from '@Business/domain/platform/audit/services/audit';
import type { AuditActivity } from '@Business/domain/platform/audit/types/contracts';
import { getAuditSearchRequest } from '../../helpers/audit.getSearchRequest';

export const useAuditPage = () => {
  const { account } = useAuth();
  const mobile = useMobileViewport();
  const [tenantScope, setTenantScope] = useState(false);
  const capabilities = useQuery({
    queryKey: ['capabilities', account?.accountId],
    queryFn: rbacOperations.capabilities,
    enabled: Boolean(account),
  });
  const canViewTenant = capabilities.data?.actions.auditSearch === true;
  const getScope = (): 'ME' | undefined => {
    if (tenantScope && canViewTenant) return undefined;
    return 'ME';
  };
  const auditScope = getScope();
  const [searchText, setSearchText] = useState('');
  const [filters, setFilters] = useState<readonly FilterCondition[]>([]);
  const [sorting, setSorting] = useState<readonly SortCondition[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(25);
  const [selected, setSelected] = useState<AuditActivity | null>(null);
  const request = getAuditSearchRequest(searchText, filters, sorting, page, pageSize);
  const mobileRequest = getAuditSearchRequest(searchText, filters, sorting, 1, 25);
  const scope = ['audit', account?.accountId, auditScope];
  const query = useQuery({
    queryKey: [...scope, request],
    enabled: Boolean(account) && !mobile,
    queryFn: () => auditApi.search({ ...request, scope: auditScope }),
  });
  const mobileQuery = useInfiniteQuery({
    queryKey: [...scope, 'mobile', mobileRequest],
    enabled: Boolean(account) && mobile,
    initialPageParam: { pageNumber: 1, cursor: undefined as string | undefined },
    queryFn: ({ pageParam }) =>
      auditApi.search({
        ...mobileRequest,
        pageInfo: { ...mobileRequest.pageInfo, ...pageParam },
        scope: auditScope,
      }),
    getNextPageParam: (last, pages) => {
      if (!last.pageInfo.hasNext || !last.pageInfo.nextCursor) return undefined;
      return { pageNumber: pages.length + 1, cursor: last.pageInfo.nextCursor };
    },
  });
  const paginate = (action: PaginationAction) => {
    if (query.isFetching) return;
    if (action.type === 'PAGE_SIZE_CHANGE') {
      setPageSize(action.pageSize);
      setPage(1);
    }
    if (action.type === 'PAGE_CHANGE' && action.page > 0) setPage(action.page);
    if (action.type === 'NEXT' && query.data?.pageInfo.hasNext) setPage(page + 1);
    if (action.type === 'PREVIOUS' && page > 1) setPage(page - 1);
  };
  const refresh = () => {
    if (mobile) {
      void mobileQuery.refetch();
      return;
    }
    void query.refetch();
  };
  const actions: TableActions<AuditActivity> = {
    SEARCH: (value) => {
      setSearchText(value);
      setPage(1);
    },
    FILTER: (value) => {
      setFilters(value);
      setPage(1);
    },
    SORT: (value) => {
      setSorting(value);
      setPage(1);
    },
    PAGINATION: paginate,
    REFRESH: refresh,
    RETRY: refresh,
    ROW_CLICK: setSelected,
    LOAD_MORE: () => {
      if (mobileQuery.hasNextPage && !mobileQuery.isFetching)
        void mobileQuery.fetchNextPage({ cancelRefetch: false });
    },
  };
  const getRows = () => {
    if (mobile) return mobileQuery.data?.pages.flatMap((result) => result.items ?? []) ?? [];
    return query.data?.items ?? [];
  };
  const getActiveQuery = () => {
    if (mobile) return mobileQuery;
    return query;
  };
  const active = getActiveQuery();
  return {
    canViewTenant,
    tenantScope,
    changeScope: (value: boolean) => {
      setTenantScope(value);
      setPage(1);
      setSelected(null);
    },
    actions,
    searchText,
    filters,
    sorting,
    selected,
    setSelected,
    rows: getRows(),
    loading: active.isPending,
    refreshing: active.isFetching,
    error: active.isError,
    pagination: {
      type: PaginationType.HYBRID,
      currentPage: page,
      pageSize,
      hasNext: Boolean(query.data?.pageInfo.hasNext),
      hasPrevious: page > 1,
      totalRecords: query.data?.pageInfo.totalRecords,
      totalPages: query.data?.pageInfo.totalPages,
    },
    mobileInfiniteProps: {
      hasNext: mobileQuery.hasNextPage,
      loading: mobileQuery.isFetching,
      error: mobileQuery.isFetchNextPageError,
      resetKey: JSON.stringify(mobileRequest),
    },
  };
};

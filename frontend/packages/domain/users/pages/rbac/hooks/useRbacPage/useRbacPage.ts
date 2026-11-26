import { useState } from 'react';
import { useQuery, useInfiniteQuery } from '@tanstack/react-query';
import { useAuth } from '@Domain/auth';
import { useMobileViewport } from '@Components/hooks/useMobileViewport';
import { PaginationType } from '@Components/components/Pagination';
import type { PaginationAction } from '@Components/components/Pagination';
import type { TableActions } from '@Components/components/Table';
import type { FilterCondition } from '@Components/components/Filters';
import type { SortCondition } from '@Components/components/Sort';
import { rbacApi } from '@Business/domain/users/rbac/services/rbac';
import type { Role } from '@Business/domain/users/rbac/types/contracts';
import { getRbacSearchRequest } from '../../helpers/rbac.getSearchRequest';

export const useRbacPage = () => {
  const { account } = useAuth();
  const mobile = useMobileViewport();
  const [searchText, setSearchText] = useState('');
  const [filters, setFilters] = useState<readonly FilterCondition[]>([]);
  const [sorting, setSorting] = useState<readonly SortCondition[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(25);
  const [selected, setSelected] = useState<Role | null>(null);
  const request = getRbacSearchRequest(searchText, filters, sorting, page, pageSize);
  const mobileRequest = getRbacSearchRequest(searchText, filters, sorting, 1, 25);
  const scope = ['rbac', account?.accountId];
  const query = useQuery({ queryKey: [...scope, request], enabled: Boolean(account) && !mobile, queryFn: () => rbacApi.search({ ...request }) });
  const mobileQuery = useInfiniteQuery({ queryKey: [...scope, 'mobile', mobileRequest], enabled: Boolean(account) && mobile, initialPageParam: { pageNumber: 1, cursor: undefined as string | undefined }, queryFn: ({ pageParam }) => rbacApi.search({ ...mobileRequest, pageInfo: { ...mobileRequest.pageInfo, ...pageParam } }), getNextPageParam: (last, pages) => { if (!last.pageInfo.hasNext || !last.pageInfo.nextCursor) return undefined; return { pageNumber: pages.length + 1, cursor: last.pageInfo.nextCursor }; } });
  const permissions = useQuery({ queryKey: [...scope, 'permissions', selected?.roleId], enabled: Boolean(account && selected), queryFn: () => rbacApi.permissions(selected?.roleId ?? '') });
  const paginate = (action: PaginationAction) => {
    if (query.isFetching) return;
    if (action.type === 'PAGE_SIZE_CHANGE') { setPageSize(action.pageSize); setPage(1); }
    if (action.type === 'PAGE_CHANGE' && action.page > 0) setPage(action.page);
    if (action.type === 'NEXT' && query.data?.pageInfo.hasNext) setPage(page + 1);
    if (action.type === 'PREVIOUS' && page > 1) setPage(page - 1);
  };
  const refresh = () => { if (mobile) { void mobileQuery.refetch(); return; } void query.refetch(); };
  const actions: TableActions<Role> = {
    SEARCH: (value) => { setSearchText(value); setPage(1); },
    FILTER: (value) => { setFilters(value); setPage(1); },
    SORT: (value) => { setSorting(value); setPage(1); },
    PAGINATION: paginate, REFRESH: refresh, RETRY: refresh, ROW_CLICK: setSelected,
    LOAD_MORE: () => { if (mobileQuery.hasNextPage && !mobileQuery.isFetching) void mobileQuery.fetchNextPage({ cancelRefetch: false }); },
  };
  const getRows = () => { if (mobile) return mobileQuery.data?.pages.flatMap((result) => result.items ?? []) ?? []; return query.data?.items ?? []; };
  const getActiveQuery = () => { if (mobile) return mobileQuery; return query; };
  const active = getActiveQuery();
  return { actions, searchText, filters, sorting, selected, setSelected, rows: getRows(), loading: active.isPending, refreshing: active.isFetching, error: active.isError, pagination: { type: PaginationType.HYBRID, currentPage: page, pageSize, hasNext: Boolean(query.data?.pageInfo.hasNext), hasPrevious: page > 1, totalRecords: query.data?.pageInfo.totalRecords, totalPages: query.data?.pageInfo.totalPages }, mobileInfiniteProps: { hasNext: mobileQuery.hasNextPage, loading: mobileQuery.isFetching, error: mobileQuery.isFetchNextPageError, resetKey: JSON.stringify(mobileRequest) }, permissions };
};

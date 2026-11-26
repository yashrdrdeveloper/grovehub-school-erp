import { useState } from 'react';
import { useQuery, useInfiniteQuery, useMutation, useQueryClient } from '@tanstack/react-query';

import { usersApi } from '@Business/domain/users';
import { useAuth } from '@Domain/auth';
import { useMobileViewport } from '@Components/hooks/useMobileViewport';
import { PaginationType } from '@Components/components/Pagination';
import type { PaginationAction, PaginationProps } from '@Components/components/Pagination';
import type { TableActions } from '@Components/components/Table';
import type { Account, CreateAccount } from '@Business/domain/users';

import { UserAction } from '../../constants/userSetup.actions';
import { USER_SETUP_DEFAULT_SORT } from '../../constants/userSetup.sort';
import { getUserSetupNextPage } from '../../helpers/userSetup.getNextPage';
import { getUserSetupMobileRows } from '../../helpers/userSetup.getMobileRows';
import { getUserSetupSearchRequest } from '../../helpers/userSetup.getSearchRequest';
import type { UserSetupDialog, UserSetupQueryState } from '../../types/userSetup.state';

export const useUserSetup = () => {
  const { account } = useAuth();
  const queryClient = useQueryClient();
  const mobile = useMobileViewport();
  const [state, setState] = useState<UserSetupQueryState>({
    searchText: '',
    filters: [],
    sorting: USER_SETUP_DEFAULT_SORT,
    quickFilterIds: ['ALL'],
    pageSize: 25,
    pageNumber: 1,
    cursors: {},
  });
  const [dialog, setDialog] = useState<UserSetupDialog | null>(null);
  const request = getUserSetupSearchRequest(state);
  const scope = ['user-setup', account?.accountId];
  const query = useQuery({
    queryKey: [...scope, request],
    queryFn: () => usersApi.search(request),
    enabled: Boolean(account) && !mobile,
  });
  const mobileRequest = getUserSetupSearchRequest({ ...state, pageNumber: 1, cursors: {} });
  const mobileQuery = useInfiniteQuery({
    queryKey: [...scope, 'mobile', mobileRequest],
    initialPageParam: { pageNumber: 1, cursor: undefined as string | undefined },
    queryFn: ({ pageParam }) => usersApi.search({ ...mobileRequest, pageInfo: { ...mobileRequest.pageInfo, ...pageParam } }),
    getNextPageParam: getUserSetupNextPage,
    enabled: Boolean(account) && mobile,
  });
  const mobileRows = getUserSetupMobileRows(mobileQuery.data?.pages ?? []);
  const getActiveQuery = () => {
    if (mobile) return mobileQuery;
    return query;
  };
  const activeQuery = getActiveQuery();
  const getData = () => {
    if (mobile) return mobileRows;
    return query.data?.items ?? [];
  };
  const getTotal = () => {
    if (mobile) return mobileQuery.data?.pages[0]?.pageInfo.totalRecords;
    return query.data?.pageInfo.totalRecords;
  };
  const getError = () => {
    if (mobile) return mobileQuery.isError && !mobileQuery.data;
    return query.isError;
  };
  const refresh = () => {
    if (mobile) { void mobileQuery.refetch(); return; }
    void query.refetch();
  };
  const mutation = useMutation({
    mutationFn: async (operation: () => Promise<unknown>) => operation(),
    onSuccess: async () => {
      setDialog(null);
      await queryClient.invalidateQueries({ queryKey: scope });
    },
  });
  const openDialog = (next: UserSetupDialog) => {
    mutation.reset();
    setDialog(next);
  };
  const closeDialog = () => {
    if (!mutation.isPending) setDialog(null);
  };
  const resetQuery = (change: Partial<UserSetupQueryState>) => {
    setState((current) => ({ ...current, ...change, pageNumber: 1, cursors: {} }));
  };
  const paginate = (action: PaginationAction) => {
    if (action.type === 'PAGE_SIZE_CHANGE') {
      resetQuery({ pageSize: action.pageSize });
      return;
    }
    if (query.isFetching) return;
    if (action.type === 'PAGE_CHANGE') {
      const totalPages = query.data?.pageInfo.totalPages;
      if (!Number.isSafeInteger(action.page) || action.page < 1 || (totalPages !== undefined && action.page > totalPages)) return;
      setState((current) => {
        const nextCursor = query.data?.pageInfo.nextCursor;
        if (action.page === current.pageNumber + 1 && nextCursor) {
          return { ...current, pageNumber: action.page, cursors: { ...current.cursors, [action.page]: nextCursor } };
        }
        return { ...current, pageNumber: action.page };
      });
      return;
    }
    if (action.type === 'PREVIOUS') {
      setState((current) => {
        if (current.pageNumber <= 1) return current;
        return { ...current, pageNumber: current.pageNumber - 1 };
      });
      return;
    }
    const nextCursor = query.data?.pageInfo.nextCursor;
    if (action.type === 'NEXT' && query.data?.pageInfo.hasNext && nextCursor) {
      setState((current) => ({
        ...current,
        pageNumber: current.pageNumber + 1,
        cursors: { ...current.cursors, [current.pageNumber + 1]: nextCursor },
      }));
    }
  };
  const actions: TableActions<Account> = {
    SEARCH: (searchText) => resetQuery({ searchText }),
    FILTER: (filters) => resetQuery({ filters }),
    SORT: (sorting) => resetQuery({ sorting }),
    QUICK_FILTER: (quickFilterIds) => resetQuery({ quickFilterIds }),
    PAGINATION: paginate,
    REFRESH: refresh,
    RETRY: refresh,
    LOAD_MORE: () => {
      if (mobile && mobileQuery.hasNextPage && !mobileQuery.isFetching) void mobileQuery.fetchNextPage({ cancelRefetch: false });
    },
    HEADER_ACTION: (id) => {
      if (id === UserAction.CREATE) openDialog({ action: UserAction.CREATE });
    },
    ROW_CLICK: (row) => openDialog({ action: UserAction.EDIT, row }),
    ROW_ACTION: ({ id, row }) => {
      if (id === UserAction.VIEW) openDialog({ action: UserAction.VIEW, row });
      if (id === UserAction.EDIT) openDialog({ action: UserAction.EDIT, row });
      if (id === UserAction.ACTIVATE) openDialog({ action: UserAction.ACTIVATE, row });
      if (id === UserAction.SUSPEND) openDialog({ action: UserAction.SUSPEND, row });
      if (id === UserAction.DEACTIVATE) openDialog({ action: UserAction.DEACTIVATE, row });
      if (id === UserAction.RESEND) openDialog({ action: UserAction.RESEND, row });
    },
  };
  const create = (input: CreateAccount, key: string) =>
    mutation.mutate(() => usersApi.create(input, key));
  const edit = (row: Account, username: string | null) =>
    mutation.mutate(() => usersApi.update(row.accountId, username));
  const confirm = () => {
    const row = dialog?.row;
    if (!row || !dialog) return;
    if (dialog.action === UserAction.ACTIVATE) mutation.mutate(() => usersApi.activate(row.accountId));
    if (dialog.action === UserAction.SUSPEND) mutation.mutate(() => usersApi.suspend(row.accountId));
    if (dialog.action === UserAction.DEACTIVATE) mutation.mutate(() => usersApi.deactivate(row.accountId));
    if (dialog.action === UserAction.RESEND) mutation.mutate(() => usersApi.resendInvitation(row.accountId));
  };
  const pagination = {
    type: PaginationType.HYBRID,
    pageSize: state.pageSize,
    currentPage: state.pageNumber,
    totalRecords: query.data?.pageInfo.totalRecords,
    totalPages: query.data?.pageInfo.totalPages,
    hasPrevious: state.pageNumber > 1,
    hasNext: Boolean(query.data?.pageInfo.hasNext && query.data.pageInfo.nextCursor),
    loading: query.isFetching,
  } satisfies Omit<Extract<PaginationProps, { type: typeof PaginationType.HYBRID }>, 'onChange'>;

  return {
    state, actions, dialog, closeDialog, create, edit, confirm, pagination,
    totalRecords: getTotal(),
    data: getData(),
    loading: activeQuery.isPending,
    refreshing: activeQuery.isFetching,
    error: getError(),
    mobileInfiniteProps: {
      hasNext: mobileQuery.hasNextPage,
      loading: mobileQuery.isFetching,
      error: mobileQuery.isFetchNextPageError,
      resetKey: JSON.stringify(mobileRequest),
    },
    submitting: mutation.isPending,
    submissionError: mutation.isError,
  };
};

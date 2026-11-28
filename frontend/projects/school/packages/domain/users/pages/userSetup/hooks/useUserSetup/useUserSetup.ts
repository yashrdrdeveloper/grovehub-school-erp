import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';

import { usersApi } from '@Business/domain/users';
import { useMobileViewport } from '@Components/hooks/useMobileViewport';
import { PaginationType } from '@Components/components/Pagination';
import type { PaginationProps } from '@Components/components/Pagination';
import type { TableActions } from '@Components/components/Table';
import type { Account, CreateAccount } from '@Business/domain/users';

import { useUserSetupQueries } from '../useUserSetupQueries';
import { useUserSetupPagination } from '../useUserSetupPagination';
import { UserAction } from '../../constants/userSetup.actions';
import { USER_SETUP_DEFAULT_SORT } from '../../constants/userSetup.sort';
import type { UserSetupDialog, UserSetupQueryState } from '../../types/userSetup.state';

export const useUserSetup = () => {
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
  const {
    scope,
    query,
    mobileQuery,
    mobileRequest,
    activeQuery,
    getData,
    getTotal,
    getError,
    refresh,
  } = useUserSetupQueries(state, mobile);
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
  const paginate = useUserSetupPagination({ query, setState, resetQuery });
  const actions: TableActions<Account> = {
    SEARCH: (searchText) => resetQuery({ searchText }),
    FILTER: (filters) => resetQuery({ filters }),
    SORT: (sorting) => resetQuery({ sorting }),
    QUICK_FILTER: (quickFilterIds) => resetQuery({ quickFilterIds }),
    PAGINATION: paginate,
    REFRESH: refresh,
    RETRY: refresh,
    LOAD_MORE: () => {
      if (mobile && mobileQuery.hasNextPage && !mobileQuery.isFetching)
        void mobileQuery.fetchNextPage({ cancelRefetch: false });
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
    if (dialog.action === UserAction.ACTIVATE)
      mutation.mutate(() => usersApi.activate(row.accountId));
    if (dialog.action === UserAction.SUSPEND)
      mutation.mutate(() => usersApi.suspend(row.accountId));
    if (dialog.action === UserAction.DEACTIVATE)
      mutation.mutate(() => usersApi.deactivate(row.accountId));
    if (dialog.action === UserAction.RESEND)
      mutation.mutate(() => usersApi.resendInvitation(row.accountId));
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
    state,
    actions,
    dialog,
    closeDialog,
    create,
    edit,
    confirm,
    pagination,
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

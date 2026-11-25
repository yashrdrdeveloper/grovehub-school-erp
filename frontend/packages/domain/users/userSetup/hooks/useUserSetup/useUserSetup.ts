import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

import { usersApi } from '@Business/domain/users';
import { useAuth } from '@Domain/auth';
import { PaginationType } from '@Components/components/Pagination';
import type { PaginationAction, PaginationProps } from '@Components/components/Pagination';
import type { TableActions } from '@Components/components/Table';
import type { Account, CreateUser, UpdateUser, UpdateUserProfile } from '@Business/domain/users';

import { UserAction } from '../../constants/userSetup.actions';
import { USER_SETUP_DEFAULT_SORT } from '../../constants/userSetup.sort';
import { getUserSetupSearchRequest } from '../../helpers/userSetup.getSearchRequest';
import type { UserSetupDialog, UserSetupQueryState } from '../../types/userSetup.state';

export const useUserSetup = () => {
  const { account } = useAuth();
  const queryClient = useQueryClient();
  const [state, setState] = useState<UserSetupQueryState>({
    searchText: '',
    filters: [],
    sorting: USER_SETUP_DEFAULT_SORT,
    quickFilterIds: ['ALL'],
    pageSize: 25,
    cursors: [undefined],
  });
  const [dialog, setDialog] = useState<UserSetupDialog | null>(null);
  const request = getUserSetupSearchRequest(state);
  const scope = ['user-setup', account?.accountId];
  const query = useQuery({
    queryKey: [...scope, request],
    queryFn: ({ signal }) => usersApi.search(request, signal),
    enabled: Boolean(account),
  });
  const roles = useQuery({
    queryKey: [...scope, 'roles'],
    queryFn: ({ signal }) => usersApi.roles(signal),
    enabled: Boolean(account && query.isSuccess),
    staleTime: 300000,
  });
  const rowId = dialog?.row?.accountId ?? '';
  const profile = useQuery({
    queryKey: [...scope, 'profile', rowId],
    queryFn: ({ signal }) => usersApi.profile(rowId, signal),
    enabled: Boolean(rowId && (dialog?.action === UserAction.EDIT || dialog?.action === UserAction.VIEW)),
  });
  const profileMutation = useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateUserProfile }) => usersApi.updateProfile(id, input),
    onSuccess: async (value, variables) => {
      queryClient.setQueryData([...scope, 'profile', variables.id], value);
      await queryClient.invalidateQueries({ queryKey: scope });
    },
  });
  const mutation = useMutation({
    mutationFn: async (operation: () => Promise<unknown>) => operation(),
    onSuccess: async () => {
      setDialog(null);
      await queryClient.invalidateQueries({ queryKey: scope });
    },
  });
  const openDialog = (next: UserSetupDialog) => {
    mutation.reset();
    profileMutation.reset();
    setDialog(next);
  };
  const closeDialog = () => {
    if (!mutation.isPending && !profileMutation.isPending) setDialog(null);
  };
  const resetQuery = (change: Partial<UserSetupQueryState>) => {
    setState((current) => ({ ...current, ...change, cursors: [undefined] }));
  };
  const paginate = (action: PaginationAction) => {
    if (action.type === 'PAGE_SIZE_CHANGE') {
      resetQuery({ pageSize: action.pageSize });
      return;
    }
    if (action.type === 'PREVIOUS') {
      setState((current) => {
        if (current.cursors.length < 2) return current;
        return { ...current, cursors: current.cursors.slice(0, -1) };
      });
    }
    const nextCursor = query.data?.pageInfo.nextCursor;
    if (action.type === 'NEXT' && query.data?.pageInfo.hasNext && nextCursor) {
      setState((current) => ({ ...current, cursors: [...current.cursors, nextCursor] }));
    }
  };
  const actions: TableActions<Account> = {
    SEARCH: (searchText) => resetQuery({ searchText }),
    FILTER: (filters) => resetQuery({ filters }),
    SORT: (sorting) => resetQuery({ sorting }),
    QUICK_FILTER: (quickFilterIds) => resetQuery({ quickFilterIds }),
    PAGINATION: paginate,
    REFRESH: () => { void query.refetch(); },
    RETRY: () => { void query.refetch(); },
    HEADER_ACTION: (id) => {
      if (id === UserAction.CREATE) openDialog({ action: UserAction.CREATE });
    },
    ROW_CLICK: (row) => openDialog({ action: UserAction.VIEW, row }),
    ROW_ACTION: ({ id, row }) => {
      if (id === UserAction.VIEW) openDialog({ action: UserAction.VIEW, row });
      if (id === UserAction.EDIT) openDialog({ action: UserAction.EDIT, row });
      if (id === UserAction.ACTIVATE) openDialog({ action: UserAction.ACTIVATE, row });
      if (id === UserAction.SUSPEND) openDialog({ action: UserAction.SUSPEND, row });
      if (id === UserAction.DEACTIVATE) openDialog({ action: UserAction.DEACTIVATE, row });
      if (id === UserAction.RESEND) openDialog({ action: UserAction.RESEND, row });
    },
  };
  const create = (input: CreateUser, key: string) =>
    mutation.mutate(() => usersApi.create(input, key));
  const edit = (row: Account, input: UpdateUser) =>
    mutation.mutate(() => usersApi.update(row.accountId, input));
  const editProfile = (row: Account, input: UpdateUserProfile) =>
    profileMutation.mutate({ id: row.accountId, input });
  const confirm = () => {
    const row = dialog?.row;
    if (!row || !dialog) return;
    if (dialog.action === UserAction.ACTIVATE) mutation.mutate(() => usersApi.activate(row.accountId));
    if (dialog.action === UserAction.SUSPEND) mutation.mutate(() => usersApi.suspend(row.accountId));
    if (dialog.action === UserAction.DEACTIVATE) mutation.mutate(() => usersApi.deactivate(row.accountId));
    if (dialog.action === UserAction.RESEND) mutation.mutate(() => usersApi.resendInvitation(row.accountId));
  };
  const pagination = {
    type: PaginationType.CURSOR,
    pageSize: state.pageSize,
    currentPage: state.cursors.length,
    hasPrevious: state.cursors.length > 1,
    hasNext: Boolean(query.data?.pageInfo.hasNext && query.data.pageInfo.nextCursor),
    loading: query.isFetching,
  } satisfies Omit<Extract<PaginationProps, { type: typeof PaginationType.CURSOR }>, 'onChange'>;

  return {
    state, actions, dialog, closeDialog, create, edit, editProfile, confirm, pagination,
    roles: roles.data ?? [],
    rolesLoading: roles.isPending,
    rolesError: roles.isError,
    retryRoles: () => { void roles.refetch(); },
    profile: profile.data,
    profileLoading: profile.isFetching,
    profileError: profile.isError,
    retryProfile: () => { void profile.refetch(); },
    profileSubmitting: profileMutation.isPending,
    profileSubmissionError: profileMutation.isError,
    profileSaved: profileMutation.isSuccess,
    preferencesKey: account?.accountId && 'groovehub:user-setup:' + account.accountId + ':columns:v1',
    data: query.data?.items ?? [],
    loading: query.isPending,
    refreshing: query.isFetching,
    error: query.isError,
    submitting: mutation.isPending,
    submissionError: mutation.isError,
  };
};

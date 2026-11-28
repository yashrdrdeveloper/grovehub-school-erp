import { useQuery, useInfiniteQuery } from '@tanstack/react-query';
import { usersApi } from '@Business/domain/users';
import { useAuth } from '@Domain/auth';

import { getUserSetupNextPage } from '../../helpers/userSetup.getNextPage';
import { getUserSetupMobileRows } from '../../helpers/userSetup.getMobileRows';
import { getUserSetupSearchRequest } from '../../helpers/userSetup.getSearchRequest';
import type { UserSetupQueryState } from '../../types/userSetup.state';

export const useUserSetupQueries = (state: UserSetupQueryState, mobile: boolean) => {
  const { account } = useAuth();
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
    queryFn: ({ pageParam }) =>
      usersApi.search({ ...mobileRequest, pageInfo: { ...mobileRequest.pageInfo, ...pageParam } }),
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
    if (mobile) {
      void mobileQuery.refetch();
      return;
    }
    void query.refetch();
  };
  return {
    scope,
    query,
    mobileQuery,
    mobileRequest,
    activeQuery,
    getData,
    getTotal,
    getError,
    refresh,
  };
};

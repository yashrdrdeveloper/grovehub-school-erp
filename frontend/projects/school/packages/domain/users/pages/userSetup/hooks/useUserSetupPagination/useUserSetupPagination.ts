import type { PaginationAction } from '@Components/components/Pagination';
import type { UserSetupPaginationProps } from './types/userSetupPagination.props';

export const useUserSetupPagination = ({
  query,
  setState,
  resetQuery,
}: UserSetupPaginationProps) => {
  const paginate = (action: PaginationAction) => {
    if (action.type === 'PAGE_SIZE_CHANGE') {
      resetQuery({ pageSize: action.pageSize });
      return;
    }
    if (query.isFetching) return;
    if (action.type === 'PAGE_CHANGE') {
      const totalPages = query.data?.pageInfo.totalPages;
      if (
        !Number.isSafeInteger(action.page) ||
        action.page < 1 ||
        (totalPages !== undefined && action.page > totalPages)
      )
        return;
      setState((current) => {
        const nextCursor = query.data?.pageInfo.nextCursor;
        if (action.page === current.pageNumber + 1 && nextCursor) {
          return {
            ...current,
            pageNumber: action.page,
            cursors: { ...current.cursors, [action.page]: nextCursor },
          };
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
  return paginate;
};

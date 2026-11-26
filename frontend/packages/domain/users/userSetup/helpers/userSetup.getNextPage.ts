import type { AccountSearchResult } from '@Business/domain/users';

type PageParam = { pageNumber: number; cursor: string | undefined };

export const getUserSetupNextPage = (
  last: AccountSearchResult,
  _pages: readonly AccountSearchResult[],
  previous: PageParam,
  params: readonly PageParam[],
): PageParam | undefined => {
  const cursor = last.pageInfo.nextCursor;
  if (!last.pageInfo.hasNext || !cursor || params.some((param) => param.cursor === cursor))
    return undefined;
  return { pageNumber: previous.pageNumber + 1, cursor };
};

import { getUserSetupNextPage } from '../userSetup.getNextPage';
import type { AccountSearchResult } from '@Business/domain/users';

describe('getUserSetupNextPage', () => {
  const previous = { pageNumber: 1, cursor: undefined };
  const page: AccountSearchResult = {
    items: [],
    pageInfo: { rows: 25, hasNext: true, nextCursor: 'next' },
  };
  it('advances using the backend cursor', () => {
    expect(getUserSetupNextPage(page, [page], previous, [previous])).toEqual({
      pageNumber: 2,
      cursor: 'next',
    });
  });
  it('stops at the last page or when the cursor is absent or repeated', () => {
    expect(
      getUserSetupNextPage({ ...page, pageInfo: { rows: 0, hasNext: false } }, [page], previous, [
        previous,
      ]),
    ).toBeUndefined();
    expect(
      getUserSetupNextPage({ ...page, pageInfo: { rows: 0, hasNext: true } }, [page], previous, [
        previous,
      ]),
    ).toBeUndefined();
    expect(
      getUserSetupNextPage(page, [page], previous, [{ pageNumber: 2, cursor: 'next' }]),
    ).toBeUndefined();
  });
});

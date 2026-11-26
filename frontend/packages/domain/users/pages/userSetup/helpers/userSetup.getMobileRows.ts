import { _flatMap } from '@Base/utilities/flatMap';
import { _map } from '@Base/utilities/map';
import type { AccountSearchResult } from '@Business/domain/users';

export const getUserSetupMobileRows = (pages: readonly AccountSearchResult[]) => {
  const rows = _flatMap(pages, (page) => page.items);
  const entries = _map(rows, (row) => [row.accountId, row] as const);
  return Array.from(new Map(entries).values());
};

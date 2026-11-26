import { _flatMap } from '@Base/utilities/flatMap';
import { _map } from '@Base/utilities/map';
import type { AccountSearchResult } from '@Business/domain/users';

export const getUserSetupMobileRows = (pages: readonly AccountSearchResult[]) =>
  Array.from(
    new Map(_map(_flatMap(pages, (page) => page.items), (row) => [row.accountId, row] as const)).values(),
  );

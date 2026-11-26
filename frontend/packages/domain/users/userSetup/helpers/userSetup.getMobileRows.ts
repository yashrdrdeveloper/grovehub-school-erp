import type { AccountSearchResult } from '@Business/domain/users';

export const getUserSetupMobileRows = (pages: readonly AccountSearchResult[]) =>
  Array.from(
    new Map(pages.flatMap((page) => page.items).map((row) => [row.accountId, row])).values(),
  );

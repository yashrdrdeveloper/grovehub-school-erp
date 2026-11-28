import type { Account, AccountSearchResult } from '@Business/domain/users';
import { getUserSetupMobileRows } from '../userSetup.getMobileRows';

const account: Account = {
  accountId: 'account-1',
  personGid: 'person-1',
  username: 'teacher01',
  email: 'teacher@example.com',
  status: 'ACTIVE',
  version: 1,
  createdAt: '2026-11-01T10:00:00Z',
  updatedAt: '2026-11-02T10:00:00Z',
};

describe('getUserSetupMobileRows', () => {
  it('retains page order while replacing duplicate account IDs with the latest row', () => {
    const second = { ...account, accountId: 'account-2' };
    const updated = { ...account, username: 'updated', version: 2 };
    const pages: AccountSearchResult[] = [
      { items: [account, second], pageInfo: { rows: 2, hasNext: true } },
      { items: [updated], pageInfo: { rows: 1, hasNext: false } },
    ];

    expect(getUserSetupMobileRows(pages)).toEqual([updated, second]);
    expect(getUserSetupMobileRows([])).toEqual([]);
  });
});

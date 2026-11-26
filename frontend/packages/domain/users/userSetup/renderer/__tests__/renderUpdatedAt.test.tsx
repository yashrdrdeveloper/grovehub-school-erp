import { render, screen } from '@testing-library/react';
import type { Account } from '@Business/domain/users';
import { getFormattedDate } from '@Base/localization/getFormattedDate';
import { renderUpdatedAt } from '../userSetup.timestamps';

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

describe('renderUpdatedAt', () => {
  it('renders the update timestamp in the shared localized format', () => {
    render(renderUpdatedAt(account));

    expect(screen.getByText(getFormattedDate(account.updatedAt, { includeTime: false }))).toBeTruthy();
    expect(screen.getByText(getFormattedDate(account.updatedAt, { timeOnly: true }))).toBeTruthy();
  });
});

import { render, screen } from '@testing-library/react';
import type { Account } from '@Business/domain/users';
import { getFormattedDate } from '@Base/localization/getFormattedDate';
import { renderCreatedAt } from '../userSetup.timestamps';

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

describe('renderCreatedAt', () => {
  it('renders the creation timestamp in the shared localized format', () => {
    render(renderCreatedAt(account));

    expect(screen.getByText(getFormattedDate(account.createdAt, { includeTime: false }))).toBeTruthy();
    expect(screen.getByText(getFormattedDate(account.createdAt, { timeOnly: true }))).toBeTruthy();
  });
});

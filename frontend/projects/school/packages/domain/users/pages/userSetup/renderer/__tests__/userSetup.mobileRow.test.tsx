import { render, screen } from '@testing-library/react';
import type { Account } from '@Business/domain/users';
import { renderUserMobileRow } from '../userSetup.mobileRow';

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

describe('renderUserMobileRow', () => {
  it('renders the username, email and plain status without a row action', () => {
    render(renderUserMobileRow(account));

    expect(screen.getByText(account.username!)).toBeTruthy();
    expect(screen.getByText(account.email)).toBeTruthy();
    expect(screen.getByText('active').getAttribute('data-plain')).toBe('true');
    expect(screen.queryByRole('button')).toBeNull();
  });

  it('falls back to the email when the username is absent', () => {
    render(renderUserMobileRow({ ...account, username: null }));

    expect(screen.getAllByText(account.email)).toHaveLength(2);
  });
});

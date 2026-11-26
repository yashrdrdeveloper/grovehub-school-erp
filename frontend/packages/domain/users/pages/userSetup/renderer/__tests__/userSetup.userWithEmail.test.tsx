import { render, screen } from '@testing-library/react';
import type { Account } from '@Business/domain/users';
import { renderUserWithEmail } from '../userSetup.userWithEmail';

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

describe('renderUserWithEmail', () => {
  it('renders the username and email', () => {
    render(renderUserWithEmail(account));

    expect(screen.getByText(account.username!)).toBeTruthy();
    expect(screen.getByText(account.email)).toBeTruthy();
  });

  it('uses the existing unknown-user label when the username is absent', () => {
    render(renderUserWithEmail({ ...account, username: null }));

    expect(screen.getByText('Unknown user')).toBeTruthy();
    expect(screen.getByText(account.email)).toBeTruthy();
  });
});

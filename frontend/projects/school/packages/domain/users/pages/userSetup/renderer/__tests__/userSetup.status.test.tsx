import { render, screen } from '@testing-library/react';
import type { Account } from '@Business/domain/users';
import { renderUserStatus } from '../userSetup.status';

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

describe('renderUserStatus', () => {
  it('renders the backend status with the shared semantic tone', () => {
    render(renderUserStatus(account));

    expect(screen.getByText('active').getAttribute('data-tone')).toBe('positive');
  });

  it('preserves pending setup status', () => {
    render(renderUserStatus({ ...account, status: 'PENDING_SETUP' }));

    expect(screen.getByText('pending setup')).toBeTruthy();
  });
});

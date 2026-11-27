import { render, screen } from '@testing-library/react';
import type { Role } from '@Business/domain/users/rbac/types/contracts';
import { renderRbacMobileRow } from '../rbac.mobileRow';

const role: Role = {
  roleId: 'role-1',
  name: 'Teacher',
  description: 'Classroom access',
  status: 'ACTIVE',
  version: 1,
  personas: [],
};

describe('renderRbacMobileRow', () => {
  it('renders the row summary without an independent action', () => {
    render(renderRbacMobileRow(role));

    expect(screen.getByText(role.name)).toBeTruthy();
    expect(screen.getByText(role.status)).toBeTruthy();
    expect(screen.getByText(role.description)).toBeTruthy();
    expect(screen.queryByRole('button')).toBeNull();
  });
});

import type { UserRole } from '@Business/domain/users';

export const getUserRoleLabel = (roles: readonly UserRole[], id?: string | null): string => {
  if (!id) return __('Unassigned');
  return roles.find((role) => role.id === id)?.name ?? id;
};

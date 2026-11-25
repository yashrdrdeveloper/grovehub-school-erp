import type { UserRole, UserRoleAssignment } from '@Business/domain/users';

export const isRoleAssignmentValid = (
  value: UserRoleAssignment,
  roles: readonly UserRole[],
): boolean => {
  const known = new Set(roles.map((role) => role.id));
  if (!known.has(value.primaryRoleId)) return false;
  const all = [value.primaryRoleId, ...value.secondaryRoleIds];
  return new Set(all).size === all.length && all.every((id) => known.has(id));
};

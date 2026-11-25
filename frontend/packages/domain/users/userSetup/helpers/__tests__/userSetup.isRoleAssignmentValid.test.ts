import { isRoleAssignmentValid } from '../userSetup.isRoleAssignmentValid';

describe('isRoleAssignmentValid', () => {
  const roles = [{ id: 'A', name: 'A' }, { id: 'B', name: 'B' }];
  it('requires a known primary and unique secondary roles', () => {
    expect(isRoleAssignmentValid({ primaryRoleId: 'A', secondaryRoleIds: ['B'] }, roles)).toBe(true);
    expect(isRoleAssignmentValid({ primaryRoleId: '', secondaryRoleIds: [] }, roles)).toBe(false);
    expect(isRoleAssignmentValid({ primaryRoleId: 'A', secondaryRoleIds: ['A'] }, roles)).toBe(false);
    expect(isRoleAssignmentValid({ primaryRoleId: 'A', secondaryRoleIds: ['B', 'B'] }, roles)).toBe(false);
  });
});

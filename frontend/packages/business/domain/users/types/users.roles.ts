export interface UserRoleAssignment {
  readonly primaryRoleId: string;
  readonly secondaryRoleIds: readonly string[];
}

export interface UserRole {
  readonly id: string;
  readonly name: string;
}

import type { Account as AccountEntity } from '@Business/auth';
import type { CreateAccount } from '@Business/accounts';

import type { UserRoleAssignment } from './users.roles';

export interface UserProfile {
  readonly firstName: string;
  readonly middleName: string | null;
  readonly lastName: string;
  readonly phoneNumber: string | null;
}

export interface Account extends AccountEntity {
  readonly primaryRoleId?: string | null;
  readonly secondaryRoleIds?: readonly string[];
  readonly profile?: UserProfile;
}

export interface CreateUser extends CreateAccount {
  readonly roleAssignment: UserRoleAssignment;
}

export interface UpdateUser {
  readonly username?: string;
  readonly roleAssignment: UserRoleAssignment;
  readonly version: number;
}

export type UpdateUserProfile = Omit<UserProfile, 'phoneNumber'>;

import type {
  Account, CreateUser, UpdateUser, UpdateUserProfile, UserRole, UserProfile,
} from '@Business/domain/users';

export interface UserFormProps {
  readonly row?: Account;
  readonly roles: readonly UserRole[];
  readonly rolesLoading: boolean;
  readonly rolesError: boolean;
  readonly retryRoles: () => void;
  readonly profile?: UserProfile;
  readonly profileLoading: boolean;
  readonly profileError: boolean;
  readonly retryProfile: () => void;
  readonly onClose: () => void;
  readonly onCreate: (input: CreateUser, key: string) => void;
  readonly onEdit: (row: Account, input: UpdateUser) => void;
  readonly onEditProfile: (row: Account, input: UpdateUserProfile) => void;
  readonly submitting: boolean;
  readonly error: boolean;
  readonly profileSubmitting: boolean;
  readonly profileSubmissionError: boolean;
  readonly profileSaved: boolean;
}

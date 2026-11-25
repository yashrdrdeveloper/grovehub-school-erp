import type { UserProfile, UpdateUserProfile } from '@Business/domain/users';

export interface EditProfileProps {
  readonly profile: UserProfile;
  readonly onSave: (input: UpdateUserProfile) => void;
  readonly submitting: boolean;
  readonly error: boolean;
  readonly saved: boolean;
}

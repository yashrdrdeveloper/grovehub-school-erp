import type { SelectOption } from '@Components/components/Select';
import type { UserRoleAssignment } from '@Business/domain/users';

export interface UserRoleFieldsProps {
  readonly options: readonly SelectOption[];
  readonly value: UserRoleAssignment;
  readonly onChange: (value: UserRoleAssignment) => void;
  readonly disabled?: boolean;
}
